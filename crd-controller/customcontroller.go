package main

import (
	"fmt"
	crdcontrolllerv1alpha1 "github.com/suaas21/go-practice/crd-controller/pkg/apis/crd.suaas21.com/v1alpha1"
	clientset "github.com/suaas21/go-practice/crd-controller/pkg/client/clientset/versioned"
	informers "github.com/suaas21/go-practice/crd-controller/pkg/client/informers/externalversions"
	listers "github.com/suaas21/go-practice/crd-controller/pkg/client/listers/crd.suaas21.com/v1alpha1"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/util/runtime"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	kubelisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"time"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Controller struct {
	kubeclientset kubernetes.Interface
	clientset clientset.Interface

	deploymentLister kubelisters.DeploymentLister
	customDeploymentLister listers.CustomDeploymentLister

	//deploymentInformaer appsinformers.DeploymentInformer
	//customDeploymentInformer custominformers.CustomDeploymentInformer

	deploymentsInformer	cache.SharedIndexInformer
	customdeploymentInformer	cache.SharedIndexInformer
	//deploymentSynced cache.InformerSynced
	//customDeploymentSync  cache.InformerSynced

	workqueue workqueue.RateLimitingInterface
	//recorder record.EventRecorder
}

func Newcontroller(kubeclientset kubernetes.Interface, clientset clientset.Interface,
	kubeInfomerFactory kubeinformers.SharedInformerFactory,
	custonInformerFactory informers.SharedInformerFactory) *Controller{

		deployment := kubeInfomerFactory.Apps().V1().Deployments()
		custom := custonInformerFactory.Crd().V1alpha1().CustomDeployments()

		controller := &Controller{
			kubeclientset: kubeclientset,
			clientset: clientset,
			deploymentLister: deployment.Lister(),
			customDeploymentLister: custom.Lister(),

			deploymentsInformer: deployment.Informer(),
			customdeploymentInformer: custom.Informer(),

			workqueue: workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(),"customdeployments"),

			//recorder: recorder,
		}
		controller.customdeploymentInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				if key, err:= cache.MetaNamespaceKeyFunc(obj); err ==nil{
					controller.workqueue.Add(key)
				}else{
					runtime.HandleError(err)
				}

			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				newcustom := newObj.(*crdcontrolllerv1alpha1.CustomDeployment)
				oldcustom := oldObj.(*crdcontrolllerv1alpha1.CustomDeployment)

				if newcustom.ResourceVersion == oldcustom.ResourceVersion{
					return

				}else{
					if key, err := cache.MetaNamespaceKeyFunc(new); err==nil{
						controller.workqueue.Add(key)
					}

				}

			},
			DeleteFunc: func(obj interface{}) {
				if key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj); err == nil{
					controller.workqueue.Add(key)
				}

			},

		})

		//controller.deploymentsInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		//	AddFunc: func(obj interface{}) {
		//		var object metav1.Object
		//		var ok bool
		//		if object, ok = obj.(metav1.Object); ok{
		//
		//		}
		//
		//	},
		//	UpdateFunc: func(oldObj, newObj interface{}) {
		//		newdeploy := newObj.(*appsv1beta2.Deployment)
		//		olddeploy := oldObj.(*appsv1beta2.Deployment)
		//		if newdeploy.ResourceVersion == olddeploy.ResourceVersion{
		//			return
		//		}else{
		//			if
		//		}
		//
		//
		//	},
		//})

    return controller
}

func (c *Controller) Run(stopch <-chan struct{})error{
	// don't let panics crash the process
    defer runtime.HandleCrash()
    defer c.workqueue.ShutDown()
    // Start the informer factories to begin populating the informer caches
	fmt.Println("Starting Something controller")
	fmt.Println("Waiting for informer caches to sync")
    if ! cache.WaitForCacheSync(stopch, c.customdeploymentInformer.HasSynced, c.customdeploymentInformer.HasSynced){
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

func (c *Controller) processNextItem() bool{
	obj, shutdown := c.workqueue.Get()

	if shutdown{
		return false
	}
	defer c.workqueue.Done(obj)
	err := c.customSyncHandler(obj.(string))
	if err == nil{
		c.workqueue.Forget(obj)
	}else if c.workqueue.NumRequeues(obj) > 10 {
		c.workqueue.AddRateLimited(obj)

	}else{
		c.workqueue.Forget(obj)
		runtime.HandleError(err)
	}

	return true

}

func (c *Controller) customSyncHandler(key string) error {
	fmt.Println("handling the something resource named \"something-exmp\"...")
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s\n", key))
		return nil
	}

	custom, err := c.customDeploymentLister.CustomDeployments(namespace).Get(name)
	if err != nil {
		// The Something resource may no longer exist, in which case we stop processing.
		if errors.IsNotFound(err) {
			runtime.HandleError(fmt.Errorf("something '%s' in work queue (somethingsQueue) no longer exists\n", key))
			return nil
		}
		fmt.Println("err in getting something resource is ", err)

		return err
	}
	fmt.Println(custom.Namespace, "/", custom.Name)

	deployment, err := c.deploymentLister.Deployments(custom.Namespace).Get(name)

	if custom.Spec.Replicas != nil && *custom.Spec.Replicas != *deployment.Spec.Replicas{
		deployment, err := c.kubeclientset.AppsV1().Deployments(custom.Namespace).Create(NewDeployment(custom))
	}


}

func NewDeployment(custom *crdcontrolllerv1alpha1.CustomDeployment) *v1.Deployment{
	labels := map[string]string{
		"app": "bookserver",
		"controller": custom.Name,
	}
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: custom.Spec.

		},
	}
}