package main

import (
	"fmt"
	"os"

	tsv "github.com/stefantds/tsv"
)

type TestExampleRow struct {
	Name            string   // 0
	Age             int      // 1
	Gender          string   // 2
	FavoriteNumbers []int    // 3
	FavoriteHeroes  []string // 4
}

type ExampleStringArrayRow struct {
	Heroes []string
}

type ExampleIntArrayRow struct {
	Heroes []int
}

type ExampleFloatArrayRow struct {
	Heroes []float64
}

type ExampleIntArrayArrayRow struct {
	Heroes [][]int
}

func main() {
	fmt.Println("running")
	file, err := os.Open("/Users/stefan.tudose/private/gitrepos/tsv/example_int_array_array.tsv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	parser, err := tsv.NewParser(file, &ExampleIntArrayArrayRow{})
	if err != nil {
		panic(err)
	}

	for {
		data := ExampleIntArrayArrayRow{}
		eof, err := parser.Next(&data)
		if eof {
			return
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(data)
	}
}
