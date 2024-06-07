package msgpack

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

func EncodeBool(value bool) string {
	if value {
		return fmt.Sprintf("%x", True)
	}
	return fmt.Sprintf("%x", False)
}

func EncodeString(value string) (string, error) {
	strlen := len(value)
	if strlen > Max32Bit {
		return "", errors.New("string exceeds maximum length")
	}

	var formatFirstByte uint8
	var header string
	var builder strings.Builder

	switch {
	case strlen <= 31:
		formatFirstByte = FixedStrLow + uint8(strlen)
		header = fmt.Sprintf("%02x", formatFirstByte)
	case strlen <= Max8Bit:
		formatFirstByte = Str8
		header = fmt.Sprintf("%02x%02x", formatFirstByte, strlen)
	case strlen <= Max16Bit:
		formatFirstByte = Str16
		header = fmt.Sprintf("%02x%04x", formatFirstByte, strlen)
	case strlen <= Max32Bit:
		formatFirstByte = Str32
		header = fmt.Sprintf("%02x%08x", formatFirstByte, strlen)
	}

	builder.WriteString(header)
	builder.WriteString(fmt.Sprintf("%x", value))

	return builder.String(), nil
}

func EncodeFloat64(value float64) string {
	bits := math.Float64bits(value)
	return fmt.Sprintf("%x", Float64) + fmt.Sprintf("%016x", bits)
}

func EncodeNil() string {
	return fmt.Sprintf("%x", Nil)
}

func EncodeArray(v []interface{}) (string, error) {
	arrlen := len(v)
	if arrlen > Max32Bit {
		return "", errors.New("string exceeds maximum length")
	}
	var formatFirstByte uint8
	var header string
	var builder strings.Builder
	switch {
	case arrlen <= 15:
		formatFirstByte := FixedArrayLow + uint8(len(v))
		header = fmt.Sprintf("%02x", formatFirstByte)
	case arrlen <= Max16Bit:
		formatFirstByte = Array16
		header = fmt.Sprintf("%02x%04x", formatFirstByte, arrlen)
	case arrlen <= Max32Bit:
		formatFirstByte = Array32
		header = fmt.Sprintf("%02x%08x", formatFirstByte, arrlen)
	}
	builder.WriteString(header)

	encodeElem := ""
	for _, elem := range v {
		elemEncoded, err := Encode(elem)
		if err != nil {
			return "", err
		}
		encodeElem = encodeElem + elemEncoded
	}
	builder.WriteString(encodeElem)

	return builder.String(), nil
}

func Encode(data interface{}) (string, error) {
	switch v := data.(type) {
	case bool:
		return EncodeBool(v), nil
	case string:
		return EncodeString(v)
	case float64:
		return EncodeFloat64(v), nil
	case map[string]interface{}:
		return EncodeStringInterface(v)
	case nil:
		return EncodeNil(), nil
	case []interface{}:
		if len(v) > (1<<32)-1 {
			return "", errors.New("array exceeds maximum number of elements")
		}
		return EncodeArray(v)
	default:
		return "", fmt.Errorf("unsupported type: %T, value: %v", v, v)
	}
}

func EncodeStringInterface(data map[string]interface{}) (string, error) {
	interfaceLen := len(data)
	if interfaceLen > Max32Bit {
		return "", errors.New("map exceeds maximum number of key-value associations")
	}
	var formatFirstByte uint8
	var header string
	var builder strings.Builder
	switch {
	case interfaceLen <= 15:
		formatFirstByte := FixedMapLow + uint8(len(data))
		header = fmt.Sprintf("%02x", formatFirstByte)
	case interfaceLen <= Max16Bit:
		formatFirstByte = Map16
		header = fmt.Sprintf("%02x%04x", formatFirstByte, interfaceLen)
	case interfaceLen <= Max32Bit:
		formatFirstByte = Map32
		header = fmt.Sprintf("%02x%08x", formatFirstByte, interfaceLen)
	}

	builder.WriteString(header)

	encodeElem := ""
	for key, value := range data {
		tmp, err := EncodeString(key)
		if err != nil {
			return "", err
		}
		encodeElem = encodeElem + tmp
		elemEncoded, err := Encode(value)
		if err != nil {
			return "", err
		}
		encodeElem = encodeElem + elemEncoded
	}
	builder.WriteString(encodeElem)

	return builder.String(), nil
}
