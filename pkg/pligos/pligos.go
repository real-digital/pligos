package pligos

import (
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type Pligos struct {
	*chart.Chart
	Flavor *chart.Chart

	ContextSpec map[string]interface{}
	Schema      map[string]interface{}
	Types       map[string]interface{}
	Instances   map[string]interface{}

	Dependencies []Pligos
}
