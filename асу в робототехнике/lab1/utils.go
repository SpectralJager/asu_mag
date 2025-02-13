package main

import "log"

func Must[T any](val T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return val
}
