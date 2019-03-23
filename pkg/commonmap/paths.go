package commonmap

import (
	"log"
	"os"
	"path/filepath"

	"github.com/kardianos/osext"
)

var binDir string
var mapservPath string
var MapfilePath string
var proj4Path string
var IndexPath string
var contentPath string
var appDir string
var fontsetPath string
var vectorTemplatePath string

func init() {
	filename, _ := osext.Executable()
	binDir = filepath.Dir(filename)

	mapservPath = GetPathFatal("bin", "mapserv.exe")
	proj4Path = GetPathFatal("bin", "nad")
	IndexPath = GetPathFatal("content", "index")
	contentPath = GetPathFatal("content")
	MapfilePath = GetPath("content", "common.map")
	fontsetPath = GetPathFatal("content", "fonts", "fontset.txt")
	appDir = GetPathFatal("content", "website")
	vectorTemplatePath = GetPathFatal("content", "vector", "Natural_Earth", "template.map")
}

func GetPath(elem ...string) string {
	elem = append([]string{binDir}, elem...)
	return filepath.Join(elem...)
}

func GetPathFatal(elem ...string) string {
	path := GetPath(elem...)
	if _, err := os.Stat(path); err != nil {
		log.Fatal("Missing required file at " + path)
	}
	return path
}

func GetIndexPath(fileName string) string {
	filePath := filepath.Join(IndexPath, fileName)
	return filePath
}
