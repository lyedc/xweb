/*
	Copyright NetFoundry, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

	https://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

package xweb

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

// WebHandlerFactoryRegistry describes a registry of binding to WebHandlerFactory registrations
type WebHandlerFactoryRegistry interface {
	Add(factory WebHandlerFactory) error
	Get(binding string) WebHandlerFactory
}

// WebHandlerFactoryRegistryImpl is a basic WebHandlerFactoryRegistry implementation backed by a simple mapping of binding (string) to WebHandlerFactories
type WebHandlerFactoryRegistryImpl struct {
	factories map[string]WebHandlerFactory
}

// NewWebHandlerFactoryRegistryImpl creates a new WebHandlerFactoryRegistryImpl
func NewWebHandlerFactoryRegistryImpl() *WebHandlerFactoryRegistryImpl {
	return &WebHandlerFactoryRegistryImpl{
		factories: map[string]WebHandlerFactory{},
	}
}

// Add adds a factory to the registry. Errors if a previous factory with the same binding is registered.
func (registry WebHandlerFactoryRegistryImpl) Add(factory WebHandlerFactory) error {
	logrus.Debugf("adding xweb factory with binding: %v", factory.Binding())
	if _, ok := registry.factories[factory.Binding()]; ok {
		return fmt.Errorf("binding [%s] already registered", factory.Binding())
	}

	registry.factories[factory.Binding()] = factory

	return nil
}

// Get retrieves a factory based on a binding or nil if no factory for the binding is registered
func (registry WebHandlerFactoryRegistryImpl) Get(binding string) WebHandlerFactory {
	return registry.factories[binding]
}

// APIBinding is an interface defines the minimum operations necessary to convert configuration into a WebHandler
// by some WebHandlerFactory. The APIBinding.Binding() value is used to map configuration data to specific
// WebHandlerFactory instances that  generate a WebHandler with the same binding value.
type APIBinding interface {
	Binding() string
	Options() map[interface{}]interface{}
}

// WebHandler is an interface that is a  http.Handler with options, binding, and information for a specific web API
// that was generated from a WebHandlerFactory.
type WebHandler interface {
	Binding() string
	Options() map[interface{}]interface{}
	RootPath() string
	IsHandler(r *http.Request) bool

	http.Handler
}

// The WebHandlerFactory interface generates WebHandler instances. Factories can use a single instance or multiple
// instances based on need. This interface allows WebHandler logic to be reused across multiple xweb.Server's while
// delegating the instance management to the factory.
type WebHandlerFactory interface {
	// Binding returns the string binding value used in configurations
	Binding() string

	// New creates a new instances of this factories API
	New(webListener *WebListener, options map[interface{}]interface{}) (WebHandler, error)

	// Validate is used for factory level configuration checks. It is not meant for individual instance validation.
	// It is meant for validating configuration that all instances generated by a factory share (i.e. sanity across
	// instances or configuration from another section).
	Validate(config *Config) error
}
