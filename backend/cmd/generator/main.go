package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"

	"github.com/hadi-projects/go-react-starter/internal/generator"
)

func main() {
	configPath := flag.String("config", "", "Path to module definition YAML")
	baseDir := flag.String("base", ".", "Base directory of the backend")
	flag.Parse()

	if *configPath == "" {
		log.Fatal("Please provide config path using -config flag")
	}

	absPath, _ := filepath.Abs(*baseDir)
	gen, err := generator.NewGenerator(*configPath, absPath)
	if err != nil {
		log.Fatalf("Failed to create generator: %v", err)
	}

	if err := gen.Generate(); err != nil {
		log.Fatalf("Failed to generate module: %v", err)
	}

	fmt.Printf("Successfully generated module: %s\n", gen.Config.ModuleName)
}
