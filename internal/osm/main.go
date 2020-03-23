package osm

import (
	"github.com/shinomontaz/tilesrv/internal/geom"
	"github.com/shinomontaz/tilesrv/internal/types"
)

type Data struct {
	Nodes  map[int64]types.IPoint
	Ways   map[int64]geom.Way
	Findex geom.S2Index
}

func (o *Data) GetFeatures(nwPt, sePt types.IPoint, zoom int) []types.IDrawable {
	ret := []types.IDrawable{}

	for _, f := range o.Findex.GetFeatures(nwPt, sePt, zoom) {
		switch f.Type {
		case types.ItemTypeNode:
			//TODO
		case types.ItemTypeWay:
			ret = append(ret, o.Ways[f.Id])
		}
	}

	return ret
}
