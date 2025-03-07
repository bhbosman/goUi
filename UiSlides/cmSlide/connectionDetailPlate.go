package cmSlide

import (
	"fmt"
	"github.com/bhbosman/gocommon/model"
	"github.com/rivo/tview"
	"strconv"
)

func ByteCountSI(b int64) string {
	return strconv.Itoa(int(b))
	//const unit = 1000
	//if b < unit {
	//	return fmt.Sprintf("%d B", b)
	//}
	//div, exp := int64(unit), 0
	//for n := b / unit; n >= unit; n /= unit {
	//	div *= unit
	//	exp++
	//}
	//return fmt.Sprintf("%.1f %cB",
	//	float64(b)/float64(div), "kMGTPE"[exp])
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

type connectionDetailPlateContent struct {
	Grid []model.LineData
}

func newConnectionPlateContent(Grid []model.LineData) *connectionDetailPlateContent {
	return &connectionDetailPlateContent{
		Grid: Grid,
	}
}

func (self *connectionDetailPlateContent) GetCell(row, column int) *tview.TableCell {
	if row == -1 || column == -1 {
		return tview.NewTableCell("")
	}

	switch column {
	case 0:
		switch row {
		case 0:
			return tview.NewTableCell("*")
		default:

			return tview.NewTableCell("")
		}
	case 1:
		switch row {
		case 0:
			return tview.NewTableCell("Name")
		default:
			return tview.NewTableCell(self.Grid[row-1].InValue.Name)
		}
	case 2:
		switch row {
		case 0:
			return tview.NewTableCell("In(Other)")
		default:
			return tview.NewTableCell(
				fmt.Sprintf("%v",
					ByteCountSI(self.Grid[row-1].InValue.OtherMsgCountIn),
				),
			).
				SetAlign(tview.AlignRight)
		}
	case 3:
		switch row {
		case 0:
			return tview.NewTableCell("In(RWS)")
		default:
			return tview.NewTableCell(
				ByteCountSI(self.Grid[row-1].InValue.RwsMsgCountIn),
			).
				SetAlign(tview.AlignRight)
		}
	case 4:
		switch row {
		case 0:
			return tview.NewTableCell("In(Bytes)")
		default:
			return tview.NewTableCell(ByteCountIEC(self.Grid[row-1].InValue.RwsBytesIn)).
				SetAlign(tview.AlignRight)
		}
	case 5:
		switch row {
		case 0:
			return tview.NewTableCell("Out(Bytes)")
		default:
			return tview.NewTableCell(ByteCountIEC(self.Grid[row-1].InValue.RwsBytesOut)).
				SetAlign(tview.AlignRight)
		}
	case 6:
		switch row {
		case 0:
			return tview.NewTableCell("Name")
		default:
			return tview.NewTableCell(self.Grid[row-1].OutValue.Name)
		}
	case 7:
		switch row {
		case 0:
			return tview.NewTableCell("In(Other)")
		default:
			return tview.NewTableCell(
				ByteCountSI(self.Grid[row-1].OutValue.OtherMsgCountOut),
			).
				SetAlign(tview.AlignRight)
		}
	case 8:
		switch row {
		case 0:
			return tview.NewTableCell("In(RWS)")
		default:
			return tview.NewTableCell(
				ByteCountSI(self.Grid[row-1].OutValue.RwsMsgCountOut),
			).
				SetAlign(tview.AlignRight)
		}
	case 9:
		switch row {
		case 0:
			return tview.NewTableCell("In(Bytes)")
		default:
			return tview.NewTableCell(ByteCountIEC(self.Grid[row-1].OutValue.RwsBytesIn)).
				SetAlign(tview.AlignRight)
		}
	case 10:
		switch row {
		case 0:
			return tview.NewTableCell("Out(Bytes)")
		default:
			return tview.NewTableCell(ByteCountIEC(self.Grid[row-1].OutValue.RwsBytesOut)).
				SetAlign(tview.AlignRight)
		}

	}
	return tview.NewTableCell("")
}

func (self *connectionDetailPlateContent) GetRowCount() int {
	return len(self.Grid) + 1
}

func (self *connectionDetailPlateContent) GetColumnCount() int {
	return 11
}

func (self *connectionDetailPlateContent) SetCell(_, _ int, _ *tview.TableCell) {
}

func (self *connectionDetailPlateContent) RemoveRow(_ int) {
}

func (self *connectionDetailPlateContent) RemoveColumn(_ int) {
}

func (self *connectionDetailPlateContent) InsertRow(_ int) {
}

func (self *connectionDetailPlateContent) InsertColumn(_ int) {
}

func (self *connectionDetailPlateContent) Clear() {
}
