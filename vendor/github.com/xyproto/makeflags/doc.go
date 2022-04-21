// Package makeflags includes functions for argument parsing.
// The goal is to support the exact same flags and arguments
// that GNU Make does, then return a Config struct.
package makeflags

// Version is the text returned by passing "-v" or "--version"
var Version = "Makeflags 1.1.0"
