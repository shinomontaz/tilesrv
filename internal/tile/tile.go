package tile

import (
	"image"
	"math"

	"github.com/golang/geo/s2"
	"github.com/shinomontaz/tilesrv/config"
	"github.com/shinomontaz/tilesrv/internal/osm"
	"github.com/shinomontaz/tilesrv/internal/types"

	"github.com/fogleman/gg"
)

type Tile struct {
	p1             s2.Point
	p2             s2.Point
	p3             s2.Point
	p4             s2.Point
	cv             *image.RGBA
	zoom, tileSize int
	nwPt, sePt     Point
	styles         map[string]map[string]config.Style
}

func (t *Tile) GetNw() types.IPoint {
	return t.nwPt
}

func (t *Tile) GetSe() types.IPoint {
	return t.sePt
}

func (t *Tile) GetRelativeXY(p types.IPoint) (float64, float64) {
	baseX, baseY := t.getXY(t.nwPt)
	nodeX, nodeY := t.getXY(p)
	x := nodeX - baseX
	y := nodeY - baseY

	return x, y
}

func (t *Tile) getXY(pt types.IPoint) (float64, float64) {
	scale := math.Pow(2, float64(t.zoom))
	x := ((pt.Lon() + 180) / 360) * scale * float64(t.tileSize)
	y := (float64(t.tileSize) / 2) - (float64(t.tileSize)*math.Log(math.Tan((math.Pi/4)+((pt.Lat()*math.Pi/180)/2)))/(2*math.Pi))*scale

	return x, y
}

func (t *Tile) IsInside(p types.IPoint) bool {
	return p.Lon() > t.nwPt.Lon() && p.Lon() < t.sePt.Lon() &&
		p.Lat() < t.nwPt.Lat() && p.Lat() > t.sePt.Lat()
}

func (t *Tile) IsCrossing(t1, t2 s2.Point) bool {
	return s2.VertexCrossing(t.p1, t.p2, t1, t2) || s2.VertexCrossing(t.p1, t.p3, t1, t2) || s2.VertexCrossing(t.p2, t.p4, t1, t2)
}

func (tl *Tile) GetZoom() int {
	return tl.zoom
}

func (tl *Tile) Draw(osmData *osm.Data) error {
	for _, feature := range osmData.GetFeatures(tl.nwPt, tl.sePt, tl.zoom) {
		feature.Draw(tl)
	}
	return nil
}

func (t *Tile) DrawPolyLine(coords [][]float64, tags map[string]string) {
	path := gg.NewContextForRGBA(t.cv)

	t.style(path, tags)

	for i, coord := range coords {
		if i == 0 {
			path.MoveTo(coord[0], coord[1])
		} else {
			path.LineTo(coord[0], coord[1])
		}
	}

	path.Stroke()
}

func (t *Tile) style(c *gg.Context, tags map[string]string) {
	styled := false

	for key, val := range tags {
		if _, exists := t.styles[key]; !exists {
			continue
		}

		for tagKey, style := range t.styles[key] {
			if tagKey != val {
				continue
			}
			style.Implement(c)
			styled = true
		}
	}

	if !styled {
		def := t.styles["default"]["default"]
		def.Implement(c)
	}
}
