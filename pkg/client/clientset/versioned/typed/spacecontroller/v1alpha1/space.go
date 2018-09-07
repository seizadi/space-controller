/*
Copyright 2018 The Kubernetes namespace-controller Authors.

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

package v1alpha1

import (
	v1alpha1 "github.com/seizadi/space-controller/pkg/apis/spacecontroller/v1alpha1"
	scheme "github.com/seizadi/space-controller/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SpacesGetter has a method to return a SpaceInterface.
// A group's client should implement this interface.
type SpacesGetter interface {
	Spaces(namespace string) SpaceInterface
}

// SpaceInterface has methods to work with Space resources.
type SpaceInterface interface {
	Create(*v1alpha1.Space) (*v1alpha1.Space, error)
	Update(*v1alpha1.Space) (*v1alpha1.Space, error)
	UpdateStatus(*v1alpha1.Space) (*v1alpha1.Space, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Space, error)
	List(opts v1.ListOptions) (*v1alpha1.SpaceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Space, err error)
	SpaceExpansion
}

// spaces implements SpaceInterface
type spaces struct {
	client rest.Interface
	ns     string
}

// newSpaces returns a Spaces
func newSpaces(c *SpacecontrollerV1alpha1Client, namespace string) *spaces {
	return &spaces{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the space, and returns the corresponding space object, and an error if there is any.
func (c *spaces) Get(name string, options v1.GetOptions) (result *v1alpha1.Space, err error) {
	result = &v1alpha1.Space{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("spaces").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Spaces that match those selectors.
func (c *spaces) List(opts v1.ListOptions) (result *v1alpha1.SpaceList, err error) {
	result = &v1alpha1.SpaceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("spaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested spaces.
func (c *spaces) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("spaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a space and creates it.  Returns the server's representation of the space, and an error, if there is any.
func (c *spaces) Create(space *v1alpha1.Space) (result *v1alpha1.Space, err error) {
	result = &v1alpha1.Space{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("spaces").
		Body(space).
		Do().
		Into(result)
	return
}

// Update takes the representation of a space and updates it. Returns the server's representation of the space, and an error, if there is any.
func (c *spaces) Update(space *v1alpha1.Space) (result *v1alpha1.Space, err error) {
	result = &v1alpha1.Space{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("spaces").
		Name(space.Name).
		Body(space).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *spaces) UpdateStatus(space *v1alpha1.Space) (result *v1alpha1.Space, err error) {
	result = &v1alpha1.Space{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("spaces").
		Name(space.Name).
		SubResource("status").
		Body(space).
		Do().
		Into(result)
	return
}

// Delete takes name of the space and deletes it. Returns an error if one occurs.
func (c *spaces) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("spaces").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *spaces) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("spaces").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched space.
func (c *spaces) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Space, err error) {
	result = &v1alpha1.Space{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("spaces").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}