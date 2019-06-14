package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func gatherProfile(componentName, profileBaseName, profileKind, dir string) error {
	profilePort, err := getPortForComponent(componentName)
	if err != nil {
		return fmt.Errorf("Profile gathering failed finding component port: %v", err)
	}
	if profileBaseName == "" {
		profileBaseName = time.Now().Format(time.RFC3339)
	}

	// Get the profile data localhost
	var cmd *exec.Cmd
	getCommand := fmt.Sprintf("curl -s http://127.0.0.1:%v/debug/pprof/%s", profilePort, profileKind)
	cmd = exec.Command("sh", "-c", getCommand)

	profilePrefix := componentName
	switch {
	case profileKind == "heap":
		profilePrefix += "_MemoryProfile_"
	case strings.HasPrefix(profileKind, "profile"):
		profilePrefix += "_CPUProfile_"
	case strings.HasPrefix(profileKind, "trace"):
		profilePrefix += "_Trace_"

	default:
		return fmt.Errorf("Unknown profile kind provided: %s", profileKind)
	}

	// Write the profile data to a file.
	rawprofileDir := path.Join(dir, "profiles")
	rawprofilePath := path.Join(dir, "profiles", profilePrefix+profileBaseName+".pprof")
	err = os.MkdirAll(rawprofileDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to create dir for the profile graph: %v", err)
	}

	rawprofile, err := os.Create(rawprofilePath)
	if err != nil {
		return fmt.Errorf("Failed to create file for the profile graph: %v", err)
	}
	defer rawprofile.Close()
	cmd.Stdout = rawprofile
	stderr := bytes.NewBuffer(nil)
	cmd.Stderr = stderr
	if err := cmd.Run(); nil != err {
		return fmt.Errorf("Failed to get profile data: %v, stderr: %#v", err, stderr.String())
	}

	return nil
}

func getPortForComponent(componentName string) (int, error) {
	switch componentName {
	case "kube-apiserver":
		return 8080, nil
	case "kube-scheduler":
		return 10251, nil
	case "kube-controller-manager":
		return 10252, nil
	}
	return -1, fmt.Errorf("Port for component %v unknown", componentName)
}

func GetRuntimeKits(kits []string) (cpuProfile, memProfile, trace bool) {
	for _, str := range kits {
		if str == "cpu" {
			cpuProfile = true
		} else if str == "mem" {
			memProfile = true
		} else if str == "trace" {
			trace = true
		}
	}
	return
}
