overseer
========

A process supervisor, written in Go. It lives within the same application space
as [supervisor](http://supervisord.org), 

Overseer can be configured using YAML, or JSON.

Right now, overseer is fairly simple; all you can do is add configuration files
to run commands. I am looking to add the following:

- HTTP status page
- API for controlling processes (with authentication)
- TOML support (maybe)

## Documentation

Documentation can be found [here](http://nesv.viewdocs.io/overseer).

## Installation

Overseer is a simple `go get` away!

	$ go get github.com/nesv/overseer
