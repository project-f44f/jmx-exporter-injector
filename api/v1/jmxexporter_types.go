/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JmxExporterSpec defines the desired state of JmxExporter
type JmxExporterSpec struct {
	// 部署注解的键，用于匹配 Deployment
	DeploymentAnnotationKey string `json:"deploymentAnnotationKey"`
	// 部署注解的值，设置为 'true' 时启用 JMX Exporter 注入
	DeploymentAnnotationValue bool `json:"deploymentAnnotationValue"`
	// 是否启用 JMX Exporter 注入到匹配的 Deployment 中
	EnableInjection bool `json:"enableInjection"`
	// 是否在 Service 中暴露 JMX Exporter 的端口
	ExposeServicePort bool `json:"exposeServicePort"`
	// 是否创建 ServiceMonitor 资源
	CreateServiceMonitor bool `json:"createServiceMonitor"`
}

// JmxExporterStatus defines the observed state of JmxExporter
type JmxExporterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// JmxExporter is the Schema for the jmxexporters API
type JmxExporter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JmxExporterSpec   `json:"spec,omitempty"`
	Status JmxExporterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JmxExporterList contains a list of JmxExporter
type JmxExporterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JmxExporter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JmxExporter{}, &JmxExporterList{})
}
