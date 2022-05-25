package siutils

import "encoding/json"

func DecodeAny(input any, output any) error {
	b, err := json.Marshal(input)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, output)

}
