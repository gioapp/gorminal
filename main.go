package main

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gop9/gorminal/mod"
	"image"
	"image/color"
	"os/exec"
	"strings"
	"time"
)

var (
	testLabel         = "testtopLabel"
	consoleInputField = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	consoleOutputList = &layout.List{
		Axis:        layout.Vertical,
		ScrollToEnd: true,
	}
	fontSize  float32 = 16
	textColor         = color.RGBA{A: 0xff, R: 0x30, G: 0xcf, B: 0x30}
)

func main() {
	go func() {
		w := app.NewWindow()
		gofont.Register()
		th := material.NewTheme()
		// START INIT OMIT
		com := mod.CommandsHistory{}
		gtx := layout.NewContext(w.Queue())
		// END INIT OMIT
		for e := range w.Events() {
			if e, ok := e.(system.FrameEvent); ok {
				gtx.Reset(e.Config, e.Size)
				fill(gtx)
				layout.Flex{}.Layout(gtx,
					layout.Flexed(1, func() {
						layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
							layout.Flex{
								Axis:    layout.Vertical,
								Spacing: layout.SpaceAround,
							}.Layout(gtx,
								layout.Flexed(1, func() {
									consoleOutputList.Layout(gtx, len(com.Commands), func(i int) {
										t := com.Commands[i]
										layout.Flex{
											Alignment: layout.End,
										}.Layout(gtx,
											layout.Rigid(func() {
												out := th.Body1(t.Out)
												out.TextSize = unit.Dp(fontSize)
												out.Color = textColor
												out.Layout(gtx)
											}),
										)
									})
								}),
								layout.Rigid(func() {
									pwd, _ := exec.Command("pwd").Output()
									layout.Flex{}.Layout(gtx,
										layout.Rigid(func() {
											p := th.Body1(out(pwd))
											p.Font.Style = text.Regular
											p.TextSize = unit.Dp(fontSize)
											p.Color = textColor
											p.Layout(gtx)
										}),
										layout.Rigid(func() {
											input := th.Editor("")
											input.Font.Style = text.Regular
											input.TextSize = unit.Dp(fontSize)
											input.Color = textColor
											input.Layout(gtx, consoleInputField)
										}),
									)
									for _, e := range consoleInputField.Events(gtx) {
										if e, ok := e.(widget.SubmitEvent); ok {
											splitted := strings.Split(e.Text, " ")
											cmd, _ := exec.Command(splitted[0], splitted[1:]...).CombinedOutput()
											com.Commands = append(com.Commands, mod.Command{
												ComID: e.Text,
												Time:  time.Time{},
												Out:   out(cmd),
											})
											consoleInputField.SetText("")
										}
									}
								}))
						})
					}),
				)
				e.Frame(gtx.Ops)
			}
		}
	}()
	app.Main()
}

func out(i []byte) string {
	return "[" + string(i) + "]$"
}

func fill(gtx *layout.Context) {
	cs := gtx.Constraints
	d := image.Point{X: cs.Width.Min, Y: cs.Height.Min}
	dr := f32.Rectangle{
		Max: f32.Point{X: float32(d.X), Y: float32(d.Y)},
	}
	paint.ColorOp{Color: color.RGBA{A: 0xff, R: 0x30, G: 0x30, B: 0x30}}.Add(gtx.Ops)
	paint.PaintOp{Rect: dr}.Add(gtx.Ops)
	gtx.Dimensions = layout.Dimensions{Size: d}
}
