package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

// State is a struct containing all results of parsing a makefile.
// All variables, all targets etc.
type State struct {
	Targets   AllTargets          // a slice of all Target structs
	TargetMap map[int]string      // map from line index to target name
	Variables map[string][]string // a map of all defined variables, from string to string.
	// TODO: Modify to point from a variable name to a slice of strings, if needed.
}

// WorkerFunc is a type of function that can be used to concurrently parse a single line
// It takes a pointer to a State, a line index, the line contents and a WaitGroup that should have
// the Done method called once the function is done, for instance with "defer wg.Done()" as the first line.
// The []string argument is a slice of all lines, so that coroutines may discover their context.
type WorkerFunc func(*State, int, string, *sync.WaitGroup, []string)

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
	lines := strings.Split(string(byteContents), "\n")
	for lineIndex, line := range lines {
		// For each line, fire off all functions in the functionCollection
		for _, f := range functionCollection {
			wg.Add(1)
			go f(state, lineIndex, line, &wg, lines)
		}
	}
	wg.Wait()
	return nil
}

// ConcurrentParsing parses the given Makefile concurrently
func (state *State) ConcurrentParsing(path string) error {

	// Using a mutex for when modifying the state
	var mut sync.Mutex

	functionCollection := []WorkerFunc{
		// .PHONY handler
		func(state *State, lineIndex int, line string, wg *sync.WaitGroup, lines []string) {
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
		func(state *State, lineIndex int, line string, wg *sync.WaitGroup, lines []string) {
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
				target.Commands = append(target.Commands, NewCommand(line))
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
	return state.ForEachLine(path, functionCollection)
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

	// Now, before starting the concurrent parsing, it might be a good idea to resolve all ifdefs first.
	// And they might depend on variables. So parsing all variables and all ifdefs first might be needed.
	// Alternatively, the ifdefs can wait for variables to be defined. But then, how does one know if it's waiting for a redifinition or not?
	// No, it's better to parse variables and ifdefs properly first, and then later parse all the details.
	// An overview over which target names applies to which lines would also be useful.
	// Then this could be a single pass, before starting on the other parsing, below.
	// By all means, the first pass could also be concurrent, but then they would need to repeately search for what the lines ahead contained.
	// No, one good old fashioned pass first is a good idea. Let's do that.
	// But! Can it be done concurrently, just for the heck of it? Yes, probably. Let's do that.

	if err := state.ConcurrentParsing(path); err != nil {
		return nil, err
	}

	return state, nil
}
