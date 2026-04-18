package settings

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"
)

func mkraw(t *testing.T, kv map[string]any) raw {
	t.Helper()
	r := raw{}
	for k, v := range kv {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		r[k] = b
	}
	return r
}

func TestDefaultSettingsValues(t *testing.T) {
	d := Default()
	if d.DefaultDurationSecs != 300 || d.DefaultIncrementSecs != 300 || d.SecondaryIncrementSecs != 60 || d.TertiaryIncrementSecs != 3600 || d.MiniWindowWidth != 70 || d.MiniWindowHeight != 70 {
		t.Fatalf("unexpected defaults: %+v", d)
	}
}

func TestValidateAcceptsBoundaries(t *testing.T) {
	if err := Validate(AppSettings{60, 60, 60, 60, 50, 50}); err != nil {
		t.Fatal(err)
	}
	if err := Validate(AppSettings{10800, 3600, 3600, 86400, 400, 400}); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRejectsDurationBelowMin(t *testing.T) {
	err := Validate(AppSettings{59, 300, 60, 3600, 70, 70})
	if err == nil || !strings.Contains(err.Error(), "default_duration_secs") {
		t.Fatalf("bad err: %v", err)
	}
}

func TestValidateRejectsDurationAboveMax(t *testing.T) {
	if err := Validate(AppSettings{10801, 300, 60, 3600, 70, 70}); err == nil {
		t.Fatal("expected err")
	}
}

func TestValidateRejectsIncrementOutOfRange(t *testing.T) {
	err := Validate(AppSettings{300, 3601, 60, 3600, 70, 70})
	if err == nil || !strings.Contains(err.Error(), "default_increment_secs") {
		t.Fatalf("bad err: %v", err)
	}
}

func TestValidateRejectsSecondaryIncrementOutOfRange(t *testing.T) {
	err := Validate(AppSettings{300, 300, 59, 3600, 70, 70})
	if err == nil || !strings.Contains(err.Error(), "secondary_increment_secs") {
		t.Fatalf("bad err: %v", err)
	}
}

func TestValidateRejectsMiniWindowWidthOutOfRange(t *testing.T) {
	err := Validate(AppSettings{300, 300, 60, 3600, 49, 70})
	if err == nil || !strings.Contains(err.Error(), "mini_window_width") {
		t.Fatalf("bad err: %v", err)
	}
	err = Validate(AppSettings{300, 300, 60, 3600, 401, 70})
	if err == nil || !strings.Contains(err.Error(), "mini_window_width") {
		t.Fatalf("bad err: %v", err)
	}
}

func TestValidateRejectsMiniWindowHeightOutOfRange(t *testing.T) {
	err := Validate(AppSettings{300, 300, 60, 3600, 70, 49})
	if err == nil || !strings.Contains(err.Error(), "mini_window_height") {
		t.Fatalf("bad err: %v", err)
	}
	err = Validate(AppSettings{300, 300, 60, 3600, 70, 401})
	if err == nil || !strings.Contains(err.Error(), "mini_window_height") {
		t.Fatalf("bad err: %v", err)
	}
}

func TestParseSettingsReadsMiniWindowDimensions(t *testing.T) {
	r := mkraw(t, map[string]any{"mini_window_width": 120, "mini_window_height": 90})
	s := ParseSettings(r)
	if s.MiniWindowWidth != 120 || s.MiniWindowHeight != 90 {
		t.Fatalf("got %+v", s)
	}
}

func TestParseSettingsFallsBackOnMiniWindowOutOfRange(t *testing.T) {
	r := mkraw(t, map[string]any{"mini_window_width": 10, "mini_window_height": 10})
	s := ParseSettings(r)
	if s.MiniWindowWidth != Default().MiniWindowWidth || s.MiniWindowHeight != Default().MiniWindowHeight {
		t.Fatalf("expected default fallback, got %+v", s)
	}
}

func TestParseSettingsMigratesLegacyMiniWindowSize(t *testing.T) {
	r := mkraw(t, map[string]any{"mini_window_size": 150})
	s := ParseSettings(r)
	if s.MiniWindowWidth != 150 || s.MiniWindowHeight != 150 {
		t.Fatalf("expected legacy migration, got %+v", s)
	}
}

func TestParseSettingsPrefersNewFieldsOverLegacy(t *testing.T) {
	r := mkraw(t, map[string]any{
		"mini_window_size":   150,
		"mini_window_width":  200,
		"mini_window_height": 100,
	})
	s := ParseSettings(r)
	if s.MiniWindowWidth != 200 || s.MiniWindowHeight != 100 {
		t.Fatalf("new fields should win, got %+v", s)
	}
}

func TestParseSettingsEmptyReturnsDefaults(t *testing.T) {
	s := ParseSettings(raw{})
	if s != Default() {
		t.Fatalf("expected defaults, got %+v", s)
	}
}

func TestParseSettingsReadsValidStoredValues(t *testing.T) {
	r := mkraw(t, map[string]any{
		"default_duration_secs":    600,
		"default_increment_secs":   120,
		"secondary_increment_secs": 90,
	})
	s := ParseSettings(r)
	if s.DefaultDurationSecs != 600 || s.DefaultIncrementSecs != 120 || s.SecondaryIncrementSecs != 90 {
		t.Fatalf("%+v", s)
	}
}

func TestParseSettingsFallsBackWhenOutOfRange(t *testing.T) {
	r := mkraw(t, map[string]any{
		"default_duration_secs":    1,
		"default_increment_secs":   99999,
		"secondary_increment_secs": 0,
	})
	s := ParseSettings(r)
	if s != Default() {
		t.Fatalf("expected defaults fallback, got %+v", s)
	}
}

func TestParseSettingsFallsBackOnWrongType(t *testing.T) {
	r := mkraw(t, map[string]any{"default_duration_secs": "not a number"})
	s := ParseSettings(r)
	if s.DefaultDurationSecs != Default().DefaultDurationSecs {
		t.Fatal("expected default fallback")
	}
}

func TestParseTimerNamesEmpty(t *testing.T) {
	if n := ParseTimerNames(raw{}); len(n) != 0 {
		t.Fatal("expected empty")
	}
}

func TestParseTimerNamesValid(t *testing.T) {
	r := mkraw(t, map[string]any{"timer_names": []string{"A", "B", "C"}})
	n := ParseTimerNames(r)
	if len(n) != 3 || n[0] != "A" || n[2] != "C" {
		t.Fatalf("%v", n)
	}
}

func TestParseTimerNamesRejectsWrongType(t *testing.T) {
	r := mkraw(t, map[string]any{"timer_names": "not an array"})
	if n := ParseTimerNames(r); len(n) != 0 {
		t.Fatal("expected empty")
	}
}

func TestAppSettingsSerializesWithExpectedFieldNames(t *testing.T) {
	b, _ := json.Marshal(Default())
	var m map[string]any
	json.Unmarshal(b, &m)
	for _, k := range []string{"default_duration_secs", "default_increment_secs", "secondary_increment_secs", "mini_window_width", "mini_window_height"} {
		if _, ok := m[k]; !ok {
			t.Fatalf("missing key %s", k)
		}
	}
}

func TestAppSettingsRoundtrip(t *testing.T) {
	in := AppSettings{600, 120, 90, 3600, 120, 80}
	b, _ := json.Marshal(in)
	var out AppSettings
	json.Unmarshal(b, &out)
	if out != in {
		t.Fatalf("roundtrip mismatch: %+v", out)
	}
}

func TestStoreFileNameIsStable(t *testing.T) {
	if StoreFileName != "settings.json" {
		t.Fatal("store filename drifted")
	}
}

func TestStoreRoundtripOnDisk(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	s, err := OpenStoreAt(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SaveSettings(AppSettings{600, 120, 90, 3600, 120, 80}); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveTimerNames([]string{"X", "Y"}); err != nil {
		t.Fatal(err)
	}
	s2, _ := OpenStoreAt(path)
	got := s2.LoadSettings()
	if got != (AppSettings{600, 120, 90, 3600, 120, 80}) {
		t.Fatalf("%+v", got)
	}
	names := s2.LoadTimerNames()
	if len(names) != 2 || names[0] != "X" {
		t.Fatalf("%v", names)
	}
}

func TestSaveSettingsStripsLegacyMiniWindowSize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "settings.json")
	s, _ := OpenStoreAt(path)
	s.data["mini_window_size"] = json.RawMessage("150")
	if err := s.SaveSettings(Default()); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.data["mini_window_size"]; ok {
		t.Fatal("legacy key should be removed")
	}
}
