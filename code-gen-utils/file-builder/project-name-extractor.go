package utils

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"os"
	"path/filepath"
	"strings"
)

func GetPath(destination string) string {
	if destination == "" {
		d, err := os.Getwd()
		if err != nil {
			otelzap.S().Fatal(err)
		}
		return d
	}
	return destination
}

func GetProjectName(projectPath string) string {
	projectPath = GetPath(projectPath)
	goModPath := filepath.Join(projectPath, "go.mod")
	f, err := os.ReadFile(goModPath)
	if err != nil {
		otelzap.S().Fatalf("cannot read go.mod in folder %s %v", projectPath, err)
	}
	splittedFile := strings.Split(string(f), "\n")
	if len(splittedFile) == 0 {
		otelzap.S().Fatalf("empty go.mod file in folder %s", projectPath)
	}
	return strings.Replace(splittedFile[0], "module ", "", 1)
}
