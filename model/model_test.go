package model_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rwxrob/keg/model"
)

func ExampleDex_json() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := model.DexEntry{
		Updated: date,
		Node:    2,
		Title:   `Some title`,
	}
	byt, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(byt))
	// Output:
	// {"updated":"2022-12-10T06:10:04Z","title":"Some title","node":2}

}

func ExampleDex_string() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := model.DexEntry{
		Updated: date,
		Node:    2,
		Title:   `Some title`,
	}
	fmt.Println(d)
	fmt.Println(d.MD())
	// Output:
	// {"updated":"2022-12-10T06:10:04Z","title":"Some title","node":2}
	// * 2022-12-10 06:10:04Z [Some title](/2)

}
