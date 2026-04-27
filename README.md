# go-netrc [![GoDoc](https://godoc.org/github.com/jdx/go-netrc?status.svg)](http://godoc.org/github.com/jdx/go-netrc) [![CircleCI](https://dl.circleci.com/status-badge/img/gh/jdx/go-netrc/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/jdx/go-netrc/tree/main)

A netrc parser for Go.

# Usage

Getting credentials for a host.

```go
usr, err := user.Current()
n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))
fmt.Println(n.Machine("api.heroku.com").Get("password"))
```

Setting credentials on a host.

```go
usr, err := user.Current()
n, err := netrc.Parse(filepath.Join(usr.HomeDir, ".netrc"))
n.Machine("api.heroku.com").Set("password", "newapikey")
n.Save()
```

# Quoting

Values that contain whitespace, double quotes, or backslashes can be wrapped in
double quotes:

```
machine example.com
  login alice
  password "my pass with spaces"
```

The following escape sequences are supported inside a quoted value: `\"`, `\\`,
`\n`, `\r`, `\t`. Any other `\x` is preserved literally on decode so
hand-edited values are never silently corrupted. `Get` returns the decoded
value, so callers see plain strings. `Set` and `AddMachine` automatically
quote and escape the value when needed; values without whitespace or special
characters are written bare so existing files round-trip unchanged.

This matches the quoting syntax adopted by curl 7.84+. Older `.netrc` parsers
may not recognize quoted values, so portability across tools varies.
