package config

import (
	"image/color"

	"github.com/fogleman/gg"
)

type Style struct {
	Color color.Color
	Width int
	Dash  float64
}

func (s *Style) Implement(c *gg.Context) {
	c.SetColor(s.Color)
	c.SetLineWidth(s.Width)
	c.SetDash(s.Dash)
}
