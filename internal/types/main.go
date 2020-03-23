package types

import (
	"github.com/golang/geo/s2"
	//	"github.com/shinomontaz/tilesrv/internal/geom"
)

const (
	ItemTypeNode = iota
	ItemTypeWay
	ItemTypeRelation
)

type Tag struct {
	Key, Val string
}

type IPoint interface {
	Lon() float64
	Lat() float64
}

type ITile interface {
	GetNw() IPoint
	GetSe() IPoint
	IsInside(p IPoint) bool
	IsCrossing(p1, p2 s2.Point) bool
	DrawPolyLine(coords [][]float64, tags map[string]string)
	GetRelativeXY(p IPoint) (float64, float64)
}

type IDrawable interface {
	Draw(t ITile)
}
