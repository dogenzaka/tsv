package tsv

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Parser has information for parser
type Parser struct {
	Reader     *csv.Reader
	structMode bool
}

// NewStructModeParser creates new TSV parser with given io.Reader as struct mode
func NewParser(reader io.Reader) (*Parser, error) {
	r := csv.NewReader(reader)
	r.Comma = '\t'
	r.LazyQuotes = true

	p := &Parser{
		Reader:     r,
		structMode: false,
	}

	return p, nil
}

// Next puts reader forward by a line
func (p *Parser) Next(data interface{}) (eof bool, err error) {

	decoderType := reflect.TypeOf((*Decoder)(nil)).Elem()

	dataReflected := reflect.ValueOf(data)

	// Get data value, while resolving pointers and interfaces
	for {
		if dataReflected.Kind() == reflect.Interface || dataReflected.Kind() == reflect.Ptr {
			dataReflected = dataReflected.Elem()
		} else {
			break
		}
	}

	records, err := p.Reader.Read()
	if err != nil {
		if err.Error() == "EOF" {
			return true, nil
		}
		return false, err
	}

	for i, record := range records {
		field := dataReflected.Field(i)

		switch field.Kind() {
		case reflect.Ptr:
			fieldType := field.Type()
			if field.Type().Implements(decoderType) {
				if field.IsZero() {
					// create a new value
					newValue := reflect.New(fieldType.Elem())
					field.Set(newValue)
				}

				fieldDecoder, _ := field.Addr().Elem().Interface().(Decoder)
				fieldDecoder.DecodeRecord(record)
			} else {
				return false, errors.New("Unsupported pointer to struct that doesn't implement the Decoder interface")
			}

		case reflect.Interface:
			fieldType := field.Type()
			if field.Type().Implements(decoderType) {
				if field.IsZero() {
					// create a new value
					newValue := reflect.New(fieldType.Elem())
					field.Set(newValue)
				}

				fieldDecoder, _ := field.Addr().Elem().Interface().(Decoder)
				fieldDecoder.DecodeRecord(record)
			} else {
				return false, errors.New("Unsupported pointer to struct that doesn't implement deserializer")
			}
		case reflect.String:
			field.SetString(record)
		case reflect.Bool:
			if record == "" {
				field.SetBool(false)
			} else {
				col, err := strconv.ParseBool(record)
				if err != nil {
					return false, err
				}
				field.SetBool(col)
			}
		case reflect.Int:
			if record == "" {
				field.SetInt(0)
			} else {
				col, err := strconv.ParseInt(record, 10, 0)
				if err != nil {
					return false, err
				}
				field.SetInt(col)
			}
		case reflect.Slice:
			fieldType := field.Type()
			elemType := fieldType.Elem()
			switch elemType.Kind() {
			case reflect.String:
				values, err := parseStringArray(record)
				if err != nil {
					return false, fmt.Errorf("could not parse %s as JSON array: %w", record, err)
				}
				appendAllString(&field, elemType, values)
			case reflect.Int:
				values, err := parseIntArray(record)
				if err != nil {
					return false, fmt.Errorf("could not parse %s as JSON array: %w", record, err)
				}
				appendAllInt(&field, elemType, values)
			case reflect.Bool:
				values, err := parseBoolArray(record)
				if err != nil {
					return false, fmt.Errorf("could not parse %s as JSON array: %w", record, err)
				}
				appendAllBool(&field, elemType, values)
			case reflect.Float64:
				values, err := parseFloatArray(record)
				if err != nil {
					return false, fmt.Errorf("could not parse %s as JSON array: %w", record, err)
				}
				appendAllFloat(&field, elemType, values)
			}
		default:
			return false, errors.New("Unsupported field type")
		}
	}

	return false, nil
}

func parseIntArray(data string) ([]int64, error) {
	var result []int64
	if err := json.NewDecoder(strings.NewReader(data)).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseStringArray(data string) ([]string, error) {
	var result []string
	if err := json.NewDecoder(strings.NewReader(data)).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseBoolArray(data string) ([]bool, error) {
	var result []bool
	if err := json.NewDecoder(strings.NewReader(data)).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func parseFloatArray(data string) ([]float64, error) {
	var result []float64
	if err := json.NewDecoder(strings.NewReader(data)).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func appendAllInt(target *reflect.Value, elemType reflect.Type, values []int64) {
	elem := reflect.New(elemType).Elem()
	for _, i := range values {
		elem.SetInt(i)
		target.Set(reflect.Append(*target, elem))
	}
}

func appendAllString(target *reflect.Value, elemType reflect.Type, values []string) {
	elem := reflect.New(elemType).Elem()
	for _, i := range values {
		elem.SetString(i)
		target.Set(reflect.Append(*target, elem))
	}
}

func appendAllBool(target *reflect.Value, elemType reflect.Type, values []bool) {
	elem := reflect.New(elemType).Elem()
	for _, i := range values {
		elem.SetBool(i)
		target.Set(reflect.Append(*target, elem))
	}
}

func appendAllFloat(target *reflect.Value, elemType reflect.Type, values []float64) {
	elem := reflect.New(elemType).Elem()
	for _, i := range values {
		elem.SetFloat(i)
		target.Set(reflect.Append(*target, elem))
	}
}
