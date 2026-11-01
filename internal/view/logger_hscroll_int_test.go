// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package view

import (
	"strings"
	"testing"

	"github.com/derailed/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func newHScrollLogger() *Logger {
	l := NewLogger(nil)
	l.SetScrollable(true).SetWrap(false)
	l.SetText(strings.Repeat("x", 200) + "\nshort")
	l.SetRect(0, 0, 40, 10)

	return l
}

func TestLoggerHScroll(t *testing.T) {
	uu := map[string]struct {
		start int
		evt   *tcell.EventKey
		ok    bool
		col   int
	}{
		"shift-right": {
			evt: tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModShift),
			ok:  true,
			col: 20,
		},
		"shift-left": {
			start: 50,
			evt:   tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModShift),
			ok:    true,
			col:   30,
		},
		"shift-left-floor": {
			start: 5,
			evt:   tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModShift),
			ok:    true,
		},
		"ctrl-right": {
			evt: tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModCtrl),
			ok:  true,
			col: 160,
		},
		"ctrl-left": {
			start: 70,
			evt:   tcell.NewEventKey(tcell.KeyLeft, 0, tcell.ModCtrl),
			ok:    true,
		},
		"shift-end": {
			evt: tcell.NewEventKey(tcell.KeyEnd, 0, tcell.ModShift),
			ok:  true,
			col: 160,
		},
		"shift-home": {
			start: 70,
			evt:   tcell.NewEventKey(tcell.KeyHome, 0, tcell.ModShift),
			ok:    true,
		},
		"plain-right": {
			start: 3,
			evt:   tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModNone),
			col:   3,
		},
		"plain-home": {
			start: 3,
			evt:   tcell.NewEventKey(tcell.KeyHome, 0, tcell.ModNone),
			col:   3,
		},
	}

	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			l := newHScrollLogger()
			l.ScrollTo(0, u.start)

			assert.Equal(t, u.ok, l.hScroll(u.evt))
			_, col := l.GetScrollOffset()
			assert.Equal(t, u.col, col)
		})
	}
}

func TestLoggerHScrollCmds(t *testing.T) {
	l := newHScrollLogger()

	assert.Nil(t, l.scrollRightCmd(nil))
	_, col := l.GetScrollOffset()
	assert.Equal(t, 20, col)

	assert.Nil(t, l.scrollLeftCmd(nil))
	_, col = l.GetScrollOffset()
	assert.Equal(t, 0, col)
}

func TestLoggerHScrollFilterActive(t *testing.T) {
	l := newHScrollLogger()
	l.cmdBuff.SetActive(true)

	assert.False(t, l.hScroll(tcell.NewEventKey(tcell.KeyRight, 0, tcell.ModShift)))
}
