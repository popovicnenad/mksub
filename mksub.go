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
    r := flag.String("r", "", "Regex to filter words from wordlist file")
    level := flag.Int("l", 1, "Subdomain level to generate (default 1)")
    output := flag.String("o", "", "Output file (optional)")
    flag.Parse()

    wordlistFile, err := os.Open(*wordlist)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer wordlistFile.Close()

    var reg *regexp.Regexp
    if *r != "" {
        reg, err = regexp.Compile(*r)
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
    }

    var outputFile *os.File
    if *output != "" {
        outputFile, err = os.Create(*output)
        if err != nil {
            fmt.Println(err.Error())
            os.Exit(1)
        }
        defer outputFile.Close()
    }

    wordSet := make(map[string]bool)
    scanner := bufio.NewScanner(wordlistFile)

    for scanner.Scan() {
        word := strings.ToLower(scanner.Text())
        if reg != nil {
            if !reg.Match([]byte(word)) {
                continue
            }
        }
        if _, isOld := wordSet[word]; word != "" && !isOld  {
            wordSet[word] = true
            //fmt.Println(word + "." + *domain)
            //if outputFile != nil {
            //    _, _ = outputFile.WriteString(word + "." + *domain + "\n")
            //}
        }
    }
    results := make([]string, 0)
    for i:=0; i<*level; i+=1{
        mergeWords := results[0:len(results)]
        if len(mergeWords) == 0 {
            for word, _ := range wordSet {
                results = append(results, word)
            }
        } else {
            for _, mw := range mergeWords {
                for word, _ := range wordSet {
                    results = append(results, fmt.Sprintf("%s.%s", word, mw))
                }
            }
        }
    }
    for _, subdomain := range results {
        fmt.Println(subdomain + "." + *domain)
        if outputFile != nil {
           _, _ = outputFile.WriteString(subdomain + "." + *domain + "\n")
        }
    }
}
