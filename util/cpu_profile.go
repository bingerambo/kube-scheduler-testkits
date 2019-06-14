package util

import (
	"fmt"
	"sync"
	"time"
)

const (
	// Default value for how long the CPU profile is gathered for.
	DefaultCPUProfileSeconds = 5
)

func StartCPUProfileGatherer(componentName string, profileBaseName string, interval time.Duration, dir string) chan struct{} {
	stopCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-time.After(interval):
				GatherCPUProfile(componentName, profileBaseName+"_"+time.Now().Format(time.RFC3339), nil, dir)
			case <-stopCh:
				return
			}
		}
	}()
	return stopCh
}

func GatherCPUProfile(componentName string, profileBaseName string, wg *sync.WaitGroup, dir string) {
	GatherCPUProfileForSeconds(componentName, profileBaseName, DefaultCPUProfileSeconds, wg, dir)
}

func GatherCPUProfileForSeconds(componentName string, profileBaseName string, seconds int, wg *sync.WaitGroup, dir string) {
	if wg != nil {
		defer wg.Done()
	}
	if err := gatherProfile(componentName, profileBaseName, fmt.Sprintf("profile?seconds=%v", seconds), dir); err != nil {
		fmt.Printf("Failed to gather %v CPU profile: %v\n", componentName, err)
	}
}
