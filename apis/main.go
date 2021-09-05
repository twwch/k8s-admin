package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"path/filepath"
)

func int32Ptr(i int32) *int32 { return &i }

func main() {
	// https://github.com/kubernetes/client-go
	path, _ := os.Getwd()
	kubeconfig := filepath.Join(
		path, "apis/.kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	coreV1 := clientset.CoreV1()
	listOptions := metav1.ListOptions{}
	ctx := context.Background()
	namespaces, err := coreV1.Namespaces().List(ctx, listOptions)
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range namespaces.Items {
		fmt.Println(item.Name)
	}
	fmt.Println(namespaces)
	// 创建Namespaces
	/*ns, err := coreV1.Namespaces().Create(ctx, &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name:"test"},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ns.Name)*/

	pod, err := coreV1.Pods("test").Create(ctx, &v1.Pod{
		TypeMeta: metav1.TypeMeta{Kind: "pods"},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "test",
			Name:      "nginx",
		},
		Spec:v1.PodSpec{
				Containers: []v1.Container{
					{Image:"nginx", Name:"nginx", Ports:[]v1.ContainerPort{{ContainerPort: 80}}},
				},
				RestartPolicy: v1.RestartPolicyAlways,
				Tolerations: []v1.Toleration{
						{Effect:v1.TaintEffectNoSchedule, Value:"no",Key:"gpu"},
				},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(pod)

	/*deploymentsClient := clientset.AppsV1().Deployments(v1.NamespaceDefault)
	//https://github.com/kubernetes/client-go/blob/master/examples/create-update-delete-deployment/main.go
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo-deployment",
			Namespace:"test",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())*/
}
