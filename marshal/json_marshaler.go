package marshal

import "encoding/json"

var (
	JSON = new(jsonMarshaler)
)

type jsonMarshaler struct {
}

func (m *jsonMarshaler) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m *jsonMarshaler) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (m *jsonMarshaler) String() string {
	return "json"
}
