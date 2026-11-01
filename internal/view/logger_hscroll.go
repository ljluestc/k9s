// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package view

import (
	"strings"

	"github.com/derailed/k9s/internal/ui"
	"github.com/derailed/tcell/v2"
)

// minHScrollStep is the smallest number of columns a fast horizontal scroll moves.
const minHScrollStep = 10

func bindHScrollKeys(l *Logger) {
	l.actions.Bulk(ui.KeyMap{
		ui.KeyLess:    ui.NewKeyAction("Scroll Left", l.scrollLeftCmd, false),
		ui.KeyGreater: ui.NewKeyAction("Scroll Right", l.scrollRightCmd, false),
	})
}

// hScroll handles modifier based horizontal scrolling keys which can not be
// expressed as plain key actions.
//
//	Shift-Left / Shift-Right: scroll half a page left/right
//	Ctrl-Left / Shift-Home:   jump to the beginning of the lines
//	Ctrl-Right / Shift-End:   jump to the end of the longest line
func (l *Logger) hScroll(evt *tcell.EventKey) bool {
	if evt == nil || l.cmdBuff.IsActive() {
		return false
	}
	mod := evt.Modifiers()
	switch evt.Key() {
	case tcell.KeyLeft:
		switch {
		case mod&tcell.ModCtrl != 0:
			l.scrollColumnTo(0)
		case mod&tcell.ModShift != 0:
			l.scrollColumnBy(-l.hScrollStep())
		default:
			return false
		}
	case tcell.KeyRight:
		switch {
		case mod&tcell.ModCtrl != 0:
			l.scrollColumnTo(l.lastColumn())
		case mod&tcell.ModShift != 0:
			l.scrollColumnBy(l.hScrollStep())
		default:
			return false
		}
	case tcell.KeyHome:
		if mod&tcell.ModShift == 0 {
			return false
		}
		l.scrollColumnTo(0)
	case tcell.KeyEnd:
		if mod&tcell.ModShift == 0 {
			return false
		}
		l.scrollColumnTo(l.lastColumn())
	default:
		return false
	}

	return true
}

func (l *Logger) scrollLeftCmd(*tcell.EventKey) *tcell.EventKey {
	l.scrollColumnBy(-l.hScrollStep())
	return nil
}

func (l *Logger) scrollRightCmd(*tcell.EventKey) *tcell.EventKey {
	l.scrollColumnBy(l.hScrollStep())
	return nil
}

// hScrollStep returns the number of columns to move, i.e. half the view width.
func (l *Logger) hScrollStep() int {
	_, _, w, _ := l.GetInnerRect()
	return max(w/2, minHScrollStep)
}

func (l *Logger) scrollColumnBy(delta int) {
	_, col := l.GetScrollOffset()
	l.scrollColumnTo(max(col+delta, 0))
}

func (l *Logger) scrollColumnTo(col int) {
	row, _ := l.GetScrollOffset()
	l.ScrollTo(row, col)
}

// lastColumn returns the column offset showing the end of the longest line.
// Byte length may overshoot the display width for multi-byte runes, in which
// case the text view clamps the offset on the next draw.
func (l *Logger) lastColumn() int {
	_, _, w, _ := l.GetInnerRect()
	longest := 0
	for line := range strings.SplitSeq(l.GetText(true), "\n") {
		longest = max(longest, len(line))
	}

	return max(longest-w, 0)
}
