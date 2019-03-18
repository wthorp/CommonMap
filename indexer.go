package main

import (
	"RPF"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type RpfBox struct {
	path string
	box  Box
}

var shapeFiles map[string]*ShpBoxWriter

// Assemble the path by looking up parent pointers in the folders map
func GetFullPath2(folders map[DWORDLONG]folderEntry, cache map[DWORDLONG]string, f folderEntry) string {
	pFrn := f.parent
	if pFrn == 0 { //check if we're at the FS root
		//		fmt.Println("FRN=0 : " + f.name)
		return f.name
	}
	if res := cache[pFrn]; res != "" {
		return res + "/" + f.name //check cached names
	}
	pPath := GetFullPath2(folders, cache, folders[pFrn]) //derive path name
	cache[pFrn] = pPath                                  //add parent to cache
	return pPath + "/" + f.name
}

func CheckPathAsNtfsDrive(path string) bool {
	//todo:  stuff
	return true
}

func Index(indexPath string) {

	t0 := time.Now()
	forShp := make(chan RpfBox) // todo: benchmark w/ pointers
	forDbf := make(chan string)
	done := make(chan bool)
	go addToShp(forShp, forDbf, done)

	totalFiles := 0

	if runtime.GOOS == "windows" { //todo:  check for NTFS drive
		//create list of files & folders, while also generating shapefiles
		folders, files := EnumFiles("\\\\.\\"+indexPath, func(rpfPath string) bool {
			totalFiles++
			isRpf, x1, y1, x2, y2 := RPF.TryGetRpfBounds(rpfPath)
			if isRpf {
				rpf := RpfBox{rpfPath, [4]float64{x1, y1, x2, y2}}
				forShp <- rpf //start building SHP / SHX / QIX now
			}
			return isRpf
		})
		close(forShp)
		fmt.Printf("Initial phase complete in %v.\n", time.Now().Sub(t0))
		<-done

		rFolders := make(map[DWORDLONG]string)
		for k, v := range folders {
			fPath := filepath.Join(indexPath, GetFullPath2(folders, rFolders, v))
			rFolders[k] = fPath
		}

		//use files and folders to generate file paths and DBFs
		emptyCount := 0
		for _, rpfData := range files {
			if rpfData.name == "" {
				fmt.Printf("File name is empty string")
				emptyCount++
				continue
			}
			pathx := rFolders[rpfData.parent] + rpfData.name
			forDbf <- pathx
		}
		close(forDbf)
	} else {
		err := filepath.Walk(indexPath, func(filepath string, f os.FileInfo, err error) error {
			totalFiles++
			_, fileName := path.Split(filepath)
			isRpf, x1, y1, x2, y2 := RPF.TryGetRpfBounds(fileName)
			if isRpf {
				rpf := RpfBox{fileName, [4]float64{x1, y1, x2, y2}}
				forShp <- rpf      // start building SHP / SHX / QIX now
				forDbf <- filepath // also generate DBF
			}
			return nil
		})
		close(forShp)
		close(forDbf)
		if err != nil {
			fmt.Println("file search error: ", err)
		}
	}
	<-done
	fmt.Printf("The call took %v to scan %d files.\n", time.Now().Sub(t0), totalFiles)
}

// add found rpf to shapefile
func addToShp(forShp chan RpfBox, forDbf chan string, done chan bool) {
	//todo:  unhardcode, possibly remove existing
	shapeFiles = make(map[string]*ShpBoxWriter)
	var shape *ShpBoxWriter
	var r RpfBox
	for r = range forShp {
		shape = getShapeFile(r.path)
		shape.WriteBox(r)
	}
	done <- true
	for shpPath := range forDbf {
		shape = getShapeFile(shpPath)
		shape.WriteDbf(shpPath)
	}
	for _, shape := range shapeFiles {
		shape.Close()
	}
	done <- true
}

//get a different shapefile for each file type
func getShapeFile(shpPath string) *ShpBoxWriter {
	seriesCode := strings.ToUpper(filepath.Ext(shpPath)[1:3])
	shape := shapeFiles[seriesCode]
	if shape == nil {
		var err error
		shape, err = Create(seriesCode)
		if err != nil {
			fmt.Println(err)
		}
		shapeFiles[seriesCode] = shape
	}
	return shape
}
