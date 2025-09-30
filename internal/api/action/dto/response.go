package dto

import (
	"bytes"
	"fmt"
)

type NextActionProbability struct {
	Data map[string]float64
	Keys []string
}

// MarshalJSON implements custom JSON marshalling to ensure the order of keys is
// preserved because map iteration order is random.
func (n *NextActionProbability) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	for i, key := range n.Keys {
		if i > 0 {
			buf.WriteString(",")
		}

		fmt.Fprintf(&buf, `"%s":%.2f`, key, n.Data[key])
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}
