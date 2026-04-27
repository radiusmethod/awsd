package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

// State is the persisted contents of ~/.awsd.
//
// Profile: empty string means "default" — the wrapper unsets AWS_PROFILE.
// Region:  paired with RegionSet to distinguish three states:
//   - RegionSet=false        → no `region=` line; wrapper leaves AWS_REGION alone.
//   - RegionSet=true, ""     → `region=` with empty value; wrapper unsets AWS_REGION.
//   - RegionSet=true, "..."  → wrapper exports AWS_REGION (and AWS_DEFAULT_REGION).
type State struct {
	Profile   string
	Region    string
	RegionSet bool
}

func statePath(loc string) string {
	return filepath.Join(loc, ".awsd")
}

func ReadState(loc string) (State, error) {
	data, err := os.ReadFile(statePath(loc))
	if err != nil {
		if os.IsNotExist(err) {
			return State{}, nil
		}
		return State{}, err
	}
	text := strings.TrimRight(string(data), "\n")
	if text == "" {
		return State{}, nil
	}

	hasKV := false
	for _, line := range strings.Split(text, "\n") {
		if strings.Contains(line, "=") {
			hasKV = true
			break
		}
	}
	if !hasKV {
		// Legacy single-line: whole file is a profile name.
		return State{Profile: strings.TrimSpace(text)}, nil
	}

	var s State
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		idx := strings.Index(line, "=")
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		switch key {
		case "profile":
			s.Profile = val
		case "region":
			s.Region = val
			s.RegionSet = true
		}
	}
	return s, nil
}

func WriteState(s State, loc string) error {
	path := statePath(loc)
	if err := TouchFile(path); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("profile=")
	b.WriteString(s.Profile)
	b.WriteString("\n")
	if s.RegionSet {
		b.WriteString("region=")
		b.WriteString(s.Region)
		b.WriteString("\n")
	}
	if err := os.WriteFile(path, []byte(b.String()), 0644); err != nil {
		log.Fatal(err)
	}
	return nil
}

// WriteFile sets the active profile, preserving any existing region.
// "default" is stored as an empty profile so the wrapper unsets AWS_PROFILE.
func WriteFile(config, loc string) error {
	s, err := ReadState(loc)
	if err != nil {
		return err
	}
	if config == "default" {
		s.Profile = ""
	} else {
		s.Profile = config
	}
	return WriteState(s, loc)
}

// WriteRegion sets the active region, preserving the existing profile.
func WriteRegion(region, loc string) error {
	s, err := ReadState(loc)
	if err != nil {
		return err
	}
	s.Region = region
	s.RegionSet = true
	return WriteState(s, loc)
}

// UnsetRegion writes an explicit empty region so the wrapper unsets AWS_REGION.
func UnsetRegion(loc string) error {
	s, err := ReadState(loc)
	if err != nil {
		return err
	}
	s.Region = ""
	s.RegionSet = true
	return WriteState(s, loc)
}

func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func CheckError(err error) {
	if err.Error() == "^D" {
		// https://github.com/manifoldco/promptui/issues/179
		log.Fatalf("<Del> not supported")
	} else if err.Error() == "^C" {
		os.Exit(1)
	} else {
		log.Fatal(err)
	}
}

func GetHomeDir() (string, error) {
	if homeDir := os.Getenv("HOME"); homeDir != "" {
		return homeDir, nil
	}
	if homeDir, err := os.UserHomeDir(); err == nil {
		return homeDir, nil
	}
	return "", fmt.Errorf("error getting user home directory: $HOME is not defined and os.UserHomeDir() failed")
}

func GetProfileFileLocation() string {
	homeDir, err := GetHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(homeDir, ".aws")
}

func GetCurrentProfileFile() string {
	homeDir, err := GetHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return GetEnv("AWS_CONFIG_FILE", filepath.Join(homeDir, ".aws/config"))
}

func IsDirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func AppendIfNotExists(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}

func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
