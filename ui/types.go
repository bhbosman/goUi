package ui

import (
	"github.com/rivo/tview"
	"io"
)

type IPrimitiveCloser interface {
	tview.Primitive
	io.Closer
	UpdateContent() error
}

type SlideCallback func(nextSlide func()) (string, IPrimitiveCloser, error)

type ISlideFactory interface {
	Title() string
	Content() SlideCallback
	OrderNumber() int
}
