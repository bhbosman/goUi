package GoFunctionCounterSlide

import (
	"github.com/rivo/tview"
	"strconv"
)

type tablePlate struct {
	emptyCell *tview.TableCell
	data      []string
}

func newTablePlate(data []string) *tablePlate {
	return &tablePlate{
		emptyCell: tview.NewTableCell(""),
		data:      data,
	}
}

func (self *tablePlate) GetCell(row, column int) *tview.TableCell {
	if row == -1 || column == -1 {
		return tview.NewTableCell("")
	}

	switch row {
	case 0:
		switch column {
		case 0:
			return tview.NewTableCell("*").SetSelectable(false)
		case 1:
			return tview.NewTableCell("Go Function").SetSelectable(false)
		}
	default:
		switch column {
		case 0:
			return tview.NewTableCell(strconv.Itoa(row)).SetSelectable(false)
		case 1:
			index := row - 1
			count := len(self.data)
			if index < count {
				return tview.NewTableCell(self.data[row-1])
			}
		}

	}
	return self.emptyCell
}

func (self *tablePlate) GetRowCount() int {
	return len(self.data) + 1
}

func (self *tablePlate) GetColumnCount() int {
	return 2
}

func (self *tablePlate) SetCell(row, column int, cell *tview.TableCell) {
}

func (self *tablePlate) RemoveRow(row int) {
}

func (self *tablePlate) RemoveColumn(column int) {
}

func (self *tablePlate) InsertRow(row int) {
}

func (self *tablePlate) InsertColumn(column int) {
}

func (self *tablePlate) Clear() {
}
