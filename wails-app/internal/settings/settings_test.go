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
	if d.DefaultDurationSecs != 300 || d.DefaultIncrementSecs != 300 || d.SecondaryIncrementSecs != 60 || d.TertiaryIncrementSecs != 3600 {
		t.Fatalf("unexpected defaults: %+v", d)
	}
}

func TestValidateAcceptsBoundaries(t *testing.T) {
	if err := Validate(AppSettings{60, 60, 60, 60}); err != nil {
		t.Fatal(err)
	}
	if err := Validate(AppSettings{10800, 3600, 3600, 86400}); err != nil {
		t.Fatal(err)
	}
}

func TestValidateRejectsDurationBelowMin(t *testing.T) {
	err := Validate(AppSettings{59, 300, 60, 3600})
	if err == nil || !strings.Contains(err.Error(), "default_duration_secs") {
		t.Fatalf("bad err: %v", err)
	}
}

func TestValidateRejectsDurationAboveMax(t *testing.T) {
	if err := Validate(AppSettings{10801, 300, 60, 3600}); err == nil {
		t.Fatal("expected err")
	}
}

func TestValidateRejectsIncrementOutOfRange(t *testing.T) {
	err := Validate(AppSettings{300, 3601, 60, 3600})
	if err == nil || !strings.Contains(err.Error(), "default_increment_secs") {
		t.Fatalf("bad err: %v", err)
	}
}

func TestValidateRejectsSecondaryIncrementOutOfRange(t *testing.T) {
	err := Validate(AppSettings{300, 300, 59, 3600})
	if err == nil || !strings.Contains(err.Error(), "secondary_increment_secs") {
		t.Fatalf("bad err: %v", err)
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
	for _, k := range []string{"default_duration_secs", "default_increment_secs", "secondary_increment_secs"} {
		if _, ok := m[k]; !ok {
			t.Fatalf("missing key %s", k)
		}
	}
}

func TestAppSettingsRoundtrip(t *testing.T) {
	in := AppSettings{600, 120, 90, 3600}
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
	if err := s.SaveSettings(AppSettings{600, 120, 90, 3600}); err != nil {
		t.Fatal(err)
	}
	if err := s.SaveTimerNames([]string{"X", "Y"}); err != nil {
		t.Fatal(err)
	}
	s2, _ := OpenStoreAt(path)
	got := s2.LoadSettings()
	if got != (AppSettings{600, 120, 90, 3600}) {
		t.Fatalf("%+v", got)
	}
	names := s2.LoadTimerNames()
	if len(names) != 2 || names[0] != "X" {
		t.Fatalf("%v", names)
	}
}
