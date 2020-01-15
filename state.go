package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

// State is a struct containing all results of parsing a makefile.
// All variables, all targets etc.
type State struct {
	Targets   AllTargets     // a slice of all Target structs
	TargetMap map[int]string // map from line index to target name
}

// WorkerFunc is a type of function that can be used to concurrently parse a single line
// It takes a pointer to a State, a line index, the line contents and a WaitGroup that should have
// the Done method called once the function is done, for instance with "defer wg.Done()" as the first line.
type WorkerFunc func(*State, int, string, *sync.WaitGroup)

func (state *State) String() string {
	return fmt.Sprintf("%v", *state)
}

// ForEachLine will call a collection of functions concurrently, per line,
// then wait for all the concurrent functions to finish after all lines has been
// iterated over.
func (state *State) ForEachLine(path string, functionCollection []WorkerFunc) error {
	byteContents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for lineIndex, byteLine := range bytes.Split(byteContents, []byte{'\n'}) {
		// For each line, fire off all functions in the functionCollection
		line := string(byteLine)
		for _, f := range functionCollection {
			wg.Add(1)
			go f(state, lineIndex, line, &wg)
		}
	}
	wg.Wait()
	return nil
}

// Parse will try to parse a Makefile into a State struct
// If there are errors, the returned state will be nil.
func Parse(path string) (*State, error) {

	// Create a state, where the results from parsing will be stored
	state := &State{}

	// Prepare 256 targets, but keep the length at 0
	state.Targets = make(AllTargets, 0, 256)

	// Prepare a map from line index to make target name
	state.TargetMap = make(map[int]string)

	// Using a mutex for when modifying the state
	var mut sync.Mutex

	functionCollection := []WorkerFunc{
		// .PHONY handler
		func(state *State, lineIndex int, line string, wg *sync.WaitGroup) {
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
		func(state *State, lineIndex int, line string, wg *sync.WaitGroup) {
			defer wg.Done()
			// Is this an indented command?
			if strings.HasPrefix(line, "\t") {
				// Count down from lineIndex until a target name is reached
				var (
					targetName string
					found      bool
				)
				// Consider using RLock instead here. Benchmark locking the mutex within the array or outside of it.
				mut.Lock()
				for i := lineIndex; i >= 0; i-- {
					targetName, found = state.TargetMap[i]
					if found {
						break
					}
				}
				// Now save this command to the target comands
				target, err := state.Targets.GetTarget(targetName)
				mut.Unlock()
				if err != nil {
					// Commands without a target, this is an error
					// TODO: Use a proper make error message
					log.Fatalf("Found indented commands without a leading target, at line %d: %s\n", (lineIndex + 1), strings.TrimSpace(line))
				}
				// Save the command to the target.Commands slice
				mut.Lock()
				target.Commands = append(target.Commands, strings.TrimSpace(line))
				mut.Unlock()
				// This is not a target declaration and the command has been saved
				return
			}
			// Is this a make target?
			fields := strings.Fields(strings.TrimSpace(line))
			if len(fields) > 0 && !strings.HasPrefix(fields[0], ".") {
				targetName := ""
				if strings.HasSuffix(fields[0], ":") {
					targetName = fields[0][:len(fields[0])-1]
				} else {
					targetName = fields[0]
				}
				mut.Lock()
				state.TargetMap[lineIndex] = targetName
				if !state.Targets.HasName(targetName) {
					newTarget := state.Targets.AddTarget(targetName)
					if len(fields) > 1 {
						// TODO: Now add the normal and order-only prerequisites
						// targets : normal-prerequsites | order-only-prerequisites
						fmt.Println("MORE INFORMATION ABOUT: "+targetName, newTarget, fields)
					}
				} else {
					//fmt.Println("Already has: " + targetName)
				}
				mut.Unlock()
			}
		},
	}

	// Perform concurrent parsing of the makefile
	if err := state.ForEachLine(path, functionCollection); err != nil {
		return nil, err
	}

	return state, nil
}
