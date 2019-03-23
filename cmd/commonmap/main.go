package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"cm/pkg/commonmap"
)

func main() {
	var useVector, doServe, genMap bool
	var IndexPath string

	flag.StringVar(&IndexPath, "index", "", "drive letter to index")
	flag.BoolVar(&useVector, "vector", true, "include vector base map")
	flag.BoolVar(&doServe, "serve", true, "start web server")
	flag.BoolVar(&genMap, "map", false, "regenerate map without reindexing")
	flag.Parse()

	if IndexPath == "" && !doServe {
		flag.PrintDefaults()
		return
	}

	if IndexPath != "" {
		fmt.Println("Indexing " + IndexPath)
		cleanIndexPath()
		commonmap.Index(IndexPath)
		genMap = true
	}

	if genMap {
		//shapeFiles fully populated... build map file
		mapfile, err := os.Create(commonmap.MapfilePath)
		if err != nil {
			fmt.Printf("CANNOT CREATE MAP FILE\n%s\n", err.Error())
		} else {
			fmt.Printf("Created map file %s\n", commonmap.MapfilePath)
		}
		defer mapfile.Close()
		commonmap.WriteMap(mapfile)
	}

	if doServe {
		commonmap.Serve()
	}
}

func cleanIndexPath() {
	//todo: error checking
	filepath.Walk(commonmap.IndexPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			os.Remove(path)
		}
		return nil
	})
}
