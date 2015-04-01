// Copyright 2015 Zack Guo <gizak@icloud.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package termui

import "strings"

// List displays []string as its items,
// it has a Overflow option (default is "hidden"), when set to "hidden",
// the item exceeding List's width is truncated, but when set to "wrap",
// the overflowed text breaks into next line.
/*
  strs := []string{
		"[0] github.com/gizak/termui",
		"[1] editbox.go",
		"[2] iterrupt.go",
		"[3] keyboard.go",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] nsf/termbox-go"}

  ls := termui.NewList()
  ls.Items = strs
  ls.ItemFgColor = termui.ColorYellow
  ls.Border.Label = "List"
  ls.Height = 7
  ls.Width = 25
  ls.Y = 0
*/
type Menu struct {
	Block
	Items         []string
	Overflow      string
	ItemFgColor   Attribute
	ItemBgColor   Attribute
	ActiveFgColor Attribute
	ActiveBgColor Attribute

	active int
	scroll int
}

// NewList returns a new *List with current theme.
func NewMenu() *Menu {
	m := &Menu{Block: *NewBlock()}
	m.Overflow = "hidden"
	m.ItemFgColor = theme.ListItemFg
	m.ItemBgColor = theme.ListItemBg
	m.ActiveFgColor = theme.ListItemBg
	m.ActiveBgColor = theme.ListItemFg
	return m
}

func (m *Menu) SelectDown() {
	m.active += 1

	if m.active >= len(m.Items) {
		m.active = len(m.Items) - 1
	}

	if m.active >= m.innerHeight+m.scroll {
		m.scroll += 1
	}
}

func (m *Menu) SelectUp() {
	m.active -= 1

	if m.active < 0 {
		m.active = 0
	}

	if m.active < m.scroll {
		m.scroll -= 1
	}
}

func (m *Menu) SelectedIndex() int {
	return m.active
}

func (m *Menu) ResetSelection() {
	m.active = 0
	m.scroll = 0
}

// Buffer implements Bufferer interface.
func (l *Menu) Buffer() []Point {
	ps := l.Block.Buffer()
	switch l.Overflow {
	case "wrap":
		rs := str2runes(strings.Join(l.Items, "\n"))
		i, j, k := 0, 0, 0
		for i < l.innerHeight && k < len(rs) {
			w := charWidth(rs[k])
			if rs[k] == '\n' || j+w > l.innerWidth {
				i++
				j = 0
				if rs[k] == '\n' {
					k++
				}
				continue
			}
			pi := Point{}
			pi.X = l.innerX + j
			pi.Y = l.innerY + i

			pi.Ch = rs[k]
			pi.Bg = l.ItemBgColor
			pi.Fg = l.ItemFgColor

			ps = append(ps, pi)
			k++
			j++
		}

	case "hidden":
		trimItems := l.Items
		if len(trimItems) > l.innerHeight {
			if len(trimItems) < l.innerHeight+l.scroll {
				trimItems = trimItems[l.scroll:]
			} else {
				trimItems = trimItems[l.scroll : l.innerHeight+l.scroll]
			}
		}

		for i, v := range trimItems {
			rs := trimStr2Runes(v, l.innerWidth)

			j := 0
			for _, vv := range rs {
				w := charWidth(vv)
				p := Point{}
				p.X = l.innerX + j
				p.Y = l.innerY + i

				p.Ch = vv

				if i+l.scroll == l.active {
					p.Bg = l.ActiveBgColor
					p.Fg = l.ActiveFgColor
				} else {
					p.Bg = l.ItemBgColor
					p.Fg = l.ItemFgColor
				}

				ps = append(ps, p)
				j += w
			}

			if i+l.scroll == l.active {
				for x := len(rs); x < l.innerWidth; x++ {
					p := Point{}
					p.X = l.innerX + x
					p.Y = l.innerY + i

					p.Ch = ' '

					p.Bg = l.ActiveBgColor
					p.Fg = l.ActiveFgColor

					ps = append(ps, p)
				}
			}
		}
	}
	return ps
}
