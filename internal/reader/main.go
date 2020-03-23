package reader

import (
	"github.com/shinomontaz/tilesrv/internal/osm"
	"github.com/shinomontaz/tilesrv/internal/types"
)

type Reader struct {
	rules    map[int][]types.Tag
	filename string
}

//New to create new reader instance
func New(filename string) *Reader {
	return &Reader{
		rules:    mapRules,
		filename: filename,
	}
}

// Init to read to memory and process provided pbf file, create spatial index
func (r *Reader) Init() *osm.Data {
	data, err := r.ParsePbf(r.filename)
	if err != nil {
		panic(err)
	}
	return data
}
