package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func upload(c *gin.Context) {
	// c.Request.ParseMultipartForm(10 << 20)

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer file.Close()

	outFile, err := ioutil.TempFile("public", "img-*.png")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = io.Copy(outFile, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var filePath string
	if os.Getenv("APP_ENV") == "docker" {
		filePath = "https://api-food-portal.systemprojects.net/file/" + filepath.Base(outFile.Name())
	} else {
		filePath = "http://localhost:7002/file/" + filepath.Base(outFile.Name())
	}

	c.JSON(http.StatusOK, gin.H{
		"file_path": filePath,
	})
}
