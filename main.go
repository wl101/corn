package main

import (
	"fmt"
	"log"

	"github.com/wl101/corn/internal"
)

func main() {
	db, err := internal.NewDB("corn.db")
	if err != nil {
		log.Fatal(err)
	}
	test := []string{"1", "2", "3", "4", "5"}
	for _, v := range test {
		err = db.Set(v, v)
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, v := range test {
		s, err := db.Get(v)
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("%s's val is %s\n", v, s)
	}
}
