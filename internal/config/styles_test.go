// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of K9s

package config_test

import (
	"testing"

	"github.com/derailed/k9s/internal/config"
	"github.com/derailed/tcell/v2"
	"github.com/derailed/tview"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStyle(t *testing.T) {
	s := config.NewStyles()

	assert.Equal(t, config.Color("black"), s.K9s.Body.BgColor)
	assert.Equal(t, config.Color("cadetblue"), s.K9s.Body.FgColor)
	assert.Equal(t, config.Color("lightskyblue"), s.K9s.Frame.Status.NewColor)
}

func TestColor(t *testing.T) {
	uu := map[string]tcell.Color{
		"blah":    tcell.ColorDefault,
		"blue":    tcell.ColorBlue.TrueColor(),
		"#ffffff": tcell.NewHexColor(16777215),
		"#ff0000": tcell.NewHexColor(16711680),
	}

	for k := range uu {
		c, u := k, uu[k]
		t.Run(k, func(t *testing.T) {
			assert.Equal(t, u, config.NewColor(c).Color())
		})
	}
}

func TestSkinHappy(t *testing.T) {
	s := config.NewStyles()
	require.NoError(t, s.Load("../../skins/black-and-wtf.yaml", false))
	s.Update()

	assert.Equal(t, "#ffffff", s.Body().FgColor.String())
	assert.Equal(t, "#000000", s.Body().BgColor.String())
	assert.Equal(t, "#000000", s.Table().BgColor.String())
	assert.Equal(t, tcell.ColorWhite.TrueColor(), s.FgColor())
	assert.Equal(t, tcell.ColorBlack.TrueColor(), s.BgColor())
	assert.Equal(t, tcell.ColorBlack.TrueColor(), tview.Styles.PrimitiveBackgroundColor)
}

func TestSkinLoad(t *testing.T) {
	uu := map[string]struct {
		f   string
		err string
	}{
		"not-exist": {
			f:   "testdata/skins/blee.yaml",
			err: "open testdata/skins/blee.yaml: no such file or directory",
		},
		"toast": {
			f: "testdata/skins/boarked.yaml",
			err: `Additional property bgColor is not allowed
Additional property fgColor is not allowed
Additional property logoColor is not allowed
Invalid type. Expected: object, given: array`,
		},
	}

	for k := range uu {
		u := uu[k]
		t.Run(k, func(t *testing.T) {
			s := config.NewStyles()
			err := s.Load(u.f, false)
			if err != nil {
				assert.Equal(t, u.err, err.Error())
			}
			assert.Equal(t, "#5f9ea0", s.Body().FgColor.String())
			assert.Equal(t, "#000000", s.Body().BgColor.String())
			assert.Equal(t, "#000000", s.Table().BgColor.String())
			assert.Equal(t, tcell.ColorCadetBlue.TrueColor(), s.FgColor())
			assert.Equal(t, tcell.ColorBlack.TrueColor(), s.BgColor())
			assert.Equal(t, tcell.ColorBlack.TrueColor(), tview.Styles.PrimitiveBackgroundColor)
		})
	}
}

func TestStylesGlobalFallback(t *testing.T) {
	s := config.NewStyles()
	g := s.Global()

	assert.Equal(t, s.K9s.Body.FgColor, g.FgColor)
	assert.Equal(t, s.K9s.Body.BgColor, g.BgColor)
	assert.Equal(t, s.K9s.Views.Yaml.KeyColor, g.KeyColor)
	assert.Equal(t, s.K9s.Views.Yaml.ValueColor, g.ValueColor)
	assert.Equal(t, s.K9s.Views.Yaml.ColonColor, g.DelimiterColor)
	assert.Equal(t, s.K9s.Frame.Status.ErrorColor, g.ErrorColor)
	assert.Equal(t, s.K9s.Frame.Status.HighlightColor, g.AccentColor)
	assert.Equal(t, s.K9s.Frame.Status.CompletedColor, g.MutedColor)
	assert.Equal(t, s.K9s.Frame.Status.NewColor, g.InfoColor)
	assert.Equal(t, s.K9s.Body.LogoColorInfo, g.SuccessColor)
	assert.Equal(t, s.K9s.Body.LogoColorWarn, g.WarnColor)
}

func TestStylesGlobalLoad(t *testing.T) {
	s := config.NewStyles()
	require.NoError(t, s.Load("testdata/skins/global.yaml", false))

	g := s.Global()
	assert.Equal(t, config.Color("red"), g.KeyColor)
	assert.Equal(t, config.Color("#00ff00"), g.DelimiterColor)
	// Unset global colors still derive from the skin.
	assert.Equal(t, s.K9s.Views.Yaml.ValueColor, g.ValueColor)

	// Resetting must drop the previous skin's global colors.
	s.Reset(false)
	assert.Equal(t, s.K9s.Views.Yaml.KeyColor, s.Global().KeyColor)
}

func TestStylesGlobalInvert(t *testing.T) {
	s := config.NewStyles()
	require.NoError(t, s.Load("testdata/skins/global.yaml", true))

	assert.Equal(t, config.Color("#00ff00").InvertColor(), s.Global().DelimiterColor)
	assert.NotEqual(t, config.Color("#00ff00"), s.Global().DelimiterColor)
}
