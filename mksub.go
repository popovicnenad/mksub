package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "regexp"
    "strings"
    "sync"
)

var (
    //flags
    domain *string
    domainFile *string
    wordlist *string
    regex *string
    level *int
    output *string

    inputDomains []string
    wordSet map[string]bool
    writeMutex sync.Mutex
)

func fileReadDomain(fileName string) {
    inputFile, err := os.Open(fileName)
    if err != nil {
        panic("Could not open file to read domains!")
    }
    defer inputFile.Close()

    scanner := bufio.NewScanner(inputFile)
    for scanner.Scan() {
        inputDomains = append(inputDomains, strings.TrimSpace(scanner.Text()))
    }
}

func prepareDomains() {
    if *domain == "" && *domainFile == "" {
        fmt.Println("No domain input provided")
        os.Exit(1)
    }

    inputDomains = make([]string, 0)
    if *domain != "" {
        inputDomains = append(inputDomains, *domain)
    } else {
        if *domainFile != "" {
            fileReadDomain(*domainFile)
        }
    }
}

func processWordList (domain string, outputFile *os.File, wg *sync.WaitGroup) {
    defer wg.Done()

    results := make([]string, 0)
    for word := range wordSet {
        results = append(results, word)
    }
    toMerge := results[0:]

    for i:=0; i< *level-1; i++ {
        toMerge = results[0:]
        for _, sd := range toMerge {
            for word := range wordSet {
                results = append(results, word + "." + sd)
            }
        }
    }

    for _, subdomain := range results {
        fmt.Println(subdomain + "." + domain)
        writeMutex.Lock()
        _, _ = outputFile.WriteString(subdomain + "." + domain + "\n")
        writeMutex.Unlock()
    }

}

func main() {
    domain = flag.String("d", "", "Input domain")
    domainFile = flag.String("df", "", "Input domain file, one domain per line")
    wordlist = flag.String("w", "", "Wordlist file")
    regex = flag.String("r", "", "Regex to filter words from wordlist file")
    level = flag.Int("l", 1, "Subdomain level to generate (default 1)")
    output = flag.String("o", "", "Output file (optional)")
    flag.Parse()

    prepareDomains()

    var reg *regexp.Regexp
    var err error
    if *regex != "" {
        reg, err = regexp.Compile(*regex)
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

    wordlistFile, err := os.Open(*wordlist)
    if err != nil {
        fmt.Println(err.Error())
        os.Exit(1)
    }
    defer wordlistFile.Close()

    wordSet = make(map[string]bool)
    scanner := bufio.NewScanner(wordlistFile)
    for scanner.Scan() {
        word := strings.ToLower(scanner.Text())
        if reg != nil {
            if !reg.Match([]byte(word)) {
                continue
            }
        }
        if word != "" {
            wordSet[word] = true
        }
    }

    var wg sync.WaitGroup
    for _, dom := range inputDomains {
        wg.Add(1)
        go processWordList(dom, outputFile, &wg)
    }

    wg.Wait()
}
