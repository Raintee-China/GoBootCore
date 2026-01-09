package common

import (
	"fmt"

	"github.com/jonas-p/go-shp"
)

var shapeTypeNames = map[shp.ShapeType]string{
	shp.NULL:        "NULL",
	shp.POINT:       "POINT",
	shp.POLYLINE:    "POLYLINE",
	shp.POLYGON:     "POLYGON",
	shp.MULTIPOINT:  "MULTIPOINT",
	shp.POINTZ:      "POINTZ",
	shp.POLYLINEZ:   "POLYLINEZ",
	shp.POLYGONZ:    "POLYGONZ",
	shp.MULTIPOINTZ: "MULTIPOINTZ",
	shp.POINTM:      "POINTM",
	shp.POLYLINEM:   "POLYLINEM",
	shp.POLYGONM:    "POLYGONM",
	shp.MULTIPOINTM: "MULTIPOINTM",
	shp.MULTIPATCH:  "MULTIPATCH",
}

// ShpFileInfo represents basic information from a SHP file
type ShpFileInfo struct {
	ShapeType string
	NumShapes int
	Fields    []shp.Field
}

// ParseShpFile opens and parses a SHP file, returning its basic information.
func ParseShpFile(filePath string) (*ShpFileInfo, error) {
	shape, err := shp.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SHP file: %w", err)
	}
	defer shape.Close()

	fileFields := shape.Fields()
	numShapes := 0
	fileShapeType := "UNKNOWN"
	if name, ok := shapeTypeNames[shape.GeometryType]; ok {
		fileShapeType = name
	}

	for shape.Next() {
		numShapes++
	}

	return &ShpFileInfo{
		ShapeType: fileShapeType,
		NumShapes: numShapes,
		Fields:    fileFields,
	}, nil
}
