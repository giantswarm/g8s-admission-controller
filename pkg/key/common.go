package key

import (
	"github.com/giantswarm/aws-admission-controller/v2/pkg/label"
)

func Cluster(getter LabelsGetter) string {
	return getter.GetLabels()[label.Cluster]
}

func ControlPlane(getter LabelsGetter) string {
	return getter.GetLabels()[label.ControlPlane]
}

func MachineDeployment(getter LabelsGetter) string {
	return getter.GetLabels()[label.MachineDeployment]
}

func Organization(getter LabelsGetter) string {
	return getter.GetLabels()[label.Organization]
}
func Release(getter LabelsGetter) string {
	return getter.GetLabels()[label.Release]
}
