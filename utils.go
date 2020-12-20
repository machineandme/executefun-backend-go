package main

import "log"

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
