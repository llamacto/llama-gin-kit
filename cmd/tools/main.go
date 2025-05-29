package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zgiai/ginext/config"
	"github.com/zgiai/ginext/pkg/storage"
)

func main() {
	toolName := flag.String("tool", "", "Tool to run (generate-url or check-file)")
	flag.Parse()

	switch *toolName {
	case "generate-url":
		GeneratePresignedURL()
	case "check-file":
		CheckR2File()
	default:
		fmt.Printf("Unknown tool: %s\n", *toolName)
		fmt.Println("Available tools: generate-url, check-file")
		os.Exit(1)
	}
}

// GeneratePresignedURL 生成预签名URL
func GeneratePresignedURL() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r2Client := storage.NewR2Client(cfg)
	url, err := r2Client.GeneratePresignedURL("test.txt", "text/plain")
	if err != nil {
		log.Fatalf("Failed to generate presigned URL: %v", err)
	}

	fmt.Printf("Presigned URL: %s\n", url)
}

// CheckR2File 检查R2文件是否存在
func CheckR2File() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r2Client := storage.NewR2Client(cfg)
	exists, err := r2Client.FileExists("test.txt")
	if err != nil {
		log.Fatalf("Failed to check file: %v", err)
	}

	fmt.Printf("File exists: %v\n", exists)
}
