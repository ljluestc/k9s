// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package model_test

import (
	"slices"
	"testing"

	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/dao"
	"github.com/derailed/k9s/internal/model"
	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestResolvePulseGVRs(t *testing.T) {
	old := dao.MetaAccess
	t.Cleanup(func() { dao.MetaAccess = old })

	dao.MetaAccess = dao.NewMeta()
	gg := model.ResolvePulseGVRs()
	assert.Equal(t, model.PulseGVRs, gg)
	assert.GreaterOrEqual(t, slices.Index(gg, client.IngGVR), 0)
	assert.Equal(t, -1, slices.Index(gg, client.HTTPRouteGVR))

	dao.MetaAccess.RegisterMeta(client.HTTPRouteGVR.String(), &metav1.APIResource{
		Name:    "httproutes",
		Group:   "gateway.networking.k8s.io",
		Version: "v1",
		Kind:    "HTTPRoute",
	})
	gg = model.ResolvePulseGVRs()
	assert.Len(t, gg, len(model.PulseGVRs))
	assert.Equal(t, -1, slices.Index(gg, client.IngGVR))
	assert.Equal(t, slices.Index(model.PulseGVRs, client.IngGVR), slices.Index(gg, client.HTTPRouteGVR))
	// The shared pulse list must not be mutated.
	assert.GreaterOrEqual(t, slices.Index(model.PulseGVRs, client.IngGVR), 0)
}
