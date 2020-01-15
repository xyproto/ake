# Plans

1. Add `execute.go` with code for executing a target.Commands slice. Respect variable expansion, trailing backslashes and leading "@" and/or "-".
2. Add functionality for changing the default make target.
3. Add support for setting variables at the top, using =, := and ?=, possibly also prefixed with `export`.
4. Add support for $(shell ) and $(wildcard ), possibly a few others too.
5. Add support for ifdef, else and endif.
6. Build a trivial project with `ake` instead of `make`.
7. Build a less trivial project with `ake` instead of `make`.

## Some things can be optimized

* A series of `cp` or `ln` commands can often be run concurrently, but overlap must be avoided. Copying in two unrelated locations should be fine, but there may be gotchas.
* A series of `echo` commands without piping can be combined to one printf, when executing, unless there are shell commands inside the echo statements.
