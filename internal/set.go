package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Separator of target elements.
const setSep = "|"

// StringSet helps encoding and decoding string slice.
type StringSet []string

func (s StringSet) MarshalJSON() ([]byte, error) {
	str := strings.Join(s, setSep)
	return json.Marshal(str)
}

func (s *StringSet) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	if len(str) == 0 {
		return nil
	}

	*s = StringSet(strings.Split(str, setSep))
	return nil
}

// Int helps encoding and decoding int slice.
type IntSet []int

func (s IntSet) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer

	for i, v := range s {
		if i != 0 {
			buf.WriteString(setSep)
		}

		buf.WriteString(strconv.Itoa(v))
	}

	return json.Marshal(buf.String())
}

func (s *IntSet) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	if len(str) == 0 {
		return nil
	}

	var r []int

	vs := strings.Split(str, setSep)
	for _, v := range vs {
		iv, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("wxwork: unmarshal int %q error: %v", v, err)
		}

		r = append(r, iv)
	}

	*s = r
	return nil
}
