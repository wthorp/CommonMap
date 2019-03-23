package rpf

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"testing"
)

const EPSILON = 0.0000001

func testEquals(a, b float64, file, name string, t *testing.T) {
	ok := (a-b) < EPSILON && (b-a) < EPSILON
	if !ok {
		t.Errorf("'%s' computed %s (%f) doesn't match expected (%f)", file, name, a, b)
	}
}

/*
func TestTryGetRpfBounds2(t *testing T){
	driver, err := gdal.GetDriverByName("GTiff")
	dataset := driver.Create(filename, 256, 256, 1, gdal.Byte, nil)
	defer dataset.Close()
	geoTransform := dataset.GeoTransform()
	fmt.Printf("%v, %v, %v, %v, %v, %v\n",
		geoTransform[0], geoTransform[1], geoTransform[2], geoTransform[3], geoTransform[4], geoTransform[5])

}
*/

func TestTryGetRpfBounds(t *testing.T) {
	file, err := os.Open("IsRPF_test.data")
	defer file.Close()
	if err != nil {
		t.Error("Cannot find source data file")
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) != 9 {
			t.Errorf("Error parsing RPF file line :\n  %s", line)
			continue
		}
		filePath := fields[0]
		ox1, err1 := strconv.ParseFloat(fields[1], 64)
		ox2, err2 := strconv.ParseFloat(fields[5], 64)
		oy1, err3 := strconv.ParseFloat(fields[4], 64)
		oy2, err4 := strconv.ParseFloat(fields[2], 64)

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			t.Errorf("Error parsing floats from line:\n  %s", line)
			continue
		}
		ok, x1, y1, x2, y2 := TryGetRpfBounds(filePath)
		if !ok {
			t.Error("Valid RPF file name '" + filePath + "' marked as invalid")
			continue
		}
		testEquals(ox1, x1, filePath, "x1", t)
		testEquals(ox2, x2, filePath, "x2", t)
		testEquals(oy1, y1, filePath, "y1", t)
		testEquals(oy2, y2, filePath, "y2", t)
	}

	if err := scanner.Err(); err != nil {
		t.Error(err)
	}
}
