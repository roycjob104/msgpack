package msgpack

import (
	"encoding/hex"
	"testing"
)

func TestDecodeBool(t *testing.T) {
	tests := []struct {
		input  string
		output bool
	}{
		{"c3", true},
		{"c2", false},
	}

	for _, test := range tests {
		data, _ := hex.DecodeString(test.input)
		result, _ := DecodeBool(data)
		if result != test.output {
			t.Errorf("DecodeBool(%s) = %v; want %v", test.input, result, test.output)
		}
	}
}

func TestDecodeFloat64(t *testing.T) {
	tests := []struct {
		input  string
		output float64
	}{
		{"cb3ff0000000000000", 1.0},
		{"cb4008000000000000", 3.0},
	}

	for _, test := range tests {
		data, _ := hex.DecodeString(test.input)
		result, _ := DecodeFloat64(data)
		if result != test.output {
			t.Errorf("DecodeFloat64(%s) = %v; want %v", test.input, result, test.output)
		}
	}
}

func TestDecodeNil(t *testing.T) {
	input := "c0"
	data, _ := hex.DecodeString(input)
	result, _ := DecodeNil(data)
	if result != nil {
		t.Errorf("DecodeNil(%s) = %v; want nil", input, result)
	}
}

func TestDecodeString(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"a161", "a"},
		{"a3616263", "abc"},
	}

	for _, test := range tests {
		data, _ := hex.DecodeString(test.input)
		result, _, err := DecodeString(data)
		if err != nil || result != test.output {
			t.Errorf("DecodeString(%s) = %v; want %v, error: %v", test.input, result, test.output, err)
		}
	}
}

func TestDecodeArray(t *testing.T) {
	tests := []struct {
		input  string
		output []interface{}
	}{
		{"93cb3ff0000000000000cb4000000000000000cb4008000000000000", []interface{}{float64(1), float64(2), float64(3)}},
	}

	for _, test := range tests {
		data, _ := hex.DecodeString(test.input)
		result, _, err := DecodeArray(data)
		if err != nil {
			t.Errorf("DecodeArray(%s) returned error: %v", test.input, err)
		}
		for i, v := range result {
			if v != test.output[i] {
				t.Errorf("DecodeArray(%s) = %v; want %v", test.input, result, test.output)
			}
		}
	}
}

func TestDecodeMap(t *testing.T) {
	tests := []struct {
		input  string
		output map[string]interface{}
	}{
		{"81a161cb3ff0000000000000", map[string]interface{}{"a": float64(1)}},
	}

	for _, test := range tests {
		data, _ := hex.DecodeString(test.input)
		result, _, err := DecodeMap(data)
		if err != nil {
			t.Errorf("DecodeMap(%s) returned error: %v", test.input, err)
		}
		for k, v := range result {
			if v != test.output[k] {
				t.Errorf("DecodeMap(%s) = %v; want %v", test.input, result, test.output)
			}
		}
	}
}
