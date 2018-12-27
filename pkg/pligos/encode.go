package pligos

import (
	"fmt"
	"path/filepath"
	"strings"
)

const identifierKey = "name"

type ValuesEncoder struct {
	config map[string]interface{}
	schema map[string]interface{}
}

type ValuesEncoderMap map[string]*ValuesEncoder

func NewValuesEncoderMap(contexts []Context, types map[string]interface{}, values map[string]interface{}) (ValuesEncoderMap, error) {
	res := make(map[string]*ValuesEncoder)
	for _, context := range contexts {
		schema, err := CreateSchema(filepath.Join(context.Flavor, "schema.yaml"), types)
		if err != nil {
			return nil, err
		}

		res[context.Name] = NewValuesEncoder(values, schema)
	}

	return ValuesEncoderMap(res), nil
}

func (v ValuesEncoderMap) Get(context string) *ValuesEncoder {
	return v[context]
}

func NewValuesEncoder(config map[string]interface{}, schema map[string]interface{}) *ValuesEncoder {
	return &ValuesEncoder{
		config: config,
		schema: schema,
	}
}

func (ve *ValuesEncoder) repeatedLength(in interface{}) int {
	if s, ok := in.([]map[string]interface{}); ok {
		return len(s)
	}

	return len(in.([]interface{}))
}

func (ve *ValuesEncoder) handleMappedRoot(context, mappedType string) (map[string]interface{}, error) {
	res := make(map[string]interface{})
	typ, err := ve.resolveType(mappedType)
	if err != nil {
		return nil, err
	}

	for _, e := range ve.config["contexts"].(map[string]interface{})[context].([]interface{}) {
		i, err := ve.resolveTypeInstance(mappedType, e.(string))
		if err != nil {
			return nil, err
		}
		n, err := ve.Encode(typ, i)
		if err != nil {
			return nil, err
		}
		res[e.(string)] = n
	}

	return res, nil
}

func (ve *ValuesEncoder) EncodeContext(context string) (map[string]interface{}, error) {
	if mappedType, ok := isOfType(ve.schema["context"], mapped); ok {
		return ve.handleMappedRoot(context, mappedType)
	}

	return ve.Encode(
		ve.schema["context"].(map[string]interface{}),
		ve.config["contexts"].(map[string]interface{})[context].(map[string]interface{}),
	)
}

func (ve *ValuesEncoder) embed(res map[string]interface{}, config map[string]interface{}) map[string]interface{} {
	for k, v := range config {
		res[k] = v
	}

	return res
}

func (ve *ValuesEncoder) Encode(schema map[string]interface{}, config map[string]interface{}) (map[string]interface{}, error) {
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

func (ve *ValuesEncoder) nextConfig(config map[string]interface{}, typ string) (map[string]interface{}, error) {
	t, err := ve.resolveType(typ)
	if err != nil {
		return nil, err
	}

	return ve.Encode(t, config)
}

func (ve *ValuesEncoder) resolveConfigs(in interface{}, t string) ([]map[string]interface{}, error) {
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

func (ve *ValuesEncoder) resolveConfig(in interface{}, t string) (map[string]interface{}, error) {
	if n, ok := in.(map[string]interface{}); ok {
		return n, nil
	}

	return ve.resolveTypeInstance(t, in.(string))
}

func (ve *ValuesEncoder) resolveType(t string) (map[string]interface{}, error) {
	for k, v := range ve.schema {
		if k == t {
			return v.(map[string]interface{}), nil
		}
	}

	return nil, fmt.Errorf("no such type defined: %s", t)
}

func (ve *ValuesEncoder) resolveTypeInstance(t, name string) (map[string]interface{}, error) {
	for k, v := range ve.config {
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
	repeated = specialType("repeated")
	mapped   = specialType("mapped")
	embedded = specialType("embedded")
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
