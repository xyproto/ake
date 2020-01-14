package main

import (
	"bytes"
	"io/ioutil"
	"sync"
)

type WorkerFunc func(*State, string, *sync.WaitGroup)

func (state *State) ForEachLine(path string, functionCollection []WorkerFunc) error {
	byteContents, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for _, byteLine := range bytes.Split(byteContents, []byte{'\n'}) {
		s := string(byteLine)
		for _, f := range functionCollection {
			wg.Add(1)
			go f(state, s, &wg)
		}
	}
	wg.Wait()
	return nil
}
