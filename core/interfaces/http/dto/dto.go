package dto

import (
	"github.com/paulmach/orb/geojson"
)

type Request struct {
	// File1Path string `json:"file1_path"`
	// File2Path string `json:"file2_path"`
	Accuracy  int    `json:"accuracy"`
}

type Response struct {
	File3 *geojson.FeatureCollection `json:"file3"`
}