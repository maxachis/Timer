//go:build !windows

package main

import "context"

func installMinSizeOverride(_ context.Context, _ string, _, _ int32) {}
