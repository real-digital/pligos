package pligos

import (
	"helm.sh/helm/v3/pkg/chart"
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
