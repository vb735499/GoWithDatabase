package main

func main() {
	router := createUploadServer()
	router.SEngine.Run(":8080")
}
