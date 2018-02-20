/*
Copyright 2018 Mike Bryant. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package main for the random secret operator
package main

import (
	"fmt"

	"crypto/rand"

	"encoding/hex"

	randomsecrets "github.com/mikebryant/random-secret-operator/pkg/apis/randomsecrets/v1"
	randomsecretsclient "github.com/mikebryant/random-secret-operator/pkg/client/clientset/versioned/typed/randomsecrets/v1"
	opkit "github.com/rook/operator-kit"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/retry"
)

// RandomSecretController represents a controller object for random secret custom resources
type RandomSecretController struct {
	context               *opkit.Context
	randomSecretClientset randomsecretsclient.RandomsecretsV1Interface
}

// newRandomSecretController create controller for watching random secret custom resources created
func newRandomSecretController(context *opkit.Context, randomSecretClientset randomsecretsclient.RandomsecretsV1Interface) *RandomSecretController {
	return &RandomSecretController{
		context:               context,
		randomSecretClientset: randomSecretClientset,
	}
}

// Watch watches for instances of RandomSecret custom resources and acts on them
func (c *RandomSecretController) StartWatch(namespace string, stopCh chan struct{}) error {

	resourceHandlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.onAdd,
		UpdateFunc: c.onUpdate,
		DeleteFunc: c.onDelete,
	}
	restClient := c.randomSecretClientset.RESTClient()
	watcher := opkit.NewWatcher(randomsecrets.RandomSecretResource, namespace, resourceHandlers, restClient)
	go watcher.Watch(&randomsecrets.RandomSecret{}, stopCh)
	return nil
}

func (c *RandomSecretController) generateData() []byte {
	src := make([]byte, 64)
	_, err := rand.Read(src)
	if err != nil {
		panic(err)
	}
	val := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(val, src)
	return val
}

func (c *RandomSecretController) updateSecret(rs *randomsecrets.RandomSecret, s *v1.Secret) {
	blockOwnerDeletion := true
	controller := true
	s.OwnerReferences = []metav1.OwnerReference{
		{
			APIVersion:         randomsecrets.RandomSecretResource.Version,
			BlockOwnerDeletion: &blockOwnerDeletion,
			Controller:         &controller,
			Kind:               randomsecrets.RandomSecretResource.Kind,
			Name:               rs.Name,
			UID:                rs.UID,
		},
	}
	if s.Data == nil {
		s.Data = make(map[string][]byte)
	}
	if val, ok := s.Data["random"]; !ok || len(val) < 128 {
		fmt.Printf("Adding data to Secret '%s'...\n", s.Name)
		s.Data["random"] = c.generateData()
	}
}

func (c *RandomSecretController) ensureSecret(obj *randomsecrets.RandomSecret) {
	// Ensure that the Secret object exists and is valid for the provided RandomSecret
	// Create if necessary, and update if the parameters are wrong
	secretsClient := c.context.Clientset.CoreV1().Secrets(obj.Namespace)
	s, err := secretsClient.Get(obj.Name, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Creating Secret '%s'...\n", obj.Name)
		s = &v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      obj.Name,
				Namespace: obj.Namespace,
			},
		}
		c.updateSecret(obj, s)
		_, err := secretsClient.Create(s)
		if err != nil {
			panic(err)
		}
	} else {
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			s, getErr := secretsClient.Get(obj.Name, metav1.GetOptions{})
			if getErr != nil {
				panic(fmt.Errorf("Failed to get Secret '%s': %v", s.Name, getErr))
			}

			c.updateSecret(obj, s)
			_, updateErr := secretsClient.Update(s)
			return updateErr
		})
		if retryErr != nil {
			panic(fmt.Errorf("Update failed: %v", retryErr))
		}
	}
}

func (c *RandomSecretController) onAdd(obj interface{}) {
	s := obj.(*randomsecrets.RandomSecret).DeepCopy()
	fmt.Printf("Add event: RandomSecret '%s/%s' with Spec=%s...\n", s.Namespace, s.Name, s.Spec)

	c.ensureSecret(s)
}

func (c *RandomSecretController) onUpdate(oldObj, newObj interface{}) {
	s := newObj.(*randomsecrets.RandomSecret).DeepCopy()
	fmt.Printf("Update event: RandomSecret '%s/%s' with Spec=%s...\n", s.Namespace, s.Name, s.Spec)

	c.ensureSecret(s)
}

func (c *RandomSecretController) onDelete(obj interface{}) {
	s := obj.(*randomsecrets.RandomSecret).DeepCopy()
	fmt.Printf("Delete event: RandomSecret '%s/%s' with Spec=%s...\n", s.Namespace, s.Name, s.Spec)
	// We ignore this, and allow garbage collection to handle it
}
