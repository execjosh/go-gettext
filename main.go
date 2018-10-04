package main

import (
	"fmt"
	"log"

	"github.com/execjosh/go-gettext/pkg/mo"
)

func main() {
	d, err := mo.Load("messages.mo")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n", d.Gettext("%.2f MBytes transferred"))
	fmt.Printf("%v\n", d.Gettext("blah doesn't exist"))
}
