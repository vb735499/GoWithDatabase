package main

func main() {
	// bucketName := "pic-image"
	// bucketClient := getClient()
	// _files := bucketClient.QueryAll(bucketName)
	// log.Println(toJson(_files))

	// current_path, _ := os.Getwd()
	// prefix_path := "/imgs/"

	// username := "user1"
	// filename := "lena.png"
	// filepath := current_path + prefix_path + username + "/" + filename

	// bucketClient.UploadFile(bucketName, username, filepath)
	// bucketClient.DownloadFile(bucketName, username, filepath)
	// bucketClient.DeleteObjects(bucketName, username, []string{filename})
	router := createUploadServer()
	router.SEngine.Run(":8080")
}
