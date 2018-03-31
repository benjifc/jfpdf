package PdfReport

import (
	"log"

	"github.com/Jeffail/gabs"
	"github.com/fatih/color"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/barcode"
)

var (
	red     = color.New(color.FgRed)
	boldRed = red.Add(color.Bold)
)

func report(jDocPdf []byte) {
	jDoc, err := gabs.ParseJSON(jDocPdf)
	if err != nil {
		boldRed.Println("Failed to parse: %v", error.Error)
		//log.Fatal("Failed to parse.")
		return
	}
	if !jDoc.Exists("pageConfig") {
		log.Fatal("Failed Config Page Not Exist.")
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.Cell(40, 10, "codebar")
	key := barcode.RegisterCode128(pdf, "codabar")
	var width float64 = 50
	var height float64 = 7
	barcode.BarcodeUnscalable(pdf, key, 25, 25, &width, &height, false)

	pdf.SetFont("Arial", "B", 16)
	pdf.OutputFileAndClose("hello.pdf")

}

/*
func addPDF(fileName string, txt string) {


	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	pdf.Output(w)
	pdf.Close()
	w.Flush()
	wr.Header().Set("Content-Type", "application/pdf")

	fmt.Println("%s.pdf", fileName)

}
*/
