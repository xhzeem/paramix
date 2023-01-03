# paramix
An enhanced replacement for [@tomnomnom/qsreplace](https://github.com/tomnomnom/qsreplace/)


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
