package utils

import (
	"gopkg.in/ini.v1"
	"log"
	"sort"
	"strings"
)

const (
	profilePrefix  = "profile"
	defaultProfile = "default"
)

func GetProfiles() []string {
	profileFileLocation := GetCurrentProfileFile()
	cfg, err := ini.Load(profileFileLocation)
	if err != nil {
		log.Fatalf("Failed to load profiles: %v", err)
	}
	sections := cfg.SectionStrings()
	profiles := make([]string, 0, len(sections)+1)
	for _, section := range sections {
		if strings.HasPrefix(section, profilePrefix) {
			trimmedProfile := strings.TrimPrefix(section, profilePrefix)
			trimmedProfile = strings.TrimSpace(trimmedProfile)
			profiles = append(profiles, trimmedProfile)
		}
	}
	profiles = AppendIfNotExists(profiles, defaultProfile)
	sort.Strings(profiles)
	return profiles
}
