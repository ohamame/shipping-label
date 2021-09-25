package label

import (
	"fmt"
	"io"
	"math"

	"github.com/signintech/gopdf"
)

type Label struct {
	ColumnCount  int
	RowCount     int
	PageSize     gopdf.Rect
	FontSize     float64
	CellXPadding float64
	CellYPadding float64
	LineHeight   float64
}

type LabelContent struct {
	OrderNumber string
	Name        string
	Address     string
	City        string
	PhoneNumber string
	PostCode    string
}

func NewLabel(columnCount int, rowCount int, pageSize gopdf.Rect, fontSize float64) Label {
	return Label{ColumnCount: columnCount, RowCount: rowCount, PageSize: pageSize, FontSize: fontSize, CellXPadding: 20, CellYPadding: 15, LineHeight: 12}
}

func (l Label) CreateShippingLabelPdf(w io.Writer, contents []LabelContent) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: l.PageSize})
	pdf.AddPage()
	err := pdf.AddTTFFont("opensans", "./fonts/opensans.ttf")
	if err != nil {
		return err
	}

	err = pdf.SetFont("opensans", "", l.FontSize)
	if err != nil {
		return err
	}

	columnWidth := l.PageSize.W / float64(l.ColumnCount)
	rowHeight := l.PageSize.H / float64(l.RowCount)
	textWidth := columnWidth - (2 * l.CellXPadding)
	fmt.Printf("Page size: %f x %f\nCol Width:%f\nRow Height:%f\nText Width: %f\n", l.PageSize.H, l.PageSize.W, columnWidth, rowHeight, textWidth)

	// TODO: paging
	for i, c := range contents {
		position := i // TODO: allow offset
		row := int(math.Floor(float64(position) / float64(l.ColumnCount)))
		column := position % l.ColumnCount
		startX := columnWidth*float64(column) + l.CellXPadding
		startY := rowHeight*float64(row) + l.CellYPadding

		fmt.Printf("Position: %d, row: %d, column: %d\n", position, row, column)
		lines, err := pdf.SplitText(c.Address, textWidth)
		if err != nil {
			return err
		}

		// Output each line as text
		for i, line := range lines {
			lineStartY := startY + (l.LineHeight * float64(i))
			pdf.SetX(startX)
			pdf.SetY(lineStartY)

			err = pdf.Cell(nil, line)
			if err != nil {
				return err
			}
		}
	}

	err = pdf.Write(w)
	if err != nil {
		return err
	}
	return nil
}
