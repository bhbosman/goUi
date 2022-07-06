package cmSlide

import (
	"github.com/rivo/tview"
)

type connectionPlate struct {
	list []IdAndName
}

func newConnectionPlace(list []IdAndName) *connectionPlate {
	return &connectionPlate{
		list: list,
	}
}

var emptyCell *tview.TableCell = tview.NewTableCell("").SetSelectable(false)

func (self *connectionPlate) GetCell(row, column int) *tview.TableCell {
	switch row {
	case 0:
		switch column {
		case 0:
			return tview.NewTableCell("*").SetSelectable(false)
		case 1:
			return tview.NewTableCell("Connection Name").SetSelectable(false)
		}
	default:
		switch column {
		case 1:
			n := row - 1
			c := len(self.list)
			if c > n {
				return tview.NewTableCell(self.list[row-1].Id)
			}
		}
	}
	return emptyCell
}

func (self *connectionPlate) GetRowCount() int {
	return len(self.list) + 1
}

func (self *connectionPlate) GetColumnCount() int {
	return 2
}

func (self *connectionPlate) SetCell(row, column int, cell *tview.TableCell) {
}

func (self *connectionPlate) RemoveRow(row int) {
}

func (self *connectionPlate) RemoveColumn(column int) {
}

func (self *connectionPlate) InsertRow(row int) {
}

func (self *connectionPlate) InsertColumn(column int) {
}

func (self *connectionPlate) Clear() {
}

func (self *connectionPlate) GetItem(row int) (IdAndName, bool) {
	index := row - 1
	count := len(self.list)
	if index >= 0 && count > index {
		return self.list[index], true
	}
	return IdAndName{}, false
}
