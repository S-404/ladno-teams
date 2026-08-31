package entity

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = JSONMap{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("cannot scan %T into JSONMap", value)
	}

	if len(bytes) == 0 {
		*j = JSONMap{}
		return nil
	}

	return json.Unmarshal(bytes, j)
}

func JSONMapFromDTO(data map[string]interface{}) JSONMap {
	if data == nil {
		return nil
	}
	return JSONMap(data)
}
