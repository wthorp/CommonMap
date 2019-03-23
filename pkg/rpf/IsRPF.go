package rpf

import (
	"path"
	"sort"
	"strings"
)

// from MIL-STD-2411-1 w/CHANGE 3 : 5.1.3
var zoneCodes = "123456789ABCDEFGHJ"

// from MIL-STD-2411-1 w/CHANGE 3 : 5.1.3.1
var zoneCodesDTED = "123456789A"

// from MIL-STD-2411-1 w/CHANGE 3 : 5.1.4
var seriesCodes = []string{"A1", "A2", "A3", "A4", "AT", "C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "CA",
	"CB", "CC", "CD", "CE", "CF", "CG", "CH", "CJ", "CK", "CL", "CM", "CN", "CO", "CP", "CQ", "CR", "CS", "CT", "D1",
	"D2", "EG", "ES", "ET", "F1", "F2", "F3", "F4", "F5", "GN", "HA", "I1", "I2", "I3", "I4", "I5", "IV", "JA", "JG",
	"JN", "JO", "JR", "K1", "K2", "K3", "K7", "K8", "KB", "KE", "KM", "KR", "KS", "KU", "L1", "L2", "L3", "L4", "L5",
	"LF", "LN", "M1", "M2", "MH", "MI", "MJ", "MM", "OA", "OH", "ON", "OW", "P1", "P2", "P3", "P4", "P5", "P6", "P7",
	"P8", "P9", "PA", "PB", "PC", "PD", "PE", "PF", "PG", "PH", "PI", "PJ", "PK", "PL", "PM", "PN", "PO", "PP", "PQ",
	"PR", "PS", "PT", "PU", "PV", "R1", "R2", "R3", "R4", "R5", "RC", "RL", "RR", "RV", "TC", "TF", "TL", "TN", "TP",
	"TQ", "TR", "TT", "UL", "V1", "V2", "V3", "V4", "VH", "VN", "VT", "WA", "WB", "WC", "WD", "WE", "WF", "WG", "WH",
	"WI", "WK", "XD", "XE", "XF", "XG", "XH", "XI", "XJ", "XK", "Y9", "YA", "YB", "YC", "YD", "YE", "YF", "YI", "YJ",
	"YZ", "Z8", "ZA", "ZB", "ZC", "ZD", "ZE", "ZF", "ZG", "ZH", "ZI", "ZJ", "ZK", "ZT", "ZV", "ZZ"}

// valid RPF extensions are two character series codes, followed by one character zone code
func isRpfExtension(extension string) bool {
	if len(extension) != 4 || extension[0] != '.' {
		return false
	}
	zone := extension[3]
	if strings.IndexByte(zoneCodes, zone) == -1 {
		return false
	}
	series := extension[1:3]
	if i := sort.SearchStrings(seriesCodes, series); i == len(seriesCodes) || seriesCodes[i] != series {
		return false
	}
	if series == "D1" || series == "D2" && strings.IndexByte(zoneCodesDTED, zone) == -1 {
		return false
	}
	return true
}

// takes a file name (excluding path) and returns
func TryGetRpfBounds(fileName string) (isValid bool, x1, y1, x2, y2 float64) {
	fileName = strings.ToUpper(fileName)
	if !isRpfExtension(path.Ext(fileName)) {
		return false, 0, 0, 0, 0
	}
	frame := NewFrameInfo(fileName)
	if frame == nil {
		return false, 0, 0, 0, 0
	}
	x1, y1, x2, y2 = GetBounds(frame)
	//todo:  establish a firmer test for X bounds
	if x1 < -181 || x2 > 181 || y1 < -90 || y2 > 90 {
		//fmt.Printf("Bad RPF bounds : %s {%f, %f, %f %f}\n", fileName, x1, y1, x2, y2)
		return false, 0, 0, 0, 0
	}
	return true, x1, y1, x2, y2
}
