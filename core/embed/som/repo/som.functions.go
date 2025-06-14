//go:build embed

package repo

import (
	"encoding/json"
	"fmt"
)

const statusOK = "OK"

type RawQuery[T any] struct {
	Status string
	Time   string
	Result T
	Detail string
}

func Unmarshal[M any](respond interface{}) (model M, err error) {
	var bytes []byte

	if arrResp, isArr := respond.([]interface{}); len(arrResp) > 0 {
		if dataMap, ok := arrResp[0].(map[string]interface{}); ok && isArr {
			if _, ok := dataMap["status"]; ok {
				if bytes, err = json.Marshal(respond); err == nil {
					var raw []RawQuery[M]
					if err = json.Unmarshal(bytes, &raw); err == nil {
						if raw[0].Status != statusOK {
							err = fmt.Errorf("%s: %s", raw[0].Status, raw[0].Detail)
						}
						model = raw[0].Result
					}
				}
				return model, err
			}
		}
	}

	if bytes, err = json.Marshal(respond); err == nil {
		err = json.Unmarshal(bytes, &model)
	}

	return model, err
}
