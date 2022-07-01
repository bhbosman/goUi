package UiService

import (
	"github.com/bhbosman/goUi/ui"
	"github.com/rivo/tview"
)

type PagePaintToggle struct {
	pages *tview.Pages
	page  string
	item  tview.Primitive
}

func NewPagePaintToggle(pages *tview.Pages) *PagePaintToggle {
	return &PagePaintToggle{
		pages: pages,
	}
}

func (self *PagePaintToggle) SetChangedFunc() {
	page, item := self.pages.GetFrontPage()
	self.setCurrent(page, item)
}

func (self *PagePaintToggle) setCurrent(page string, item tview.Primitive) {
	if self.item != nil {
		if screenDrawToggle, ok := self.item.(ui.IScreenDrawToggle); ok {
			screenDrawToggle.Toggle(false)
		}
	}
	self.page = page
	self.item = item

	if self.item != nil {
		if screenDrawToggle, ok := self.item.(ui.IScreenDrawToggle); ok {
			screenDrawToggle.Toggle(true)
		}
	}
}
