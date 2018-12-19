package customdeploy

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"log"

	//"flag"
	//"fmt"
	core_util "github.com/appscode/kutil/core/v1"
	csappV1 "github.com/suaas21/go-practice/custom-deployment/pkg/apis/crd.suaas21.com/v1alpha1"
	clientset "github.com/suaas21/go-practice/custom-deployment/pkg/client/clientset/versioned"
	"github.com/tamalsaha/go-oneliners"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/kubernetes"
	//"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/util/homedir"
	//"path/filepath"
)

//var kubeclient *kubernetes.Clientset

func CreateCustomDeployment(client *clientset.Clientset) {

	customdeployment := &csappV1.CustomDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "bookserver",
			Namespace: "default",
			Labels: map[string]string{
				"app": "bookserver",
			},
		},
		Spec: csappV1.CustomDeploymentSpec{
			Name:     "bookserver",
			Replicas: int32Ptr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "bookserver",
				},
			},
			Template: csappV1.CustomPodTemplate{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "bookserver",
					},
				},

				Spec: csappV1.PodSpec{
					Containers: []csappV1.Container{
						{
							Name:  "book-server-cli",
							Image: "suaas21/book-server-cli:part1",
							Ports: []csappV1.ContainerPort{
								{
									Name:          "http",
									Protocol:      csappV1.ProtocolTCP,
									ContainerPort: 8081,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := client.CrdV1alpha1().CustomDeployments("default").Create(customdeployment)
	if err != nil {
		panic(err)
	}
	oneliners.PrettyJson(result)
	fmt.Println(result.GetObjectMeta().GetName())
	oneliners.PrettyJson(result)

}

func DeleteAll(client *clientset.Clientset, kubeClient kubernetes.Interface) {
	//Delete Customdeployment
	err := client.CrdV1alpha1().CustomDeployments("default").Delete("bookserver", &metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}

	// delete service
	kubeclient := kubeClient.CoreV1().Services(corev1.NamespaceDefault)
	log.Println("Deleting service...")
	if err := kubeclient.Delete("svc-bookserver", &metav1.DeleteOptions{
		//PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Fatal(err)
	}
	log.Println("service deleted")

}

func CreateService(kubeClient kubernetes.Interface) {
	fmt.Println("In CreateService function--->")

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "svc-bookserver",
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": "bookserver",
			},
			Ports: []corev1.ServicePort{
				{
					Protocol:   corev1.ProtocolTCP,
					Port:       80,
					TargetPort: intstr.FromInt(8081),
				},
			},
			Type: corev1.ServiceTypeLoadBalancer,
		},
	}


	//Cerate Service
	resService,_, err:= core_util.CreateOrPatchService(kubeClient,svc.ObjectMeta, func(in *corev1.Service) *corev1.Service {
		in.Spec = svc.Spec

		return in

	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Printf("Created service %q.\n", resService.GetObjectMeta().GetName())

	// The url at which we can access the now
	//node, err := kubeclient.corev1().Nodes().Get("minikube", metav1.GetOptions{})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Print("The url at which we can access the now is")
	//fmt.Printf("%v:%v\n", node.Status.Addresses[0].Address, result.Spec.Ports[0].NodePort)

}

func int32Ptr(i int32) *int32 { return &i }
