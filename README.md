# paramix
An enhanced version of [@tomnomnom/qsreplace](https://github.com/tomnomnom/qsreplace/) that support parameter-by-parameter modification


## Install
```bash
go install github.com/xhzeem/paramix@latest
```

## Flags

```python
-r  [Relace Mode]  Replace the value instead of appending it.
-s  [Single Mode]  Modify the parameters one by one.
```

## Usage
```bash
➜ xhzeem $ echo "http://xhzeem.me/?x=1&y=2&z=3" | paramix -s xhzeem

http://xhzeem.me/?x=1xhzeem&y=2&z=3
http://xhzeem.me/?x=1&y=2xhzeem&z=3
http://xhzeem.me/?x=1&y=2&z=3xhzeem
```
