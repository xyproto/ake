package main

import (
	"errors"
)

type AllTargets []Target

type Target struct {
	ID        int       // ID, a counter
	Name      string    // Can be a regular name or it can be something like $(OBJDIR)/%.o
	Normal    []*Target // Before "|"
	OrderOnly []*Target // After "|"
	Phony     bool      // Is it .PHONY ?
}

func (all AllTargets) HasName(name string) bool {
	for _, t := range all {
		if t.Name == name {
			return true
		}
	}
	return false
}

//func (all AllTargets) SetPhony(name string) error {
//	for i, t := range all {
//		if t.Name == name {
//			all[i].Phony = true
//			return nil
//		}
//	}
//	return errors.New("could not find " + name)
//}

func (all AllTargets) HasTarget(target *Target) bool {
	for _, t := range all {
		if t.ID == target.ID {
			return true
		}
	}
	return false
}

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

func (all *AllTargets) AddTarget(name string) *Target {
	t := &Target{}
	t.ID = len(*all)
	t.Name = name
	*all = append(*all, *t)
	// Return a pointer to the element in the list
	// (instead of a pointer to the local variable)
	return &((*all)[t.ID])
}
