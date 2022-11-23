package keg_test

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/rwxrob/keg"
)

func ExampleDex_json() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := keg.DexEntry{U: date, N: 2, T: `Some title`}
	byt, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(byt))
	// Output:
	// {"U":"2022-12-10T06:10:04Z","T":"Some title","N":2}

}

func ExampleDex_string() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := keg.DexEntry{U: date, N: 2, T: `Some title`}
	fmt.Println(d)
	fmt.Println(d.MD())
	// Output:
	// * 2022-12-10 06:10:04Z [Some title](/2)
	// * 2022-12-10 06:10:04Z [Some title](/2)
}

func ExampleDex_tsv() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d := keg.DexEntry{U: date, N: 2, T: `Some title`}
	fmt.Println(d.TSV())
	// Output:
	// 2	2022-12-10 06:10:04Z	Some title
}

/*
func ExampleDex_Pretty() {
	date := time.Date(2022, 12, 10, 6, 10, 4, 0, time.UTC)
	d1 := keg.DexEntry{U: date, N: 2000, T: `Some title`}
	d2 := keg.DexEntry{U: date, N: 1, T: `Another title`}
	dex := keg.Dex{d1, d2}
	fmt.Println(dex.Pretty())
	// Output:
	// ignored
}
*/
