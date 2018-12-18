package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	crdcontrolllerv1alpha1 "github.com/suaas21/go-practice/crd-controller/pkg/apis/crd.suaas21.com/v1alpha1"
	clientset "github.com/suaas21/go-practice/crd-controller/pkg/client/clientset/versioned"
	informers "github.com/suaas21/go-practice/crd-controller/pkg/client/informers/externalversions"
	listers "github.com/suaas21/go-practice/crd-controller/pkg/client/listers/crd.suaas21.com/v1alpha1"
	"k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	kubelisters "k8s.io/client-go/listers/apps/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

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
	podLabel          string

	PreviousPodPhase map[string]string
	PodOwnerKey      map[string]string
}

func Newcontroller(kubeclientset kubernetes.Interface, clientset clientset.Interface,
	//kubeInfomerFactory kubeinformers.SharedInformerFactory,
	custonInformerFactory informers.SharedInformerFactory) *Controller {

	//deployment := kubeInfomerFactory.Apps().V1().Deployments()
	custom := custonInformerFactory.Crd().V1alpha1().CustomDeployments()

	controller := &Controller{
		kubeclientset:          kubeclientset,
		clientset:              clientset,
		//deploymentLister:       deployment.Lister(),
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

func (c *Controller) Run(stopch <-chan struct{}) error {
	// don't let panics crash the process
	defer runtime.HandleCrash()
	defer c.workqueue.ShutDown()
	// Start the informer factories to begin populating the informer caches
	fmt.Println("Starting CustomDeployment controller")
	fmt.Println("Waiting for informer caches to sync")
	if !cache.WaitForCacheSync(stopch, c.customdeploymentInformer.HasSynced, c.customdeploymentInformer.HasSynced) {
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
	fmt.Println("handling the customdeployment resource named \"something-exmp\"...")
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
			runtime.HandleError(fmt.Errorf("something '%s' in work queue (workqueue) no longer exists\n", key))
			return nil
		}
		fmt.Println("err in getting something resource is ", err)

		return err
	}
	customdeployment := custom.DeepCopy()

	if customdeployment.Status.AvailableReplicas + customdeployment.Status.CreatingReplicas < *customdeployment.Spec.Replicas{
		pod, err := c.CreateNewPod(customdeployment.Spec.Template, customdeployment)

		label := ""
		if err != nil{
			fmt.Printf("Can't create pod. Reason: %v\n", err.Error())
			return err
		}
		podName := string(pod.GetName())
		c.PodOwnerKey[podName] = key
		c.PreviousPodPhase[podName] = "Creating"

		mp := pod.GetLabels()
        for key, value := range mp{
			label += key + "=" + value
		}
        c.podLabel =label
		err2 := c.updateCustomStatus(customdeployment)

		if err2 != nil {
			fmt.Printf("Pod created but failed to update DeploymentStatus.")
			return err2
		}

	}else if customdeployment.Status.AvailableReplicas+customdeployment.Status.CreatingReplicas > *customdeployment.Spec.Replicas{
		err := c.DeletePod(customdeployment.Status.AvailableReplicas + customdeployment.Status.CreatingReplicas - *customdeployment.Spec.Replicas)

		if err != nil {
			fmt.Println("Can't Delete Pod. Reason: ", err)
			return err
		}

		err = c.updateCustomStatus(customdeployment)
		if err != nil {
			fmt.Println("Failed to update DeploymentStatus.")
			return err
		}

	}else{
		//everything ok
	}


	//fmt.Println(custom.Namespace, "/", custom.Name)
	//deploymentName := custom.Spec.Name
	//if deploymentName == "" {
	//
	//	runtime.HandleError(fmt.Errorf("%s: deployment name must be specified\n", key))
	//	return nil
	//}
	//fmt.Printf("deployment Name -> %s\n", deploymentName)
	//
	//deployment, err := c.deploymentLister.Deployments(custom.Namespace).Get(deploymentName)
	//
	//fmt.Println(err)
	//
	//if err != nil {
	//	if errors.IsNotFound(err) {
	//		deployment, err = c.kubeclientset.AppsV1().Deployments(custom.Namespace).Create(NewDeployment(custom))
	//		if err != nil {
	//			fmt.Errorf("====================", err)
	//		}
	//	}
	//
	//	log.Println("============Error===========", err)
	//	return err
	//}
	//
	//oneliners.PrettyJson(custom, "custom")
	//oneliners.PrettyJson(deployment, "deployment")
	//
	//if custom.Spec.Replicas != nil && *custom.Spec.Replicas != *deployment.Spec.Replicas {
	//	fmt.Println("customdeployment: %d, deployR: %d", *custom.Spec.Replicas, *deployment.Spec.Replicas)
	//	deployment, err = c.kubeclientset.AppsV1().Deployments(custom.Namespace).Update(NewDeployment(custom))
	//}
	//if err != nil {
	//	fmt.Println("error occured in updating deployment", deployment.Name, "is", err)
	//	return err
	//}
	//fmt.Println("no error in updating deployment", deployment.Name)
	//err = c.updateCustomStatus(custom, deployment)
	//if err != nil {
	//	fmt.Println("error occured in updating status of something", custom.Name, "is", err)
	//	return err
	//}

	return nil
}

func (c *Controller) CreateNewPod(podTemplate crdcontrolllerv1alpha1.CustomPodTemplate, customdeployment *crdcontrolllerv1alpha1.CustomDeployment)(*apiv1.Pod, error){
	podClient := c.kubeclientset.CoreV1().Pods(apiv1.NamespaceDefault)

	pod := &apiv1.Pod{
		ObjectMeta:metav1.ObjectMeta{
		   Name:customdeployment.GetName()+"-"+strconv.Itoa(rand.Int()),
		   Labels: podTemplate.GetObjectMeta().GetLabels(),
		},
		//Spec: podTemplate.Spec,

	}
	newPod, err := podClient.Create(pod)
	if err == nil{
		fmt.Printf("New pod with name %v has been created.\n", newPod.GetName())
	}

	return newPod, err
}
func (c *Controller) DeletePod(deletionLimit int32) error{

	podClient := c.kubeclientset.CoreV1().Pods(apiv1.NamespaceDefault)

	podList, err := c.kubeclientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{LabelSelector: c.podLabel})
	if err != nil {
		fmt.Println("Can't get pod list. Reason: ", err)
	}

	deletedPod := int32(0)

	for _, pod := range podList.Items {

		delErr := podClient.Delete(pod.GetName(), &metav1.DeleteOptions{})

		if delErr != nil {
			return delErr
		} else {
			c.PreviousPodPhase[pod.GetName()] = "Terminating"
			deletedPod++
			if deletedPod >= deletionLimit {
				break
			}
		}
	}

	return nil
}

func (c *Controller) updateCustomStatus(customdeployment *crdcontrolllerv1alpha1.CustomDeployment) error {

	creating := 0

	podList, err := c.kubeclientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{LabelSelector: c.podLabel})
	if err != nil {
		fmt.Println("Can't get pod list. Reason: ", err)
	}

	for _, pod := range podList.Items {

		if c.PreviousPodPhase[pod.GetName()] == "Creating" {
			creating++
		}
	}

	//Don't modify cache. Work on it's copy
	customdeploymentCopy := customdeployment.DeepCopy()

	customdeploymentCopy.Spec.Replicas = customdeployment.Spec.Replicas
	customdeploymentCopy.Status.CreatingReplicas = int32(creating)


	//Now update the cache
	_, err = c.clientset.CrdV1alpha1().CustomDeployments(apiv1.NamespaceDefault).Update(customdeploymentCopy)

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
				Spec: custom.Spec.Template.Spec,
				//Spec: apiv1.PodSpec{
					//Containers: []apiv1.Container{
					//	{
					//		Name:  custom.Spec.Template.Spec.Containers[0].Name,
					//		Image: custom.Spec.Template.Spec.Containers[0].Image,
					//		Ports: []apiv1.ContainerPort{
					//			{
					//				Name:          custom.Spec.Template.Spec.Containers[0].Ports[0].Name,
					//				Protocol:      apiv1.ProtocolTCP,
					//				ContainerPort: custom.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort,
					//			},
					//		},
					//	},
					//},
				//},
			},
		},
	}
}