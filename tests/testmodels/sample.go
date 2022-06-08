package testmodels

import "encoding/json"

type Embedded struct {
	NilTypo string `json:"embedded_nil_"`
	Nil     string `json:"embedded_nil"`
}

type Sample struct {
	*Embedded
	Nil  string `json:"nil"`
	Int2 int    `json:"int2_"`
	Int3 *int   `json:"int3_"`
}

func (s *Sample) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}

type SampleList []Sample

func (s *SampleList) String() string {
	b, _ := json.Marshal(s)
	return string(b)
}
