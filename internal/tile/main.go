package tile

import (
	"image"
	"image/draw"
	"image/png"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/geo/s2"
	"github.com/shinomontaz/tilesrv/config"
	"github.com/shinomontaz/tilesrv/internal/osm"
)

type Server struct {
	prefix string
	gIndex *osm.Data
	styles map[string]map[string]config.Style
}

// New to create new tile server instance
func New(prefix string, gIndex *osm.Data, styles map[string]map[string]config.Style) *Server {
	return &Server{
		prefix: prefix,
		gIndex: gIndex,
		styles: styles,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len(s.prefix):]

	if !(strings.HasPrefix(path, "/") && strings.HasSuffix(path, ".png")) {
		w.Write([]byte("404"))
		return
	}

	xyz := strings.Split(path[1:len(path)-4], "/")
	if len(xyz) != 3 {
		w.Write([]byte("404"))
		return
	}

	xyz_ := []int{}
	for _, value := range xyz {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			w.Write([]byte("404"))
			return
		}

		xyz_ = append(xyz_, intVal)
	}

	zoom := xyz_[0]
	x := xyz_[1]
	y := xyz_[2]

	nwPt := getPointByCoords(x, y, zoom)
	sePt := getPointByCoords(x+1, y+1, zoom)
	const tileSize = 256

	img := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	draw.Draw(img, img.Bounds(), image.White, image.ZP, draw.Src)

	tl := &Tile{
		cv:       img,
		zoom:     zoom,
		tileSize: tileSize,
		nwPt:     nwPt,
		sePt:     sePt,
		p1:       s2.PointFromLatLng(s2.LatLngFromDegrees(nwPt.Lat(), nwPt.Lon())),
		p2:       s2.PointFromLatLng(s2.LatLngFromDegrees(sePt.Lat(), nwPt.Lon())),
		p3:       s2.PointFromLatLng(s2.LatLngFromDegrees(nwPt.Lat(), sePt.Lon())),
		p4:       s2.PointFromLatLng(s2.LatLngFromDegrees(sePt.Lat(), sePt.Lon())),
		styles:   s.styles,
	}

	err := tl.Draw(s.gIndex)
	if err != nil {
		panic(err)
	}

	png.Encode(w, tl.cv)
}
