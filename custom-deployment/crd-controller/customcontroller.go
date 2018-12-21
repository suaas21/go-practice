package main

import (
	"flag"
	"fmt"
	"github.com/appscode/go/signals"
	"github.com/tamalsaha/go-oneliners"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"
	"github.com/golang/glog"
	"log"
	"path/filepath"
	"time"

	crdcontrolllerv1alpha1 "github.com/suaas21/go-practice/custom-deployment/pkg/apis/crd.suaas21.com/v1alpha1"
	clientset "github.com/suaas21/go-practice/custom-deployment/pkg/client/clientset/versioned"
	informers "github.com/suaas21/go-practice/custom-deployment/pkg/client/informers/externalversions"
	listers "github.com/suaas21/go-practice/custom-deployment/pkg/client/listers/crd.suaas21.com/v1alpha1"
	"k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubelisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
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
		//panic(err)
		klog.Fatalf("unexpected error occured: %v", err)
	}

	kubeclient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	client, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	kubeInformerFactory := kubeinformers.NewSharedInformerFactory(kubeclient, time.Second*30)
	customInformerFactory := informers.NewSharedInformerFactory(client, time.Second*30)

	controller := Newcontroller(kubeclient, client, kubeInformerFactory, customInformerFactory)

	go kubeInformerFactory.Start(stopCh)
	go customInformerFactory.Start(stopCh)

	if err = controller.Run(stopCh); err != nil {
		log.Fatal(err)
	}

}

type Controller struct {
	kubeclientset kubernetes.Interface
	clientset     clientset.Interface

	deploymentLister       kubelisters.DeploymentLister
	customDeploymentLister listers.CustomDeploymentLister

	//deploymentInformaer appsinformers.DeploymentInformer
	//customDeploymentInformer custominformers.CustomDeploymentInformer

	deploymentsInformer      cache.SharedIndexInformer
	customdeploymentInformer cache.SharedIndexInformer
	//deploymentSynced cache.InformerSynced
	//customDeploymentSync  cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	//recorder record.EventRecorder
	podLabel string

	PreviousPodPhase map[string]string
	PodOwnerKey      map[string]string
}

func Newcontroller(kubeclientset kubernetes.Interface, clientset clientset.Interface,
	kubeInfomerFactory kubeinformers.SharedInformerFactory,
	custonInformerFactory informers.SharedInformerFactory) *Controller {

	deployment := kubeInfomerFactory.Apps().V1().Deployments()
	custom := custonInformerFactory.Crd().V1alpha1().CustomDeployments()

	controller := &Controller{
		kubeclientset:          kubeclientset,
		clientset:              clientset,
		deploymentLister:       deployment.Lister(),
		customDeploymentLister: custom.Lister(),

		//deploymentsInformer:      deployment.Informer(),
		customdeploymentInformer: custom.Informer(),

		workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "customdeployments"),

		//recorder: recorder,
	}
	controller.customdeploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			if key, err := cache.MetaNamespaceKeyFunc(obj); err == nil {
				controller.workqueue.Add(key)
			} else {
				fmt.Println("AddFunc")
				runtime.HandleError(err)
			}

		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			newcustom := newObj.(*crdcontrolllerv1alpha1.CustomDeployment)
			oldcustom := oldObj.(*crdcontrolllerv1alpha1.CustomDeployment)

			if newcustom.ResourceVersion == oldcustom.ResourceVersion {
				return

			} else {
				if key, err := cache.MetaNamespaceKeyFunc(newObj); err == nil {
					controller.workqueue.Add(key)
				}

			}

		},
		DeleteFunc: func(obj interface{}) {
			if key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err == nil {
				controller.workqueue.Add(key)
			}

		},
	})
	return controller
}

func (c *Controller) Run(stopch <-chan struct{}) error {
	// don't let panics crash the process
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()
	// Start the informer factories to begin populating the informer caches
	fmt.Println("Starting CustomDeployment controller")
	fmt.Println("Waiting for informer caches to sync")
	if !cache.WaitForCacheSync(stopch, c.customdeploymentInformer.HasSynced) {
		fmt.Println("Timed out waiting for caches to sync")
		return fmt.Errorf("%s", "Timed out waiting for caches to sync")
	}
	fmt.Println("Starting workers")
	go wait.Until(c.runWorker, time.Second, stopch)

	fmt.Println("Started workers")
	<-stopch
	fmt.Println("Shutting down workers")

	return nil
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}
	defer c.workqueue.Done(obj)
	err := c.customSyncHandler(obj.(string))
	if err == nil {
		c.workqueue.Forget(obj)
	} else if c.workqueue.NumRequeues(obj) < 10 {
		c.workqueue.AddRateLimited(obj)

	} else {
		c.workqueue.Forget(obj)
		runtime.HandleError(err)
	}

	return true

}

func (c *Controller) customSyncHandler(key string) error {
	fmt.Println("handling the customdeployment resource named ...")
	//glog.Infoln("handling the customdeployment resource named ...")
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		//glog.Infoln("error in spliting")
		runtime.HandleError(fmt.Errorf("invalid resource key: %s\n", key))
		return nil
	}

	//glog.Infoln("splited into", namespace, name)
	custom, err := c.customDeploymentLister.CustomDeployments(namespace).Get(name)
	if err != nil {
		// The Something resource may no longer exist, in which case we stop processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("something '%s' in work queue (workqueue) no longer exists\n", key))
			return nil
		}
		fmt.Println("err in getting something resource is ", err)
		//glog.Infoln("err in getting something resource: %v", err)

		return err
	}

	fmt.Println(custom.Namespace, "/", custom.Name)
	deploymentName := custom.Spec.Name
	if deploymentName == "" {

		//glog.Infoln(fmt.Errorf("%s: deployment name must be specified\n", key))
		runtime.HandleError(fmt.Errorf("%s: deployment name must be specified\n", key))
		return nil
	}
	fmt.Printf("deployment Name -> %s\n", deploymentName)
	//glog.Infoln("deployment Name -> %s\n", deploymentName)
	//oneliners.PrettyJson(custom, "custom")

	deployment, err := c.deploymentLister.Deployments(custom.Namespace).Get(deploymentName)

	//glog.Infoln(err)

	if err != nil {
		if errors.IsNotFound(err) {
			deployment, err = c.kubeclientset.AppsV1().Deployments(custom.Namespace).Create(NewDeployment(custom))
			if err != nil {
				log.Fatal(err)
			}
		}

		log.Println(err)
		//glog.Infoln(err)
		return err
	}

	oneliners.PrettyJson(custom, "custom")
	//oneliners.PrettyJson(deployment, "deployment")

	if custom.Spec.Replicas != nil && *custom.Spec.Replicas != *deployment.Spec.Replicas {
		fmt.Println("customdeployment: %d, deployR: %d", *custom.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1().Deployments(custom.Namespace).Update(NewDeployment(custom))
	}
	if err != nil {
		fmt.Println("error occured in updating deployment", deployment.Name, "is", err)
		return err
	}
	fmt.Println("no error in updating deployment", deployment.Name)
	glog.Infoln("no error in updating deployment", deployment.Name)
	err = c.updateCustomStatus(custom, deployment)
	if err != nil {
		fmt.Println("error occured in updating status of something", custom.Name, "is", err)
		return err
	}

	return nil
}

func (c *Controller) updateCustomStatus(custom *crdcontrolllerv1alpha1.CustomDeployment, deployment *v1.Deployment) error {
	customcopy := custom.DeepCopy()
	customcopy.Status.AvailableReplicas = deployment.Status.AvailableReplicas

	_, err := c.clientset.CrdV1alpha1().CustomDeployments(custom.Namespace).Update(customcopy)

	return err
}

func NewDeployment(custom *crdcontrolllerv1alpha1.CustomDeployment) *v1.Deployment {
	//labels := map[string]string{
	//	"app": "bookserver",
	//	"controller": custom.Name,
	//}
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      custom.Spec.Name,
			Namespace: custom.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: custom.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: custom.Spec.Selector.MatchLabels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: custom.Spec.Template.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  custom.Spec.Template.Spec.Containers[0].Name,
							Image: custom.Spec.Template.Spec.Containers[0].Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          custom.Spec.Template.Spec.Containers[0].Ports[0].Name,
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: custom.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}
}
