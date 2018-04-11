package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	lister_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}

var Version string = "none" // Set at compile time

func main() {
	// When running as a pod in-cluster, a kubeconfig is not needed. Instead
	// this will make use of the service account injected into the pod.
	// However, allow the use of a local kubeconfig as this can make local
	// development & testing easier.
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig file")

	// We log to stderr because glog will default to logging to a file.
	// By setting this debugging is easier via `kubectl logs`
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Build the client config - optionally using a provided kubeconfig file.
	config, err := GetClientConfig(*kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to load client config: %v", err)
	}

	// Construct the Kubernetes client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %v", err)
	}

	stopCh := make(chan struct{})
	defer close(stopCh)

	newDotmeshController(client).Run(stopCh)
}

type dotmeshController struct {
	client     kubernetes.Interface
	nodeLister lister_v1.NodeLister
	informer   cache.Controller
}

func newDotmeshController(client kubernetes.Interface) *dotmeshController {
	rc := &dotmeshController{
		client: client,
	}

	// TODO: listen for node_list_changed, dotmesh_pod_list_changed,
	// config_changed or pv_list_changed

	indexer, informer := cache.NewIndexerInformer(
		&cache.ListWatch{
			ListFunc: func(lo meta_v1.ListOptions) (runtime.Object, error) {
				// We do not add any selectors because we want to watch all
				// nodes.  This is so we can determine the total count of
				// "unavailable" nodes.  However, this could also be
				// implemented using multiple informers (or better,
				// shared-informers)
				return client.Core().Nodes().List(lo)
			},
			WatchFunc: func(lo meta_v1.ListOptions) (watch.Interface, error) {
				return client.Core().Nodes().Watch(lo)
			},
		},
		// The types of objects this informer will return
		&v1.Node{},
		// The resync period of this object. This will force a re-update of all
		// cached objects at this interval.  Every object will trigger the
		// `Updatefunc` even if there have been no actual updates triggered.
		// In some cases you can set this to a very high interval - as you can
		// assume you will see periodic updates in normal operation.  The
		// interval is set low here for demo purposes.
		10*time.Second,
		// Callback Functions to trigger on add/update/delete
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				glog.Info("NODE ADD %+v", obj)
				rc.scheduleUpdate()
			},
			UpdateFunc: func(old, new interface{}) {
				glog.Info("NODE UPDATE %+v -> %+v", old, new)
				rc.scheduleUpdate()
			},
			DeleteFunc: func(obj interface{}) {
				glog.Info("NODE DELETE %+v", obj)
				rc.scheduleUpdate()
			},
		},
		cache.Indexers{},
	)

	rc.informer = informer
	// NodeLister avoids some boilerplate code (e.g. convert runtime.Object to
	// *v1.node)
	rc.nodeLister = lister_v1.NewNodeLister(indexer)

	return rc
}

func (c *dotmeshController) scheduleUpdate() {
	// FIXME: Set a flag so that a goroutine somewhere will notice we need to do The Algorithm
	// because something may have changed.
}

func (c *dotmeshController) Run(stopCh chan struct{}) {
	glog.Info("Starting RebootController")

	go c.informer.Run(stopCh)

	// Wait for all caches to be synced, before processing is started
	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		glog.Error(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	go wait.Until(c.runWorker, time.Second, stopCh)

	<-stopCh
	glog.Info("Stopping Reboot Controller")
}

func (c *dotmeshController) runWorker() {
	// FIXME: Block until the flag set by scheduleUpdate is set and then run The Algorithm
	// because something may have changed.
	// while not time to quit {
	// block on flag
	// run process()
}

func (c *dotmeshController) process() error {
	/*
		TODO: implement the following algorithm in pseudocode.

		nodes = list_of_eligible_nodes_in_az(config.node_eligibility_filter)

		for node in nodes:
			if not labelled(node):
				label(node)

		pvs = list_of_dotmesh_pvs_in_az()
		dotmeshes = list_of_dotmesh_instances_in_az()
		orphan_dotmeshes = filter(node_of_dotmesh not in nodes, dotmeshes)
		unused_pvs = pvs - map(pv_of_dotmesh, dotmesh_nodes)
		unnoded_dotmeshes = filter(node_of_dotmesh == null, dotmeshes)

		undotted_nodes = nodes - map(node_of_dotmesh, dotmeshes)
		for each node in undotted_nodes:
			 if config.use_pv_storage:
		if len(unused_pvs) == 0:
			unused_pvs.push(make_new_pv(config.new_pv_size))
				storage = unused_pvs.pop()
			else:
				storage = local()
			create_dotmesh_with_existing_pv(storage, node)

		for each dotmesh in unnoded_dotmeshes:
			delete_dotmesh(dotmesh)

	*/
	//	node, err := c.nodeLister.Get(key)
	//	if err != nil {
	//		return fmt.Errorf("failed to retrieve node by key %q: %v", key, err)
	//	}

	// glog.V(4).Infof("Received update of node: %s", node.GetName())
	// if node.Annotations == nil {
	// 	return nil // If node has no annotations, then it doesn't need a reboot
	// }

	// if _, ok := node.Annotations[RebootNeededAnnotation]; !ok {
	// 	return nil // Node does not need reboot
	// }

	// // Determine if we should reboot based on maximum number of unavailable nodes
	// unavailable := 0
	// if err != nil {
	// 	return fmt.Errorf("Failed to determine number of unavailable nodes: %v", err)
	// }

	// if unavailable >= maxUnavailable {
	// 	// TODO(aaron): We might want this case to retry indefinitely. Could
	// 	// create a specific error an check in handleErr()
	// 	return fmt.Errorf(
	// 		"Too many nodes unvailable (%d/%d). Skipping reboot of %s",
	// 		unavailable, maxUnavailable, node.Name,
	// 	)
	// }

	// // We should not modify the cache object directly, so we make a copy first
	// nodeCopy := node.DeepCopy()

	// glog.Infof("Marking node %s for reboot", node.Name)
	// nodeCopy.Annotations[RebootAnnotation] = ""
	// if _, err := c.client.Core().Nodes().Update(nodeCopy); err != nil {
	// 	return fmt.Errorf("Failed to set %s annotation: %v", RebootAnnotation, err)
	// }
	return nil
}
