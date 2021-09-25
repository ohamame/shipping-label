package label

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/mitchellh/go-wordwrap"
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
	OrderNumber  string `csv:"Order Number"`
	Name         string `csv:"Full Name (Shipping)"`
	Address      string `csv:"Address (Shipping)"`
	City         string `csv:"City (Shipping)"`
	PhoneNumber  string `csv:"Phone (Billing)"`
	PostCode     string `csv:"Postcode (Shipping)"`
	CustomerNote string `csv:"Customer Note"`
}

func (c LabelContent) GetText() string {
	lines := []string{}
	if c.Name != "" {
		lines = append(lines, c.Name)
	}
	if c.Address != "" {
		lines = append(lines, c.Address)
	}
	if c.City != "" {
		lines = append(lines, c.City)
	}
	if c.PostCode != "" {
		lines = append(lines, c.PostCode)
	}
	if c.PhoneNumber != "" {
		lines = append(lines, c.PhoneNumber)
	}
	if c.CustomerNote != "" {
		lines = append(lines, "Customer Note:\n"+c.CustomerNote)
	}
	return strings.Join(lines, "\n")
}

func NewLabel(columnCount int, rowCount int, pageSize gopdf.Rect, fontSize float64) Label {
	return Label{ColumnCount: columnCount, RowCount: rowCount, PageSize: pageSize, FontSize: fontSize, CellXPadding: 20, CellYPadding: 20, LineHeight: 14}
}

func (l Label) CreateShippingLabelPdf(w io.Writer, contents []LabelContent) error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: l.PageSize})

	// sort by order number
	sort.Slice(contents, func(i, j int) bool {
		num1, err := strconv.ParseUint(contents[i].OrderNumber, 10, 32)
		if err != nil {
			fmt.Println(err)
			return false
		}
		num2, err := strconv.ParseUint(contents[j].OrderNumber, 10, 32)
		if err != nil {
			fmt.Println(err)
			return false
		}
		return num1 < num2
	})

	labelsPerPage := l.ColumnCount * l.RowCount
	totalPages := int(math.Ceil(float64(len(contents)) / float64(labelsPerPage)))

	for page := 0; page < totalPages; page++ {
		startIndex := page * labelsPerPage
		endIndex := int(math.Min(float64(page*labelsPerPage+labelsPerPage), float64(len(contents))))
		pageContents := contents[startIndex:endIndex]
		fmt.Printf("Page size: %d", len(pageContents))

		pdf.AddPage()
		err := pdf.AddTTFFont("opensans", "./fonts/opensans.ttf")
		if err != nil {
			return err
		}

		err = pdf.SetFont("opensans", "", l.FontSize)
		if err != nil {
			return err
		}

		err = l.drawGrids(&pdf)
		if err != nil {
			return err
		}

		columnWidth := l.PageSize.W / float64(l.ColumnCount)
		rowHeight := l.PageSize.H / float64(l.RowCount)
		textWidth := (columnWidth - (2 * l.CellXPadding))
		fmt.Printf("Page size: %f x %f\nCol Width:%f\nRow Height:%f\nText Width: %f\n", l.PageSize.H, l.PageSize.W, columnWidth, rowHeight, textWidth)

		for i, c := range pageContents {
			position := i // TODO: allow offset
			row := int(math.Floor(float64(position) / float64(l.ColumnCount)))
			column := position % l.ColumnCount
			startX := columnWidth * float64(column)
			startY := rowHeight * float64(row)

			fmt.Printf("Position: %d, row: %d, column: %d\n", position, row, column)

			// Note: Wraps word by a character limit, not particularly accurate as
			// not all characters are equal.
			wrappedText := wordwrap.WrapString(c.GetText(), 40)
			lines, err := pdf.SplitText(wrappedText, textWidth)
			if err != nil {
				return err
			}

			// Order Number at top left
			pdf.SetX(startX + l.CellXPadding)
			pdf.SetY(startY + l.CellYPadding)
			pdf.Cell(nil, "Order ")
			pdf.Cell(nil, c.OrderNumber)
			pdf.SetX(startX + l.CellXPadding)
			pdf.SetY(startY + l.CellYPadding + l.LineHeight)
			pdf.Cell(nil, "Thank you for your order!")

			// Logo at top right
			imageSize := &gopdf.Rect{H: 60, W: 60}
			logoX := startX + columnWidth - l.CellXPadding - imageSize.W
			logoY := startY + l.CellYPadding
			pdf.Image("./images/logo.png", logoX, logoY, imageSize)

			// Output each line as text at bottom left
			for i, line := range lines {
				lineStartX := startX + l.CellXPadding
				lineStartY := startY + rowHeight - l.CellYPadding - (l.LineHeight * float64(len(lines)-i))

				pdf.SetX(lineStartX)
				pdf.SetY(lineStartY)

				err = pdf.Cell(nil, line)
				if err != nil {
					return err
				}
			}
		}
	}

	err := pdf.Write(w)
	if err != nil {
		return err
	}
	return nil
}

func (l Label) drawGrids(pdf *gopdf.GoPdf) error {
	pdf.SetLineWidth(0.5)
	pdf.SetLineType("dashed")

	columnWidth := l.PageSize.W / float64(l.ColumnCount)
	rowHeight := l.PageSize.H / float64(l.RowCount)

	// draw vertical lines
	for i := 0; i < l.ColumnCount; i++ {
		x := float64(columnWidth) * float64(i+1)
		pdf.Line(x, 0, x, l.PageSize.H)
	}

	// draw horizontal lines
	for i := 0; i < l.RowCount; i++ {
		y := float64(rowHeight) * float64(i+1)
		pdf.Line(0, y, l.PageSize.W, y)
	}

	return nil
}
