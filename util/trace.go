package util

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Default value for how long the Trace is gathered for.
	DefaultTraceSeconds = 5
)

func StartTraceGatherer(componentName string, profileBaseName string, interval time.Duration, dir string) chan struct{} {
	stopCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(interval):
				GatherTrace(componentName, profileBaseName+"_"+time.Now().Format(time.RFC3339), nil, dir)
			case <-stopCh:
				return
			}
		}
	}()
	return stopCh
}

func GatherTrace(componentName string, profileBaseName string, wg *sync.WaitGroup, dir string) {
	GatherTraceForSeconds(componentName, profileBaseName, DefaultTraceSeconds, wg, dir)
}

func GatherTraceForSeconds(componentName string, profileBaseName string, seconds int, wg *sync.WaitGroup, dir string) {
	if wg != nil {
		defer wg.Done()
	}
	if err := gatherProfile(componentName, profileBaseName, fmt.Sprintf("trace?seconds=%v", seconds), dir); err != nil {
		fmt.Printf("Failed to gather %v trace: %v\n", componentName, err)
	}
}
