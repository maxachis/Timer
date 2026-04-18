package settings

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	DurationMin          uint64 = 60
	DurationMax          uint64 = 10800
	IncrementMin         uint64 = 60
	IncrementMax         uint64 = 3600
	TertiaryIncrementMax uint64 = 86400
	MiniWindowSizeMin    uint64 = 50
	MiniWindowSizeMax    uint64 = 400

	StoreFileName = "settings.json"
)

type AppSettings struct {
	DefaultDurationSecs    uint64 `json:"default_duration_secs"`
	DefaultIncrementSecs   uint64 `json:"default_increment_secs"`
	SecondaryIncrementSecs uint64 `json:"secondary_increment_secs"`
	TertiaryIncrementSecs  uint64 `json:"tertiary_increment_secs"`
	MiniWindowSize         uint64 `json:"mini_window_size"`
}

func Default() AppSettings {
	return AppSettings{
		DefaultDurationSecs:    300,
		DefaultIncrementSecs:   300,
		SecondaryIncrementSecs: 60,
		TertiaryIncrementSecs:  3600,
		MiniWindowSize:         70,
	}
}

func Validate(s AppSettings) error {
	if s.DefaultDurationSecs < DurationMin || s.DefaultDurationSecs > DurationMax {
		return errors.New("default_duration_secs must be between 60 and 10800")
	}
	if s.DefaultIncrementSecs < IncrementMin || s.DefaultIncrementSecs > IncrementMax {
		return errors.New("default_increment_secs must be between 60 and 3600")
	}
	if s.SecondaryIncrementSecs < IncrementMin || s.SecondaryIncrementSecs > IncrementMax {
		return errors.New("secondary_increment_secs must be between 60 and 3600")
	}
	if s.TertiaryIncrementSecs < IncrementMin || s.TertiaryIncrementSecs > TertiaryIncrementMax {
		return errors.New("tertiary_increment_secs must be between 60 and 86400")
	}
	if s.MiniWindowSize < MiniWindowSizeMin || s.MiniWindowSize > MiniWindowSizeMax {
		return errors.New("mini_window_size must be between 50 and 400")
	}
	return nil
}

// raw is the on-disk JSON map; values may be missing or out of range.
type raw map[string]json.RawMessage

func (r raw) u64(key string) (uint64, bool) {
	v, ok := r[key]
	if !ok {
		return 0, false
	}
	var n uint64
	if err := json.Unmarshal(v, &n); err != nil {
		return 0, false
	}
	return n, true
}

func ParseSettings(r raw) AppSettings {
	d := Default()
	out := d
	if v, ok := r.u64("default_duration_secs"); ok && v >= DurationMin && v <= DurationMax {
		out.DefaultDurationSecs = v
	}
	if v, ok := r.u64("default_increment_secs"); ok && v >= IncrementMin && v <= IncrementMax {
		out.DefaultIncrementSecs = v
	}
	if v, ok := r.u64("secondary_increment_secs"); ok && v >= IncrementMin && v <= IncrementMax {
		out.SecondaryIncrementSecs = v
	}
	if v, ok := r.u64("tertiary_increment_secs"); ok && v >= IncrementMin && v <= TertiaryIncrementMax {
		out.TertiaryIncrementSecs = v
	}
	if v, ok := r.u64("mini_window_size"); ok && v >= MiniWindowSizeMin && v <= MiniWindowSizeMax {
		out.MiniWindowSize = v
	}
	return out
}

func ParseTimerNames(r raw) []string {
	v, ok := r["timer_names"]
	if !ok {
		return nil
	}
	var names []string
	if err := json.Unmarshal(v, &names); err != nil {
		return nil
	}
	return names
}

// Store is a JSON-file-backed key/value store at settings.json.
type Store struct {
	path string
	data raw
}

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "Timer", StoreFileName), nil
}

func OpenStore() (*Store, error) {
	p, err := configPath()
	if err != nil {
		return nil, err
	}
	return OpenStoreAt(p)
}

func OpenStoreAt(path string) (*Store, error) {
	s := &Store{path: path, data: raw{}}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s, nil
		}
		return nil, err
	}
	_ = json.Unmarshal(b, &s.data)
	return s, nil
}

func (s *Store) Raw() raw { return s.data }

func (s *Store) LoadSettings() AppSettings { return ParseSettings(s.data) }

func (s *Store) LoadTimerNames() []string { return ParseTimerNames(s.data) }

func (s *Store) SaveSettings(a AppSettings) error {
	if err := Validate(a); err != nil {
		return err
	}
	s.set("default_duration_secs", a.DefaultDurationSecs)
	s.set("default_increment_secs", a.DefaultIncrementSecs)
	s.set("secondary_increment_secs", a.SecondaryIncrementSecs)
	s.set("tertiary_increment_secs", a.TertiaryIncrementSecs)
	s.set("mini_window_size", a.MiniWindowSize)
	return s.flush()
}

func (s *Store) SaveTimerNames(names []string) error {
	s.set("timer_names", names)
	return s.flush()
}

func (s *Store) set(key string, v any) {
	b, _ := json.Marshal(v)
	s.data[key] = b
}

func (s *Store) flush() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
