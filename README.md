# paramix
An enhanced version of @tomnomnom/qsreplace that support parameter-by-parameter modification


## Install
```bash
go install github.com/xhzeem/paramix@latest
```

## Flags

```python
-r    Replace the value instead of appending it.
-f    Modify the parameters one by one.
```

## Usage
```bash
➜ xhzeem $ echo "http://xhzeem.me/?x=1&y=2&z=3" | paramix -f xhzeem

http://xhzeem.me/?x=1xhzeem&y=2&z=3
http://xhzeem.me/?x=1&y=2xhzeem&z=3
http://xhzeem.me/?x=1&y=2&z=3xhzeem
```
