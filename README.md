# Paramix
Paramix is Golang project built to help in bug bounty and penetration testing for modifying and cleaning the URLs of some targets focusing on the parameters of URLs. It reads a list of URLs from stdin, performs the specified actions on them, and then returns the modified URLs in stdout. It can add, remove, or modify the values of parameters, and it can also URL-decode the values. It can operate on all parameters at once, or on specific ones specified by the user. It can also keep URLs with no parameters, replace the values of parameters or append to them.


## Install
```bash
go install github.com/xhzeem/paramix@latest
```

## Flags

```python
  -a *str   Add custom parameters to the URLs, comma seprated
  -r *str   Remove a parameter from the URLs, comma seprated
  -v *str   Value to modify the parameters upon
  -d	    URLdecode the values of the parameters
  -k	    Keep the URLs with no parameters
  -m	    Modify all parameters at once
  -o	    Replace the value instead of appending
```

## Usage
```bash
➜ xhzeem $ echo "http://xhzeem.me/?x=1&y=2&z=3" | paramix -a a -r y,z -v xss  

http://xhzeem.me/?a=xss&x=1
http://xhzeem.me/?a=xss&x=1xss
```

## Thanks
-  @tomnomnom This project uses some code from his tool [qsreplace](https://github.com/tomnomnom/qsreplace/).
