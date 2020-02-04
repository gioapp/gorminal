package main

import (
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gop9/gorminal/mod"
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
												sat := th.Body1("ds://" + t.ComID + "_" + t.Out)
												sat.Font.Size = unit.Dp(16)
												sat.Layout(gtx)
											}),
										)
									})
								}),
								layout.Rigid(func() {
									layout.UniformInset(unit.Dp(8)).Layout(gtx, func() {
										e := th.Editor("Run command")
										e.Font.Style = text.Regular
										e.Font.Size = unit.Dp(16)
										e.Layout(gtx, consoleInputField)
										for _, e := range consoleInputField.Events(gtx) {
											if e, ok := e.(widget.SubmitEvent); ok {

												splitted := strings.Split(e.Text, " ")

												cmd, _ := exec.Command(splitted[0], splitted[1:]...).Output()

												com.Commands = append(com.Commands, mod.Command{
													ComID: e.Text,
													Time:  time.Time{},
													Out:   string(cmd),
												})
												consoleInputField.SetText("")
											}
										}
									})
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
