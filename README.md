# rip
[![Travis CI](https://img.shields.io/travis/gesquive/rip/master.svg?style=flat-square)](https://travis-ci.org/gesquive/rip)
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/rip/blob/master/LICENSE.md)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/gesquive/rip)

Sends a text file line by line to a remote host/port.


### Why?
 - Because bash has the ability to send data directly to a port, but other shells and OSs do not.
 - Because netcat is sometimes blocked because it is a "hacking" tool.
 - Because sometimes you need to throttle a data stream

## Installing

### Compile
This project requires go 1.6+ to compile. Just run `go get -u github.com/gesquive/rip` and the executable should be built for you automatically in your `$GOPATH`.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/rip/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/rip`.

## Usage

```console
Sends a text file line by line to a remote host/port.

Usage:
  rip [flags] <host>[:<port>] <tcp|udp> <file_path> [<file_path>...]

Flags:
  -r, --rate-limit int   Message rate allowed per second, use -1 for no limit (default: -1)
  -V, --version   Show the version and exit
```
Optionally, a hidden debug flag is available in case you need additional output.
```console
Hidden Flags:
  -D, --debug                  Include debug statements in log output
```

You can also pipe in input in addition to specified files on the command line:

```console
$ rip server:3333 tcp massive.log
$ cat massive.log | rip server:3333 tcp
$ app-with-output | rip server:3333 tcp
```

## Documentation

This documentation can be found at github.com/gesquive/rip

## License

This package is made available under an MIT-style license. See LICENSE.

## Contributing

PRs are always welcome!
