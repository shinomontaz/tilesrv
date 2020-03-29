package tile

import "math"

type Point struct {
	Lon_, Lat_ float64
}

func (p Point) Lon() float64 { return p.Lon_ }
func (p Point) Lat() float64 { return p.Lat_ }

// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Tile_numbers_to_lon..2Flat.
func getPointByCoords(x, y, zoom int) Point {
	n := math.Pow(2, float64(zoom))
	lon := (float64(x) / n * 360) - 180
	latRad := math.Atan(math.Sinh(math.Pi * (1 - (2 * float64(y) / n))))
	lat := latRad * 180 / math.Pi

	return Point{lon, lat}
}
