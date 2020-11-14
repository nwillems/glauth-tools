package chpass

import (
	"context"
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/util/retry"
)

// GetConfigMap finds and return the data of a given configmap in a given namespace
func GetConfigMap(configmapName, namespace string) (map[string]string, error) {
	config, err := restclient.InClusterConfig() // TODO: Make configureable
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	api := clientset.CoreV1()

	getOpts := metav1.GetOptions{}
	fetchedConfigmap, err := api.ConfigMaps(namespace).Get(context.TODO(), configmapName, getOpts)
	if err != nil {
		return nil, err
	}

	return fetchedConfigmap.Data, nil
}

// UpdateConfigMap stores the given configmap data, does not consider changes that has happened between the original data was fetched and now
func UpdateConfigMap(configmapName, namespace string, data map[string]string) error {
	config, err := restclient.InClusterConfig()
	if err != nil {
		return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}

	api := clientset.CoreV1()
	configMapClient := api.ConfigMaps(namespace)

	log.Println("Updating configmap - config, clientset and apiclient created")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		log.Println("Updating configmap - in retryOnConflict")
		result, getErr := configMapClient.Get(context.TODO(), configmapName, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of ConfigMap: %v", getErr)
		}

		result.Data = data
		_, updateErr := configMapClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		log.Printf("Updated configmap - %v", updateErr)
		return updateErr
	})
	if retryErr != nil {
		return fmt.Errorf("Update failed: %v", retryErr)
	}

	return nil
}
