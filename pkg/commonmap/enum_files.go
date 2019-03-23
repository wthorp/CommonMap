package commonmap

import (
	"fmt"
	"math"
	"reflect"
	"syscall"
	"unsafe"
)

// Build map of folder names using MFT (based on PopulateMethod2)
func EnumFiles(basePath string, isRpf func(string) bool) (map[DWORDLONG]folderEntry, []folderEntry) {
	//build map of folders
	folders := make(map[DWORDLONG]folderEntry)
	//build list of files
	files := make([]folderEntry, 0, 5000)

	volumeHandle, _ := open(basePath, syscall.O_RDONLY, FILE_ATTRIBUTE_NORMAL)
	//fmt.Println("dir,err", volumeHandle, err)

	med := MFT_ENUM_DATA{0, 0, math.MaxInt64}

	for {
		data, done, err := enumUsnData(volumeHandle, &med)
		if err != nil {
			if err.Error() != "Reached the end of the file." {
				fmt.Println(err)
			}
		}
		if done == 0 {
			return folders, files
		}

		var usn USN = *(*USN)(unsafe.Pointer(&data[0]))
		var ur *USN_RECORD
		for i := unsafe.Sizeof(usn); i < uintptr(done); i += uintptr(ur.RecordLength) {
			ur = (*USN_RECORD)(unsafe.Pointer(&data[i]))
			if ur.FileAttributes&FILE_ATTRIBUTE_DIRECTORY != 0 {
				nameLength := uintptr(ur.FileNameLength) / unsafe.Sizeof(ur.FileName[0])
				fnp := unsafe.Pointer(&data[i+uintptr(ur.FileNameOffset)])
				fnUtf := (*[10000]uint16)(fnp)[:nameLength]
				fn := syscall.UTF16ToString(fnUtf)
				(*reflect.SliceHeader)(unsafe.Pointer(&fn)).Cap = int(nameLength)
				// fmt.Println("len", ur.FileNameLength, ur.FileNameOffset, "fn", fn)
				folders[ur.FileReferenceNumber] = folderEntry{fn, ur.ParentFileReferenceNumber}
			} else { // file not folder
				nameLength := uintptr(ur.FileNameLength) / unsafe.Sizeof(ur.FileName[0])
				if nameLength == 0 {
					fmt.Println("zero length file name!!!")
				}
				fnp := unsafe.Pointer(&data[i+uintptr(ur.FileNameOffset)])
				fnUtf := (*[10000]uint16)(fnp)[:nameLength]
				fn := syscall.UTF16ToString(fnUtf)
				if isRpf(fn) {
					files = append(files, folderEntry{fn, ur.ParentFileReferenceNumber})
				}
			}
		}
		med.StartFileReferenceNumber = DWORDLONG(usn)
	}
}
