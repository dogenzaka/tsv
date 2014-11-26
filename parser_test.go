package tsv

import (
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

func TestParser(t *testing.T) {

	file, err := os.Open("example.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	data := TestRow{}
	parser, err := NewParser(file, &data)
	if err != nil {
		t.Error(err)
		return
	}

	i := 0

	for {
		eof, err := parser.Next()
		if eof {
			return
		}
		if i == 0 {
			if data.Name != "alex" ||
				data.Age != 10 ||
				data.Gender != "male" ||
				data.Active != true {
				t.Error("Record does not match index:0")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 1 {
			if data.Name != "john" ||
				data.Age != 24 ||
				data.Gender != "male" ||
				data.Active != false {
				t.Error("Record does not match index:1")
				if err != nil {
					t.Error(err)
				}
			}
		}
		if i == 2 {
			if data.Name != "sara" ||
				data.Age != 30 ||
				data.Gender != "female" ||
				data.Active != true {
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

func TestParserWithTag(t *testing.T) {

	file, err := os.Open("example.tsv")
	if err != nil {
		t.Error(err)
		return
	}
	defer file.Close()

	data := TestTaggedRow{}
	parser, err := NewParser(file, &data)
	if err != nil {
		t.Error(err)
		return
	}

	i := 0

	for {
		eof, err := parser.Next()
		if eof {
			return
		}
		if i == 0 {
			if err != nil {
				t.Error(err)
			}
			if data.Name != "alex" ||
				data.Age != 10 ||
				data.Gender != "male" ||
				data.Active != true {
				t.Error("Record does not match index:0")
			}
		}
		if i == 1 {
			if err != nil {
				t.Error(err)
			}
			if data.Name != "john" ||
				data.Age != 24 ||
				data.Gender != "male" ||
				data.Active != false {
				t.Error("Record does not match index:1")
			}
		}
		if i == 2 {
			if err != nil {
				t.Error(err)
			}
			if data.Name != "sara" ||
				data.Age != 30 ||
				data.Gender != "female" ||
				data.Active != true {
				t.Error("Record does not match index:2")
			}
		}
		i++
	}

}
