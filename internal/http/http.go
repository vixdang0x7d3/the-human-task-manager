package http

import "encoding/json"

func parseTagsJSON(data []byte) (tags []string, err error) {
	if len(data) == 0 {
		return
	}

	var items []struct {
		Value string `json:"value"`
	}
	if err = json.Unmarshal(data, &items); err != nil {
		return []string{}, err
	}

	for _, item := range items {
		tags = append(tags, item.Value)
	}
	return
}
