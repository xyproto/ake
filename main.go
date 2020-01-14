package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

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

	// Concurrent parsing, using a mutex for when changing the state
	var mut sync.Mutex

	functionCollection := []WorkerFunc{
		// .PHONY handler
		func(state *State, line string, wg *sync.WaitGroup) {
			defer wg.Done()
			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) > 1 {
				if fields[0] == ".PHONY:" {
					for _, name := range fields[1:] {
						mut.Lock()
						// Check if the target already exists
						if existingTarget, err := state.Targets.GetTarget(name); err != nil {
							// Target does not exist, create a new one
							newTarget := state.Targets.AddTarget(name)
							// Updating this variable directly is possible,
							// since it points into the list of targets.
							newTarget.Phony = true
						} else {
							// Target does exists, set it to phony.
							// Updating this variable directly is possible,
							// since it points into the list of targets.
							existingTarget.Phony = true
						}
						mut.Unlock()
					}
				}
			}
		},
		// Target handler
		// TODO: Also store commands in the Target variable
		func(state *State, line string, wg *sync.WaitGroup) {
			defer wg.Done()
			if strings.HasPrefix(line, "\t") {
				// This is not a target declaration
				return
			}
			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) > 0 && !strings.HasPrefix(fields[0], ".") {
				targetName := ""
				if strings.HasSuffix(fields[0], ":") {
					targetName = fields[0][:len(fields[0])-1]
				} else {
					targetName = fields[0]
				}
				mut.Lock()
				if !state.Targets.HasName(targetName) {
					newTarget := state.Targets.AddTarget(targetName)
					// TODO: Now add the normal and order-only prerequisites
					// targets : normal-prerequsites | order-only-prerequisites
					fmt.Println("New target: "+targetName, newTarget)
				} else {
					fmt.Println("Already has: " + targetName)
				}
				mut.Unlock()
			}
		},
	}

	// Create a state, where the results from parsing will be stored
	state := &State{}
	// Prepare 256 targets, but keep the length at 0
	state.Targets = make(AllTargets, 0, 256)

	if err := state.ForEachLine(config.Makefile, functionCollection); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(*state)
}
