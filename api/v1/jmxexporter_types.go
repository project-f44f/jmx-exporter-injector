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
	// 部署注解的值
	DeploymentAnnotationValue string `json:"deploymentAnnotationValue"`
	// 是否启用 JMX Exporter 注入到匹配的 Deployment 中
	EnableInjection bool `json:"enableInjection"`
	// InitContainer 使用的基础镜像
	JmxExporterImage string `json:"jmxExporterImage"`
	// servicemonitor 使用的label
	ServiceMonitorLabelValue string `json:"serviceMonitorLabelValue"`
	// 指定暴露端口的 service 名称
	ServiceName string `json:"serviceName"`
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
