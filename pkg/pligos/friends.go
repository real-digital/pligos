package pligos

import (
	"bytes"

	"gopkg.in/yaml.v2"
)

type FriendType string

const (
	Dependency  FriendType = "local"
	Requirement FriendType = "public"
)

type FriendConfig struct {
	Version        string                 `yaml:"version"`
	Name           string                 `yaml:"name"`
	Credentials    interface{}            `yaml:"credentials"`
	Config         interface{}            `yaml:"config"`
	Values         map[string]interface{} `yaml:"values"`
	ChartReference interface{}            `yaml:"chartReference"`
	FriendType     FriendType             `yaml:"type"`
}

type Friend struct {
	friends    []FriendConfig
	Normalizer *Normalizer
}

func NewFriend(input [][]byte) (*Friend, error) {
	if len(input) == 0 {
		return nil, nil
	}

	specSeperator := []byte("---")

	buf := make([]byte, 0)
	for _, e := range input {
		buf = append(buf, e...)
		buf = append(buf, specSeperator...)
	}

	buf = buf[:len(buf)-len(specSeperator)]

	specs := bytes.Split(buf, specSeperator)
	friends := make([]FriendConfig, 0, len(specs))
	for _, e := range specs {
		var friend FriendConfig
		err := yaml.Unmarshal(e, &friend)
		if err != nil {
			return nil, err
		}

		friends = append(friends, friend)
	}

	return &Friend{friends: friends, Normalizer: &Normalizer{}}, nil
}

func (f *Friend) EnrichValues(values map[string]interface{}) map[string]interface{} {
	if _, ok := values["dependencies"]; !ok {
		values["dependencies"] = make(map[string]interface{})
	}

	if _, ok := values["credentials"]; !ok {
		values["credentials"] = make(map[string]interface{})
	}

	for _, c := range f.friends {
		values["dependencies"].(map[string]interface{})[c.Name] = f.Normalizer.Normalize(c.Config)
		values["credentials"].(map[string]interface{})[c.Name] = f.Normalizer.Normalize(c.Credentials)
	}

	return values
}
