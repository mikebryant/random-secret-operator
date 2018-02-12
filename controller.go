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

	randomsecrets "github.com/mikebryant/random-secret-operator/pkg/apis/randomsecrets/v1"
	randomsecretsclient "github.com/mikebryant/random-secret-operator/pkg/client/clientset/versioned/typed/randomsecrets/v1"
	opkit "github.com/rook/operator-kit"
	"k8s.io/client-go/tools/cache"
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

func (c *RandomSecretController) onAdd(obj interface{}) {
	s := obj.(*randomsecrets.RandomSecret).DeepCopy()

	fmt.Printf("Added RandomSecret '%s' with Spec=%s\n", s.Name, s.Spec)
}

func (c *RandomSecretController) onUpdate(oldObj, newObj interface{}) {
	oldRandomSecret := oldObj.(*randomsecrets.RandomSecret).DeepCopy()
	newRandomSecret := newObj.(*randomsecrets.RandomSecret).DeepCopy()

	fmt.Printf("Updated RandomSecret '%s' from %s to %s\n", newRandomSecret.Name, oldRandomSecret.Spec, newRandomSecret.Spec)
}

func (c *RandomSecretController) onDelete(obj interface{}) {
	s := obj.(*randomsecrets.RandomSecret).DeepCopy()

	fmt.Printf("Deleted RandomSecret '%s' with Spec=%s\n", s.Name, s.Spec)
}
