package geom

import (
	"fmt"
	"strconv"

	"github.com/golang/geo/s2"
	"github.com/shinomontaz/tilesrv/internal/types"
)

type Way struct {
	NodeIDs []int64
	Id      int64
	Tags    map[string]string
	Nodes   []types.IPoint
}

func (w *Way) SetNodes(nodes map[int64]types.IPoint) {
	w.Nodes = make([]types.IPoint, 0)
	for _, id := range w.NodeIDs {
		if _, ok := nodes[id]; !ok {
			fmt.Println(id)
			panic("cannot find node!")
		}
		w.Nodes = append(w.Nodes, nodes[id])
	}
}

func (w Way) match(tags []types.Tag) bool {
	for key, val := range w.Tags {
		for _, tag := range tags {
			if key == tag.Key && (val == tag.Val || tag.Val == "*") {
				return true
			}
		}
	}
	return false
}

func (w Way) MatchAny(rules map[int][]types.Tag) (int, bool) {
	if min_zoom, ok := w.Tags["min_zoom"]; ok {
		mz, err := strconv.ParseFloat(min_zoom, 64)
		if err == nil {
			return int(mz), true
		}
	}
	for zoom, tags := range rules {
		if w.match(tags) {
			return zoom, true
		}
	}

	return -1, false
}

func (w Way) Draw(t types.ITile) {
	var (
		prevCoords  []float64 // previous  plotting point, nil if already added to coords or first time
		prevS2Point s2.Point  // previous  s2 Point. Always present, except first time
	)
	coords := [][]float64{}
	prevWithinBounds := false
	first := true
	for _, node := range w.Nodes {
		x, y := t.GetRelativeXY(node)
		s2Point := s2.PointFromLatLng(s2.LatLngFromDegrees(node.Lat(), node.Lon()))

		if first {
			first = false
			if t.IsInside(node) {
				coords = append(coords, []float64{x, y})
				prevWithinBounds = true
			} else {
				prevS2Point = s2Point
				prevCoords = []float64{x, y}
				prevWithinBounds = false
			}
			continue
		}

		if t.IsInside(node) {
			if prevWithinBounds == false {
				if len(prevCoords) > 0 {
					coords = append(coords, prevCoords)
				}
			}
			coords = append(coords, []float64{x, y})
			prevWithinBounds = true
			prevCoords = nil
		} else { // Очередная точка вне тайла, следовательно рисуем до границы тайла
			if prevWithinBounds == true {
				coords = append(coords, []float64{x, y})
				prevCoords = nil
			} else {

				if t.IsCrossing(s2Point, prevS2Point) {
					if len(prevCoords) > 0 {
						coords = append(coords, prevCoords)
					}
					coords = append(coords, []float64{x, y})
					prevCoords = nil
				} else {
					if len(coords) > 0 {
						t.DrawPolyLine(coords, w.Tags)
						coords = coords[:0]
					}
					prevCoords = []float64{x, y}
				}
			}
			prevS2Point = s2Point
			prevWithinBounds = false
		}
	}

	if len(coords) > 0 {
		t.DrawPolyLine(coords, w.Tags)
	}
}
