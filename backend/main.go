package main

import (
	"backend/analyzer"
	"backend/commands"
	"backend/global"
	"backend/structures"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"net/http"
	"strings"
)

type ExecuteRequest struct {
	Content string `json:"content" binding:"required"`
}

type ExecuteResponse struct {
	Result string `json:"result"`
}

type FolderResponse struct {
	Result []structures.FolderElement `json:"result"`
}

type FileResponse struct {
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
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World!",
		})
	})

	app.Post("/execute", func(c *fiber.Ctx) error {
		var req ExecuteRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		result := processContent(req.Content)
		return c.JSON(ExecuteResponse{
			Result: result,
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		result := login(req.PartitionId, req.Username, req.Password)
		return c.JSON(LoginResponse{
			Result: result,
		})
	})

	app.Post("/logout", func(c *fiber.Ctx) error {
		_, err := commands.ParserLogout([]string{})
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{
			"result": true,
		})
	})

	app.Get("/disks", func(c *fiber.Ctx) error {
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

		return c.JSON(DiskListResponse{
			Result: partitions,
		})
	})

	app.Get("/partitions/:diskId", func(c *fiber.Ctx) error {
		partitions := getPartitionsListDisk(c.Params("diskId"))

		return c.JSON(PartitionsListResponse{
			Result: partitions,
		})
	})

	app.Get("/filesystem/:partitionId", func(c *fiber.Ctx) error {
		isFile := c.Query("type") == "file"

		if isFile {
			text, err := commands.ParserCat([]string{"-file1=" + c.Query("path")})
			if err != nil {
				return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}

			return c.JSON(FileResponse{
				Result: text,
			})
		}

		data, err := getElementsInFolder(c.Params("partitionId"), c.Query("path"))
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(FolderResponse{
			Result: data,
		})
	})

	err := app.Listen(":5000")
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

func getElementsInFolder(partitionID, path string) ([]structures.FolderElement, error) {
	mountedPartition, partitionPath, err := global.GetMountedPartition(partitionID)
	if err != nil {
		return nil, err
	}

	sb := &structures.SuperBlock{}
	if err := sb.ReadSuperBlock(partitionPath, int64(mountedPartition.PartStart)); err != nil {
		return nil, err
	}

	array := strings.Split(path, "/")
	var result []string
	for _, part := range array {
		if part != "" {
			result = append(result, part)
		}
	}

	indexInode := sb.GetInodeReference(partitionPath, 0, result)
	folderElement, err := sb.GetInodeElements(partitionPath, indexInode)
	if err != nil {
		return nil, err
	}

	return folderElement, nil
}
