package som

import (
	"encoding/json"
)

func toMap(val any) (map[string]any, error) {
	data, err := json.Marshal(val)
	if err != nil {
		return nil, err
	}

	var res map[string]any
	err = json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
