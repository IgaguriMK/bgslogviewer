package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/IgaguriMK/bgslogviewer/api"
)

func main() {

	f, err := os.Open("sol.json")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	var v api.FactionsStatus
	err = json.NewDecoder(f).Decode(&v)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%v\n", v)
}
