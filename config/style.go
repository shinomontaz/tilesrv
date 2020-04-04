package config

import (
	"encoding/json"
	"fmt"

	"github.com/fogleman/gg"
)

type Color struct {
	R, G, B, A float64
}

type Style struct {
	Color Color
	Width int
	Dash  float64
}

type JSONStyle struct {
	Color   string  `json:"color"`
	Width   int     `json:"width"`
	Dash    float64 `json:"dash"`
	Opacity float64 `json:"opacity"`
}

func (s *Style) UnmarshalJSON(b []byte) error {
	var jStyle JSONStyle
	err := json.Unmarshal(b, &jStyle)
	if err != nil {
		return err
	}

	s.Color, err = parseHexColor(jStyle.Color, jStyle.Opacity)
	if err != nil {
		return err
	}

	s.Width = jStyle.Width
	s.Dash = jStyle.Dash

	return nil
}

func (s *Style) Implement(c *gg.Context) {
	c.SetRGBA(s.Color.R, s.Color.G, s.Color.B, s.Color.A)
	c.SetLineWidth(float64(s.Width))
	c.SetDash()
	if s.Dash > 0 {
		c.SetDash(256.0*s.Dash, 256.0*s.Dash*2.0)
	}
}

func parseHexColor(s string, opacity float64) (c Color, err error) {
	if opacity <= 0.0 {
		opacity = 1.0
	}
	var R, G, B uint32

	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &R, &G, &B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &R, &G, &B)
		// Double the hex digits:
		R *= 17
		G *= 17
		B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")
	}
	c = Color{R: float64(R) / 255.0, G: float64(G) / 255.0, B: float64(B) / 255.0, A: opacity}

	return
}
