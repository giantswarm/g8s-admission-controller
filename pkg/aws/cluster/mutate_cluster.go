// Package cluster intercepts write activity to Cluster objects.
package cluster

import (
	"github.com/blang/semver"
	"github.com/giantswarm/k8sclient/v4/pkg/k8sclient"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	admissionv1 "k8s.io/api/admission/v1"
	capiv1alpha2 "sigs.k8s.io/cluster-api/api/v1alpha2"

	"github.com/giantswarm/aws-admission-controller/v2/config"
	"github.com/giantswarm/aws-admission-controller/v2/pkg/aws"
	"github.com/giantswarm/aws-admission-controller/v2/pkg/key"
	"github.com/giantswarm/aws-admission-controller/v2/pkg/label"
	"github.com/giantswarm/aws-admission-controller/v2/pkg/mutator"
)

type Config struct {
	K8sClient k8sclient.Interface
	Logger    micrologger.Logger
}

// Mutator for Cluster object.
type Mutator struct {
	k8sClient k8sclient.Interface
	logger    micrologger.Logger
}

func NewMutator(config config.Config) (*Mutator, error) {
	if config.K8sClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.K8sClient must not be empty", config)
	}
	if config.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "%T.Logger must not be empty", config)
	}

	mutator := &Mutator{
		k8sClient: config.K8sClient,
		logger:    config.Logger,
	}

	return mutator, nil
}

// Mutate is the function executed for every matching webhook request.
func (m *Mutator) Mutate(request *admissionv1.AdmissionRequest) ([]mutator.PatchOperation, error) {
	var result []mutator.PatchOperation

	if request.DryRun != nil && *request.DryRun {
		return result, nil
	}
	if request.Operation == admissionv1.Create {
		return m.MutateCreate(request)
	}
	return result, nil
}

// MutateCreate is the function executed for every create webhook request.
func (m *Mutator) MutateCreate(request *admissionv1.AdmissionRequest) ([]mutator.PatchOperation, error) {
	var result []mutator.PatchOperation
	var patch []mutator.PatchOperation
	var err error

	// Parse incoming object
	cluster := &capiv1alpha2.Cluster{}
	if _, _, err := mutator.Deserializer.Decode(request.Object.Raw, nil, cluster); err != nil {
		return nil, microerror.Maskf(parsingFailedError, "unable to parse Cluster: %v", err)
	}

	releaseVersion, err := aws.ReleaseVersion(cluster, patch)
	if err != nil {
		return nil, microerror.Maskf(parsingFailedError, "unable to parse release version from Cluster")
	}

	patch, err = m.MutateOperatorVersion(*cluster, releaseVersion)
	if err != nil {
		return nil, microerror.Mask(err)
	}
	result = append(result, patch...)

	return result, nil
}

func (m *Mutator) MutateOperatorVersion(cluster capiv1alpha2.Cluster, releaseVersion *semver.Version) ([]mutator.PatchOperation, error) {
	var result []mutator.PatchOperation
	var patch []mutator.PatchOperation
	var err error

	if key.ClusterOperator(&cluster) != "" {
		return result, nil
	}
	// Retrieve the `Release` CR.
	release, err := aws.FetchRelease(&aws.Handler{K8sClient: m.k8sClient, Logger: m.logger}, releaseVersion)
	if err != nil {
		return nil, microerror.Mask(err)
	}

	// mutate the operator label
	patch, err = aws.MutateLabelFromRelease(&aws.Handler{K8sClient: m.k8sClient, Logger: m.logger}, &cluster, *release, label.ClusterOperatorVersion, "cluster-operator")
	if err != nil {
		return nil, microerror.Mask(err)
	}
	result = append(result, patch...)

	return result, nil
}

func (m *Mutator) Log(keyVals ...interface{}) {
	m.logger.Log(keyVals...)
}

func (m *Mutator) Resource() string {
	return "cluster"
}
