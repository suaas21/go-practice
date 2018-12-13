package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//CustomDeployment Specification
type CustomDeployment struct{
	metav1.TypeMeta		`json:",inline"`
	metav1.ObjectMeta	`json:"metadata,omitempty"`

	Spec CustomDeploymentSpec		`json:"spec"`
	Status CustomDeploymentStatus	`json:"status"`
}

//CustomDeployment Specification
type CustomDeploymentSpec struct {
	Replicas *int32	`json:"replicas"`
	Selector *metav1.LabelSelector `json:"selector"`
	Template CustomPodTemplate `json:"template"`
}

//PodTemplate Specification
type CustomPodTemplate struct{
	metav1.ObjectMeta 	`json:"metadata,omitempty"`
	Spec              PodSpec `json:"spec"`
}
type PodSpec struct{
	Containers []Container `json:"containers" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=containers"`
}
type Container struct{
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	Image string `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	Ports []ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
}
type ContainerPort struct {
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	HostPort int32 `json:"hostPort,omitempty" protobuf:"varint,2,opt,name=hostPort"`
	ContainerPort int32 `json:"containerPort" protobuf:"varint,3,opt,name=containerPort"`
	Protocol Protocol `json:"protocol,omitempty" protobuf:"bytes,4,opt,name=protocol,casttype=Protocol"`
	HostIP string `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
}
type Protocol string
const (
	ProtocolTCP Protocol = "TCP"
	ProtocolUDP Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

//Status of the CustomDeployment
type CustomDeploymentStatus struct {
	AvailableReplicas   int32 `json:"available_replicas"`
	CreatingReplicas    int32 `json:"creating_replicas"`
	TerminatingReplicas int32 `json:"terminating_replicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

//List of CustomDeployment
type CustomDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []CustomDeployment `json:"items"`
}