package awscluster

import (
	"context"
	"strconv"
	"testing"

	"github.com/giantswarm/micrologger/microloggertest"

	"github.com/giantswarm/aws-admission-controller/v2/pkg/aws"
	"github.com/giantswarm/aws-admission-controller/v2/pkg/mutator"
	"github.com/giantswarm/aws-admission-controller/v2/pkg/unittest"
)

func TestAWSClusterPodCIDR(t *testing.T) {
	testCases := []struct {
		ctx  context.Context
		name string

		currentPodCIDR  string
		expectedPodCIDR string
	}{
		{
			// Don't default the Pod CIDR if it is set
			name: "case 0",
			ctx:  context.Background(),

			currentPodCIDR:  unittest.DefaultPodCIDR,
			expectedPodCIDR: "",
		},
		{
			// Default the Pod CIDR if it is not set
			name: "case 1",
			ctx:  context.Background(),

			currentPodCIDR:  "",
			expectedPodCIDR: unittest.DefaultPodCIDR,
		},
	}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error
			var updatedCIDR map[string]string

			fakeK8sClient := unittest.FakeK8sClient()
			mutate := &Mutator{
				podCIDRBlock: unittest.DefaultPodCIDR,
				k8sClient:    fakeK8sClient,
				logger:       microloggertest.New(),
			}

			// run admission request to default AWSCluster Pod CIDR
			var patch []mutator.PatchOperation
			request, err := unittest.CustomAdmissionRequestAWSCluster(tc.currentPodCIDR)
			if err != nil {
				t.Fatal(err)
			}
			patch, err = mutate.Mutate(&request)
			if err != nil {
				t.Fatal(err)
			}
			// parse patches
			for _, p := range patch {
				if p.Path == "/spec/provider/pods" {
					updatedCIDR = p.Value.(map[string]string)
				}
			}
			// check if the pod CIDR is as expected
			if tc.expectedPodCIDR != updatedCIDR["cidrBlock"] {
				t.Fatalf("expected %#q to be equal to %#q", tc.expectedPodCIDR, updatedCIDR)
			}
		})
	}
}

func TestAWSClusterDescription(t *testing.T) {
	testCases := []struct {
		ctx  context.Context
		name string

		currentDescription string
		expectedPatch      string
	}{
		{
			// Don't default the Cluster Description if it is set
			name: "case 0",
			ctx:  context.Background(),

			currentDescription: "My cluster",
			expectedPatch:      "",
		},
		{
			// Default the Cluster Description if it is not set
			name: "case 1",
			ctx:  context.Background(),

			currentDescription: "",
			expectedPatch:      aws.DefaultClusterDescription,
		},
	}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error
			var updatedDescription string

			fakeK8sClient := unittest.FakeK8sClient()
			mutate := &Mutator{
				podCIDRBlock: unittest.DefaultPodCIDR,
				k8sClient:    fakeK8sClient,
				logger:       microloggertest.New(),
			}

			// run mutate function to default AWSCluster Description
			var patch []mutator.PatchOperation
			awscluster := unittest.DefaultAWSCluster()
			awscluster.Spec.Cluster.Description = tc.currentDescription
			patch, err = mutate.MutateDescription(awscluster)
			if err != nil {
				t.Fatal(err)
			}
			// parse patches
			for _, p := range patch {
				if p.Path == "/spec/cluster/description" {
					updatedDescription = p.Value.(string)
				}
			}
			// check if the pod CIDR is as expected
			if tc.expectedPatch != updatedDescription {
				t.Fatalf("expected %#q to be equal to %#q", tc.expectedPatch, updatedDescription)
			}
		})
	}
}
func TestAWSClusterDomain(t *testing.T) {
	testCases := []struct {
		ctx  context.Context
		name string

		currentDomain string
		expectedPatch string
	}{
		{
			// Don't default the Cluster DNS Domain if it is set
			name: "case 0",
			ctx:  context.Background(),

			currentDomain: unittest.DefaultClusterDNSDomain,
			expectedPatch: "",
		},
		{
			// Default the Cluster DNS Domain if it is not set
			name: "case 1",
			ctx:  context.Background(),

			currentDomain: "",
			expectedPatch: unittest.DefaultClusterDNSDomain,
		},
	}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error
			var updatedDomain string

			fakeK8sClient := unittest.FakeK8sClient()
			mutate := &Mutator{
				dnsDomain: unittest.DefaultClusterDNSDomain,
				k8sClient: fakeK8sClient,
				logger:    microloggertest.New(),
			}

			// run mutate function to default AWSCluster Description
			var patch []mutator.PatchOperation
			awscluster := unittest.DefaultAWSCluster()
			awscluster.Spec.Cluster.DNS.Domain = tc.currentDomain
			patch, err = mutate.MutateDomain(awscluster)
			if err != nil {
				t.Fatal(err)
			}
			// parse patches
			for _, p := range patch {
				if p.Path == "/spec/cluster/dns/domain" {
					updatedDomain = p.Value.(string)
				}
			}
			// check if the pod CIDR is as expected
			if tc.expectedPatch != updatedDomain {
				t.Fatalf("expected %#q to be equal to %#q", tc.expectedPatch, updatedDomain)
			}
		})
	}
}

func TestAWSClusterRegion(t *testing.T) {
	testCases := []struct {
		ctx  context.Context
		name string

		currentRegion string
		expectedPatch string
	}{
		{
			// Don't default the Cluster Region if it is set
			name: "case 0",
			ctx:  context.Background(),

			currentRegion: unittest.DefaultClusterRegion,
			expectedPatch: "",
		},
		{
			// Default the Cluster Region if it is not set
			name: "case 1",
			ctx:  context.Background(),

			currentRegion: "",
			expectedPatch: unittest.DefaultClusterRegion,
		},
	}
	for i, tc := range testCases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var err error
			var updatedRegion string

			fakeK8sClient := unittest.FakeK8sClient()
			mutate := &Mutator{
				region:    unittest.DefaultClusterRegion,
				k8sClient: fakeK8sClient,
				logger:    microloggertest.New(),
			}

			// run mutate function to default AWSCluster Description
			var patch []mutator.PatchOperation
			awscluster := unittest.DefaultAWSCluster()
			awscluster.Spec.Provider.Region = tc.currentRegion
			patch, err = mutate.MutateRegion(awscluster)
			if err != nil {
				t.Fatal(err)
			}
			// parse patches
			for _, p := range patch {
				if p.Path == "/spec/provider/region" {
					updatedRegion = p.Value.(string)
				}
			}
			// check if the pod CIDR is as expected
			if tc.expectedPatch != updatedRegion {
				t.Fatalf("expected %#q to be equal to %#q", tc.expectedPatch, updatedRegion)
			}
		})
	}
}
