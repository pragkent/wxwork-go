package internal

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestStringSet_Marshal(t *testing.T) {
	tests := []struct {
		s    StringSet
		want []byte
	}{
		{
			StringSet{},
			[]byte(`""`),
		},
		{
			StringSet{"id0"},
			[]byte(`"id0"`),
		},
		{
			StringSet{"id0", "id1", "id2"},
			[]byte(`"id0|id1|id2"`),
		},
	}

	for _, tt := range tests {
		b, err := json.Marshal(tt.s)
		if err != nil {
			t.Errorf("json.Marshal returned error: %v", err)
		}

		if !reflect.DeepEqual(b, tt.want) {
			t.Errorf("json.Marshal returned %q; want %q", b, tt.want)
		}
	}
}

func TestStringSet_Unmarshal(t *testing.T) {
	tests := []struct {
		b    []byte
		want StringSet
	}{
		{
			[]byte(`""`),
			nil,
		},
		{
			[]byte(`"id0"`),
			StringSet{"id0"},
		},
		{
			[]byte(`"id0|id1|id2"`),
			StringSet{"id0", "id1", "id2"},
		},
		{
			[]byte(`"id0|id1||id2"`),
			StringSet{"id0", "id1", "", "id2"},
		},
	}

	for _, tt := range tests {
		var v StringSet
		if err := json.Unmarshal(tt.b, &v); err != nil {
			t.Errorf("json.Unmarshal returned error: %v", err)
		}

		if !reflect.DeepEqual(v, tt.want) {
			t.Errorf("json.Marshal returned %#v; want %#v", v, tt.want)
		}
	}
}

func TestStringSet_UnmarshalError(t *testing.T) {
	tests := []struct {
		b []byte
	}{
		{
			[]byte(`id0`),
		},
		{
			[]byte(`"id0|id1|id2`),
		},
	}

	for _, tt := range tests {
		var v StringSet
		if err := json.Unmarshal(tt.b, &v); err == nil {
			t.Errorf("json.Unmarshal returned no error")
		}
	}
}

func TestIntSet_Marshal(t *testing.T) {
	tests := []struct {
		s    IntSet
		want []byte
	}{
		{
			IntSet{},
			[]byte(`""`),
		},
		{
			IntSet{0},
			[]byte(`"0"`),
		},
		{
			IntSet{0, 1, 2},
			[]byte(`"0|1|2"`),
		},
	}

	for _, tt := range tests {
		b, err := json.Marshal(tt.s)
		if err != nil {
			t.Errorf("json.Marshal returned error: %v", err)
		}

		if !reflect.DeepEqual(b, tt.want) {
			t.Errorf("json.Marshal returned %q; want %q", b, tt.want)
		}
	}
}

func TestIntSet_Unmarshal(t *testing.T) {
	tests := []struct {
		b    []byte
		want IntSet
	}{
		{
			[]byte(`""`),
			nil,
		},
		{
			[]byte(`"0"`),
			IntSet{0},
		},
		{
			[]byte(`"0|1|2"`),
			IntSet{0, 1, 2},
		},
	}

	for _, tt := range tests {
		var v IntSet
		if err := json.Unmarshal(tt.b, &v); err != nil {
			t.Errorf("json.Unmarshal returned error: %v", err)
		}

		if !reflect.DeepEqual(v, tt.want) {
			t.Errorf("json.Marshal returned %v; want %v", v, tt.want)
		}
	}
}

func TestIntSet_UnmarshalError(t *testing.T) {
	tests := []struct {
		b []byte
	}{
		{
			[]byte(`id0`),
		},
		{
			[]byte(`"id0|id1|id2`),
		},
		{
			[]byte(`"|`),
		},
		{
			[]byte(`"id"`),
		},
	}

	for _, tt := range tests {
		var v IntSet
		if err := json.Unmarshal(tt.b, &v); err == nil {
			t.Errorf("json.Unmarshal returned no error")
		}
	}
}
