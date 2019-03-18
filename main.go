package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var useVector, doServe, genMap bool
	var indexPath string

	flag.StringVar(&indexPath, "index", "", "drive letter to index")
	flag.BoolVar(&useVector, "vector", true, "include vector base map")
	flag.BoolVar(&doServe, "serve", true, "start web server")
	flag.BoolVar(&genMap, "map", false, "regenerate map without reindexing")
	flag.Parse()

	if indexPath == "" && !doServe {
		flag.PrintDefaults()
		return
	}

	if indexPath != "" {
		fmt.Println("Indexing " + indexPath)
		cleanIndexPath()
		Index(indexPath)
		genMap = true
	}

	if genMap {
		//shapeFiles fully populated... build map file
		mapfile, err := os.Create(mapfilePath)
		if err == nil {
			fmt.Errorf("CANNOT CREATE MAP FILE")
		}
		defer mapfile.Close()
		WriteMap(mapfile)
	}

	if doServe {
		serve()
	}
}

func cleanIndexPath() {
	//todo: error checking
	filepath.Walk(indexPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			os.Remove(path)
		}
		return nil
	})
}
