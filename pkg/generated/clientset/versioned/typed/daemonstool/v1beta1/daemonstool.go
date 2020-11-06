/*
Copyright The Kubernetes Authors.

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

package v1beta1

import (
	"context"
	v1beta1 "my-controller/pkg/apis/daemonstool/v1beta1"
	scheme "my-controller/pkg/generated/clientset/versioned/scheme"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DaemonstoolsGetter has a method to return a DaemonstoolInterface.
// A group's client should implement this interface.
type DaemonstoolsGetter interface {
	Daemonstools(namespace string) DaemonstoolInterface
}

// DaemonstoolInterface has methods to work with Daemonstool resources.
type DaemonstoolInterface interface {
	Create(ctx context.Context, daemonstool *v1beta1.Daemonstool, opts v1.CreateOptions) (*v1beta1.Daemonstool, error)
	Update(ctx context.Context, daemonstool *v1beta1.Daemonstool, opts v1.UpdateOptions) (*v1beta1.Daemonstool, error)
	UpdateStatus(ctx context.Context, daemonstool *v1beta1.Daemonstool, opts v1.UpdateOptions) (*v1beta1.Daemonstool, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.Daemonstool, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.DaemonstoolList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Daemonstool, err error)
	DaemonstoolExpansion
}

// daemonstools implements DaemonstoolInterface
type daemonstools struct {
	client rest.Interface
	ns     string
}

// newDaemonstools returns a Daemonstools
func newDaemonstools(c *UalterV1beta1Client, namespace string) *daemonstools {
	return &daemonstools{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the daemonstool, and returns the corresponding daemonstool object, and an error if there is any.
func (c *daemonstools) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Daemonstool, err error) {
	result = &v1beta1.Daemonstool{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("daemonstools").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Daemonstools that match those selectors.
func (c *daemonstools) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.DaemonstoolList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.DaemonstoolList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("daemonstools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested daemonstools.
func (c *daemonstools) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("daemonstools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a daemonstool and creates it.  Returns the server's representation of the daemonstool, and an error, if there is any.
func (c *daemonstools) Create(ctx context.Context, daemonstool *v1beta1.Daemonstool, opts v1.CreateOptions) (result *v1beta1.Daemonstool, err error) {
	result = &v1beta1.Daemonstool{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("daemonstools").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(daemonstool).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a daemonstool and updates it. Returns the server's representation of the daemonstool, and an error, if there is any.
func (c *daemonstools) Update(ctx context.Context, daemonstool *v1beta1.Daemonstool, opts v1.UpdateOptions) (result *v1beta1.Daemonstool, err error) {
	result = &v1beta1.Daemonstool{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("daemonstools").
		Name(daemonstool.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(daemonstool).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *daemonstools) UpdateStatus(ctx context.Context, daemonstool *v1beta1.Daemonstool, opts v1.UpdateOptions) (result *v1beta1.Daemonstool, err error) {
	result = &v1beta1.Daemonstool{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("daemonstools").
		Name(daemonstool.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(daemonstool).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the daemonstool and deletes it. Returns an error if one occurs.
func (c *daemonstools) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("daemonstools").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *daemonstools) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("daemonstools").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched daemonstool.
func (c *daemonstools) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Daemonstool, err error) {
	result = &v1beta1.Daemonstool{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("daemonstools").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}