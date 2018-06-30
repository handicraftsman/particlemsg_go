package particlemsg

import "reflect"

func typeof(v interface{}) string {
	return reflect.TypeOf(v).String()
}

func isSupersetOf(what, of map[string]interface{}) bool {
	for k, v := range of {
		if v2, ok := what[k]; ok {
			if typeof(v) != typeof(v2) {
				return false
			}
			switch v.(type) {
			case map[string]interface{}:
				if !isSupersetOf(v2.(map[string]interface{}), v.(map[string]interface{})) {
					return false
				}
				break
			}
		} else {
			return false
		}
	}
	return true
}
