package main

import (
	"bufio"
	"fmt"
	"github.com/manifoldco/promptui"
	"log"
	"os"
	"regexp"
	"sort"
)

const (
	NoticeColor = "\033[0;38m%s\u001B[0m"
	PromptColor = "\033[1;38m%s\u001B[0m"
	KermitColor = "\033[1;32m%s\033[0m"
	CyanColor   = "\033[0;36m%s\033[0m"
)

func main() {

	home := os.Getenv("HOME")
	profileFileLocation := getenv("AWS_CONFIG_FILE", fmt.Sprintf("%s/.aws/config", home))
	profiles := getProfiles(profileFileLocation)

	fmt.Printf(NoticeColor, "AWS Profile Switcher\n")
	prompt := promptui.Select{
		Label:    fmt.Sprintf(PromptColor, "Choose a profile"),
		Items:    profiles,
		HideHelp: true,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}?",
			Active:   fmt.Sprintf("%s {{ . | cyan }}", promptui.IconSelect),
			Inactive: "  {{.}}",
			Selected: "  {{ . | cyan }}",
		},
		Stdout: &bellSkipper{},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf(KermitColor, "? ")
	fmt.Printf(PromptColor, "Choose a profile ")
	fmt.Printf(CyanColor, result)
	writeFile(result, home)
}

func writeFile(profile, loc string) {
	s := []byte("")
	if profile != "default" {
		s = []byte(profile)
	}
	err := os.WriteFile(fmt.Sprintf("%s/.awsd", loc), s, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getProfiles(profileFileLocation string) []string {
	profiles := make([]string, 0)

	file, err := os.Open(profileFileLocation)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	r, err := regexp.Compile(`\[profile .*]`) // this can also be a regex

	if err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		if r.MatchString(scanner.Text()) {
			s := scanner.Text()
			reg := regexp.MustCompile(`(\[profile )|(\])`)
			res := reg.ReplaceAllString(s, "")
			profiles = append(profiles, res)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	profiles = append(profiles, "default")
	sort.Strings(profiles)
	return profiles
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

type bellSkipper struct{}

func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}
