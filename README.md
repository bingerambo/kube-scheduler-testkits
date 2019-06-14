Kube-scheduler-testkits is a performance testing tool for kube-scheduler's throughput，which can also grab cpu-profile, mem-profile and go-trace. It can be easier used and extended than kubernetes e2e test.

By default, it grabs cpu-profile and go-trace for five seconds every one second.

The best practice:
1. Use Kubemark to make a simulated cluster which can be much bigger than a real one. ([kubemark-outside-gce](https://github.com/snowplayfire/kubemark-outside-gce))
2. Use kube-scheduler-testkits to run density test which can tell us the kube-scheduler's throughput and grab profiles for us.

## Build
``
GOPROXY=https://athens.azurefd.net GO111MODULE=on go build
``
>Precondition: Go >= 1.11

## Deploy
Put the binary `schedtest` on the leader-scheduler's host.

## Use
```
[root@xxxx]# ./schedtest -h
Usage of ./schedtest:
  -kubeConfigFile string
        Kubernetest config file used to connect kube-apiserver. (default "/etc/kubernetes/admin.kubeconfig")
  -name string
        Statefulset name (default "app-density-test")
  -namespace string
        Statefulset namespace (default "app-density-test")
  -replicas int
        Statefulset replicas (default 1)
  -reportDir string
        Profiles dir (default "/tmp/schedtest/")
  -runtimeKit string
        cpu,mem,trace
```

example:

```
[root@xxxx]# ./schedtest -namespace=e2e-tests-density-50-1 -replicas=1000 -runtimeKit=cpu,mem,trace
Start time: 2019-06-12 20:54:31.339970709 +0800 CST m=+0.009257840
End Time: 2019-06-12 20:54:50.406348385 +0800 CST m=+19.075635466
1000 scheduled for 1000 pods
E2E startup time for 1000 pods: 19.066378
Throughput (pods/s) during cluster saturation phase: 52.44835

root@xxxx]# ll /tmp/schedtest/profiles/
total 67772
-rw-r--r-- 1 root root    56509 Jun 12 20:54 kube-scheduler_CPUProfile_density_2019-06-12T20:54:32+08:00.pprof
-rw-r--r-- 1 root root    56699 Jun 12 20:54 kube-scheduler_CPUProfile_density_2019-06-12T20:54:38+08:00.pprof
-rw-r--r-- 1 root root    49617 Jun 12 20:54 kube-scheduler_CPUProfile_density_2019-06-12T20:54:44+08:00.pprof
-rw-r--r-- 1 root root    52125 Jun 12 20:54 kube-scheduler_MemoryProfile_density_2019-06-12T20:54:50+08:00.pprof
-rw-r--r-- 1 root root 23851725 Jun 12 20:54 kube-scheduler_Trace_density_2019-06-12T20:54:32+08:00.pprof
-rw-r--r-- 1 root root 24613493 Jun 12 20:54 kube-scheduler_Trace_density_2019-06-12T20:54:38+08:00.pprof
-rw-r--r-- 1 root root 20703746 Jun 12 20:54 kube-scheduler_Trace_density_2019-06-12T20:54:44+08:00.pprof
```

Next, you can use "go tool pprof", "go-torch", "go tool trace" to analyze the above pprof.

