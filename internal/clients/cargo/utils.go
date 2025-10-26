package cargo

import "time"

func parseString(v interface{}) string {
	if v == nil {
		return ""
	}
	if str, ok := v.(string); ok {
		return str
	}
	return ""
}

func parseBool(v interface{}) *bool {
	if v == nil {
		b := false
		return &b
	}
	if str, ok := v.(string); ok {
		if str == "1" || str == "true" {
			b := true
			return &b
		}
		b := false
		return &b
	}
	if bool, ok := v.(bool); ok {
		return &bool
	}
	b := false
	return &b
}

func parseTime(v interface{}) *time.Time {
	if v == nil {
		t := time.Time{}
		return &t
	}
	if time, ok := v.(time.Time); ok {
		return &time
	}
	if str, ok := v.(string); ok {
		t, err := time.Parse(time.RFC3339, str)
		if err == nil {
			return &t
		}
	}
	t := time.Time{}
	return &t
}

func parseStringSlice(v interface{}) []string {
	if v == nil {
		return []string{}
	}
	if arr, ok := v.([]interface{}); ok {
		var result []string
		for _, item := range arr {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return []string{}
}

func parseInt(v interface{}) *int {
	if v == nil {
		return nil
	}
	if num, ok := (v).(int); ok {
		return &num
	}
	return nil
}
