# Reggi

Interactively validate regex against test data in your console

![Example 1](./assets/example1.png)

![Example 2](./assets/example2.png)

![Example 3](./assets/example3.png)

# Install

```
$ go get github.com/byxorna/reggi/cmd/reggi
```

Then, `reggi some.txt files.txt` to interactively test your regexp against your files.

# Use

## Input Mode

Enter a regex to match against the currently focused buffer. Matches are highlighted.

- `ctrl-y` enable match all expressions
- `ctrl-l` enable multiline match: ^ and $ match begin/end line
- `ctrl-s` enable span line: let . match \n
- `ctrl-i` enable insensitive matching
- Press `esc` to enter pager

## Pager Mode

- `i`,`a` to go back to the regex editor
- `H`,`L` to change buffers (if multiple files are open)
- Normal pagination (`hjkl`, `ctrl-f`, `ctrl-b`, `g`, `G`)
- `q`,`ctrl-c` to quit

# Dev

```
$ make dev # opens a fixture
```

# About

I use [rubular.com](rubular.com) constantly, and wanted to make something similar that I could keep closer at hand in the console.

# TODO

- [ ] implement different color for submatches vs matches
- [ ] implement a submatch expression explorer to visualise submatches as a tree

