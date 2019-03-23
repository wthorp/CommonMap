package commonmap

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"cm/pkg/rpf"
)

type SeriesRes struct {
	seriesCode string
	bestName   string
	scale      float64
}
type AllSeries []SeriesRes

func (slice AllSeries) Len() int {
	return len(slice)
}
func (slice AllSeries) Less(i, j int) bool {
	return slice[i].scale < slice[j].scale
}
func (slice AllSeries) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func EscapeSlashes(input string) string {
	return strings.Replace(input, "\\", "\\\\", -1)
}

func WriteMap(w io.Writer) {

	WriteHeader(w, proj4Path, IndexPath, "http://localhost:7070/wms")
	WriteVector(w, vectorTemplatePath, contentPath)
	allSeries := make(AllSeries, 0)

	//scan shapepath for existing RPF shapefiles
	//fmt.Printf("walking " + IndexPath)
	filepath.Walk(IndexPath, func(path string, f os.FileInfo, err error) error {
		path = strings.ToUpper(filepath.Base(path))
		if len(path) == 6 && filepath.Ext(path) == ".SHP" {
			seriesCode := path[0:2]
			//fmt.Println("SERIES " + seriesCode)
			series := rpf.DataSeries[seriesCode]
			scale := series.Scale
			if scale == -1 {
				//fmt.Println("SKIPPING BAD SERIES TYPE : " + seriesCode)
				return nil
			}
			if series.Type == rpf.CIB {
				scale *= 15000.0
			}
			bestName := series.GroupCode
			if bestName == "" {
				bestName = seriesCode
			}
			allSeries = append(allSeries, SeriesRes{seriesCode, bestName, scale})
		}
		return nil
	})

	//sort by resolution
	sort.Sort(allSeries)
	//for _, series := range allSeries {
	//	WriteTileLayer(w, series)
	//}
	for _, series := range allSeries {
		WriteShapeLayer(w, series)
	}
	WriteFooter(w)
}

func WriteVector(w io.Writer, vectorTemplatePath, shapePath string) {
	bytes, err := ioutil.ReadFile(vectorTemplatePath)
	if err != nil {
		fmt.Println("error reading map template")
	}
	vectorText := strings.Replace(string(bytes), "{shpPath}", EscapeSlashes(contentPath), -1)
	w.Write([]byte(vectorText))
}

func WriteHeader(w io.Writer, proj4Path, shapePath string, wmsLink string) {
	header := `
MAP
  NAME "CommonMap"
  IMAGETYPE png
  SIZE 1600 800
  UNITS DD
  DEFRESOLUTION 72
  EXTENT -180 -90 180 90
  CONFIG "MS_ERRORFILE" "stderr"
  CONFIG "PROJ_LIB" "` + EscapeSlashes(proj4Path) + `"
  CONFIG "ON_MISSING_DATA" "IGNORE"
  PROJECTION
    'init=epsg:4326'
  END
  SHAPEPATH "` + EscapeSlashes(IndexPath) + `"
  MAXSIZE 4096
  FONTSET "` + EscapeSlashes(fontsetPath) + `"
  

  OUTPUTFORMAT
    NAME "png8"
    DRIVER AGG/PNG8
    MIMETYPE "image/png; mode=8bit"
    EXTENSION "png"
    TRANSPARENT ON
    IMAGEMODE RGBA
    FORMATOPTION "QUANTIZE_FORCE=off"
    FORMATOPTION "QUANTIZE_COLORS=256"
    FORMATOPTION "INTERLACE=ON"
  END

  WEB
    METADATA
      OWS_ENABLE_REQUEST "*"
      WMS_SRS "EPSG:4326 EPSG:3857"
      WMS_EXTENT "-180 -90 180 90"
      WMS_ONLINERESOURCE "` + EscapeSlashes(wmsLink) + `"
      LABELCACHE_MAP_EDGE_BUFFER "-10"
      WMS_TITLE "CommonMap"
    END
  END

`
	w.Write([]byte(header))
}

func WriteFooter(w io.Writer) {
	footer := `
END
`
	w.Write([]byte(footer))
}

func WriteTileLayer(w io.Writer, series SeriesRes) {
	layer := `
  LAYER
    NAME "RPF-` + series.seriesCode + `"
	GROUP "RPF"
	METADATA
		"WMS_TITLE" "RPF-` + series.seriesCode + `"
		"WMS_SRS" "EPSG:4326 EPSG:3857"
	END
	MAXSCALEDENOM ` + strconv.FormatFloat(series.scale*2.0, 'f', 0, 64) + `
	MINSCALEDENOM ` + strconv.FormatFloat(series.scale/3.0, 'f', 0, 64) + `
    STATUS ON
    TYPE RASTER
    TILEINDEX "` + series.seriesCode + `.shp"
    TILEITEM "LOCATION"
  END

`
	w.Write([]byte(layer))
}

func WriteShapeLayer(w io.Writer, series SeriesRes) {
	layer := `
  LAYER
    NAME "RPF-` + series.seriesCode + `-index"
	GROUP "RPF-index"
	METADATA
		"WMS_TITLE" "RPF-` + series.seriesCode + `-index"
		"WMS_SRS" "EPSG:4326 EPSG:3857"
	END
	MAXSCALEDENOM ` + strconv.FormatFloat(series.scale*10, 'f', 0, 64) + `
    STATUS ON
    TYPE POLYGON
    DATA "` + series.seriesCode + `.shp"
    CLASS
		LABEL
		  TEXT "` + series.bestName + `"
		END
      STYLE
        WIDTH 0.5
        OUTLINECOLOR "#006600"
      END
    END
  END

`
	w.Write([]byte(layer))
}
