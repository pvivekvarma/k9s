// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package dao

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestExtractMeta(t *testing.T) {
	c := load(t, "dr")
	m, ee := extractMeta(c)

	assert.Equal(t, 0, len(ee))
	assert.Equal(t, "destinationrules", m.Name)
	assert.Equal(t, "destinationrule", m.SingularName)
	assert.Equal(t, "DestinationRule", m.Kind)
	assert.Equal(t, "networking.istio.io", m.Group)
	assert.Equal(t, "v1alpha3", m.Version)
	assert.Equal(t, true, m.Namespaced)
	assert.Equal(t, []string{"dr"}, m.ShortNames)
	var vv metav1.Verbs
	assert.Equal(t, vv, m.Verbs)
}

func TestExtractInterfaceSlice(t *testing.T) {
	uu := map[string]struct {
		m  map[string]interface{}
		n  string
		nn []interface{}
		ee []error
	}{
		"plain": {
			m:  map[string]interface{}{"versions": []interface{}{map[string]interface{}{"name": "v1alpha1", "served": false}, map[string]interface{}{"name": "v1beta1", "served": true}}},
			n:  "versions",
			nn: []interface{}{map[string]interface{}{"name": "v1alpha1", "served": false}, map[string]interface{}{"name": "v1beta1", "served": true}},
		},
		"empty": {
			m: map[string]interface{}{},
			n: "versions",
		},
	}

	var ee []error
	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			ss, e := extractInterfaceSlice(u.m, u.n, ee)
			assert.Equal(t, u.ee, e)
			assert.Equal(t, u.nn, ss)
		})
	}
}

func TestExtractStrSlice(t *testing.T) {
	uu := map[string]struct {
		m  map[string]interface{}
		n  string
		nn []string
		ee []error
	}{
		"plain": {
			m:  map[string]interface{}{"shortNames": []string{"a", "b", "c"}},
			n:  "shortNames",
			nn: []string{"a", "b", "c"},
		},
		"empty": {
			m: map[string]interface{}{},
			n: "shortNames",
		},
	}

	var ee []error
	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			ss, e := extractStrSlice(u.m, u.n, ee)
			assert.Equal(t, u.ee, e)
			assert.Equal(t, u.nn, ss)
		})
	}
}

func TestExtractString(t *testing.T) {
	uu := map[string]struct {
		m  map[string]interface{}
		n  string
		s  string
		ee []error
	}{
		"plain": {
			m: map[string]interface{}{"blee": "fred"},
			n: "blee",
			s: "fred",
		},
		"missing": {
			m:  map[string]interface{}{},
			n:  "blee",
			ee: []error{fmt.Errorf("failed to extract string blee")},
		},
	}

	var ee []error
	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			as, ae := extractStr(u.m, u.n, ee)
			assert.Equal(t, u.ee, ae)
			assert.Equal(t, u.s, as)
		})
	}
}

func TestExtractBool(t *testing.T) {
	uu := map[string]struct {
		m  map[string]interface{}
		n  string
		b  bool
		ee []error
	}{
		"plain": {
			m: map[string]interface{}{"served": true},
			n: "served",
			b: true,
		},
		"missing": {
			m:  map[string]interface{}{},
			n:  "served",
			ee: []error{fmt.Errorf("failed to extract bool served")},
		},
	}

	var ee []error
	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			as, ae := extractBool(u.m, u.n, ee)
			assert.Equal(t, u.ee, ae)
			assert.Equal(t, u.b, as)
		})
	}
}

// Helpers...

func load(t *testing.T, n string) *unstructured.Unstructured {
	raw, err := os.ReadFile(fmt.Sprintf("testdata/%s.json", n))
	assert.Nil(t, err)

	var o unstructured.Unstructured
	err = json.Unmarshal(raw, &o)
	assert.Nil(t, err)

	return &o
}
