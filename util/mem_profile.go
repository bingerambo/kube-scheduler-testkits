package util

import (
	"fmt"
	"sync"
	"time"
)

func GatherMemoryProfile(componentName string, profileBaseName string, wg *sync.WaitGroup, dir string) {
	if wg != nil {
		defer wg.Done()
	}
	if err := gatherProfile(componentName, profileBaseName+"_"+time.Now().Format(time.RFC3339), "heap", dir); err != nil {
		fmt.Printf("Failed to gather %v memory profile: %v\n", componentName, err)
	}
}
