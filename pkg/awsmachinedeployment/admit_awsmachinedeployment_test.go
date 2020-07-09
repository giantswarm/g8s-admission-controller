package awsmachinedeployment

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"

	infrastructurev1alpha2 "github.com/giantswarm/apiextensions/pkg/apis/infrastructure/v1alpha2"
	"k8s.io/api/admission/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/giantswarm/admission-controller/pkg/unittest"
)

var (
	machineDeploymentID        = "mymd1"
	machineDeploymentNamespace = "mymd1"
)

func TestAWSMachineDeploymentAdmit(t *testing.T) {
	testCases := []struct {
		name                    string
		ctx                     context.Context
		currentAvailabilityZone []string
		// expectAvailabilityZones needs to be in order
		expectAvailabilityZones []string
		validAvailabilityZones  []string
	}{
		{
			name: "case 0",
			ctx:  context.Background(),
		},
	}

	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error

			fakeK8sClient := unittest.FakeK8sClient()
			admit := &Admitter{
				k8sClient: fakeK8sClient,
			}

			// create AWSMachineDeployment
			err = fakeK8sClient.CtrlClient().Create(tc.ctx, awsMachineDeployment())
			if err != nil {
				t.Fatal(err)
			}

			// run admission request to update AWSControlPlane AZ's
			_, err = admit.Admit(awsMachineDeploymentAdmissionRequest())
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func awsMachineDeploymentAdmissionRequest() *v1beta1.AdmissionRequest {
	req := &v1beta1.AdmissionRequest{
		Kind: metav1.GroupVersionKind{
			Version: "infrastructure.giantswarm.io/v1alpha2",
			Kind:    "AWSMachineDeployment",
		},
		Resource: metav1.GroupVersionResource{
			Version:  "infrastructure.giantswarm.io/v1alpha2",
			Resource: "awsmachinedeployments",
		},
		Operation: v1beta1.Update,
		Object: runtime.RawExtension{
			Raw:    awsMachineDeploymentRawByte(),
			Object: nil,
		},
		OldObject: runtime.RawExtension{
			Raw:    awsMachineDeploymentRawByte(),
			Object: nil,
		},
	}
	return req
}

func awsMachineDeploymentRawByte() []byte {
	cr := infrastructurev1alpha2.AWSMachineDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSMachineDeployment",
			APIVersion: "infrastructure.giantswarm.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineDeploymentID,
			Namespace: machineDeploymentNamespace,
			Labels: map[string]string{
				"giantswarm.io/control-plane":   machineDeploymentID,
				"giantswarm.io/organization":    "giantswarm",
				"release.giantswarm.io/version": "11.5.0",
			},
		},
		Spec: infrastructurev1alpha2.AWSMachineDeploymentSpec{
			NodePool: infrastructurev1alpha2.AWSMachineDeploymentSpecNodePool{
				Description: "Some friendly name",
				Scaling: infrastructurev1alpha2.AWSMachineDeploymentSpecNodePoolScaling{
					Min: 1,
					Max: 30,
				},
			},
			Provider: infrastructurev1alpha2.AWSMachineDeploymentSpecProvider{
				InstanceDistribution: infrastructurev1alpha2.AWSMachineDeploymentSpecInstanceDistribution{
					OnDemandBaseCapacity:                10,
					OnDemandPercentageAboveBaseCapacity: 10,
				},
			},
		},
	}
	byt, _ := json.Marshal(cr)
	return byt
}

func awsMachineDeployment() *infrastructurev1alpha2.AWSMachineDeployment {
	return &infrastructurev1alpha2.AWSMachineDeployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "AWSMachineDeployment",
			APIVersion: "infrastructure.giantswarm.io/v1alpha2",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      machineDeploymentID,
			Namespace: machineDeploymentNamespace,
			Labels: map[string]string{
				"giantswarm.io/control-plane":   machineDeploymentID,
				"giantswarm.io/organization":    "giantswarm",
				"release.giantswarm.io/version": "11.5.0",
			},
		},
		Spec: infrastructurev1alpha2.AWSMachineDeploymentSpec{
			NodePool: infrastructurev1alpha2.AWSMachineDeploymentSpecNodePool{
				Description: "Some friendly name",
				Scaling: infrastructurev1alpha2.AWSMachineDeploymentSpecNodePoolScaling{
					Min: 1,
					Max: 30,
				},
			},
			Provider: infrastructurev1alpha2.AWSMachineDeploymentSpecProvider{
				InstanceDistribution: infrastructurev1alpha2.AWSMachineDeploymentSpecInstanceDistribution{
					OnDemandBaseCapacity:                10,
					OnDemandPercentageAboveBaseCapacity: 10,
				},
			},
		},
	}
}
