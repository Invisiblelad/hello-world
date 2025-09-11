package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const kubeconfigDir = "/loca/dev"

func main() {
	// Find all kubeconfig files
	files, err := filepath.Glob(filepath.Join(kubeconfigDir, "*.yaml"))
	if err != nil {
		log.Fatalf("❌ Failed to read directory: %v", err)
	}
	if len(files) == 0 {
		log.Fatalf("❌ No kubeconfig files found in %s", kubeconfigDir)
	}

	// Show choices
	fmt.Println("📌 Select a kubeconfig:")
	for i, file := range files {
		fmt.Printf("%d) %s\n", i+1, filepath.Base(file))
	}

	// Get user choice
	fmt.Print("Enter choice: ")
	var choice string
	fmt.Scanln(&choice)
	index, err := strconv.Atoi(choice)
	if err != nil || index < 1 || index > len(files) {
		log.Fatalf("❌ Invalid choice")
	}
	selected := files[index-1]

	// Connect to cluster
	connectToCluster(selected)
}

func connectToCluster(kubeconfig string) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("❌ Error loading kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("❌ Error creating client: %v", err)
	}

	// List pods in default namespace
	pods, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("❌ Error listing pods: %v", err)
	}

	fmt.Printf("\n✅ Connected using %s\n", kubeconfig)
	fmt.Printf("Found %d pods in 'default' namespace:\n", len(pods.Items))
	for _, pod := range pods.Items {
		fmt.Printf(" - %s (%s)\n", pod.Name, pod.Status.Phase)
	}
}
