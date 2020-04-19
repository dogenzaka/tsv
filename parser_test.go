package tsv

import (
	"strings"
	"testing"
)

type MyDecoder struct {
	value string
}

func (d *MyDecoder) DecodeRecord(value string) error {
	d.value = value + "!"
	return nil
}

func (d *MyDecoder) DecodedValue() string {
	return d.value
}

func TestDecoder(t *testing.T) {
	type TestRow struct {
		Field      *MyDecoder
		OtherField int
	}

	for _, tc := range []struct {
		name      string
		RowStruct TestRow
		data      string
		expected  string
	}{
		{
			name:      "should work for a struct with uninitialized pointers",
			RowStruct: TestRow{},
			data:      "record1\t1\n",
			expected:  "record1!",
		},
		{
			name: "should work for a struct with initialized pointers",
			RowStruct: TestRow{
				Field: &MyDecoder{},
			},
			data:     "record2\n2\n",
			expected: "record2!",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data))
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for {
				eof, err := parser.Next(&tc.RowStruct)
				if err != nil {
					t.Error(err)
				}

				if tc.RowStruct.Field.DecodedValue() != tc.expected {
					t.Errorf("expected value '%s' got '%s'", tc.expected, tc.RowStruct.Field.DecodedValue())
				}

				if eof {
					break
				}
			}
		})
	}
}

func TestDecoderInterface(t *testing.T) {
	type TestRowWithInterface struct {
		Field      Decoder
		OtherField int
	}

	for _, tc := range []struct {
		name      string
		RowStruct TestRowWithInterface
		data      string
		expected  string
	}{
		{
			name: "should work for a struct using the Decoder interface",
			RowStruct: TestRowWithInterface{
				Field: &MyDecoder{},
			},
			data: "record1	1",
			expected: "record1!",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			parser, err := NewParser(strings.NewReader(tc.data))
			if err != nil {
				t.Fatalf("could not create parser: %w", err)
			}

			for {
				eof, err := parser.Next(&tc.RowStruct)
				if err != nil {
					t.Error(err)
				}

				if tc.RowStruct.Field.(*MyDecoder).DecodedValue() != tc.expected {
					t.Errorf("expected value '%s' got '%s'", tc.expected, tc.RowStruct.Field.(*MyDecoder).DecodedValue())
				}

				if eof {
					break
				}
			}
		})
	}
}
