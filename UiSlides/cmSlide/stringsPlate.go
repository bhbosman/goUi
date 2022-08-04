package cmSlide

import (
	"github.com/bhbosman/gocommon/model"
	"github.com/rivo/tview"
	"strconv"
)

type stringsPlate struct {
	data []model.KeyValue
}

func newStringsPlate(data []model.KeyValue) *stringsPlate {
	return &stringsPlate{data: data}
}

func (self *stringsPlate) GetCell(row, column int) *tview.TableCell {
	if row == -1 || column == -1 {
		return emptyCell
	}

	switch row {
	case 0:
		switch column {
		case 0:
			return tview.NewTableCell("*").SetSelectable(false)
		case 1:
			return tview.NewTableCell("Key").SetSelectable(false)
		case 2:
			return tview.NewTableCell("Value").SetSelectable(false)
		}
	default:
		switch column {
		case 0:
			return tview.NewTableCell(strconv.Itoa(row)).SetSelectable(false).SetAlign(tview.AlignRight)
		case 1:
			n := row - 1
			c := len(self.data)
			if c > n {
				return tview.NewTableCell(self.data[row-1].Key)
			}
			break
		case 2:
			n := row - 1
			c := len(self.data)
			if c > n {
				return tview.NewTableCell(self.data[row-1].Value)
			}
			break
		}

	}
	return emptyCell
}

func (self *stringsPlate) GetRowCount() int {
	return len(self.data) + 1
}

func (self *stringsPlate) GetColumnCount() int {
	return 3
}

func (self *stringsPlate) SetCell(row, column int, cell *tview.TableCell) {
}

func (self *stringsPlate) RemoveRow(row int) {
}

func (self *stringsPlate) RemoveColumn(column int) {
}

func (self *stringsPlate) InsertRow(row int) {
}

func (self *stringsPlate) InsertColumn(column int) {
}

func (self *stringsPlate) Clear() {
}

func (self *stringsPlate) GetItem(row int) (model.KeyValue, bool) {

	if row == -1 {
		return model.KeyValue{}, false
	}

	index := row - 1
	count := len(self.data)
	if index >= 0 && count > index {
		return self.data[index], true
	}
	return model.KeyValue{}, false
}
