package compiler

import (
	"fmt"
	"strings"
)

const identifierKey = "name"

type Compiler struct {
	config    map[string]interface{} // describes what instances are used
	schema    map[string]interface{} // describes the schema for the output
	types     map[string]interface{} // describes all available types
	instances map[string]interface{} // describes all instances
}

func New(config, schema, types, instances map[string]interface{}) *Compiler {
	return &Compiler{
		config:    config,
		schema:    schema,
		types:     types,
		instances: instances,
	}
}

func (ve *Compiler) repeatedLength(in interface{}) int {
	if s, ok := in.([]map[string]interface{}); ok {
		return len(s)
	}

	return len(in.([]interface{}))
}

func (ve *Compiler) Compile() (map[string]interface{}, error) {
	return ve.compile(
		ve.schema,
		ve.config,
	)
}

func (ve *Compiler) embed(res map[string]interface{}, config map[string]interface{}) map[string]interface{} {
	for k, v := range config {
		res[k] = v
	}

	return res
}

func (ve *Compiler) compile(schema map[string]interface{}, config map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	for k, typ := range schema {
		if _, ok := config[k]; !ok {
			continue
		}

		if isPrimitive(typ) {
			if _, ok := isOfType(typ, embedded); ok {
				res = ve.embed(res, config[k].(map[string]interface{}))
				continue
			}

			res[k] = config[k]
			continue
		}

		if typEmbeddedMapped, ok := isOfType(typ, embeddedMapped); ok {
			cs, err := ve.resolveConfigs(config[k], typEmbeddedMapped)
			if err != nil {
				return nil, err
			}

			for _, e := range cs {
				n, err := ve.nextConfig(e, typEmbeddedMapped)
				if err != nil {
					return nil, err
				}

				res[e[identifierKey].(string)] = n
			}
			continue
		}

		if typEmbedded, ok := isOfType(typ, embedded); ok {
			c, err := ve.resolveConfig(config[k], typEmbedded)
			if err != nil {
				return nil, err
			}

			toEmbed, err := ve.nextConfig(c, typEmbedded)
			if err != nil {
				return nil, err
			}

			res = ve.embed(res, toEmbed)
			continue
		}

		if typMapped, ok := isOfType(typ, mapped); ok {
			configs := make(map[string]interface{})

			cs, err := ve.resolveConfigs(config[k], typMapped)
			if err != nil {
				return nil, err
			}

			for _, e := range cs {
				n, err := ve.nextConfig(e, typMapped)
				if err != nil {
					return nil, err
				}

				configs[e[identifierKey].(string)] = n
			}

			res[k] = configs
			continue
		}

		if typRepeated, ok := isOfType(typ, repeated); ok {
			configs := make([]map[string]interface{}, 0, ve.repeatedLength(config[k]))

			cs, err := ve.resolveConfigs(config[k], typRepeated)
			if err != nil {
				return nil, err
			}

			for _, e := range cs {
				n, err := ve.nextConfig(e, typRepeated)
				if err != nil {
					return nil, err
				}

				configs = append(configs, n)
			}

			res[k] = configs
			continue
		}

		if _, ok := typ.(map[string]interface{}); ok {
			next, err := ve.compile(schema[k].(map[string]interface{}), config[k].(map[string]interface{}))
			if err != nil {
				return nil, err
			}

			res[k] = next
			continue
		}

		c, err := ve.resolveConfig(config[k], typ.(string))
		if err != nil {
			return nil, err
		}

		res[k], err = ve.nextConfig(c, typ.(string))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func (ve *Compiler) nextConfig(config map[string]interface{}, typ string) (map[string]interface{}, error) {
	t, err := ve.resolveType(typ)
	if err != nil {
		return nil, err
	}

	return ve.compile(t, config)
}

func (ve *Compiler) resolveConfigs(in interface{}, t string) ([]map[string]interface{}, error) {
	if n, ok := in.([]map[string]interface{}); ok {
		return n, nil
	}

	res := make([]map[string]interface{}, 0, len(in.([]interface{})))
	for _, e := range in.([]interface{}) {
		n, err := ve.resolveTypeInstance(t, e.(string))
		if err != nil {
			return nil, err
		}
		res = append(res, n)
	}

	return res, nil
}

func (ve *Compiler) resolveConfig(in interface{}, t string) (map[string]interface{}, error) {
	if n, ok := in.(map[string]interface{}); ok {
		return n, nil
	}

	return ve.resolveTypeInstance(t, in.(string))
}

func (ve *Compiler) resolveType(t string) (map[string]interface{}, error) {
	for k, v := range ve.types {
		if k == t {
			return v.(map[string]interface{}), nil
		}
	}

	return nil, fmt.Errorf("no such type defined: %s", t)
}

func (ve *Compiler) resolveTypeInstance(t, name string) (map[string]interface{}, error) {
	for k, v := range ve.instances {
		if k != t {
			continue
		}

		configs, ok := v.([]map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("bad configuration for type: %s", t)
		}

		for _, e := range configs {
			if e[identifierKey] == name {
				return e, nil
			}
		}

		break
	}

	return nil, fmt.Errorf("no instance of type %s defined", t)
}

func isPrimitive(t interface{}) bool {
	if _, ok := t.(string); !ok {
		return false
	}

	for _, e := range []specialType{repeated, mapped, embedded} {
		special, ok := isOfType(t, e)
		if ok {
			t = special
			break
		}
	}

	for _, e := range []string{"string", "numeric", "bool", "object"} {
		if t == e {
			return true
		}
	}

	return false
}

const (
	repeated       = specialType("repeated")
	mapped         = specialType("mapped")
	embedded       = specialType("embedded")
	embeddedMapped = specialType("embedded mapped")
)

type specialType string

func isOfType(t interface{}, typ specialType) (string, bool) {
	if _, ok := t.(string); !ok {
		return "", false
	}

	if !strings.HasPrefix(t.(string), string(typ)) {
		return "", false
	}

	return t.(string)[len(typ+" "):], true
}
