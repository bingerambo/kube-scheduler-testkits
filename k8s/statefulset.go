package k8s

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	apiresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func RetrieveSts(namespace, name string) (*appsv1.StatefulSet, error) {
	cs := GetClient()

	sts, err := cs.AppsV1().StatefulSets(namespace).Get(name, metav1.GetOptions{})

	return sts, err
}

func CreateSts(stsName, ns string, replicas int) error {
	cs := GetClient()

	stsImage := "k8s.gcr.io/pause:3.1"

	cpuFloat := 0.01
	memFloat := 40265318

	lables := map[string]string{
		"app":  stsName,
		"name": stsName,
		"type": "densityPod",
	}

	replicasCopy := int32(replicas)

	sts := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name: stsName,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas:            &replicasCopy,
			PodManagementPolicy: appsv1.ParallelPodManagement,
			UpdateStrategy: appsv1.StatefulSetUpdateStrategy{
				Type: "OnDelete",
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": stsName,
				},
			},
			ServiceName: stsName,
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: lables,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  stsName,
							Image: stsImage,
							Ports: []apiv1.ContainerPort{{ContainerPort: 80}},
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									apiv1.ResourceCPU:    *apiresource.NewMilliQuantity(int64(cpuFloat*1000), apiresource.DecimalSI),
									apiv1.ResourceMemory: *apiresource.NewQuantity(int64(memFloat), apiresource.DecimalSI),
								},
							},
						},
					},
					DNSPolicy: apiv1.DNSDefault,
					Affinity:  genAffinity(stsName),
				},
			},
		},
	}

	sts, err := cs.AppsV1().StatefulSets(ns).Create(sts)
	if err != nil {
		return err
	}

	return nil
}

func DeleteSts(namespace, name string) error {
	cs := GetClient()

	deletePolicy := metav1.DeletePropagationForeground
	err := cs.AppsV1().StatefulSets(namespace).Delete(name, &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})

	if errors.IsNotFound(err) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func genAffinity(name string) *apiv1.Affinity {
	affinity := &apiv1.Affinity{
		PodAntiAffinity: &apiv1.PodAntiAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []apiv1.WeightedPodAffinityTerm{
				{
					Weight: 100,
					PodAffinityTerm: apiv1.PodAffinityTerm{
						LabelSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app": name,
							},
						},
						TopologyKey: "kubernetes.io/hostname",
					},
				},
			},
		},
	}
	return affinity
}
