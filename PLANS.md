# Plans

## Before parsing targets

Before parsing, there are two things that should be done:

* Joining lines ending with "\"
* Going through all variable definitions and ifdefs
  * Add support for setting variables at the top, using =, := and ?=, possibly also prefixed with `export`.
  * Add support for $(shell ) and $(wildcard ), possibly a few others too.
  * Add support for ifdef, else and endif.

## When parsing targets

The goroutines that parse targets should trigger another goroutine that has access to the entire file,
that will iterate from the current lineIndex to either the end or until another target name has been reached. Alternatively, look at the lack of indentation ("\t") to determine the end of a target.

## Executing targets

* Add `execute.go` with code for executing a target.Commands slice. Respect variable expansion, trailing backslashes and leading "@" and/or "-".

## Special variables

* Add functionality for changing the default make target.

## General plans

* Build a trivial project with `ake` instead of `make`.
* Build increasingly less trivial project with `ake` instead of `make`.

## Some functions would be nice to have

* Determine if a shell command contains "`" or "$("
* Determine if a shell command contains "$$"

## Some things can be optimized

* A series of `cp` or `ln` commands can often be run concurrently, but overlap must be avoided. Copying in two unrelated locations should be fine, but there may be gotchas.
* A series of `echo` commands without piping can be combined to one printf, when executing, unless there are shell commands inside the echo statements.
* Starting to execute some targets before all targets have been parsed would be sweet.
