package main

import (
	"context"
	//"fmt"
	"time"
	"strings"
	"log"
	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	//"github.com/kr/pretty"
	//"k8s.io/api/apps/v1"
	//"k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {

		var guildID string = "554446372619681803"

		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		//pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		deployList, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		for i := range deployList.Items {
			if strings.Contains(deployList.Items[i].Labels["custguild"], guildID) {
				log.Printf(deployList.Items[i].Labels["custguild"])
				log.Printf(deployList.Items[i].Labels["gamename"])

			}
		}

		time.Sleep(10 * time.Second)
	}
}
