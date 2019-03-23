package rpf

import (
	"fmt"
	"math"
)

// This is an implementation of the Equal Arc-Second Raster Chart/Map System
// See MIL-A-89007

type ArcZone struct {
	Equatorward, Poleward float64
	AParam                float64 // East-West pixel spacing at scale 1:S in zone Z
	IsPolar               bool
}

var ArcZones = map[byte]ArcZone{
	'1': {0, 32, 369664, false},
	'2': {32, 48, 302592, false},
	'3': {48, 56, 245760, false},
	'4': {56, 64, 199168, false},
	'5': {64, 68, 163328, false},
	'6': {68, 72, 137216, false},
	'7': {72, 76, 110080, false},
	'8': {76, 80, 82432, false},
	'9': {80, 90, 400384, true},
	'A': {0, -32, 369664, false},
	'B': {-32, -48, 302592, false},
	'C': {-48, -56, 245760, false},
	'D': {-56, -64, 199168, false},
	'E': {-64, -68, 163328, false},
	'F': {-68, -72, 137216, false},
	'G': {-72, -76, 110080, false},
	'H': {-76, -80, 82432, false},
	'J': {-80, -90, 400384, true},
}

const pixelsPerFrame = 1536 // 256x256 pixel subframes, stacked 6x6

func GetBounds(frame *FrameInfo) (x1, y1, x2, y2 float64) {
	zone := ArcZones[frame.ArcZone]
	series := DataSeries[frame.SeriesCode]
	isCADRG := series.Type == CADRG
	scale := series.Scale
	_, cols := CalculateNumRowsCols(frame.ArcZone, scale, isCADRG)
	latDpp, lonDpp := CalculateDegreesPerPixel(frame.ArcZone, scale, isCADRG)

	//calc row and column for this frame
	row := frame.FrameNumber / cols
	column := frame.FrameNumber - row*cols

	x1 = float64(column)*lonDpp*pixelsPerFrame - 180.0
	y1 = math.Min(zone.Equatorward, zone.Poleward) + float64(row)*latDpp*pixelsPerFrame
	x2 = x1 + lonDpp*pixelsPerFrame
	y2 = y1 + latDpp*pixelsPerFrame

	if math.IsInf(x1, 0) || math.IsInf(y1, 0) || math.IsInf(x2, 0) || math.IsInf(y2, 0) {
		fmt.Println("RPF calc error - something is Infinite")
	}
	if math.IsNaN(x1) || math.IsNaN(y1) || math.IsNaN(x2) || math.IsNaN(y2) {
		fmt.Println("RPF calc error - something is NaN")
	}
	if x1 > x2 {
		x1, x2 = x2, x1
	}
	if y1 > y2 {
		y1, y2 = y2, y1
	}
	return x1, y1, x2, y2
}

// Computes the geographic width and height of a frame in the given zone for the given map type
func CalculateGeoFrameWidthHeight(zone byte, scale float64, isCADRG bool) (geoFrameWidth float64, geoFrameHeight float64) {
	latDpp, lonDpp := CalculateDegreesPerPixel(zone, scale, isCADRG)
	return lonDpp * pixelsPerFrame, latDpp * pixelsPerFrame
}

// Computes the number of degrees per pixel - only called by CalculateGeoFrameWidthHeight
func CalculateDegreesPerPixel(zone byte, scale float64, isCADRG bool) (float64, float64) {
	if ArcZones[zone].IsPolar {
		dppLatLon := 360.0 / calcPolarPixConst(scale)
		return dppLatLon, dppLatLon
	} // else...
	dppLat := 90.0 / CalcNsPixConst(scale, isCADRG)
	dppLon := 360.0 / CalcEwPixConst(zone, scale, isCADRG)
	return dppLat, dppLon
}

// Computes the number of rows and columns of frames in the given zone for the given map scale
func CalculateNumRowsCols(zone byte, scale float64, isCADRG bool) (int, int) {
	if ArcZones[zone].IsPolar {
		numFrames := int(math.Ceil(calcPolarPixConst(scale) / 18.0 / pixelsPerFrame))
		if numFrames%2 == 0 { // round up to the next odd number of frames
			numFrames++
		}
		return numFrames, numFrames
	} // else...
	rows := int(math.Ceil(calcLatFrameRows(zone, scale, isCADRG)))
	cols := int(math.Ceil(CalcEwPixConst(zone, scale, isCADRG) / pixelsPerFrame))
	return rows, cols
}

func calcPolarPixConst(scale float64) float64 {
	const BParam = 400384.0
	AdrgPixConst := roundUp(BParam*1000000.0/scale, 512)
	return roundUp(AdrgPixConst/27.0, 512) * 18.0
}

func CalcNsPixConst(scale float64, isCADRG bool) float64 {
	const BParam = 400384.0 // fixed value
	if isCADRG {            // 1:N scale
		AdrgPixConst := roundUp(BParam*1000000.0/scale, 512)
		return roundUp(AdrgPixConst/6.0, 256)
	}
	// is CIB/DTED meters scale
	AdrgPixConst := roundUp(BParam*100.0/scale, 512)
	return roundUp(AdrgPixConst/4.0, 256)
}

func CalcEwPixConst(zone byte, scale float64, isCADRG bool) float64 {
	AParam := float64(ArcZones[zone].AParam)
	if isCADRG { // 1:N scale
		AdrgPixConst := roundUp(AParam*1000000.0/scale, 512)
		return roundUp(AdrgPixConst/1.5, 256)
	}
	// is CIB/DTED meters scale
	AdrgPixConst := roundUp(AParam*100.0/scale, 512)
	return roundUp(AdrgPixConst, 256)
}

// determine the actual equatorward extent of a zone as described in the CADRG standard
func GetEquatorwardExtent(zone byte, scale float64, isCADRG bool) float64 {
	dLatPixConst := CalcNsPixConst(scale, isCADRG)

	// determine the number of frames needed to reach the nominal zone boundary.
	absNominalEquatorwardZoneBoundary := math.Abs(ArcZones[zone].Equatorward)
	pixelsPerDegree := dLatPixConst / 90.0
	numberOfFrames := math.Trunc(pixelsPerDegree * absNominalEquatorwardZoneBoundary * (1.0 / pixelsPerFrame))

	// use the number of frames to get the poleward zone extent
	equatorwardZoneLat := numberOfFrames * pixelsPerFrame * (1.0 / pixelsPerDegree)

	if zone < 0 {
		equatorwardZoneLat = -equatorwardZoneLat
	}

	return equatorwardZoneLat
}

// determine the actual poleward extent of a zone as described in the CADRG standard
func GetPolewardExtent(zone byte, scale float64, isCADRG bool) float64 {
	dLatPixConst := CalcNsPixConst(scale, isCADRG)

	// determine the number of frames needed to reach the nominal zone boundary.
	absNominalPolewardZoneBoundary := math.Abs(ArcZones[zone].Poleward)
	pixelsPerDegree := dLatPixConst / 90.0
	numberOfFrames := math.Ceil(pixelsPerDegree * absNominalPolewardZoneBoundary * (1.0 / pixelsPerFrame))

	// use the number of frame to get the poleward zone extent =
	polewardZoneLat := numberOfFrames * pixelsPerFrame * (1.0 / pixelsPerDegree)

	if zone < 0 {
		polewardZoneLat = -polewardZoneLat
	}

	return polewardZoneLat
}

func calcLatFrameRows(zone byte, scale float64, isCADRG bool) float64 {
	// check if we are at polar zones, and constrain it
	// to the non-polar zones if necessary
	arcZone := ArcZones[zone]

	// calculate the number of pixels per degree
	PixelsPerDegreeLat := CalcNsPixConst(scale, isCADRG) / 90.0

	// calculate the number of frames needed to reach each of the nominal
	// zone boundaries.  This number is the pixels per degree lat multiplied
	// by the nominal zone boundary (in degrees), divided by 1536 (the number of
	// pixel rows in a frame)
	numFramesPoleward := PixelsPerDegreeLat * arcZone.Poleward / 1536.0
	numFramesEquatorward := PixelsPerDegreeLat * arcZone.Equatorward / 1536.0

	// the exact poleward zone extent is calculated by multiplying the number of frames
	// (rounded up) by 1536 and dividing by the number of pixels in a degree of latitude
	exactPolewardZoneExtent := math.Ceil(numFramesPoleward) * 1536.0 / PixelsPerDegreeLat

	// equatorward zone extent is calculated the same way, except we round the number
	// of frames down rather than up
	exactEquatorwardZoneExtent := math.Trunc(numFramesEquatorward) * 1536.0 / PixelsPerDegreeLat

	// The number of latitudinal frames is the difference (in degrees)
	// between the exact poleward extent and exact equatorward zone extent,
	// multiplied by the number of pixels per degree, and divided by 1536 (the
	// number of pixel rows per frame).
	return (exactPolewardZoneExtent - exactEquatorwardZoneExtent) * PixelsPerDegreeLat / 1536.0
}

// Round a float up to the nearest multiple of an integer factor
func roundUp(number float64, factor int) float64 {
	ceiling := int(math.Ceil(number))
	return float64(ceiling + ((factor - (ceiling % factor)) % factor))
}
