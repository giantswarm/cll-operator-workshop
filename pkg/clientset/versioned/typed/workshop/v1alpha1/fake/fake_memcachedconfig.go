/*
Copyright 2018 Giant Swarm GmbH.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/giantswarm/cll-operator-workshop/pkg/apis/workshop/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeMemcachedConfigs implements MemcachedConfigInterface
type FakeMemcachedConfigs struct {
	Fake *FakeWorkshopV1alpha1
	ns   string
}

var memcachedconfigsResource = schema.GroupVersionResource{Group: "workshop.continuouslifecycle.london", Version: "v1alpha1", Resource: "memcachedconfigs"}

var memcachedconfigsKind = schema.GroupVersionKind{Group: "workshop.continuouslifecycle.london", Version: "v1alpha1", Kind: "MemcachedConfig"}

// Get takes name of the memcachedConfig, and returns the corresponding memcachedConfig object, and an error if there is any.
func (c *FakeMemcachedConfigs) Get(name string, options v1.GetOptions) (result *v1alpha1.MemcachedConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(memcachedconfigsResource, c.ns, name), &v1alpha1.MemcachedConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MemcachedConfig), err
}

// List takes label and field selectors, and returns the list of MemcachedConfigs that match those selectors.
func (c *FakeMemcachedConfigs) List(opts v1.ListOptions) (result *v1alpha1.MemcachedConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(memcachedconfigsResource, memcachedconfigsKind, c.ns, opts), &v1alpha1.MemcachedConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.MemcachedConfigList{}
	for _, item := range obj.(*v1alpha1.MemcachedConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested memcachedConfigs.
func (c *FakeMemcachedConfigs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(memcachedconfigsResource, c.ns, opts))

}

// Create takes the representation of a memcachedConfig and creates it.  Returns the server's representation of the memcachedConfig, and an error, if there is any.
func (c *FakeMemcachedConfigs) Create(memcachedConfig *v1alpha1.MemcachedConfig) (result *v1alpha1.MemcachedConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(memcachedconfigsResource, c.ns, memcachedConfig), &v1alpha1.MemcachedConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MemcachedConfig), err
}

// Update takes the representation of a memcachedConfig and updates it. Returns the server's representation of the memcachedConfig, and an error, if there is any.
func (c *FakeMemcachedConfigs) Update(memcachedConfig *v1alpha1.MemcachedConfig) (result *v1alpha1.MemcachedConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(memcachedconfigsResource, c.ns, memcachedConfig), &v1alpha1.MemcachedConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MemcachedConfig), err
}

// Delete takes name of the memcachedConfig and deletes it. Returns an error if one occurs.
func (c *FakeMemcachedConfigs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(memcachedconfigsResource, c.ns, name), &v1alpha1.MemcachedConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeMemcachedConfigs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(memcachedconfigsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.MemcachedConfigList{})
	return err
}

// Patch applies the patch and returns the patched memcachedConfig.
func (c *FakeMemcachedConfigs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.MemcachedConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(memcachedconfigsResource, c.ns, name, data, subresources...), &v1alpha1.MemcachedConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.MemcachedConfig), err
}
