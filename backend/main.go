package main

import (
	"backend/analyzer"
	"backend/commands"
	"backend/global"
	"backend/structures"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ExecuteRequest struct {
	Content string `json:"content" binding:"required"`
}

type ExecuteResponse struct {
	Result string `json:"result"`
}

type LoginRequest struct {
	PartitionId string `json:"partitionId" binding:"required"`
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Result bool `json:"result"`
}

type DiskListResponse struct {
	Result []map[string]string `json:"result"`
}

type PartitionsListResponse struct {
	Result []map[string]string `json:"result"`
}

func main() {

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

	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result := login(req.PartitionId, req.Username, req.Password)

		c.JSON(http.StatusOK, LoginResponse{
			Result: result,
		})
	})

	r.POST("/logout", func(c *gin.Context) {
		_, err := commands.ParserLogout([]string{})

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"result": true,
		})
	})

	r.GET("/disks", func(c *gin.Context) {
		var partitions []map[string]string
		pathSet := make(map[string]bool)

		for id, path := range global.MountedPartitions {
			if _, exists := pathSet[path]; exists {
				continue
			}

			pathSet[path] = true

			parts := strings.Split(path, "/")
			name := parts[len(parts)-1]

			partitions = append(partitions, map[string]string{
				"id":   id,
				"name": name,
			})
		}

		c.JSON(http.StatusOK, DiskListResponse{
			Result: partitions,
		})
	})

	r.GET("/partitions/:diskId", func(c *gin.Context) {
		partitions := getPartitionsListDisk(c.Param("diskId"))

		c.JSON(http.StatusOK, PartitionsListResponse{
			Result: partitions,
		})
	})

	err := r.Run(":5000")

	if err != nil {
		return
	}
}

func processContent(content string) string {
	return analyzer.Analyzer(content)
}

func login(id, username, password string) bool {
	loginStrings := []string{
		"-id=" + id,
		"-user=" + username,
		"-pass=" + password,
	}

	_, err := commands.ParserLogin(loginStrings)

	if err != nil {
		return false
	}
	return true
}

func getPartitionsListDisk(diskId string) []map[string]string {
	var partitions []map[string]string

	_, path, err := global.GetMountedPartition(diskId)
	if err != nil {
		return []map[string]string{}
	}

	mbr := &structures.MBR{}
	if err := mbr.ReadMBR(path); err != nil {
		return []map[string]string{}
	}

	for _, part := range mbr.MbrPartition {
		if part.PartStatus == '1' {
			partitions = append(partitions, map[string]string{
				"id":   string(part.PartId[:]),
				"name": string(part.PartName[:]),
			})
		}
	}

	return partitions
}
