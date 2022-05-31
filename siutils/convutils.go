package siutils

import (
	jsoniter "github.com/json-iterator/go"
)

// func DecodeAny(input any, output any) error {
// 	b, err := json.Marshal(input)
// 	if err != nil {
// 		return err
// 	}

// 	return json.Unmarshal(b, output)

// }

var ji = jsoniter.Config{
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

// var ji = jsoniter.ConfigCompatibleWithStandardLibrary
// var ji = jsoniter.ConfigFastest

func DecodeAny(input any, output any) error {
	// js := jsoniter.ConfigDefault
	if b, err := ji.Marshal(input); err != nil {
		return err
	} else {
		if err = ji.Unmarshal(b, output); err != nil {
			return err
		}
	}
	return nil
}
