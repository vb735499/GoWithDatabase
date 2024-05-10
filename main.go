package main

func main() {
	router := createUploadServer()
	err := router.SEngine.Run(":8080")
	ErrorOccurMsg(err)
}
