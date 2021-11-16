package main

import (
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/hako/durafmt"
	"image"
	"math"
	"time"
)

// EditContactPage is the page for adding a new contact
type EditContactPage struct {
	a        *App
	nickname string
	back     *widget.Clickable
	apply    *widget.Clickable
	avatar   *gesture.Click
	clear    *widget.Clickable
	expiry   *widget.Float
	rename   *widget.Clickable
	remove   *widget.Clickable
	settings *layout.List
	widgets  []layout.Widget
	duration time.Duration
}

const (
	minExpiration = 0.0  // never delete messages
	maxExpiration = 14.0 // 2 weeks
)

// Layout returns the contact options menu
func (p *EditContactPage) Layout(gtx layout.Context) layout.Dimensions {
	bg := Background{
		Color: th.Bg,
		Inset: layout.Inset{},
	}

	return bg.Layout(gtx, func(gtx C) D {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd, Alignment: layout.Start}.Layout(gtx,
			// topbar
			layout.Rigid(func(gtx C) D {
				return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(button(th, p.back, backIcon).Layout),
					layout.Flexed(1, fill{th.Bg}.Layout),
					layout.Rigid(material.H6(th, "Edit Contact").Layout),
					layout.Flexed(1, fill{th.Bg}.Layout))
			}),
			// settings list
			layout.Flexed(1, func(gtx C) D {
				in := layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(12), Right: unit.Dp(12)}
				return in.Layout(gtx, func(gtx C) D {
					return p.settings.Layout(gtx, len(p.widgets), func(gtx C, i int) layout.Dimensions {
						// Layout the widget in the list. can wrap and inset, etc, here...
						return p.widgets[i](gtx)
					})
				})
			}),
		)
	})
}

type EditContactComplete struct {
	nickname string
}

type ChooseAvatar struct {
	nickname string
}

type RenameContact struct {
	nickname string
}

// Event catches the widget submit events and calls catshadow.NewContact
func (p *EditContactPage) Event(gtx layout.Context) interface{} {
	if p.back.Clicked() {
		return BackEvent{}
	}
	for _, e := range p.avatar.Events(gtx.Queue) {
		if e.Type == gesture.TypeClick {
			return ChooseAvatar{nickname: p.nickname}
		}
	}
	if p.clear.Clicked() {
		// TODO: confirmation dialog
		p.a.c.WipeConversation(p.nickname)
		return EditContactComplete{nickname: p.nickname}
	}
	if p.expiry.Changed() {
		p.expiry.Value = float32(math.Round(float64(p.expiry.Value)))
	}
	// update duration
	p.duration = time.Duration(int64(p.expiry.Value)) * time.Minute * 60 * 24
	if p.rename.Clicked() {
		return RenameContact{nickname: p.nickname}
	}
	if p.remove.Clicked() {
		// TODO: confirmation dialog
		p.a.c.RemoveContact(p.nickname)
		return EditContactComplete{nickname: p.nickname}
	}
	if p.apply.Clicked() {
		p.a.c.ChangeExpiration(p.nickname, p.duration)
		return BackEvent{}
	}
	return nil
}

func (p *EditContactPage) Start(stop <-chan struct{}) {
}

func newEditContactPage(a *App, contact string) *EditContactPage {
	expiry, _ := a.c.GetExpiration(contact)
	p := &EditContactPage{a: a, nickname: contact, back: &widget.Clickable{},
		avatar: &gesture.Click{}, clear: &widget.Clickable{},
		expiry: &widget.Float{}, rename: &widget.Clickable{},
		remove: &widget.Clickable{}, apply: &widget.Clickable{},
		settings: &layout.List{Axis: layout.Vertical},
	}
	p.expiry.Value = float32(math.Round(float64(expiry) / float64(time.Minute*60*24)))
	p.widgets = []layout.Widget{
		func(gtx C) D {
			dims := layout.Center.Layout(gtx, func(gtx C) D {
				return layoutAvatar(gtx, p.a.c, p.nickname)
			})
			a := pointer.Rect(image.Rectangle{Max: dims.Size})
			t := a.Push(gtx.Ops)
			p.avatar.Add(gtx.Ops)
			t.Pop()
			return dims
		},
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		func(gtx C) D {
			var label string
			if p.expiry.Value < 1.0 {
				label = "Delete after: never"
			} else {
				label = "Delete after: " + durafmt.Parse(p.duration).Format(units)
			}
			return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle, Spacing: layout.SpaceBetween}.Layout(gtx,
				layout.Rigid(material.Body2(th, "Message deletion").Layout),
				layout.Rigid(material.Slider(th, p.expiry, minExpiration, maxExpiration).Layout),
				layout.Rigid(material.Caption(th, label).Layout),
			)
		},
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		material.Button(th, p.clear, "Clear History").Layout,
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		material.Button(th, p.rename, "Rename Contact").Layout,
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		material.Button(th, p.remove, "Delete Contact").Layout,
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		material.Button(th, p.apply, "Apply Changes").Layout,
	}
	return p
}
