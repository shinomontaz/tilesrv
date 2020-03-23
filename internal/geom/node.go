package geom

import (
	"github.com/shinomontaz/tilesrv/internal/types"
)

type Node struct {
	Lon_, Lat_ float64
	Id         int64
}

func (n *Node) Lon() float64 { return n.Lon_ }
func (n *Node) Lat() float64 { return n.Lat_ }

func (n Node) Draw(t types.ITile, nodes map[int64]types.IPoint) {

}
