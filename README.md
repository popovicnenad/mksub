mksub
-----
Make subdomains using a wordlist

Read a wordlist file and generate subdomains for given domain or list of domains.
Input from wordlist file is lowercased and unique words are processed. Additionally, wordlist can be
filtered using regex. 

```
Usage of mksub:
  -d string
        Input domain
  -df string
        Input domain file, one domain per line
  -w string
    	Wordlist file
  -r string
        Regex to filter words from wordlist file
  -o string
    	Output file (optional)
```

### Example

##### wordlist.txt
```
dev
DEV
*
foo.bar
prod
```
```shell script
> go run mksub.go -d example.com -w input.txt -r "^[a-zA-Z0-9\.-_]+$"
dev.example.com
foo.bar.example.com
prod.example.com
```
