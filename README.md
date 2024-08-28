# md

A no-nonsense markdown viewer for the terminal based on
[glamour](https://github.com/charmbracelet/glamour).

## Installation

```sh
go install moehl.dev/md@latest
```

## Tips

Monitoring files for changes is not supported. The easiest way to iterate quickly is this
one-liner:

```sh
while :; do md -p ./README.md; done
```

Pressing `q` quits less and causes the file to be reloaded, CTRL+C will cause `md` to fail which
also terminates less and breaks the loop.
