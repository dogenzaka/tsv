package tsv

import (
	"golang.org/x/text/unicode/norm"
	"os"
	"testing"
)

type TestRow struct {
	Name   string
	Age    int
	Gender string
	Active bool
}

type TestTaggedRow struct {
	Age    int    `tsv:"age"`
	Active bool   `tsv:"active"`
	Gender string `tsv:"gender"`
	Name   string `tsv:"name"`
}

func TestParserWithoutHeader(t *testing.T) {

	file, err := os.Open("example_simple.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	parser := NewParserWithoutHeader(file, &TestRow{})

	i := 0

	for {
		item := TestRow{}
		eof, err := parser.Next(&item)
		if eof {
			return
		}
		if i == 0 {
			if item.Name != "alex" ||
				item.Age != 10 ||
				item.Gender != "male" ||
				item.Active != true {
				t.Error("Record does not match index:0")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 1 {
			if item.Name != "john" ||
				item.Age != 24 ||
				item.Gender != "male" ||
				item.Active != false {
				t.Error("Record does not match index:1")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 2 {
			if item.Name != "sara" ||
				item.Age != 30 ||
				item.Gender != "female" ||
				item.Active != true {
				t.Error("Record does not match index:2")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 3 {
			if err == nil {
				t.Error("Error should be caused")
				return
			}
		}
		if i == 4 {
			if err == nil {
				t.Error("Error should be caused")
				return
			}
		}
		if i == 5 {
			if err == nil {
				t.Error("Error should be caused")
			}
		}
		i++
	}

}

func TestParserTaggedStructure(t *testing.T) {

	file, err := os.Open("example.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	parser, err := NewParser(file, &TestTaggedRow{})
	if err != nil {
		t.Error(err)
		return
	}

	i := 0

	for {
		item := TestTaggedRow{}
		eof, err := parser.Next(&item)
		if eof {
			return
		}
		if i == 0 {
			if err != nil {
				t.Error(err)
			}
			if item.Name != "alex" ||
				item.Age != 10 ||
				item.Gender != "male" ||
				item.Active != true {
				t.Error("Record does not match index:0")
			}
		}
		if i == 1 {
			if err != nil {
				t.Error(err)
			}
			if item.Name != "john" ||
				item.Age != 24 ||
				item.Gender != "male" ||
				item.Active != false {
				t.Error("Record does not match index:1")
			}
		}
		if i == 2 {
			if err != nil {
				t.Error(err)
			}
			if item.Name != "sara" ||
				item.Age != 30 ||
				item.Gender != "female" ||
				item.Active != true {
				t.Error("Record does not match index:2")
			}
		}
		i++
	}

}

func TestParserNormalize(t *testing.T) {

	file, err := os.Open("example_norm.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	parser, err := NewParser(file, &TestRow{})
	if err != nil {
		t.Error(err)
		return
	}
	// Use NFC as normalization
	parser.normalize = norm.NFKC

	i := 0

	for {
		item := TestRow{}
		eof, err := parser.Next(&item)
		if eof {
			return
		}
		if err != nil {
			t.Error(err)
		}
		if i == 0 && item.Name != "アレックス" {
			t.Errorf("name is not normalized %v", item.Name)
		}
		if i == 1 && item.Name != "デボラ" {
			t.Errorf("name is not normalized %v", item.Name)
		}
		if i == 2 && item.Name != "デボラ" {
			t.Errorf("name is not normalized %v", item.Name)
		}
		if i == 3 && item.Name != "(テスト)" {
			t.Errorf("name is not normalized %v", item.Name)
		}
		if i == 4 && item.Name != "/" {
			t.Errorf("name is not normalized %v", item.Name)
		}
		i++
	}

}
