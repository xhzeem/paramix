# paramix
An enhanced replacement for [@tomnomnom/qsreplace](https://github.com/tomnomnom/qsreplace/)


## Install
```bash
go install github.com/xhzeem/paramix@latest
```

## Flags

```python
  -v  [str]   Value to modify the parameters upon
  -p  [str]   Add a custom parameter to the URLs 
  -r          Replace the value instead of appending
  -m          Modify all parameters at once
  -d          URLdecode the values of the parameters
  -k          Keep the URLs with no parameters
```

## Usage
```bash
➜ xhzeem $ echo "http://xhzeem.me/?x=1&y=2&z=3" | paramix -v xhzeem -p new

http://xhzeem.me/?new=xhzeem&x=1&y=2&z=3
http://xhzeem.me/?new=xhzeem&x=1xhzeem&y=2&z=3
http://xhzeem.me/?new=xhzeem&x=1&y=2xhzeem&z=3
http://xhzeem.me/?new=xhzeem&x=1&y=2&z=3xhzeem
```
