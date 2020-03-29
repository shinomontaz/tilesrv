package geom

import (
	"github.com/golang/geo/s2"
	"github.com/shinomontaz/tilesrv/internal/types"
)

type FeatureRef struct {
	Id   int64
	Type int
	Zoom int
}

type S2Index map[s2.CellID]([]FeatureRef)

func (si S2Index) AddWay(w Way, zoom int) {
	cap := s2.EmptyCap()

	for _, node := range w.Nodes {
		cap = cap.AddPoint(s2.PointFromLatLng(s2.LatLngFromDegrees(node.Lat(), node.Lon())))
	}

	if cap.IsEmpty() {
		return
	}

	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 10}
	cellUnion := rc.FastCovering(cap)
	for _, cid := range cellUnion {
		si[cid] = append(si[cid], FeatureRef{w.Id, types.ItemTypeWay, zoom})
		for l := cid.Level(); l > 0; l-- {
			cid = cid.Parent(l - 1)
			if _, ok := si[cid]; !ok {
				si[cid] = make([]FeatureRef, 0)
			}
		}
	}
}

func (si S2Index) GetFeatures(nwPt, sePt types.IPoint, zoom int) []FeatureRef {

	r := s2.RectFromLatLng(s2.LatLngFromDegrees(nwPt.Lat(), nwPt.Lon()))
	r = r.AddPoint(s2.LatLngFromDegrees(sePt.Lat(), sePt.Lon()))

	rc := &s2.RegionCoverer{MaxLevel: 30, MaxCells: 10}

	cellUnion := rc.Covering(r)

	visitCid := make(map[s2.CellID]bool)
	visitF := make(map[int64]bool)
	ret := make([]FeatureRef, 0)

	for _, cid := range cellUnion {
		if v, ok := si[cid]; ok {
			ret = si.visitDown(cid, v, visitF, ret)
		}
		for l := cid.Level(); l > 0; l-- { // по уровню идем вниз (т.е. увеличивается размер клетки)
			cid = cid.Parent(l - 1)
			ret = si.visitUp(cid, visitCid, visitF, ret)
		}
	}

	i := 0
	for _, fr := range ret {
		if fr.Zoom > 0 && fr.Zoom < zoom {
			ret[i] = fr
			i++
		}
	}

	return ret[:i]
}

func (si S2Index) visitUp(cid s2.CellID, visitCid map[s2.CellID]bool, visitF map[int64]bool, ret []FeatureRef) []FeatureRef {
	fr, ok := si[cid]
	if !ok {
		return ret
	}
	if visitCid[cid] {
		return ret
	}
	visitCid[cid] = true

	for _, f := range fr {
		if !visitF[f.Id] {
			ret = append(ret, f)
			visitF[f.Id] = true
		}
	}
	return ret
}

func (si S2Index) visitDown(cid s2.CellID, fr []FeatureRef, visitF map[int64]bool, ret []FeatureRef) []FeatureRef {

	for _, f := range fr {
		if !visitF[f.Id] {
			ret = append(ret, f)
			visitF[f.Id] = true
		}
	}

	if !cid.IsLeaf() { // Если не самый нижний уровень этой клетки то проходимся по потомкам
		chs := cid.Children()
		for _, cellID := range chs {
			if v, ok := si[cellID]; ok {
				ret = si.visitDown(cellID, v, visitF, ret)
			}
		}
	}

	return ret
}
