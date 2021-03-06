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

// Code generated by informer-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	daemonstoolv1beta1 "my-controller/pkg/apis/daemonstool/v1beta1"
	versioned "my-controller/pkg/generated/clientset/versioned"
	internalinterfaces "my-controller/pkg/generated/informers/externalversions/internalinterfaces"
	v1beta1 "my-controller/pkg/generated/listers/daemonstool/v1beta1"
	time "time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// DaemonstoolInformer provides access to a shared informer and lister for
// Daemonstools.
type DaemonstoolInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1beta1.DaemonstoolLister
}

type daemonstoolInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewDaemonstoolInformer constructs a new informer for Daemonstool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewDaemonstoolInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredDaemonstoolInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredDaemonstoolInformer constructs a new informer for Daemonstool type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredDaemonstoolInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.UalterV1beta1().Daemonstools(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.UalterV1beta1().Daemonstools(namespace).Watch(context.TODO(), options)
			},
		},
		&daemonstoolv1beta1.Daemonstool{},
		resyncPeriod,
		indexers,
	)
}

func (f *daemonstoolInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredDaemonstoolInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *daemonstoolInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&daemonstoolv1beta1.Daemonstool{}, f.defaultInformer)
}

func (f *daemonstoolInformer) Lister() v1beta1.DaemonstoolLister {
	return v1beta1.NewDaemonstoolLister(f.Informer().GetIndexer())
}
