/*
	Copyright NetFoundry Inc.

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
	"gitee.com/zhaochuninhefei/gmgo/gmhttp"
)

// ApiBinding is an interface defines the minimum operations necessary to convert configuration into a ApiHandler
// by some ApiHandlerFactory. The ApiBinding.Binding() value is used to map configuration data to specific
// ApiHandlerFactory instances that  generate a ApiHandler with the same binding value.
type ApiBinding interface {
	Binding() string
	Options() map[interface{}]interface{}
}

// ApiHandler is an interface that is a  http.Handler with options, binding, and information for a specific ApiConfig
// that was generated from a ApiHandlerFactory.
type ApiHandler interface {
	Binding() string
	Options() map[interface{}]interface{}
	RootPath() string
	IsHandler(r *gmhttp.Request) bool

	gmhttp.Handler
}

// The ApiHandlerFactory interface generates ApiHandler instances. Factories can use a single instance or multiple
// instances based on need. This interface allows ApiHandler logic to be reused across multiple xweb.Server's while
// delegating the instance management to the factory.
type ApiHandlerFactory interface {
	// Binding returns the string binding value used in configurations
	Binding() string

	// New creates a new instances of this factories ApiConfig
	New(serverConfig *ServerConfig, options map[interface{}]interface{}) (ApiHandler, error)

	// Validate is used for factory level configuration checks. It is not meant for individual instance validation.
	// It is meant for validating configuration that all instances generated by a factory share (i.e. sanity across
	// instances or configuration from another section).
	Validate(config *InstanceConfig) error
}
