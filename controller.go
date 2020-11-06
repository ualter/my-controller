package main

import (
	"context"
	"fmt"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	appsinformers "k8s.io/client-go/informers/apps/v1"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	daemonstoolv1beta1 "my-controller/pkg/apis/daemonstool/v1beta1"
	clientset "my-controller/pkg/generated/clientset/versioned"
	samplescheme "my-controller/pkg/generated/clientset/versioned/scheme"
	informers "my-controller/pkg/generated/informers/externalversions/daemonstool/v1beta1"
	listers "my-controller/pkg/generated/listers/daemonstool/v1beta1"

	appslisters "k8s.io/client-go/listers/apps/v1"
)

// Controller is the controller :-)
type Controller struct {
	// kubeclientset is a standard kubernetes clientset
	kubeclientset kubernetes.Interface
	// sampleclientset is a clientset for our own API group
	sampleclientset   clientset.Interface
	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced
	daemonstoolLister listers.DaemonstoolLister
	daemonstoolSynced cache.InformerSynced
	workqueue         workqueue.RateLimitingInterface
	recorder          record.EventRecorder
}

const controllerAgentName = "my-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

// NewController returns a new sample controller
func NewController(
	kubeclientset kubernetes.Interface,
	sampleclientset clientset.Interface,
	deploymentInformer appsinformers.DeploymentInformer,
	daemonstoolInformer informers.DaemonstoolInformer) *Controller {

	// Create event broadcaster
	// Add my-controller types to the default Kubernetes Scheme so Events can be
	// logged for my-controller types.
	utilruntime.Must(samplescheme.AddToScheme(scheme.Scheme))
	klog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartStructuredLogging(0)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	controller := &Controller{
		kubeclientset:     kubeclientset,
		sampleclientset:   sampleclientset,
		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,
		daemonstoolLister: daemonstoolInformer.Lister(),
		daemonstoolSynced: daemonstoolInformer.Informer().HasSynced,
		workqueue:         workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Daemonstool"),
		recorder:          recorder,
	}

	klog.Info("Setting up event handlers")
	// Set up an event handler for when Daemonstool resources change
	daemonstoolInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueueFoo,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueueFoo(new)
		},
	})

	// Set up an event handler for when Deployment resources change. This
	// handler will lookup the owner of the given Deployment, and if it is
	// owned by a Daemonstool resource will enqueue that Daemonstool resource for
	// processing. This way, we don't need to implement custom logic for
	// handling Deployment resources. More info on this pattern:
	// https://github.com/kubernetes/community/blob/8cafef897a22026d42f5e5bb3f104febe7e29830/contributors/devel/controllers.md
	deploymentInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.handleObject,
		UpdateFunc: func(old, new interface{}) {
			newDepl := new.(*appsv1.Deployment)
			oldDepl := old.(*appsv1.Deployment)
			if newDepl.ResourceVersion == oldDepl.ResourceVersion {
				// Periodic resync will send update events for all known Deployments.
				// Two different versions of the same Deployment will always have different RVs.
				return
			}
			controller.handleObject(new)
		},
		DeleteFunc: controller.handleObject,
	})

	return controller
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	klog.Info("Starting Daemonstool controller")

	klog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.deploymentsSynced, c.daemonstoolSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Info("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Info("Started workers")
	<-stopCh
	klog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two. It then updates the Status block of the Foo resource
// with the current status of the resource.
func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}

	// Get the Daemonstool resource with this namespace/name
	daemonstool, err := c.daemonstoolLister.Daemonstools(namespace).Get(name)
	if err != nil {
		// The Daemonstool resource may no longer exist, in which case we stop
		// processing.
		if errors.IsNotFound(err) {
			utilruntime.HandleError(fmt.Errorf("daemonstool '%s' in work queue no longer exists", key))
			return nil
		}

		return err
	}

	deploymentName := daemonstool.Spec.Image
	if deploymentName == "" {
		// We choose to absorb the error here as the worker would requeue the
		// resource otherwise. Instead, the next time the resource is updated
		// the resource will be queued again.
		utilruntime.HandleError(fmt.Errorf("%s: deployment name must be specified", key))
		return nil
	}

	// Get the deployment with the name specified in Foo.spec
	deployment, err := c.deploymentsLister.Deployments(daemonstool.Namespace).Get(deploymentName)
	// If the resource doesn't exist, we'll create it
	if errors.IsNotFound(err) {
		deployment, err = c.kubeclientset.AppsV1().Deployments(daemonstool.Namespace).Create(context.TODO(), newDeployment(daemonstool), metav1.CreateOptions{})
	}

	// If an error occurs during Get/Create, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// If the Deployment is not controlled by this Daemonstool resource, we should log
	// a warning to the event recorder and return error msg.
	if !metav1.IsControlledBy(deployment, daemonstool) {
		msg := fmt.Sprintf(MessageResourceExists, deployment.Name)
		c.recorder.Event(daemonstool, corev1.EventTypeWarning, ErrResourceExists, msg)
		return fmt.Errorf(msg)
	}

	// If this number of the replicas on the Daemonstool resource is specified, and the
	// number does not equal the current desired replicas on the Deployment, we
	// should update the Deployment resource.
	/*if daemonstool.Spec.Replicas != nil && *foo.Spec.Replicas != *deployment.Spec.Replicas {
		klog.V(4).Infof("Foo %s replicas: %d, deployment replicas: %d", name, *foo.Spec.Replicas, *deployment.Spec.Replicas)
		deployment, err = c.kubeclientset.AppsV1().Deployments(foo.Namespace).Update(context.TODO(), newDeployment(foo), metav1.UpdateOptions{})
	}*/

	// If an error occurs during Update, we'll requeue the item so we can
	// attempt processing again later. This could have been caused by a
	// temporary network failure, or any other transient reason.
	if err != nil {
		return err
	}

	// Finally, we update the status block of the Foo resource to reflect the
	// current state of the world
	err = c.updateDaemonstoolStatus(daemonstool, deployment)
	if err != nil {
		return err
	}

	c.recorder.Event(daemonstool, corev1.EventTypeNormal, SuccessSynced, MessageResourceSynced)

	return nil
}

func (c *Controller) updateDaemonstoolStatus(daemonstool *daemonstoolv1beta1.Daemonstool, deployment *appsv1.Deployment) error {
	// NEVER modify objects from the store. It's a read-only, local cache.
	// You can use DeepCopy() to make a deep copy of original object and modify this copy
	// Or create a copy manually for better performance
	daemonstoolCopy := daemonstool.DeepCopy()
	daemonstoolCopy.Status.Installed = 1
	//.AvailableReplicas = deployment.Status.AvailableReplicas
	// If the CustomResourceSubresources feature gate is not enabled,
	// we must use Update instead of UpdateStatus to update the Status block of the Foo resource.
	// UpdateStatus will not allow changes to the Spec of the resource,
	// which is ideal for ensuring nothing other than resource status has been updated.
	_, err := c.sampleclientset.UalterV1beta1().Daemonstools(daemonstool.Namespace).Update(context.TODO(), daemonstoolCopy, metav1.UpdateOptions{})
	return err
}

// enqueueFoo takes a Daemonstool resource and converts it into a namespace/name
// string which is then put onto the work queue. This method should *not* be
// passed resources of any type other than Daemonstool.
func (c *Controller) enqueueFoo(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// handleObject will take any resource implementing metav1.Object and attempt
// to find the Daemonstool resource that 'owns' it. It does this by looking at the
// objects metadata.ownerReferences field for an appropriate OwnerReference.
// It then enqueues that Daemonstool resource to be processed. If the object does not
// have an appropriate OwnerReference, it will simply be skipped.
func (c *Controller) handleObject(obj interface{}) {
	var object metav1.Object
	var ok bool
	if object, ok = obj.(metav1.Object); !ok {
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object, invalid type"))
			return
		}
		object, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			utilruntime.HandleError(fmt.Errorf("error decoding object tombstone, invalid type"))
			return
		}
		klog.V(4).Infof("Recovered deleted object '%s' from tombstone", object.GetName())
	}
	klog.V(4).Infof("Processing object: %s", object.GetName())
	if ownerRef := metav1.GetControllerOf(object); ownerRef != nil {
		// If this object is not owned by a Daemonstool, we should not do anything more
		// with it.
		if ownerRef.Kind != "Daemonstool" {
			return
		}

		foo, err := c.daemonstoolLister.Daemonstools(object.GetNamespace()).Get(ownerRef.Name)
		if err != nil {
			klog.V(4).Infof("ignoring orphaned object '%s' of foo '%s'", object.GetSelfLink(), ownerRef.Name)
			return
		}

		c.enqueueFoo(foo)
		return
	}
}

// newDeployment creates a new Deployment for a Daemonstool resource. It also sets
// the appropriate OwnerReferences on the resource so handleObject can discover
// the Foo resource that 'owns' it.
func newDeployment(daemonstool *daemonstoolv1beta1.Daemonstool) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "nginx",
		"controller": daemonstool.Name,
	}
	replicas := int32(1)
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonstool.Spec.Image,
			Namespace: daemonstool.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(daemonstool, daemonstoolv1beta1.SchemeGroupVersion.WithKind("Daemonstool")),
			},
		},
		Spec: appsv1.DeploymentSpec{
			//Replicas: daemonstool.Spec.Replicas,
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "nginx",
							Image: "nginx:latest",
						},
					},
				},
			},
		},
	}
}
