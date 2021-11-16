package main

import (
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/notify"
)

// SettingsPage is for user settings
type SettingsPage struct {
	a                 *App
	back              *widget.Clickable
	submit            *widget.Clickable
	switchUseTor      *widget.Bool
	switchAutoConnect *widget.Bool
}

var (
	inset = layout.UniformInset(unit.Dp(8))
)

const (
	settingNameColumnWidth    = .3
	settingDetailsColumnWidth = 1 - settingNameColumnWidth
)

// Layout returns a simple centered layout prompting to update settings
func (p *SettingsPage) Layout(gtx layout.Context) layout.Dimensions {
	bg := Background{
		Color: th.Bg,
		Inset: layout.Inset{},
	}

	return bg.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical, Alignment: layout.End}.Layout(gtx,
			// topbar
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Baseline}.Layout(gtx,
					layout.Rigid(button(th, p.back, backIcon).Layout),
					layout.Flexed(1, fill{th.Bg}.Layout),
					layout.Rigid(material.H6(th, "Settings").Layout),
					layout.Flexed(1, fill{th.Bg}.Layout))
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(settingNameColumnWidth, func(gtx C) D {
						return inset.Layout(gtx, material.Body1(th, "Use Tor").Layout)
					}),
					layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
						return inset.Layout(gtx, material.Switch(th, p.switchUseTor).Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(settingNameColumnWidth, func(gtx C) D {
						return inset.Layout(gtx, material.Body1(th, "Connect Automatically").Layout)
					}),
					layout.Flexed(settingDetailsColumnWidth, func(gtx C) D {
						return inset.Layout(gtx, material.Switch(th, p.switchAutoConnect).Layout)
					}),
				)
			}),
			layout.Rigid(func(gtx C) D {
				return material.Button(th, p.submit, "Apply Settings").Layout(gtx)
			}),
		)
	})
}

type restartClient struct{}

// Event catches the widget submit events and calls Settings
func (p *SettingsPage) Event(gtx layout.Context) interface{} {
	if p.back.Clicked() {
		return BackEvent{}
	}
	if p.switchUseTor.Changed() {
		if p.switchUseTor.Value && !hasTor() {
			p.switchUseTor.Value = false
			p.a.c.DeleteBlob("UseTor")
			warnNoTor()
			return nil
		}
		if p.switchUseTor.Value {
			p.a.c.AddBlob("UseTor", []byte{1})
		} else {
			p.a.c.DeleteBlob("UseTor")
		}
	}
	if p.switchAutoConnect.Changed() {
		if p.switchAutoConnect.Value {
			p.a.c.AddBlob("AutoConnect", []byte{1})
		} else {
			p.a.c.DeleteBlob("AutoConnect")
		}
	}
	if p.submit.Clicked() {
		go func() {
			if n, err := notify.Push("Restarting", "Catchat is restarting"); err == nil {
				<-time.After(notificationTimeout)
				n.Cancel()
			}
		}()
		p.a.c.Shutdown()
		return restartClient{}
	}
	return nil
}

func (p *SettingsPage) Start(stop <-chan struct{}) {
}

func newSettingsPage(a *App) *SettingsPage {
	p := &SettingsPage{a: a}
	p.back = &widget.Clickable{}
	p.submit = &widget.Clickable{}
	if _, err := a.c.GetBlob("UseTor"); err == nil {
		p.switchUseTor = &widget.Bool{Value: true}
	} else {
		p.switchUseTor = &widget.Bool{Value: false}
	}
	if _, err := a.c.GetBlob("AutoConnect"); err == nil {
		p.switchAutoConnect = &widget.Bool{Value: true}
	} else {
		p.switchAutoConnect = &widget.Bool{Value: false}
	}
	return p
}

func warnNoTor() {
	go func() {
		if n, err := notify.Push("Failure", "Tor requested, but not available on port 9050. Disable in settings to connect."); err == nil {
			<-time.After(notificationTimeout)
			n.Cancel()
		}
	}()
}
