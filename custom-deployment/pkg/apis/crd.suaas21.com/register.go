package crd_suaas21_com

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "crd.suaas21.com"
)

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1alpha1"}
