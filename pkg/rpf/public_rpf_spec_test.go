package rpf

import (
	"math"
	"strings"
	"testing"
)

const publicSpecEpsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= publicSpecEpsilon
}

func TestDecodeBase34PublicAlphabet(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{name: "zero", input: "0", want: 0},
		{name: "nine", input: "9", want: 9},
		{name: "a", input: "A", want: 10},
		{name: "h", input: "H", want: 17},
		{name: "j", input: "J", want: 18},
		{name: "z", input: "Z", want: 33},
		{name: "base rollover", input: "10", want: 34},
		{name: "double max", input: "ZZ", want: 34*33 + 33},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeBase34(tt.input); got != tt.want {
				t.Fatalf("DecodeBase34(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestArcZonesPublicZoneSet(t *testing.T) {
	validZones := []byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J'}
	for _, zone := range validZones {
		arcZone, ok := ArcZones[zone]
		if !ok {
			t.Fatalf("ArcZones is missing documented zone %q", zone)
		}
		if arcZone.AParam <= 0 {
			t.Fatalf("ArcZones[%q].AParam = %v, want positive", zone, arcZone.AParam)
		}
		if arcZone.Equatorward < -90 || arcZone.Equatorward > 90 {
			t.Fatalf("ArcZones[%q].Equatorward = %v, want latitude range", zone, arcZone.Equatorward)
		}
		if arcZone.Poleward < -90 || arcZone.Poleward > 90 {
			t.Fatalf("ArcZones[%q].Poleward = %v, want latitude range", zone, arcZone.Poleward)
		}
	}

	for _, zone := range []byte{'I', 'O'} {
		if _, ok := ArcZones[zone]; ok {
			t.Fatalf("ArcZones unexpectedly contains excluded zone %q", zone)
		}
	}
}

func TestDataSeriesPublicFamilies(t *testing.T) {
	tests := []struct {
		code      string
		wantType  Type
		wantScale float64
	}{
		{code: "A1", wantType: CADRG, wantScale: 10000},
		{code: "A2", wantType: CADRG, wantScale: 25000},
		{code: "I1", wantType: CIB, wantScale: 10},
		{code: "I2", wantType: CIB, wantScale: 5},
		{code: "I3", wantType: CIB, wantScale: 2},
		{code: "I4", wantType: CIB, wantScale: 1},
		{code: "I5", wantType: CIB, wantScale: 0.5},
		{code: "D1", wantType: CDTED, wantScale: 100},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			series, ok := DataSeries[tt.code]
			if !ok {
				t.Fatalf("DataSeries is missing documented series %q", tt.code)
			}
			if series.Type != tt.wantType {
				t.Fatalf("DataSeries[%q].Type = %v, want %v", tt.code, series.Type, tt.wantType)
			}
			if series.Scale != tt.wantScale {
				t.Fatalf("DataSeries[%q].Scale = %v, want %v", tt.code, series.Scale, tt.wantScale)
			}
			if series.SeriesCode != tt.code {
				t.Fatalf("DataSeries[%q].SeriesCode = %q, want %q", tt.code, series.SeriesCode, tt.code)
			}
			if strings.TrimSpace(series.ScaleText) == "" {
				t.Fatalf("DataSeries[%q].ScaleText is empty", tt.code)
			}
			if strings.TrimSpace(series.Name) == "" {
				t.Fatalf("DataSeries[%q].Name is empty", tt.code)
			}
		})
	}
}

func TestCADRGFrameGeometryUses1536Pixels(t *testing.T) {
	tests := []struct {
		name       string
		seriesCode string
		zone       byte
	}{
		{name: "combat chart zone 1", seriesCode: "A1", zone: '1'},
		{name: "combat chart zone 5", seriesCode: "A1", zone: '5'},
		{name: "combat chart polar zone", seriesCode: "A1", zone: '9'},
		{name: "25k chart zone 1", seriesCode: "A2", zone: '1'},
		{name: "25k chart polar zone", seriesCode: "A2", zone: '9'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			series := DataSeries[tt.seriesCode]
			width, height := CalculateGeoFrameWidthHeight(tt.zone, series.Scale, true)
			degPerPixelX, degPerPixelY := CalculateDegreesPerPixel(tt.zone, series.Scale, true)

			if degPerPixelX <= 0 || degPerPixelY <= 0 {
				t.Fatalf("CalculateDegreesPerPixel(%q, %v, true) = (%v, %v), want positive values", tt.zone, series.Scale, degPerPixelX, degPerPixelY)
			}

			directMatch := almostEqual(width/degPerPixelX, 1536) && almostEqual(height/degPerPixelY, 1536)
			swappedMatch := almostEqual(width/degPerPixelY, 1536) && almostEqual(height/degPerPixelX, 1536)
			if !directMatch && !swappedMatch {
				t.Fatalf(
					"frame geometry for %s did not match 1536x1536 pixels: width=%v height=%v degPerPixel=(%v,%v)",
					tt.name,
					width,
					height,
					degPerPixelX,
					degPerPixelY,
				)
			}
		})
	}
}

func TestTryGetRpfBoundsMatchesGetBoundsForValidFrames(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		series   string
		zone     byte
		edition  int
	}{
		{name: "cib 1m frame", fileName: "0REF5K4A.I41", series: "I4", zone: '1', edition: 4},
		{name: "cib 1m adjacent frame", fileName: "0RF1ZE2A.I41", series: "I4", zone: '1', edition: 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := NewFrameInfo(tt.fileName)
			if frame == nil {
				t.Fatalf("NewFrameInfo(%q) returned nil", tt.fileName)
			}
			if frame.SeriesCode != tt.series {
				t.Fatalf("NewFrameInfo(%q).SeriesCode = %q, want %q", tt.fileName, frame.SeriesCode, tt.series)
			}
			if frame.ArcZone != tt.zone {
				t.Fatalf("NewFrameInfo(%q).ArcZone = %q, want %q", tt.fileName, frame.ArcZone, tt.zone)
			}
			if frame.Edition != tt.edition {
				t.Fatalf("NewFrameInfo(%q).Edition = %d, want %d", tt.fileName, frame.Edition, tt.edition)
			}

			ok, x1, y1, x2, y2 := TryGetRpfBounds(tt.fileName)
			if !ok {
				t.Fatalf("TryGetRpfBounds(%q) reported invalid", tt.fileName)
			}

			gx1, gy1, gx2, gy2 := GetBounds(frame)
			if !almostEqual(x1, gx1) || !almostEqual(y1, gy1) || !almostEqual(x2, gx2) || !almostEqual(y2, gy2) {
				t.Fatalf("TryGetRpfBounds(%q) = (%v, %v, %v, %v), GetBounds(NewFrameInfo(...)) = (%v, %v, %v, %v)", tt.fileName, x1, y1, x2, y2, gx1, gy1, gx2, gy2)
			}
		})
	}
}

func TestTryGetRpfBoundsRejectsMalformedNames(t *testing.T) {
	tests := []string{
		"",
		"INVALID",
		"0REF5K4A",
		"0REF5K4A.I4",
		"0REF5K4A.I411",
		"0REI5K4A.I41",
		"123456789.1234",
		"0REF5K4A.ZZ1",
	}

	for _, fileName := range tests {
		t.Run(fileName, func(t *testing.T) {
			ok, _, _, _, _ := TryGetRpfBounds(fileName)
			if ok {
				t.Fatalf("TryGetRpfBounds(%q) reported valid for malformed name", fileName)
			}
		})
	}
}
