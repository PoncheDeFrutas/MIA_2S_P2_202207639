package main

import (
	"backend/analyzer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

type ExecuteRequest struct {
	Content string `json:"content" binding:"required"`
}

type ExecuteResponse struct {
	Result string `json:"result"`
}

func main() {

	content, err := ioutil.ReadFile("initial.txt")
	if err != nil {
		log.Fatalf("Error al leer el archivo: %v", err)
	}

	processContent(string(content))

	r := gin.Default()

	r.Use(cors.Default())

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	r.POST("/execute", func(c *gin.Context) {
		var req ExecuteRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := processContent(req.Content)

		c.JSON(http.StatusOK, ExecuteResponse{
			Result: result,
		})
	})

	err = r.Run(":5000")

	if err != nil {
		return
	}
}

func processContent(content string) string {
	return analyzer.Analyzer(content)
}
