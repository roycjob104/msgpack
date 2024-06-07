package msgpack

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
)

func IsFixedMap(value byte) bool {
	return value >= FixedMapLow && value <= FixedMapHigh
}

func IsMap16(value byte) bool {
	return value == Map16
}

func IsMap32(value byte) bool {
	return value == Map32
}

func IsFixedArray(value byte) bool {
	return value >= FixedArrayLow && value <= FixedArrayHigh
}

func IsArray16(value byte) bool {
	return value == Array16
}

func IsArray32(value byte) bool {
	return value == Array32
}

func IsBool(value byte) bool {
	return value == True || value == False
}

func IsFloat64(value byte) bool {
	return value == Float64
}

func IsNil(value byte) bool {
	return value == Nil
}

func IsFixString(value byte) bool {
	return value >= FixedStrLow && value <= FixedStrHigh
}

func IsStr8(value byte) bool {
	return value == Str8
}

func IsStr16(value byte) bool {
	return value == Str16
}

func IsStr32(value byte) bool {
	return value == Str32
}

func DecodeBool(data []byte) (bool, int) {
	if IsBool(data[0]) {
		return data[0] == True, 1
	}
	return false, 0
}

func DecodeFloat64(data []byte) (float64, int) {
	if !IsFloat64(data[0]) {
		return 0, 0
	}

	if len(data) < 9 {
		return 0, 0
	}

	bits := binary.BigEndian.Uint64(data[1:9])
	return math.Float64frombits(bits), 9
}

func DecodeNil(data []byte) (interface{}, int) {
	if IsNil(data[0]) {
		return nil, 1
	}
	return nil, 0
}

func DecodeFixString(data []byte) (string, int, error) {
	length := int(data[0] - FixedStrLow)
	if len(data) < length+1 {
		return "", 0, errors.New("data length is insufficient")
	}
	return string(data[1 : length+1]), length + 1, nil
}

func DecodeStr8(data []byte) (string, int, error) {
	if len(data) < 2 {
		return "", 0, errors.New("data length is insufficient")
	}
	length := int(data[1])
	if len(data) < length+2 {
		return "", 0, errors.New("data length is insufficient")
	}
	return string(data[2 : length+2]), length + 2, nil
}

func DecodeStr16(data []byte) (string, int, error) {
	if len(data) < 3 {
		return "", 0, errors.New("data length is insufficient")
	}
	length := int(binary.BigEndian.Uint16(data[1:3]))
	if len(data) < length+3 {
		return "", 0, errors.New("data length is insufficient")
	}
	return string(data[3 : length+3]), length + 3, nil
}

func DecodeStr32(data []byte) (string, int, error) {
	if len(data) < 5 {
		return "", 0, errors.New("data length is insufficient")
	}
	length := int(binary.BigEndian.Uint32(data[1:5]))
	if len(data) < length+5 {
		return "", 0, errors.New("data length is insufficient")
	}
	return string(data[5 : length+5]), length + 5, nil
}

func DecodeString(data []byte) (string, int, error) {
	if IsFixString(data[0]) {
		return DecodeFixString(data)
	} else if IsStr8(data[0]) {
		return DecodeStr8(data)
	} else if IsStr16(data[0]) {
		return DecodeStr16(data)
	} else if IsStr32(data[0]) {
		return DecodeStr32(data)
	}
	return "", 0, fmt.Errorf("unsupported string format: 0x%x", data[0])
}

func InitDecode(encoded string) (interface{}, error) {
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	result, _ := DecodeInterface(decoded)
	return result, nil
}

func DecodeArray(data []byte) ([]interface{}, int, error) {
	length := 0
	var n int

	switch {
	case IsFixedArray(data[0]):
		length = int(data[0] - FixedArrayLow)
		n = 1
	case IsArray16(data[0]):
		if len(data) < 3 {
			return nil, 0, errors.New("data length is insufficient")
		}
		length = int(binary.BigEndian.Uint16(data[1:3]))
		n = 3
	case IsArray32(data[0]):
		if len(data) < 5 {
			return nil, 0, errors.New("data length is insufficient")
		}
		length = int(binary.BigEndian.Uint32(data[1:5]))
		n = 5
	}

	array := make([]interface{}, length)
	for i := 0; i < length; i++ {
		elem, consumed := DecodeInterface(data[n:])
		array[i] = elem
		n += consumed
	}

	return array, n, nil
}

func DecodeMap(data []byte) (map[string]interface{}, int, error) {
	length := 0
	var n int

	switch {
	case IsFixedMap(data[0]):
		length = int(data[0] - FixedMapLow)
		n = 1
	case IsMap16(data[0]):
		if len(data) < 3 {
			return nil, 0, errors.New("data length is insufficient")
		}
		length = int(binary.BigEndian.Uint16(data[1:3]))
		n = 3
	case IsMap32(data[0]):
		if len(data) < 5 {
			return nil, 0, errors.New("data length is insufficient")
		}
		length = int(binary.BigEndian.Uint32(data[1:5]))
		n = 5
	}

	m := make(map[string]interface{}, length)
	for i := 0; i < length; i++ {
		key, consumed, err := DecodeString(data[n:])
		if err != nil {
			return nil, 0, err
		}
		n += consumed
		value, consumed := DecodeInterface(data[n:])
		m[key] = value
		n += consumed
	}

	return m, n, nil
}

func DecodeInterface(data []byte) (interface{}, int) {
	switch {
	case IsBool(data[0]):
		return DecodeBool(data)
	case IsNil(data[0]):
		return DecodeNil(data)
	case IsFixString(data[0]), IsStr8(data[0]), IsStr16(data[0]), IsStr32(data[0]):
		str, n, _ := DecodeString(data)
		return str, n
	case IsFloat64(data[0]):
		return DecodeFloat64(data)
	case IsFixedArray(data[0]), IsArray16(data[0]), IsArray32(data[0]):
		array, n, _ := DecodeArray(data)
		return array, n
	case IsFixedMap(data[0]), IsMap16(data[0]), IsMap32(data[0]):
		m, n, _ := DecodeMap(data)
		return m, n
	default:
		return nil, 0
	}
}
