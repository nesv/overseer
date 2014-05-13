# Overseer's documentation

Welcome to the documentation for Overseer: yet another process supervisor!

Most process supervisors are written in an interpreted language of some sort
(for example, [supervisor](http://supervisord.org) is written in Python), and
it seems a bit silly to rely on an interpreted language to do something like
process supervision. There are always dependencies to pull in that may not be
packaged by your preferred operating system vendor, or there may be a version
mismatch between the version of the interpreter required to run the program and
the version supplied by your OS vendor.

Enter Overseer.

It is written in [Go](http://golang.org), and once it is compiled, it has no
runtime dependencies.

## Table of contents

1. Configuration
2. Processes
3. ...
