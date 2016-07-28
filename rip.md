#rip

A small utility to send a file to a specific port

### Why?
Bash has the ability to echo to `/dev/tcp/server-name/8080` and it opens a port to that box to send the data
In windows this does not exist
In zsh this does not exist
Sometimes netcat is blocked because it is a "hacking" tool

## Usage

```console
rip <host:port|ip:port> [flags] <file_paths>...

-f file path to file to send (allow multiple)
-p sending protocol (default:tcp|udp)

```
