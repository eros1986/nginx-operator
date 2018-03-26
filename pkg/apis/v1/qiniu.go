package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Qiniu struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   QiniuSpec   `json:"spec"`
	Status QiniuStatus `json:"status"`
}

type QiniuList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Qiniu `json:"items"`
}

type QiniuSpec struct {
	Name string `json:"name"`
}

type QiniuStatus struct {
	Msg string `json:"msg"`
}
