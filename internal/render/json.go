package render

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v3"
)

// writeJSON writes the given value v as JSON to the writer w.
func writeJSON(w io.Writer, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", " ")
	return enc.Encode(v)
}

func writeYaml(w io.Writer, v any) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()
	return enc.Encode(v)
}
