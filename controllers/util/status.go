package util

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	runtimeclient "sigs.k8s.io/controller-runtime/pkg/client"

	kedav1alpha1 "github.com/kedacore/keda/api/v1alpha1"
)

func SetStatusConditions(client runtimeclient.Client, logger logr.Logger, object interface{}, conditions *kedav1alpha1.Conditions) error {
	var patch runtimeclient.Patch

	runtimeObj := object.(runtime.Object)
	switch obj := runtimeObj.(type) {
	case *kedav1alpha1.ScaledObject:
		patch = runtimeclient.MergeFrom(obj.DeepCopy())
		obj.Status.Conditions = *conditions
	case *kedav1alpha1.ScaledJob:
		patch = runtimeclient.MergeFrom(obj.DeepCopy())
		obj.Status.Conditions = *conditions
	default:
		err := fmt.Errorf("Unknown scalable object type %v", obj)
		logger.Error(err, "Failed to patch Objects Status with Conditions")
		return err
	}

	err := client.Status().Patch(context.TODO(), runtimeObj, patch)
	if err != nil {
		logger.Error(err, "Failed to patch Objects Status with Conditions")
	}
	return err
}

func UpdateScaledObjectStatus(client runtimeclient.Client, logger logr.Logger, scaledObject *kedav1alpha1.ScaledObject, status *kedav1alpha1.ScaledObjectStatus) error {
	patch := runtimeclient.MergeFrom(scaledObject.DeepCopy())
	scaledObject.Status = *status
	err := client.Status().Patch(context.TODO(), scaledObject, patch)
	if err != nil {
		logger.Error(err, "Failed to patch ScaledObjects Status")
	}
	return err
}
