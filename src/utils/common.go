package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func TouchFile(name string) error {
	file, err := os.OpenFile(name, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func WriteFile(config, loc string) error {
	homeDir, err := GetHomeDir()
	if err != nil {
		return err
	}
	if err := TouchFile(fmt.Sprintf("%s/.awsd", homeDir)); err != nil {
		return err
	}
	s := []byte("")
	if config != "default" {
		s = []byte(config)
	}
	err = os.WriteFile(fmt.Sprintf("%s/.awsd", loc), s, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return nil
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
