package config

import (
	"encoding/json"
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
)

type Style struct {
	Color color.Color
	Width int
	Dash  float64
}

type JSONStyle struct {
	Color string
	Width int
	Dash  float64
}

func (s *Style) UnmarshalJSON(b []byte) error {
	var jStyle JSONStyle
	err := json.Unmarshal(b, &jStyle)
	if err != nil {
		return err
	}

	s.Color, err = parseHexColor(jStyle.Color)
	if err != nil {
		return err
	}

	s.Width = jStyle.Width
	s.Dash = jStyle.Dash

	return nil
}

func (s *Style) Implement(c *gg.Context) {
	c.SetColor(s.Color)
	c.SetLineWidth(float64(s.Width))
	c.SetDash()
	if s.Dash > 0 {
		c.SetDash(s.Dash)
	}
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	return
}
