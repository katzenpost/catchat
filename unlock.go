package main

import (
	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/katzenpost/katzenpost/catshadow"
)

type unlockPage struct {
	result chan interface{}
}

func (p *unlockPage) Layout(gtx layout.Context) layout.Dimensions {
	bg := Background{
		Color: th.Bg,
		Inset: layout.Inset{},
	}

	return bg.Layout(gtx, func(gtx C) D {
		return layout.Center.Layout(gtx, material.Caption(th, "Decrypting statefile...").Layout)
	})
}

func (p *unlockPage) Start(stop <-chan struct{}) {
}

type unlockError struct {
	err error
}

type unlockSuccess struct {
	client *catshadow.Client
}

func (p *unlockPage) Event(gtx layout.Context) interface{} {
	select {
	case r := <-p.result:
		switch r := r.(type) {
		case error:
			return unlockError{err: r}
		case *catshadow.Client:
			return unlockSuccess{client: r}
		}
	default:
	}
	return nil
}

func newUnlockPage(result chan interface{}) *unlockPage {
	p := new(unlockPage)
	p.result = result
	return p
}
