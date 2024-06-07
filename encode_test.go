package msgpack

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var runLargeTests = flag.Bool("runLargeTests", false, "Set to true to run large tests")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestEncodeArrayMaxElements(t *testing.T) {
	if !*runLargeTests {
		t.Skip("Skipping large test")
	}

	const maxElements2 = (1 << 32) - 1

	// Preallocate the slice with its maximum capacity
	array := make([]interface{}, maxElements2)

	for i := 0; i < maxElements2; i++ {
		array[i] = float64(i)
	}
	_, err := Encode(array)
	assert.NoError(t, err)

	array = append(array, float64(100)) // Adding one more element
	_, err = Encode(array)
	assert.Error(t, err)
	assert.EqualError(t, err, "array exceeds maximum number of elements")
}

func TestEncodeMapMaxAssociations(t *testing.T) {
	if !*runLargeTests {
		t.Skip("Skipping large test")
	}

	data := make(map[string]interface{})
	for i := 0; i < (1<<32)-1; i++ {
		key := fmt.Sprintf("key%d", i)
		data[key] = i
	}
	_, err := EncodeStringInterface(data)
	assert.NoError(t, err)

	data["newKey"] = 100 // Adding one more key-value pair
	_, err = EncodeStringInterface(data)
	assert.Error(t, err)
	assert.EqualError(t, err, "map exceeds maximum number of key-value associations")
}

func TestEncodeBool(t *testing.T) {
	tests := []struct {
		input    bool
		expected string
	}{
		{true, "c3"},
		{false, "c2"},
	}

	for _, test := range tests {
		result := EncodeBool(test.input)
		if result != test.expected {
			t.Errorf("EncodeBool(%v) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestEncodeString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      string
		expectErr bool
	}{
		{
			name:  "Empty string",
			input: "",
			want:  "a0",
		},
		{
			name:  "Short string",
			input: "abc",
			want:  "a3616263",
		},
		{
			name:  "String of length 16",
			input: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:  "bf61616161616161616161616161616161616161616161616161616161616161",
		},
		{
			name:  "String of length 32",
			input: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:  "d9206161616161616161616161616161616161616161616161616161616161616161",
		},
		{
			name:      "String exceeding max length",
			input:     string(make([]byte, Max32Bit+1)),
			want:      "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeString(tt.input)
			if (err != nil) != tt.expectErr {
				t.Errorf("EncodeString() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodeFloat64(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{1.23, "cb3ff3ae147ae147ae"},
		{0.0, "cb0000000000000000"},
	}

	for _, test := range tests {
		result := EncodeFloat64(test.input)
		if result != test.expected {
			t.Errorf("EncodeFloat64(%v) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestEncodeNil(t *testing.T) {
	expected := fmt.Sprintf("%x", Nil)
	result := EncodeNil()
	if result != expected {
		t.Errorf("EncodeNil() = %s; expected %s", result, expected)
	}
}

func TestEncodeStringInterface(t *testing.T) {
	jsonData := `{"a":{"a":10}}`
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		t.Fatalf("Error unmarshalling JSON: %v", err)
	}

	expected := "81a16181a161cb4024000000000000" // Adjust based on your actual encoding logic
	// Capture the output of EncodeStringInterface
	out, _ := EncodeStringInterface(data)

	if out != expected {
		t.Errorf("EncodeStringInterface(%v) = %s; expected %s", data, out, expected)
	}
}
