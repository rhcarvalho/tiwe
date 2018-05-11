package main

import (
	"log"
	"sync"
)

type SyncGroup struct {
	wg sync.WaitGroup
}

func (sg *SyncGroup) Go(f func()) {
	sg.wg.Add(1)
	go func() {
		f := f
		defer sg.wg.Done()
		log.Printf("begin %v", f)
		defer log.Printf("end %v", f)
		f()
	}()
}

func (sg *SyncGroup) Wait() {
	sg.wg.Wait()
}
