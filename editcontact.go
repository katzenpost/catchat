package main

import (
	"bytes"
	"gioui.org/gesture"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/png"
)

// EditContactPage is the page for adding a new contact
type EditContactPage struct {
	a        *App
	nickname string
	back     *widget.Clickable
	avatar   *gesture.Click
	clear    *widget.Clickable
	expiry   *widget.Clickable
	rename   *widget.Clickable
	remove   *widget.Clickable
	settings *layout.List
	widgets  []layout.Widget
	//avatar // select an avatar image
}

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
			ct := Contactal{}
			ct.Reset()
			sz := image.Point{X: gtx.Px(unit.Dp(96)), Y: gtx.Px(unit.Dp(96))}
			i := ct.Render(sz)
			b := new(bytes.Buffer)
			if err := png.Encode(b, i); err == nil {
				p.a.c.AddBlob("avatar://"+p.nickname, b.Bytes())
				return RedrawEvent{}
			}
		}
	}
	if p.clear.Clicked() {
		// TODO: confirmation dialog
		p.a.c.WipeConversation(p.nickname)
		return EditContactComplete{nickname: p.nickname}
	}
	if p.expiry.Clicked() {
		// TODO: add message expiry configuration to catshadow
	}
	if p.rename.Clicked() {
		return RenameContact{nickname: p.nickname}
	}
	if p.remove.Clicked() {
		// TODO: confirmation dialog
		p.a.c.RemoveContact(p.nickname)
		return EditContactComplete{nickname: p.nickname}
	}
	return nil
}

func (p *EditContactPage) Start(stop <-chan struct{}) {
}

func newEditContactPage(a *App, contact string) *EditContactPage {
	p := &EditContactPage{a: a, nickname: contact, back: &widget.Clickable{},
		avatar: &gesture.Click{}, clear: &widget.Clickable{},
		expiry: &widget.Clickable{}, rename: &widget.Clickable{},
		remove:   &widget.Clickable{},
		settings: &layout.List{Axis: layout.Vertical},
	}
	p.widgets = []layout.Widget{
		func(gtx C) D {
			dims := layout.Center.Layout(gtx, func(gtx C) D {
				gtx.Constraints.Max.X = gtx.Constraints.Max.X / 4
				return layoutAvatar(gtx, p.a.c.GetContacts()[p.nickname])
			})
			a := pointer.Rect(image.Rectangle{Max: dims.Size})
			t := a.Push(gtx.Ops)
			p.avatar.Add(gtx.Ops)
			t.Pop()
			return dims
		},
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		material.Button(th, p.clear, "Clear History").Layout,
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		material.Button(th, p.rename, "Rename Contact").Layout,
		layout.Spacer{Height: unit.Dp(8)}.Layout,
		material.Button(th, p.remove, "Delete Contact").Layout,
	}
	return p
}
