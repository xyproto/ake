# makeflags

[Online Documentation](https://godoc.org/github.com/xyproto/makeflags)

Package for handling the same flags and arguments that GNU Make parses, then return a Config struct.

It can be used as a starting-point for writing a drop-in replacement for GNU Make.

## Example use

```go
package main

import (
	"fmt"
	"os"

	"github.com/xyproto/makeflags"
)

func main() {
	config := makeflags.New()

	if config.VersionInfoAndExit {
		fmt.Println(makeflags.Version)
		os.Exit(0)
	}

	if len(config.Targets) == 0 && config.Makefile == "" {
		fmt.Println("make: *** No targets specified and no makefile found.  Stop.")
		os.Exit(2)
	}

	fmt.Printf("[%s] %v\n", config.Makefile, config.Targets)
}
```

## General info

* Version: 1.1.0
* License: MIT
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
