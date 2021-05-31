package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)


func main() {
	domain := flag.String("d", "", "Domain")
	wordlist := flag.String("w", "", "Wordlist file")
	flag.Parse()

	wordlistFile, err := os.Open(*wordlist)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer wordlistFile.Close()

	wordSet := make(map[string]bool)
	reg, _ := regexp.Compile("[^a-zA-Z0-9-_.]+")
	scanner := bufio.NewScanner(wordlistFile)

	for scanner.Scan() {
		word := reg.ReplaceAllString(strings.ToLower(scanner.Text()), "")
		if _, isOld := wordSet[word]; word != "" && !isOld  {
			wordSet[word] = true
			fmt.Println(word + "." + *domain)
		}
	}
}