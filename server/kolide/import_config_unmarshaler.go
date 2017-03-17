package kolide

import (
	"encoding/json"
	"errors"
	"strconv"
)

// UnmarshalJSON custom unmarshaling for PackNameMap will determine whether
// the pack section of an osquery config file refers to a file path, or
// pack details.  Pack details are unmarshalled into into PackDetails structure
// as oppossed to nested map[string]interface{}
func (pnm PackNameMap) UnmarshalJSON(b []byte) error {
	var temp map[string]interface{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}
	for key, val := range temp {
		switch t := val.(type) {
		case string:
			pnm[key] = t
		case map[string]interface{}:
			pnm[key], err = unmarshalPackDetails(t)
			if err != nil {
				return err
			}
		default:
			return errors.New("can't unmarshal json")
		}
	}
	return nil
}

func strptr(v interface{}) *string {
	if v == nil {
		return nil
	}
	s := new(string)
	*s = v.(string)
	return s
}

func boolptr(v interface{}) *bool {
	if v == nil {
		return nil
	}
	b := new(bool)
	*b = v.(bool)
	return b
}

func uintptr(v interface{}) *uint {
	if v == nil {
		return nil
	}
	i := new(uint)
	*i = uint(v.(float64))
	return i
}

func unmarshalPackDetails(v map[string]interface{}) (PackDetails, error) {
	queries, err := unmarshalQueryDetails(v["queries"])
	if err != nil {
		return PackDetails{}, nil
	}
	return PackDetails{
		Queries:   queries,
		Shard:     uintptr(v["shard"]),
		Version:   strptr(v["version"]),
		Platform:  v["platform"].(string),
		Discovery: unmarshalDiscovery(v["discovery"]),
	}, nil
}

func unmarshalDiscovery(val interface{}) []string {
	var result []string
	if val == nil {
		return result
	}
	v := val.([]interface{})
	for _, val := range v {
		result = append(result, val.(string))
	}
	return result
}

func unmarshalQueryDetails(v interface{}) (QueryNameToQueryDetailsMap, error) {
	var err error
	result := make(QueryNameToQueryDetailsMap)
	if v == nil {
		return result, nil
	}
	for qn, details := range v.(map[string]interface{}) {
		result[qn], err = unmarshalQueryDetail(details)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func unmarshalQueryDetail(val interface{}) (QueryDetails, error) {
	v := val.(map[string]interface{})
	interval, err := unmarshalInterval(v["interval"])
	if err != nil {
		return QueryDetails{}, err
	}
	return QueryDetails{
		Query:    v["query"].(string),
		Interval: interval,
		Removed:  boolptr(v["removed"]),
		Platform: strptr(v["platform"]),
		Version:  strptr(v["version"]),
		Shard:    uintptr(v["shard"]),
		Snapshot: boolptr(v["snapshot"]),
	}, nil
}

// It is valid for the interval can be a string that is convertable to an int,
// or an float64. The float64 is how all numbers in JSON are represented, so
// we need to convert to uint
func unmarshalInterval(val interface{}) (uint, error) {
	// if interval is nil return zero value
	if val == nil {
		return uint(0), nil
	}
	switch v := val.(type) {
	case string:
		i, err := strconv.ParseUint(v, 10, 64)
		return uint(i), err
	case float64:
		return uint(v), nil
	default:
		return uint(0), errors.New("type mismatch for interval value")
	}
}
