package utils

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

func WriteToFile(path string, offset int64, maxSize int64, data interface{}) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	dataSize := int64(binary.Size(data))
	if dataSize < 0 {
		return fmt.Errorf("failed to calculate data size")
	}

	if offset+dataSize > maxSize {
		return fmt.Errorf("data exceeds allowed space: offset=%d, dataSize=%d, maxSize=%d", offset, dataSize, maxSize)
	}

	if _, err = file.Seek(offset, 0); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	if err = binary.Write(file, binary.LittleEndian, data); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}

func ReadFromFile(path string, offset int64, data interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}
	fileSize := fileInfo.Size()

	if offset >= fileSize {
		return fmt.Errorf("offset is beyond the file size")
	}

	if _, err = file.Seek(offset, 0); err != nil {
		return fmt.Errorf("failed to seek file: %v", err)
	}

	dataSize := int64(binary.Size(data))
	if fileSize-offset < dataSize {
		return fmt.Errorf("not enough data to read: file size is %d but need %d bytes from offset %d", fileSize, dataSize, offset)
	}

	if err = binary.Read(file, binary.LittleEndian, data); err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	return nil
}

func ReadFromBitMap(path string, offset int64, end int64) (string, error) {
	if end <= offset {
		return "", fmt.Errorf("end must be greater than offset")
	}

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	if _, err = file.Seek(offset, 0); err != nil {
		return "", fmt.Errorf("failed to seek file: %v", err)
	}

	length := end - offset
	buffer := make([]byte, length)

	if _, err = io.ReadFull(file, buffer); err != nil {
		return "", fmt.Errorf("failed to read from file: %v", err)
	}

	allowedChars := "01OX"
	var result strings.Builder

	for _, char := range buffer {
		if strings.ContainsRune(allowedChars, rune(char)) {
			result.WriteRune(rune(char))
			result.WriteRune(' ')
		}
	}

	finalString := result.String()
	lines := splitIntoLines(finalString, 40)

	return lines, nil
}

func splitIntoLines(s string, n int) string {
	var builder strings.Builder
	for i, char := range s {
		builder.WriteRune(char)
		if (i+1)%n == 0 {
			builder.WriteRune('\n')
		}
	}
	return builder.String()
}
