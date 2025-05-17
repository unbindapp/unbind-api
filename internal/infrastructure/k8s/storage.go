package k8s

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type StorageMetadata struct {
	StorageClassName          string `json:"storage_class_name"`
	AllocatableGB             string `json:"allocatable_bytes_gb"`
	UnableToDetectAllocatable bool   `json:"unable_to_detect_allocatable"`
	MinimumStorageGB          string `json:"minimum_storage_gb"`
	StorageStepGB             string `json:"storage_step_gb"`
}

// AvailableStorageBytes inspects the default StorageClass and returns
// capacity / sizing metadata
//
// • Longhorn  – sums .status.diskStatus[*].storageAvailable live
// • Hetzner   – 10 TiB max, 10 GiB min, 1 GiB step
// • AWS EBS   – 64 TiB max, 1 GiB  min, 1 GiB step
// • Azure Disk – 64 TiB max, 1 GiB  min, 1 GiB step
// • GCP PD    – 64 TiB max, 1 GiB  min, 1 GiB step
// • DigitalOcean Volumes – 16 TiB max, 1 GiB min, 1 GiB step
// • Vultr Block Storage  – 10 TiB max, 10 GiB min, 1 GiB step
// • Linode Block Storage – 16 TiB max, 10 GiB min, 1 GiB step
// • OpenStack Cinder     – 12 TiB max, 10 GiB min, 1 GiB step
//
// Anything else falls through with UnableToDetectAllocatable=true.
func (self *KubeClient) AvailableStorageBytes(ctx context.Context) (*StorageMetadata, error) {
	const (
		giB = 1 << 30 // 1 GiB  = 1 073 741 824 bytes
		tiB = 1 << 40 // 1 TiB  = 1 099 511 627 776 bytes
	)
	resp := &StorageMetadata{UnableToDetectAllocatable: true}

	scList, err := self.clientset.StorageV1().StorageClasses().List(ctx, meta.ListOptions{})
	if err != nil {
		return resp, err
	}

	for _, sc := range scList.Items {
		if isDefault(sc) {
			resp.StorageClassName = sc.Name
		}

		switch sc.Provisioner {

		// * Longhorn - we will query the Longhorn node for available storage
		case "driver.longhorn.io":
			gvr := schema.GroupVersionResource{
				Group:    "longhorn.io",
				Version:  "v1beta2",
				Resource: "nodes",
			}
			list, err := self.client.Resource(gvr).
				Namespace("longhorn-system").
				List(ctx, meta.ListOptions{})
			if err != nil {
				return resp, err
			}

			// Track the biggest node-level free capacity
			maxFree := resource.NewQuantity(0, resource.BinarySI)

			for _, u := range list.Items {
				nodeTotal := resource.NewQuantity(0, resource.BinarySI)

				// .status.diskStatus is a map keyed by disk UUID
				disks, _, _ := unstructured.NestedMap(u.Object, "status", "diskStatus")
				for _, v := range disks {
					if disk, ok := v.(map[string]interface{}); ok {
						switch x := disk["storageAvailable"].(type) {
						case int64:
							nodeTotal.Add(*resource.NewQuantity(x, resource.BinarySI))
						case float64:
							nodeTotal.Add(*resource.NewQuantity(int64(x), resource.BinarySI))
						case string:
							if i, err := strconv.ParseInt(x, 10, 64); err == nil {
								nodeTotal.Add(*resource.NewQuantity(i, resource.BinarySI))
							}
						}
					}
				}

				// Keep the largest node total seen so far
				if nodeTotal.Cmp(*maxFree) > 0 {
					maxFree = nodeTotal
				}
			}

			gbValue := float64(maxFree.Value()) / (1024 * 1024 * 1024)
			sizeGBValue := fmt.Sprintf("%.2f", gbValue) // Format to 2 decimal places
			// if gbValue is a whole number, remove .00
			if strings.HasSuffix(sizeGBValue, ".00") {
				sizeGBValue = strings.TrimSuffix(sizeGBValue, ".00")
			}
			resp.AllocatableGB = sizeGBValue
			resp.UnableToDetectAllocatable = false
			resp.MinimumStorageGB = "0.10"
			resp.StorageStepGB = "0.25" // 256 MiB
			return resp, nil

		// * Hetzner - predefined limits
		case "csi.hetzner.cloud":
			resp.AllocatableGB = "10000" // 10 TiB
			resp.MinimumStorageGB = "10" // 10 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil

		// * AWS EBS - predefined limits
		case "ebs.csi.aws.com", "kubernetes.io/aws-ebs":
			resp.AllocatableGB = "64000" // 64 TiB
			resp.MinimumStorageGB = "1"  // 1 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil

		// * Azure Disk - predefined limits
		case "disk.csi.azure.com", "kubernetes.io/azure-disk":
			resp.AllocatableGB = "64000" // 64 TiB
			resp.MinimumStorageGB = "1"  // 1 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil

		// * GCP PD - predefined limits
		case "pd.csi.storage.gke.io", "pd.csi.storage.k8s.io":
			resp.AllocatableGB = "64000" // 64 TiB
			resp.MinimumStorageGB = "1"  // 1 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil

		// * DigitalOcean Volumes - predefined limits
		case "dobs.csi.digitalocean.com":
			resp.AllocatableGB = "16000" // 16 TiB
			resp.MinimumStorageGB = "1"  // 1 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil

		// * Vultr Block Storage - predefined limits
		case "csi.vultr.com":
			resp.AllocatableGB = "10000" // 10 TiB
			resp.MinimumStorageGB = "10" // 10 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil

		// * Linode Block Storage - predefined limits
		case "linodebs.csi.linode.com":
			resp.AllocatableGB = "16000" // 16 TiB
			resp.MinimumStorageGB = "10" // 10 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil

		// * OpenStack Cinder - predefined limits (OVH and others)
		case "cinder.csi.openstack.org":
			resp.AllocatableGB = "12000" // 12 TiB
			resp.MinimumStorageGB = "10" // 10 GiB
			resp.StorageStepGB = "1"     // 1 GiB
			resp.UnableToDetectAllocatable = false
			return resp, nil
		}
	}
	// No recognised driver – leave UnableToDetectAllocatable = true
	return resp, nil
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
