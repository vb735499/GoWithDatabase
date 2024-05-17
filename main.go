package main

import "log"

func main() {
	router := createUploadServer()
	err := router.SEngine.Run(":8080")
	if err != nil {
		log.Println(err)
	}
}
