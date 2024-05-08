package main

import (
	"fmt"
	"testing"
)

func TestPass(t *testing.T) {
	router := createUploadServer()
	err := router.SEngine.Run(":8080")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
