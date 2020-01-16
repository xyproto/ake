package main

import (
	"errors"
)

// AllTargets is a slice of all make targets, as discovered when parsing
type AllTargets []Target

// Target represents a make target, like "all", "clean" or "main.o"
type Target struct {
	ID        int        // ID, a counter
	Name      string     // Can be a regular name or it can be something like $(OBJDIR)/%.o
	Normal    []*Target  // Before "|"
	OrderOnly []*Target  // After "|"
	Phony     bool       // Is it .PHONY ?
	Commands  []*Command // Commands to run (not ifdef etc, just the ones indented with tab)
}

// HasName checks if the given name exists in the slice of targets
func (all AllTargets) HasName(name string) bool {
	for _, t := range all {
		if t.Name == name {
			return true
		}
	}
	return false
}

// HasTarget checks if the given target thas the same ID as one in the slice of targets
func (all AllTargets) HasTarget(target *Target) bool {
	for _, t := range all {
		if t.ID == target.ID {
			return true
		}
	}
	return false
}

// GetTarget returns a pointer to a Target within the
func (all AllTargets) GetTarget(name string) (*Target, error) {
	for i, t := range all {
		if t.Name == name {
			//return &t, nil
			// Return a pointer to the target directly in the list
			return &((all)[i]), nil
		}
	}
	return nil, errors.New("could not find " + name)
}

// AddTarget creates a new Target struct with the given name
// and returns a pointer directly into the AllTargets slice,
// that can be used for modifying the target later.
func (all *AllTargets) AddTarget(name string) *Target {
	t := &Target{}
	t.ID = len(*all)
	t.Name = name
	*all = append(*all, *t)
	// Return a pointer to the element in the list
	// (instead of a pointer to the local variable)
	return &((*all)[t.ID])
}
