package label

import (
	"io"

	"github.com/signintech/gopdf"
)

type Label struct {
	ColumnCount int
	RowCount    int
	PageSize    gopdf.Rect
	FontSize    float64
}

func NewLabel(columnCount int, rowCount int, pageSize gopdf.Rect, fontSize float64) Label {
	return Label{ColumnCount: columnCount, RowCount: rowCount, PageSize: pageSize, FontSize: fontSize}
}

func (l Label) CreateShippingLabelPdf(w io.Writer) error {
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

	err = pdf.Cell(nil, "hello world")
	if err != nil {
		return err
	}

	err = pdf.Write(w)
	if err != nil {
		return err
	}
	return nil
}
