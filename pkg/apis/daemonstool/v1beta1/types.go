package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Daemonstool is a specification for a Generic Daemon resource
type Daemonstool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DaemonstoolSpec   `json:"spec"`
	Status DaemonstoolStatus `json:"status"`
}

// DaemonstoolSpec is the spec for a Daemonstool resource
type DaemonstoolSpec struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Image string `json:"image"`
}

// DaemonstoolStatus is the status for a Daemonstool resource
type DaemonstoolStatus struct {
	Installed int32 `json:"installed"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DaemonstoolList is a list of Daemonstool resources
type DaemonstoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Daemonstool `json:"items"`
}
