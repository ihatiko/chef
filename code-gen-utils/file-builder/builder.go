package utils

import (
	"bytes"
	"embed"
	"errors"
	"io/fs"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Builder struct {
	Prefix    string
	structure embed.FS
}

func NewBuilder(prefix string, structure embed.FS) *Builder {
	return &Builder{Prefix: prefix, structure: structure}
}
func (b *Builder) Build(prefix, destination string, obj any) {
	b.process(prefix, destination, obj)
}
func (b *Builder) process(prefix, destination string, obj any) {
	dir, err := b.structure.ReadDir(prefix)
	if err != nil {
		log.Fatal(err)
	}
	Mkdir(destination)
	for _, fl := range dir {
		secondPath := filepath.Join(prefix, fl.Name())
		if fl.IsDir() {
			directory := b.cleanPath(secondPath)
			newDirectory := b.RewritePath(directory, obj)
			if strings.HasSuffix(newDirectory, b.OSSlash("/")) {
				continue
			}
			newDirectoryPath := filepath.Join(destination, newDirectory)
			Mkdir(newDirectoryPath)

			if b.checkWindows() {
				secondPath = filepath.ToSlash(secondPath)
			}
			b.process(secondPath, destination, obj)
			b.CleanEmptyDir(newDirectoryPath)
			continue
		}

		b.buildFile(secondPath, destination, obj)
	}
}

func (b *Builder) CleanEmptyDir(destination string) {
	dir, err := os.ReadDir(destination)
	if err != nil {
		slog.Error(err.Error())
		panic("failed to clean empty directory")
	}
	if len(dir) == 0 {
		err = os.Remove(destination)
		if err != nil {
			slog.Error(err.Error())
			panic("failed to clean empty directory")
		}
	}
}
func (b *Builder) checkWindows() bool {
	return os.Getenv("GOOS") == "windows" ||
		strings.Contains(strings.ToLower(os.Getenv("OS")), "windows")
}

func (b *Builder) cleanPath(path string) string {
	path = b.OSSlash(path)
	result, _ := strings.CutPrefix(path, b.Prefix)
	return result
}

func (b *Builder) OSSlash(path string) string {
	if b.checkWindows() {
		path = filepath.ToSlash(path)
	}
	return path
}

func (b *Builder) buildFile(secondPath string, folder string, obj any) {
	secondPath = b.OSSlash(secondPath)
	rF, err := b.structure.ReadFile(secondPath)
	if err != nil {
		slog.Error(err.Error())
		panic("failed to read file")
	}
	t, err := template.New("").Parse(string(rF))
	if err != nil {
		slog.Error(err.Error())
		panic("failed to parse template")
	}
	parsedPath, _ := strings.CutSuffix(secondPath, ".tmpl")
	parsedPath = b.cleanPath(parsedPath)

	filePath := filepath.Join(
		folder,
		b.RewritePath(parsedPath, obj),
	)
	f, err := os.Create(filePath)
	if err != nil {
		slog.Error(err.Error())
		panic("failed to create file")
	}
	err = t.ExecuteTemplate(f, "", obj)
	if err != nil {
		slog.Error(err.Error())
		panic("failed to execute template")
	}
}
func (b *Builder) RewritePath(folder string, obj any) string {
	t, err := template.New("").Parse(folder)
	if err != nil {
		slog.Error(err.Error())
		panic("failed to parse template")
	}
	buffer := bytes.NewBufferString("")
	err = t.ExecuteTemplate(buffer, "", obj)
	if err != nil {
		slog.Error(err.Error())
		panic("failed to execute template")
	}
	return strings.ReplaceAll(buffer.String(), " ", "")
}

func Mkdir(path string) {
	var err error
	if path == "" {
		return
	}
	_, err = os.ReadDir(path)
	if errors.Is(err, fs.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			slog.Error(err.Error())
			panic("failed to create directory")
		}
	}
}
