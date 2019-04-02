package pligos

import (
	"github.com/golang/protobuf/ptypes/any"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type Pligos struct {
	Metadata           *chart.Metadata
	Flavor             *chart.Chart
	ChartDependencies  []*chart.Chart
	ConfigurationFiles []*any.Any

	Contexts  map[string]interface{}
	Schema    map[string]interface{}
	Types     map[string]interface{}
	Instances map[string]interface{}

	Dependencies []Pligos
}
