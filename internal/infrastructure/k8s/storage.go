package k8s

import (
	"context"
	"fmt"
	"strconv"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	NotLonghornError = fmt.Errorf("not longhorn")
)

// isLonghornDefaultStorageClass returns (true, scName, nil) when the default SC is Longhorn.
func (self *KubeClient) isLonghornDefaultStorageClass(ctx context.Context) (bool, string, error) {
	scList, err := self.clientset.StorageV1().StorageClasses().List(ctx, meta.ListOptions{})
	if err != nil {
		return false, "", err
	}

	for _, sc := range scList.Items {
		if isDefault(sc) {
			return sc.Provisioner == "driver.longhorn.io", sc.Name, nil
		}
	}
	return false, "", fmt.Errorf("no default StorageClass found")
}

func isDefault(sc storagev1.StorageClass) bool {
	if v, ok := sc.Annotations["storageclass.kubernetes.io/is-default-class"]; ok && v == "true" {
		return true
	}
	if v, ok := sc.Annotations["storageclass.beta.kubernetes.io/is-default-class"]; ok && v == "true" {
		return true
	}
	return false
}

// AvailableStorageBytes returns the total available storage in bytes as a string.
// Returns a string representation of the total bytes available across all Longhorn nodes.
func (self *KubeClient) AvailableStorageBytes(ctx context.Context) (storageBytes string, storageClass string, err error) {
	sc, className, err := self.isLonghornDefaultStorageClass(ctx)
	if err != nil {
		return "", className, err
	}
	if !sc {
		return "", className, NotLonghornError
	}

	gvr := schema.GroupVersionResource{
		Group:    "longhorn.io",
		Version:  "v1beta2",
		Resource: "nodes",
	}

	// All Longhorn CRs live in the longhorn-system namespace.
	list, err := self.client.Resource(gvr).Namespace("longhorn-system").List(
		ctx, meta.ListOptions{})
	if err != nil {
		return "", className, err
	}

	total := resource.NewQuantity(0, resource.BinarySI)

	for _, u := range list.Items {
		// diskStatus is a map keyed by disk UUID
		disks, _, _ := unstructured.NestedMap(u.Object, "status", "diskStatus")
		for _, v := range disks {
			if disk, ok := v.(map[string]interface{}); ok {
				switch x := disk["storageAvailable"].(type) {
				case int64:
					total.Add(*resource.NewQuantity(x, resource.BinarySI))
				case float64: // JSON numbers default to float64
					total.Add(*resource.NewQuantity(int64(x), resource.BinarySI))
				case string: // some Longhorn versions emit strings
					if i, err := strconv.ParseInt(x, 10, 64); err == nil {
						total.Add(*resource.NewQuantity(i, resource.BinarySI))
					}
				}
			}
		}
	}

	return total.String(), className, nil
}
