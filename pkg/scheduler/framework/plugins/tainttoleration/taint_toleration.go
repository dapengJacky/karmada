package tainttoleration

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1helper "k8s.io/component-helpers/scheduling/corev1"

	clusterv1alpha1 "github.com/karmada-io/karmada/pkg/apis/cluster/v1alpha1"
	policyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	workv1alpha2 "github.com/karmada-io/karmada/pkg/apis/work/v1alpha2"
	"github.com/karmada-io/karmada/pkg/scheduler/framework"
)

const (
	// Name is the name of the plugin used in the plugin registry and configurations.
	Name = "TaintToleration"
)

// TaintToleration is a plugin that checks if a propagation policy tolerates a cluster's taints.
type TaintToleration struct{}

var _ framework.FilterPlugin = &TaintToleration{}

// New instantiates the TaintToleration plugin.
func New() (framework.Plugin, error) {
	return &TaintToleration{}, nil
}

// Name returns the plugin name.
func (p *TaintToleration) Name() string {
	return Name
}

// Filter checks if the given tolerations in placement tolerate cluster's taints.
func (p *TaintToleration) Filter(ctx context.Context, placement *policyv1alpha1.Placement, resource *workv1alpha2.ObjectReference, cluster *clusterv1alpha1.Cluster) *framework.Result {
	filterPredicate := func(t *corev1.Taint) bool {
		return t.Effect == corev1.TaintEffectNoSchedule || t.Effect == corev1.TaintEffectNoExecute
	}

	taint, isUntolerated := v1helper.FindMatchingUntoleratedTaint(cluster.Spec.Taints, placement.ClusterTolerations, filterPredicate)
	if !isUntolerated {
		return framework.NewResult(framework.Success)
	}

	return framework.NewResult(framework.Unschedulable, fmt.Sprintf("cluster had taint {%s: %s}, that the propagation policy didn't tolerate",
		taint.Key, taint.Value))
}
