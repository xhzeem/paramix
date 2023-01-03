# paramix
An enhanced version of [@tomnomnom/qsreplace](https://github.com/tomnomnom/qsreplace/) that support parameter-by-parameter modification


## Install
```bash
go install github.com/xhzeem/paramix@latest
```

## Flags

```python
  -v  [val]  Set the custom value to modify the URLs upon
  -p  [val]  Add a custom parameter to the URL
  -d         URL decode the values of the paramters
  -r         Replace the value instead of appending it
  -s         Modify the single parameter at a time
```

## Usage
```bash
➜ xhzeem $ echo "http://xhzeem.me/?x=1&y=2&z=3" | paramix -s -v xhzeem -p new param

http://xhzeem.me/?added=xhzeem&x=1&y=2&z=3
http://xhzeem.me/?added=&x=1xhzeem&y=2&z=3
http://xhzeem.me/?added=&x=1&y=2xhzeem&z=3
http://xhzeem.me/?added=&x=1&y=2&z=3xhzeem
```
