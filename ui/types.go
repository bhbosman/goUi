package ui

import (
	"github.com/rivo/tview"
	"io"
)

type IPrimitiveCloser interface {
	tview.Primitive
	io.Closer
}

type SlideCallback func(nextSlide func()) (string, IPrimitiveCloser)

type ISlideFactory interface {
	Title() string
	Content() SlideCallback
	OrderNumber() int
}
