# rip
[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/gesquive/rip/blob/master/LICENSE)
[![Build Status](https://img.shields.io/circleci/build/github/gesquive/rip?style=flat-square)](https://circleci.com/gh/gesquive/rip)
[![Coverage Report](https://img.shields.io/codecov/c/gh/gesquive/rip?style=flat-square)](https://codecov.io/gh/gesquive/rip)

Sends a text file line by line to a remote host/port.


### Why?
 - Because bash has the ability to send data directly to a port, but other shells and OSs do not.
 - Because netcat is sometimes blocked because it is a "hacking" tool.
 - Because sometimes you need to throttle a data stream

## Installing

### Compile
This project has only been tested with go1.11+. To compile just run `go get -u github.com/gesquive/rip` and the executable should be built for you automatically in your `$GOPATH`. This project uses go mods, so you might need to set `GO111MODULE=on` in order for `go get` to complete properly.

Optionally you can run `make install` to build and copy the executable to `/usr/local/bin/` with correct permissions.

### Download
Alternately, you can download the latest release for your platform from [github](https://github.com/gesquive/rip/releases).

Once you have an executable, make sure to copy it somewhere on your path like `/usr/local/bin` or `C:/Program Files/`.
If on a \*nix/mac system, make sure to run `chmod +x /path/to/rip`.

### Homebrew
This app is also avalable from this [homebrew tap](https://github.com/gesquive/homebrew-tap). Just install the tap and then the app will be available.
```shell
$ brew tap gesquive/tap
$ brew install rip
```

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

```shell
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
