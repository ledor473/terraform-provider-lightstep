package lightstep

import "encoding/json"

func marshalRules(rules map[string][]string) (string, error) {
	b, err := json.MarshalIndent(rules, "", " ")

	if err != nil {
		return "", err
	}
	return string(b), nil
}

func unmarshalRules(rules string) (map[string][]string, error) {
	var r = make(map[string][]string)
	err := json.Unmarshal([]byte(rules), &r)
	return r, err
}
