package main

import (
	"flag"
	"fmt"

	apps_util "github.com/appscode/kutil/apps/v1"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/tamalsaha/go-oneliners"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "book-server-cli",
			Namespace: "default",
			Labels: map[string]string{
				"app" : "bookserver",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "bookserver",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "bookserver",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "book-server-cli",
							Image: "suaas21/book-server-cli:part1",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 8081,
								},
							},
						},
					},
				},
			},
		},
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-service",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app" : "bookserver",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port: 80,
					TargetPort: intstr.FromInt(8081),

				},
			},
			Type: corev1.ServiceTypeNodePort,
		},


	}

   //Create Deployment
   resDeploy, _, err := apps_util.CreateOrPatchDeployment(clientset, deployment.ObjectMeta, func(in *appsv1.Deployment) *appsv1.Deployment {
	   in.Spec = deployment.Spec

	   return in
   })
   //Display Deployment
   oneliners.PrettyJson(resDeploy)

   if err != nil {
	   fmt.Println(err.Error())
   }

   //Cerate Service
   resService,_, err:= core_util.CreateOrPatchService(clientset,service.ObjectMeta, func(in *corev1.Service) *corev1.Service {
   	in.Spec = service.Spec

   	return in

   })
   if err != nil {
   		fmt.Println(err.Error())
   }
   //Display Service
   oneliners.PrettyJson(resService)

   //Update Deployment
   patchDeploymet,_,err := apps_util.PatchDeployment(clientset, deployment, func(in *appsv1.Deployment) *appsv1.Deployment {
	  in.Spec.Replicas = int32Ptr(2)
	  return in
   })
	//Display Updated Deployment
	oneliners.PrettyJson(patchDeploymet)

	//Update Service
	patchService,_,err := core_util.PatchService(clientset, service, func(in *corev1.Service) *corev1.Service {
		in.Spec.Type = corev1.ServiceTypeNodePort
		return in
	})
	//Display Updated Service
	oneliners.PrettyJson(patchService)

	//Delete deployment
	er := apps_util.DeleteDeployment(clientset,deployment.ObjectMeta)
	if err != nil {
		fmt.Println(er.Error())
	}
	//Delete Service



}
func int32Ptr(i int32) *int32 { return &i }
