package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	"schedtest/k8s"
	"schedtest/util"
)

var replicas int
var name string
var ns string
var reportDir string
var runtimeKit string

func init() {
	flag.StringVar(&k8s.KubeConfigFile, "kubeConfigFile", "/etc/kubernetes/admin.kubeconfig", "Kubernetest config file used to connect kube-apiserver.")
	flag.IntVar(&replicas, "replicas", 1, "Statefulset replicas")
	flag.StringVar(&name, "name", "app-density-test", "Statefulset name")
	flag.StringVar(&ns, "namespace", "app-density-test", "Statefulset namespace")
	flag.StringVar(&reportDir, "reportDir", "/tmp/schedtest/", "Profiles dir")
	flag.StringVar(&runtimeKit, "runtimeKit", "", "cpu,mem,trace")
}

func main() {
	flag.Parse()
	k8s.MustInit()

	// Solve runtimeKit
	kits := strings.Split(runtimeKit, ",")
	cpuProfile, memProfile, trace := util.GetRuntimeKits(kits)

	// Start scheduler CPU profile-gatherer before we begin cluster saturation.
	profileGatheringDelay := 1 * time.Second
	schedulerProfilingStopCh := make(chan struct{})
	schedulerTraceStopCh := make(chan struct{})
	if cpuProfile {
		schedulerProfilingStopCh = util.StartCPUProfileGatherer("kube-scheduler", "density", profileGatheringDelay, reportDir)
	}
	// Start scheduler trace-gatherer before we begin cluster saturation.
	if trace {
		schedulerTraceStopCh = util.StartTraceGatherer("kube-scheduler", "density", profileGatheringDelay, reportDir)
	}
	startTime := time.Now()

	err := k8s.CreateSts(name, ns, replicas)
	if err != nil {
		fmt.Printf("Failed to create sts, err: %v", err)
		return
	}

	label := labels.SelectorFromSet(labels.Set(map[string]string{"app": name}))

	cs := k8s.GetClient()

	ps, err := k8s.NewPodStore(cs, ns, label, fields.Everything())
	if err != nil {
		fmt.Printf("Failed to new podStore, err: %v", err)
		return
	}
	defer ps.Stop()

	timeout := 3 * time.Minute

	oldPods := make([]*v1.Pod, 0)
	oldRunning := 0
	lastChange := time.Now()

	var startupStatus k8s.PodsStartupStatus

	for oldRunning != replicas {
		//time.Sleep(interval)

		pods := ps.List()
		startupStatus = k8s.ComputePodsStartupStatus(pods, replicas)

		if len(pods) > len(oldPods) || startupStatus.Scheduled > oldRunning {
			lastChange = time.Now()
		}
		oldPods = pods
		oldRunning = startupStatus.Scheduled

		if time.Since(lastChange) > timeout {
			break
		}
	}

	if oldRunning != replicas {
		fmt.Printf("Only %d pods scheduled out of %d", oldRunning, replicas)
		return
	}

	endTime := time.Now()

	startupTime := endTime.Sub(startTime)
	close(schedulerProfilingStopCh)
	close(schedulerTraceStopCh)

	// Grabbing scheduler memory profile after cluster saturation finished.
	if memProfile {
		wg := sync.WaitGroup{}
		wg.Add(1)
		util.GatherMemoryProfile("kube-scheduler", "density", &wg, reportDir)
		wg.Wait()
	}

	fmt.Printf("Start time: %v\n", startTime)
	fmt.Printf("End Time: %v\n", endTime)

	fmt.Printf("%d scheduled for %d pods\n", startupStatus.Scheduled, replicas)

	fmt.Printf("E2E startup time for %d pods: %v\n", replicas, float32(startupTime.Seconds()))
	fmt.Printf("Throughput (pods/s) during cluster saturation phase: %v\n", float32(replicas)/float32(startupTime.Seconds()))
}
