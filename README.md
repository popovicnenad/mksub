mksub
-----
Make subdomains using a wordlist

Read a wordlist file (lowercase, remove `[^a-zA-Z0-9-_.]+`), filter unique words and generate subdomains.

```
Usage of mksub:
  -d string
    	Domain
  -w string
    	Wordlist file
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
> go run mksub.go -d example.com -w wordlist.txt
dev.domain.com
foo.bar.domain.com
prod.domain.com
```
