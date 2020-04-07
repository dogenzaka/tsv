package tsv

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Parser has information for parser
type Parser struct {
	Headers    []string
	Reader     *csv.Reader
	ref        reflect.Value
	indices    []int // indices is field index list of header array
	structMode bool
	normalize  norm.Form
}

// NewStructModeParser creates new TSV parser with given io.Reader as struct mode
func NewParser(reader io.Reader, data interface{}) (*Parser, error) {
	r := csv.NewReader(reader)
	r.Comma = '\t'
	r.LazyQuotes = true

	// first line should be fields
	headers, err := r.Read()

	if err != nil {
		return nil, err
	}

	for i, header := range headers {
		headers[i] = header
	}

	p := &Parser{
		Reader:     r,
		Headers:    headers,
		ref:        reflect.ValueOf(data).Elem(),
		indices:    make([]int, len(headers)),
		structMode: false,
		normalize:  -1,
	}

	// get type information
	t := p.ref.Type()

	for i := 0; i < t.NumField(); i++ {
		// get TSV tag
		tsvtag := t.Field(i).Tag.Get("tsv")
		if tsvtag != "" {
			// find tsv position by header
			for j := 0; j < len(headers); j++ {
				if headers[j] == tsvtag {
					// indices are 1 start
					p.indices[j] = i + 1
					p.structMode = true
				}
			}
		}
	}

	if !p.structMode {
		for i := 0; i < len(headers); i++ {
			p.indices[i] = i + 1
		}
	}

	return p, nil
}

// NewParserWithoutHeader creates new TSV parser with given io.Reader
func NewParserWithoutHeader(reader io.Reader, data interface{}) *Parser {
	r := csv.NewReader(reader)
	r.Comma = '\t'

	p := &Parser{
		Reader:    r,
		ref:       reflect.ValueOf(data).Elem(),
		normalize: -1,
	}

	return p
}

// Next puts reader forward by a line
func (p *Parser) Next(data interface{}) (eof bool, err error) {

	// Get data reflect value
	dataReflected := reflect.ValueOf(data).Elem()

	// Get next record
	var records []string

	for {
		// read until valid record
		records, err = p.Reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				return true, nil
			}
			return false, err
		}
		if len(records) > 0 {
			break
		}
	}

	if len(p.indices) == 0 {
		p.indices = make([]int, len(records))
		// mapping simple index
		for i := 0; i < len(records); i++ {
			p.indices[i] = i + 1
		}
	}

	// record should be a pointer
	for i, record := range records {
		idx := p.indices[i]
		if idx == 0 {
			// skip empty index
			continue
		}
		// get target field
		field := dataReflected.Field(idx - 1)
		switch field.Kind() {
		case reflect.String:
			// Normalize text
			if p.normalize >= 0 {
				record = p.normalize.String(record)
			}
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
