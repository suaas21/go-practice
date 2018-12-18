package main

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"log"
	"time"

	clientset "github.com/suaas21/go-practice/crd-controller/pkg/client/clientset/versioned"
	informers "github.com/suaas21/go-practice/crd-controller/pkg/client/informers/externalversions"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/sample-controller/pkg/signals"
	"path/filepath"
)
var (
	kubeconfig *string
)

func main() {

	stopCh := signals.SetupSignalHandler()

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

	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err)
	}
	client, err := clientset.NewForConfig(config)
	if err != nil{
		panic(err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeclient, time.Second*30)
	customInformerFactory := informers.NewSharedInformerFactory(client, time.Second*30)

	controller := Newcontroller(kubeclient, client, customInformerFactory)

	go kubeInformerFactory.Start(stopCh)
	go customInformerFactory.Start(stopCh)


    if err = controller.Run(stopCh); err != nil{
    	log.Fatal(err)
	}


	fmt.Println("go-client")
   //
	//customdeployment := &csappV1.CustomDeployment{
	//	ObjectMeta: metav1.ObjectMeta{
	//	  Name: "bookserver",
	//	  Namespace: "default",
	//	  Labels: map[string]string{
	//	  	"app" : "bookserver",
	//	  },
	//	},
	//	Spec:csappV1.CustomDeploymentSpec{
	//		Name: "bookserver",
	//	  Replicas: int32Ptr(3),
	//	  Selector: &metav1.LabelSelector{
	//	  	MatchLabels: map[string]string{
	//	  		"app" : "bookserver",
	//		},
	//	  },
	//	  Template:csappV1.CustomPodTemplate{
	//	  	ObjectMeta: metav1.ObjectMeta{
	//	  		Labels: map[string]string{
	//	  			"app" : "bookserver",
	//			},
	//		},
	//
	//		//Spec:csappV1.PodSpec{
	//		//	Containers:[]csappV1.Container{
	//		//		{
	//		//			Name:  "book-server-cli",
	//		//			Image: "suaas21/book-server-cli:part1",
	//		//			Ports: []csappV1.ContainerPort{
	//		//				{
	//		//					Name:          "http",
	//		//					Protocol:      csappV1.ProtocolTCP,
	//		//					ContainerPort: 8081,
	//		//				},
	//		//			},
	//		//		},
	//		//	},
	//		//},
	//	  },
	//	},
   //
	//}
	//fmt.Println("go-client")
   // //Create Customdeployment
   // result, err := client.CrdV1alpha1().CustomDeployments("default").Create(customdeployment)
	//if err != nil{
	//	panic(err)
	//}
	//fmt.Println(result.GetObjectMeta().GetName())
	//oneliners.PrettyJson(result)
   // //Update Customdeployment
   //// update, err := cs.CrdV1alpha1().CustomDeployments("default").Update()
	////if err != nil{
	////	panic(err)
	////}
	////fmt.Println("Updated Cuntomdeployment %q.\n",update.GetObjectMeta().GetName())
   //
   // //Delete Customdeployment
   // er := client.CrdV1alpha1().CustomDeployments("default").Delete(customdeployment.Name, &metav1.DeleteOptions{})
	//if er != nil{
	//	panic(err)
	//}

}
func int32Ptr(i int32) *int32 { return &i }
