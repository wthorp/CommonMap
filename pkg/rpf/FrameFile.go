package rpf

import (
	"path"
	"strings"
)

type FrameInfo struct {
	FrameNumber int
	Edition     int
	Producer    int
	SeriesCode  string
	ArcZone     byte
}

func NewFrameInfo(filePath string) *FrameInfo {
	_, fileName := path.Split(filePath)
	i := strings.LastIndex(fileName, ".")
	file, ext := filePath[:i], filePath[i+1:]
	if len(ext) != 3 || len(file) != 8 {
		return nil
	}
	seriesCode := ext[0:2]
	dataSeries, validDataSeries := DataSeries[seriesCode]
	if !validDataSeries {
		return nil
	}
	arcZone := ext[2]
	_, validArcZone := ArcZones[arcZone]
	if !validArcZone {
		return nil
	}

	var frameNumber int
	var edition, producer int
	if dataSeries.Type == CIB { //  MIL-PRF-89041 - ffffffvp.ccz
		frameNumber = DecodeBase34(fileName[0:6])
		edition = DecodeBase34(fileName[6:7])
		producer = DecodeBase34(fileName[8:8])
	} else { //  MIL-PRF 89038 - fffffvvp.ccz format
		frameNumber = DecodeBase34(fileName[0:5])
		edition = DecodeBase34(fileName[5:7])
		producer = DecodeBase34(fileName[8:8])
	}
	if frameNumber == -1 || edition == -1 || producer == -1 {
		return nil
	}
	f := FrameInfo{frameNumber, edition, producer, seriesCode, arcZone}
	return &f
}

func DecodeBase34(input string) int {
	i := 0
	//todo:  replace bytes with runes, loop with range?
	for digit := 0; digit < len(input); digit++ {
		var index int
		char := input[digit]
		switch {
		case char >= '0' && char <= '9':
			index = int(char - '0')
		case char >= 'A' && char <= 'H':
			index = int(10 + char - 'A')
		case char >= 'J' && char <= 'N':
			index = int(18 + char - 'J')
		case char >= 'P' && char <= 'Z':
			index = int(23 + char - 'P')
		default:
			return -1
		}
		i = (i * 34) + index
	}
	return i
}
