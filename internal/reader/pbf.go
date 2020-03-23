package reader

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/qedus/osmpbf"
	"github.com/shinomontaz/tilesrv/internal/geom"
	"github.com/shinomontaz/tilesrv/internal/osm"
	"github.com/shinomontaz/tilesrv/internal/types"
)

func NodeFromPbf(n *osmpbf.Node) geom.Node {
	return geom.Node{
		Lon_: n.Lon,
		Lat_: n.Lat,
		Id:   n.ID,
	}
}

func WayFromPbf(w *osmpbf.Way) geom.Way {
	return geom.Way{
		NodeIDs: w.NodeIDs,
		Id:      w.ID,
		Tags:    w.Tags,
	}
}

func (r *Reader) ParsePbf(path string) (*osm.Data, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	d := osmpbf.NewDecoder(f)
	err = d.Start(runtime.GOMAXPROCS(0) - 1)
	if err != nil {
		return nil, err
	}

	data := &osm.Data{
		Nodes:  map[int64]types.IPoint{},
		Ways:   map[int64]geom.Way{},
		Findex: make(geom.S2Index),
	}

	for {
		if v, err := d.Decode(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else {
			switch v := v.(type) {
			case *osmpbf.Node:
				node := NodeFromPbf(v)
				data.Nodes[node.Id] = &node
			case *osmpbf.Way:
				way := WayFromPbf(v)
				zoom, ok := way.MatchAny(r.rules)
				if ok {
					way.SetNodes(data.Nodes)
					data.Ways[way.Id] = way
					data.Findex.AddWay(way, zoom)
				}
			case *osmpbf.Relation:
				// Ignore
			default:
				return nil, fmt.Errorf("unknown type %T", v)
			}
		}
	}

	log.Println("Num s2Cells", len(data.Findex))
	log.Println("Num ways", len(data.Ways))
	log.Println("Num nodes", len(data.Nodes))
	return data, nil
}
