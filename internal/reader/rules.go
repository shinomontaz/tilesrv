package reader

import "github.com/shinomontaz/tilesrv/internal/types"

var mapRules = map[int][]types.Tag{
	0: []types.Tag{
		{"natural", "coastline"},
		{"admin_level", "2"},
	},
	6: []types.Tag{
		{"admin_level", "3"},
		{"highway", "motorway"},
		{"highway", "trunk"},
		{"railway", "rail"},
	},
	8: []types.Tag{
		{"admin_level", "4"},
		{"highway", "primary"},
		{"highway", "secondary"},
		{"highway", "tertiary"},
	},
	10: []types.Tag{
		{"admin_level", "5"},
		{"highway", "motorway_link"},
		{"highway", "trunk_link"},
		{"highway", "primary_link"},
		{"highway", "secondary_link"},
		{"highway", "road"},
		{"railway", "light_rail"},
		{"railway", "monorail"},
	},
	12: []types.Tag{
		{"admin_level", "6"},
		{"admin_level", "7"},
		{"highway", "unclassified"},
		{"highway", "residential"},
		{"highway", "living_street"},
		{"highway", "bus_guideway"},
		{"highway", "raceway"},
	},
	13: []types.Tag{
		{"admin_level", "8"},
		{"highway", "service"},
	},
	14: []types.Tag{
		{"admin_level", "9"},
		{"highway", "track"},
	},
}
