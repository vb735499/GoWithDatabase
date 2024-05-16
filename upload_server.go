package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadServer struct {
	SEngine *gin.Engine
}

func removeFile(filePath string) {
	e := os.Remove(filePath)
	if e != nil {
		log.Fatal(e)
	}
}

func toJson(_files map[string][]string) []map[string]string {
	out := []map[string]string{}

	for user, files := range _files {
		for _, filename := range files {
			title := strings.Split(filename, ".")
			tmp := map[string]string{
				"id":       uuid.New().String(),
				"username": user,
				"date":     "123",
				"title":    title[len(title)-2],
				"image":    filename,
			}
			out = append(out, tmp)
		}
	}
	log.Println(out)
	return out
}

func jsonMsg(statusCode string, Msg string) []map[string]string {
	out := []map[string]string{}
	tmp := map[string]string{
		"statusCode": statusCode,
		"Msg":        Msg,
	}
	out = append(out, tmp)
	return out
}

func createUploadServer() UploadServer {
	bucketName := "pic-image"
	bucketClient := getClient()

	validation := []string{"jpg", "png", "jpeg", "gif"}

	router := UploadServer{
		SEngine: gin.Default(),
	}

	router.SEngine.Static("/imgs", "./imgs")
	router.SEngine.LoadHTMLGlob("./test/*")

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	router.SEngine.MaxMultipartMemory = 8 << 20 // 8 MiB

	router.SEngine.GET("/api", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html", gin.H{})
	})

	router.SEngine.GET("/api/query", func(c *gin.Context) {
		_files := bucketClient.QueryAll(bucketName)
		c.JSON(200, toJson(_files))
	})

	router.SEngine.POST("/api/upload", func(c *gin.Context) {
		uploadInfo := ""
		uploadSucces := 0

		form, _ := c.MultipartForm()
		files := form.File["upload[]"]
		memberId := form.Value["memberId"][0]
		log.Println(form)

		for _, file := range files {
			fileType := strings.Split(file.Filename, ".")[1]
			if !slices.Contains(validation, fileType) {
				uploadInfo += fmt.Sprintf("'%s' are not vaild image types.\n", file.Filename)
				log.Println(file.Filename, "are not valid image types.")
				continue
			}
			log.Println(file.Filename)
			uploadSucces += 1

			filePath := "./imgs/" + memberId + "/" + file.Filename

			// Upload the file to specific dst.
			err := c.SaveUploadedFile(file, filePath)
			ErrorOccurMsg(err)
			if err != nil {
				removeFile(filePath)
				c.JSON(http.StatusBadRequest, jsonMsg(fmt.Sprint(http.StatusBadRequest), "Failed to save file."))
			}

			// Upload to the AWS S3 bucket.
			upload_err := bucketClient.UploadFile(bucketName, memberId, filePath)
			if upload_err != nil {
				removeFile(filePath)
				c.JSON(http.StatusBadRequest, jsonMsg(fmt.Sprint(http.StatusBadRequest), "Upload file to AWS S3 bucket failed."))
			}
			// Remove file from web server.
			removeFile(filePath)
		}
		uploadInfo += fmt.Sprintf("'%d' files uploaded!\n", uploadSucces)
		log.Printf("%v", uploadInfo)

		c.JSON(http.StatusOK, jsonMsg(fmt.Sprint(http.StatusOK), uploadInfo))
		// c.String(http.StatusOK, uploadInfo)

		// c.Redirect(http.StatusOK, c.GetHeader("authority")+"/")
	})
	return router
}
