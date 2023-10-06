package utils

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"sort"
)

func GetProfiles() []string {
	profileFileLocation := GetCurrentProfileFile()
	profiles := make([]string, 0)

	file, err := os.Open(profileFileLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	r, err := regexp.Compile(`\[profile .*]`)

	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			s := scanner.Text()
			reg := regexp.MustCompile(`(\[profile )|(])`)
			res := reg.ReplaceAllString(s, "")
			profiles = append(profiles, res)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	profiles = AppendIfNotExists(profiles, "default")
	sort.Strings(profiles)
	return profiles
}
