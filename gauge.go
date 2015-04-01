// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "strconv"

// Gauge is a progress bar like widget.
// A simple example:
/*
  g := termui.NewGauge()
  g.Percent = 40
  g.Width = 50
  g.Height = 3
  g.Border.Label = "Slim Gauge"
  g.BarColor = termui.ColorRed
  g.PercentColor = termui.ColorBlue
*/
type Gauge struct {
	Block
	Percent      int
	BarColor     Attribute
	PercentColor Attribute

	LeftMargin  int
	RightMargin int

	LeftText  string
	RightText string
}

// NewGauge return a new gauge with current theme.
func NewGauge() *Gauge {
	g := &Gauge{
		Block:        *NewBlock(),
		PercentColor: theme.GaugePercent,
		BarColor:     theme.GaugeBar,
	}

	g.Width = 12
	g.Height = 5

	return g
}

// Buffer implements Bufferer interface.
func (g *Gauge) Buffer() []Point {
	block := g.Block.Buffer()

	maxWidth := g.innerWidth - g.RightMargin - g.LeftMargin

	width := g.Percent * maxWidth / 100

	mtext := str2runes(strconv.Itoa(g.Percent) + "%")
	ltext := str2runes(g.LeftText)
	rtext := str2runes(g.RightText)

	midy := g.innerY + g.innerHeight/2
	midx := g.innerX + g.LeftMargin + maxWidth/2 - len(rtext)/2

	// plot left text
	for i, v := range ltext {
		p := Point{}
		p.X = g.innerX + i
		p.Y = midy
		p.Ch = v
		p.Bg = g.Block.BgColor
		p.Fg = g.PercentColor
		block = append(block, p)
	}

	// plot right text
	for i, v := range rtext {
		p := Point{}
		p.X = g.innerX + g.innerWidth - g.RightMargin + i
		p.Y = midy
		p.Ch = v
		p.Bg = g.Block.BgColor
		p.Fg = g.PercentColor
		block = append(block, p)
	}

	// plot bar
	for i := 0; i < g.innerHeight; i++ {
		for j := 0; j < width; j++ {
			p := Point{}
			p.X = g.innerX + g.LeftMargin + j
			p.Y = g.innerY + i
			p.Ch = ' '
			p.Bg = g.BarColor
			if p.Bg == ColorDefault {
				p.Bg |= AttrReverse
			}
			block = append(block, p)
		}
	}

	// plot percentage
	for i, v := range mtext {
		p := Point{}
		p.X = midx + i
		p.Y = midy
		p.Ch = v
		p.Fg = g.PercentColor
		if g.innerX+g.LeftMargin+width > p.X {
			p.Bg = g.BarColor
			if p.Bg == ColorDefault {
				p.Bg |= AttrReverse
			}

		} else {
			p.Bg = g.Block.BgColor
		}
		block = append(block, p)
	}
	return block
}
