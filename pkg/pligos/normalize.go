package pligos

type Normalizer struct{}

func (n *Normalizer) Normalize(in interface{}) map[string]interface{} {
	return n.normalize(in).(map[string]interface{})
}

func (n *Normalizer) normalize(in interface{}) interface{} {
	switch t := in.(type) {
	case map[string]interface{}:
		res := make(map[string]interface{})
		for k, v := range t {
			res[k] = n.normalize(v)
		}
		return res
	case map[interface{}]interface{}:
		res := make(map[string]interface{})
		for k, v := range t {
			res[k.(string)] = n.normalize(v)
		}
		return res
	case []map[string]interface{}:
		res := make([]map[string]interface{}, 0, len(t))
		for _, e := range t {
			res = append(res, n.normalize(e).(map[string]interface{}))
		}
		return res
	case []map[interface{}]interface{}:
		res := make([]map[string]interface{}, 0, len(t))
		for _, e := range t {
			res = append(res, n.normalize(e).(map[string]interface{}))
		}
		return res
	case []interface{}:
		if len(t) == 0 {
			return make(map[string]interface{})
		}

		if _, ok := t[0].(map[interface{}]interface{}); ok {
			res := make([]map[string]interface{}, 0, len(t))
			for _, e := range t {
				res = append(res, n.normalize(e).(map[string]interface{}))
			}
			return res
		}

		if _, ok := t[0].(map[string]interface{}); ok {
			res := make([]map[string]interface{}, 0, len(t))
			for _, e := range t {
				res = append(res, n.normalize(e).(map[string]interface{}))
			}
			return res
		}

		res := make([]interface{}, 0, len(t))
		for _, e := range t {
			res = append(res, n.normalize(e))
		}
		return res
	default:
		return t
	}
}
