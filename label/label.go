package label

import (
	"github.com/signintech/gopdf"
)

func CreateShippingLabelPdf() error {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()
	err := pdf.AddTTFFont("opensans", "./fonts/opensans.ttf")
	if err != nil {
		return err
	}

	err = pdf.SetFont("opensans", "", 10)
	if err != nil {
		return err
	}
	pdf.Cell(nil, "hello world")
	pdf.WritePdf("output.pdf")
	return nil
}
