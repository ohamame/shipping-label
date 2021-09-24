package label

import (
	"io"

	"github.com/signintech/gopdf"
)

func CreateShippingLabelPdf(w io.Writer) error {
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
