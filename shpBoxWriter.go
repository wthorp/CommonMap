package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

type ShpBoxWriter struct {
	shp, shx, dbf, qix, prj *os.File
	shpW, shxW, dbfW, qixW  *bufio.Writer
	bbox                    Box
	n                       int32
	shxBuffer               []byte
	shpBuffer               []byte
	qixData                 *qixTree
}

func Create(seriesCode string) (*ShpBoxWriter, error) {
	//filename = filename[0 : len(filename)-3]
	shp, err := os.Create(GetIndexPath(seriesCode + ".shp"))
	if err != nil {
		return nil, err
	}
	shx, err := os.Create(GetIndexPath(seriesCode + ".shx"))
	if err != nil {
		return nil, err
	}
	dbf, err := os.Create(GetIndexPath(seriesCode + ".dbf"))
	if err != nil {
		return nil, err
	}
	qix, err := os.Create(GetIndexPath(seriesCode + ".qix"))
	if err != nil {
		return nil, err
	}
	prj, err := os.Create(GetIndexPath(seriesCode + ".prj"))
	if err != nil {
		return nil, err
	}

	//skip headers at first, we'll write them on Close()
	shp.Seek(100, os.SEEK_SET)
	shx.Seek(100, os.SEEK_SET)
	dbf.Seek(65, os.SEEK_SET)
	//qix.Seek(40, os.SEEK_SET)

	//create our w740000000. r+++++   riter
	s := &ShpBoxWriter{
		shp:       shp,
		shx:       shx,
		dbf:       dbf,
		qix:       qix,
		prj:       prj,
		shpW:      bufio.NewWriterSize(shp, 8192),
		shxW:      bufio.NewWriterSize(shx, 8192),
		dbfW:      bufio.NewWriterSize(dbf, 8192),
		qixW:      bufio.NewWriterSize(qix, 8192),
		bbox:      Box{0.0, 0.0, 0.0, 0.0},
		shxBuffer: make([]byte, 8, 8),
		shpBuffer: make([]byte, 136, 136),
		qixData:   CreateQixTree(),
	}

	//pre-populate fixed record values
	putBigInt32(s.shxBuffer, int32(64), 4) //Content Length (64 16-bit words)
	putBigInt32(s.shpBuffer, int32(64), 4) //Content Length (64 16-bit words)
	putLilInt32(s.shpBuffer, int32(5), 8)  // Shape Type (Polygon)
	putLilInt32(s.shpBuffer, int32(1), 44) // number of parts
	putLilInt32(s.shpBuffer, int32(5), 48) // number of points
	putLilInt32(s.shpBuffer, int32(0), 52) // index to part 0

	return s, nil
}

func (s *ShpBoxWriter) Close() {
	s.shpW.Flush()
	s.shxW.Flush()
	s.dbfW.Flush()
	s.shp.Seek(0, os.SEEK_SET)
	s.shx.Seek(0, os.SEEK_SET)
	s.dbf.Seek(0, os.SEEK_SET)
	s.writeHeader(s.shx)
	s.writeHeader(s.shp)
	s.writeDbfHeader(s.dbf)

	s.writePrjContent(s.prj)
	s.writeQixContent(s.qix)
	s.shp.Close()
	s.shx.Close()
	s.dbf.Close()
	s.prj.Close()
	s.qix.Close()
}

//grow a bbox by another bbpx
func extendBbox(old, new *Box) {
	if old[0] > new[0] { //minX
		old[0] = new[0]
	}
	if old[1] > new[1] { //minY
		old[1] = new[1]
	}
	if old[2] < new[2] { //maxX
		old[2] = new[2]
	}
	if old[3] < new[3] { //maxY
		old[3] = new[3]
	}
}

//Byte 0 Record Number Record Number Integer Big
//Byte 4 Content Length Content Length Integer Big
func (s *ShpBoxWriter) WriteBox(rpf RpfBox) {
	// housekeeping
	var bbox = rpf.box
	if s.n == 0 {
		s.bbox = bbox
	} else {
		extendBbox(&s.bbox, &bbox)
	}
	s.n++ //(begins at 1)

	// write shp
	shpBuffer := s.shpBuffer
	putBigInt32(shpBuffer, s.n, 0)
	putLilFloat64(shpBuffer, bbox[0], 12)
	putLilFloat64(shpBuffer, bbox[1], 20)
	putLilFloat64(shpBuffer, bbox[2], 28)
	putLilFloat64(shpBuffer, bbox[3], 36)
	putLilFloat64(shpBuffer, bbox[0], 56)
	putLilFloat64(shpBuffer, bbox[3], 64)
	putLilFloat64(shpBuffer, bbox[2], 72)
	putLilFloat64(shpBuffer, bbox[3], 80)
	putLilFloat64(shpBuffer, bbox[2], 88)
	putLilFloat64(shpBuffer, bbox[1], 96)
	putLilFloat64(shpBuffer, bbox[0], 104)
	putLilFloat64(shpBuffer, bbox[1], 112)
	putLilFloat64(shpBuffer, bbox[0], 120)
	putLilFloat64(shpBuffer, bbox[3], 128)
	_, err := s.shpW.Write(shpBuffer)
	if err != nil {
		fmt.Println("error writing shp: ", err)
	}

	// write shx
	putBigInt32(s.shxBuffer, int32(-18+(68*s.n)), 0) // start index
	s.shxW.Write(s.shxBuffer)

	//build in memory QIX tree
	s.qixData.Insert(s.n, &bbox)
}

func (s *ShpBoxWriter) WriteDbf(path string) {
	//write dbf
	s.dbfW.WriteString(pad(path))
}

// Writes SHP/SHX headers to specified file.
func (s *ShpBoxWriter) writeHeader(file *os.File) {
	filelength, _ := file.Seek(0, os.SEEK_END)
	if filelength == 0 {
		filelength = 100
	}
	file.Seek(0, os.SEEK_SET)
	// file code
	Write(file, binary.BigEndian, []int32{9994, 0, 0, 0, 0, 0})
	// file length
	Write(file, binary.BigEndian, int32(filelength/2))
	// version and shape type
	Write(file, binary.LittleEndian, []int32{1000, 5})
	// bounding box
	Write(file, binary.LittleEndian, s.bbox)
	// elevation, measure
	Write(file, binary.LittleEndian, []float64{0.0, 0.0, 0.0, 0.0})
}

// Write DBF header.
func (s *ShpBoxWriter) writeDbfHeader(file *os.File) {
	file.Seek(0, os.SEEK_SET)
	// version, year (YEAR-1990), month, day
	Write(file, binary.LittleEndian, []byte{3, 24, 5, 3})
	// number of records
	Write(file, binary.LittleEndian, s.n)
	// header length (#fields * 32 + 33), record length (field sizes + 1)
	Write(file, binary.LittleEndian, []int16{65, 255})
	// padding
	Write(file, binary.LittleEndian, make([]byte, 20))
	//location field
	Write(file, binary.LittleEndian, []byte("location   ")) //Name
	Write(file, binary.LittleEndian, []byte("C"))           //Fieldtype
	Write(file, binary.LittleEndian, make([]byte, 4))       // Addr
	Write(file, binary.LittleEndian, uint8(255))            // Size
	Write(file, binary.LittleEndian, uint8(0))              // Precision
	Write(file, binary.LittleEndian, make([]byte, 14))      // Padding
	// end with return
	Write(file, binary.LittleEndian, []byte("\r"))
}

func pad(s string) string {
	return " " + s + strings.Repeat(" ", 254-len(s))
}

func (s *ShpBoxWriter) writePrjContent(file *os.File) {
	file.WriteString(`GEOGCS["GCS_WGS_1984",DATUM["D_WGS_1984",SPHEROID["WGS_1984",6378137,298.257223563]],PRIMEM["Greenwich",0],UNIT["Degree",0.017453292519943295]]`)
}

func (s *ShpBoxWriter) writeQixContent(file *os.File) {
	//fmt.Println("starting writeQixHeader")
	header := [8]byte{'S', 'Q', 'T', 1, 1, 0, 0, 0}
	Write(file, binary.BigEndian, header)
	Write(file, binary.LittleEndian, s.qixData.numFeatures)
	Write(file, binary.LittleEndian, int32(12))
	s.writeQixNode(file, s.qixData.root)
}

// calculate # of bytes to skip this node and subnodes
func qixGetNodeSkipOffset(node *qixNode) int32 {
	var offset int32
	for i := int32(0); i < node.numSubNodes; i++ {
		if node.SubNodes[i] != nil {
			offset += 44 + (node.SubNodes[i].numFeatures * 4)
			offset += qixGetNodeSkipOffset(node.SubNodes[i])
		}
	}
	return (offset)
}

func (s *ShpBoxWriter) writeQixNode(file *os.File, node *qixNode) {
	//fmt.Println("starting writeQixNode")
	offset := qixGetNodeSkipOffset(node)
	size := 44 + (node.numFeatures * 4)
	qixBuffer := make([]byte, size, size)

	putLilInt32(qixBuffer, offset, 0)
	putLilFloat64(qixBuffer, node.bbox[0], 4)
	putLilFloat64(qixBuffer, node.bbox[1], 12)
	putLilFloat64(qixBuffer, node.bbox[2], 20)
	putLilFloat64(qixBuffer, node.bbox[3], 28)
	putLilInt32(qixBuffer, node.numFeatures, 36)
	index := int32(40)
	for i := int32(0); i < node.numFeatures; i++ {
		putLilInt32(qixBuffer, node.FeatureIds[i], index)
		index += 4
	}
	putLilInt32(qixBuffer, node.numSubNodes, index)
	Write(file, binary.LittleEndian, qixBuffer)

	for i := int32(0); i < node.numSubNodes; i++ {
		if node.SubNodes[i] != nil {
			s.writeQixNode(file, node.SubNodes[i])
		}
	}
}

func Write(w io.Writer, order binary.ByteOrder, data interface{}) {
	err := binary.Write(w, order, data)
	if err != nil {
		fmt.Println(err, data)
	}
	return
}

func putBigInt32(b []byte, v int32, index int32) {
	b[index+0] = byte(v >> 24)
	b[index+1] = byte(v >> 16)
	b[index+2] = byte(v >> 8)
	b[index+3] = byte(v)
}

func putLilInt32(b []byte, v int32, index int32) {
	b[index+0] = byte(v)
	b[index+1] = byte(v >> 8)
	b[index+2] = byte(v >> 16)
	b[index+3] = byte(v >> 24)
}

func putLilFloat64(b []byte, v float64, index int32) {
	u := math.Float64bits(v)
	b[index+0] = byte(u)
	b[index+1] = byte(u >> 8)
	b[index+2] = byte(u >> 16)
	b[index+3] = byte(u >> 24)
	b[index+4] = byte(u >> 32)
	b[index+5] = byte(u >> 40)
	b[index+6] = byte(u >> 48)
	b[index+7] = byte(u >> 56)
}
