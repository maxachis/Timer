//go:build windows

package main

import (
	"context"
	"sync"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	wmGetMinMaxInfo = 0x0024
	gwlpWndProc     = ^uintptr(3) // -4 as uintptr
	defaultDPI      = 96
)

type point struct{ X, Y int32 }

type minMaxInfo struct {
	Reserved     point
	MaxSize      point
	MaxPosition  point
	MinTrackSize point
	MaxTrackSize point
}

var (
	user32              = windows.NewLazySystemDLL("user32.dll")
	procFindWindowW     = user32.NewProc("FindWindowW")
	procSetWindowLongW  = user32.NewProc("SetWindowLongPtrW")
	procCallWindowProcW = user32.NewProc("CallWindowProcW")
	procGetDpiForWindow = user32.NewProc("GetDpiForWindow")

	subclassOnce sync.Once
	origWndProc  uintptr
	minTrackW    int32 = 50
	minTrackH    int32 = 50
)

func installMinSizeOverride(_ context.Context, title string, minW, minH int32) {
	subclassOnce.Do(func() {
		hwnd := findWindow(title)
		if hwnd == 0 {
			return
		}
		minTrackW, minTrackH = minW, minH
		cb := syscall.NewCallback(wndProc)
		prev, _, _ := procSetWindowLongW.Call(uintptr(hwnd), uintptr(gwlpWndProc), cb)
		origWndProc = prev
	})
}

func findWindow(title string) windows.HWND {
	p, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return 0
	}
	r, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(p)))
	return windows.HWND(r)
}

func wndProc(hwnd windows.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procCallWindowProcW.Call(origWndProc, uintptr(hwnd), uintptr(msg), wParam, lParam)
	if msg == wmGetMinMaxInfo && lParam != 0 {
		dpi := getDPI(hwnd)
		scale := float32(dpi) / float32(defaultDPI)
		mmi := (*minMaxInfo)(unsafe.Pointer(lParam))
		mmi.MinTrackSize.X = int32(float32(minTrackW) * scale)
		mmi.MinTrackSize.Y = int32(float32(minTrackH) * scale)
	}
	return ret
}

func getDPI(hwnd windows.HWND) uint32 {
	r, _, _ := procGetDpiForWindow.Call(uintptr(hwnd))
	if r == 0 {
		return defaultDPI
	}
	return uint32(r)
}
