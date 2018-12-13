package main

import (
	"flag"
	"fmt"

	csappV1 "github.com/suaas21/go-practice/crd-controller/pkg/apis/crd.suaas21.com/v1alpha1"
	cs "github.com/suaas21/go-practice/crd-controller/pkg/client/clientset/versioned"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	_ "k8s.io/code-generator"
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

	cs, err := cs.NewForConfig(config)

	if err != nil{
		panic(err)
	}

	customdeployment := &csappV1.CustomDeployment{
		ObjectMeta: metav1.ObjectMeta{
		  Name: "bookserver",
		  Namespace: "default",
		  Labels: map[string]string{
		  	"app" : "bookserver",
		  },
		},
		Spec:csappV1.CustomDeploymentSpec{
		  Replicas: int32Ptr(3),
		  Selector: &metav1.LabelSelector{
		  	MatchLabels: map[string]string{
		  		"app" : "bookserver",
			},
		  },
		  Template:csappV1.CustomPodTemplate{
		  	ObjectMeta: metav1.ObjectMeta{
		  		Labels: map[string]string{
		  			"app" : "bookserver",
				},
			},
			Spec:csappV1.PodSpec{
				Containers:[]csappV1.Container{
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
    //Create Customdeployment
    result, err := cs.CrdV1alpha1().CustomDeployments("default").Create(customdeployment)
	if err != nil{
		panic(err)
	}
    fmt.Println("Created Cuntomdeployment %q.\n",result.GetObjectMeta().GetName())

    //Update Customdeployment
    update, err := cs.CrdV1alpha1().CustomDeployments("default").Update(customdeployment)
	if err != nil{
		panic(err)
	}
	fmt.Println("Updated Cuntomdeployment %q.\n",update.GetObjectMeta().GetName())

    //Delete Customdeployment
    er := cs.CrdV1alpha1().CustomDeployments("default").Delete(customdeployment.Name, &metav1.DeleteOptions{})
	if er != nil{
		panic(err)
	}



}
func int32Ptr(i int32) *int32 { return &i }
