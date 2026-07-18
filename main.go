package main

import (
	"fmt"
	"sync"
	"time"
)

// ---- 1. TYPES OF EVENTS ----
type EventType string

const (
	EventPodAdded   EventType = "ADDED"
	EventPodUpdated EventType = "UPDATED"
)

type WatchEvent struct {
	Type EventType
	Pod  *Pod
}

type Pod struct {
	Name   string
	Image  string
	Node   string
	Status string
}

// ---- 2. THE EVENT-DRIVEN API SERVER ----
type APIserver struct {
	mu          sync.Mutex
	pods        map[string]*Pod
	subscribers []chan WatchEvent // Broadcast channels for components
}

func NewAPIserver() *APIserver {
	return &APIserver{
		pods:        make(map[string]*Pod),
		subscribers: make([]chan WatchEvent, 0),
	}
}

// Watch allows Scheduler and Kubelets to subscribe to real-time updates
func (api *APIserver) Watch() <-chan WatchEvent {
	api.mu.Lock()
	defer api.mu.Unlock()
	
	// Create a buffered channel so slower components don't block the API server
	ch := make(chan WatchEvent, 100)
	api.subscribers = append(api.subscribers, ch)
	return ch
}

// broadcast sends the event to all active watchers asynchronously
func (api *APIserver) broadcast(event WatchEvent) {
	for _, sub := range api.subscribers {
		select {
		case sub <- event:
		default:
			// Drop event or handle slow consumer to prevent deadlocks
		}
	}
}

func (api *APIserver) CreatePod(name, image string) {
	api.mu.Lock()
	pod := &Pod{Name: name, Image: image, Status: "Pending"}
	api.pods[name] = pod
	api.mu.Unlock()

	fmt.Printf("[API-Server] Pod '%s' created\n", name)
	api.broadcast(WatchEvent{Type: EventPodAdded, Pod: pod})
}

func (api *APIserver) UpdatePod(pod *Pod) {
	api.mu.Lock()
	api.pods[pod.Name] = pod
	api.mu.Unlock()

	api.broadcast(WatchEvent{Type: EventPodUpdated, Pod: pod})
}

// ---- 3. THE SCALABLE SCHEDULER ----
// Instead of loops, it sits and blocks on the channel until a Pod event arrives
func StartScheduler(api *APIserver, availableNodes []string) {
	events := api.Watch()
	var nodeIndex int

	go func() {
		for event := range events {
			// Scalability Win: We ONLY evaluate the single pod that triggered the event
			pod := event.Pod
			if pod.Node == "" {
				// Round-robin distribution across many worker nodes
				assignedNode := availableNodes[nodeIndex%len(availableNodes)]
				nodeIndex++

				// Create a copy to modify safely
				updatedPod := *pod
				updatedPod.Node = assignedNode
				
				fmt.Printf("[Scheduler] Assigned '%s' to '%s'\n", pod.Name, assignedNode)
				api.UpdatePod(&updatedPod)
			}
		}
	}()
}

// ---- 4. THE SCALABLE KUBELET ----
type Kubelet struct {
	NodeName string
	api      *APIserver
}

func StartKubelet(nodeName string, api *APIserver) {
	k := &Kubelet{NodeName: nodeName, api: api}
	events := api.Watch()

	go func() {
		for event := range events {
			pod := event.Pod
			// Scalability Win: Kubelet instantly ignores pods belonging to other nodes
			if pod.Node == k.NodeName && pod.Status == "Pending" {
				// Run the container launch in its own goroutine so it doesn't block the queue!
				go k.reconcile(pod)
			}
		}
	}()
}

func (k *Kubelet) reconcile(pod *Pod) {
	fmt.Printf("[%s Kubelet] Launching '%s' using image '%s'...\n", k.NodeName, pod.Name, pod.Image)
	
	// Simulate container startup time
	time.Sleep(1 * time.Second)

	updatedPod := *pod
	updatedPod.Status = "Running"
	k.api.UpdatePod(&updatedPod)
	fmt.Printf("[%s Kubelet] Pod '%s' is now Running\n", k.NodeName, pod.Name)
}

// ---- 5. THE SCALE TEST ----
func main() {
	api := NewAPIserver()

	// Scale up our infrastructure: 3 active worker nodes
	nodes := []string{"worker-node-1", "worker-node-2", "worker-node-3"}

	StartScheduler(api, nodes)
	for _, node := range nodes {
		StartKubelet(node, api)
	}

	time.Sleep(100 * time.Millisecond)

	fmt.Println("\n--- Simulating High Load: Deploying 10 Pods at once ---")
	for i := 1; i <= 10; i++ {
		podName := fmt.Sprintf("app-replica-%d", i)
		imageName := fmt.Sprintf("custom-image-%d:v1", i)
		api.CreatePod(podName, imageName)
	}

	// Wait for processing to finish
	time.Sleep(3 * time.Second)
}