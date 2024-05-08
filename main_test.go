package main

import (
	"testing"
)

func TestPass(t *testing.T) {
	router := createUploadServer()
	router.SEngine.Run(":8080")
}
