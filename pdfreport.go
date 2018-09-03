package pdfreport

import (
	"strings"

	"strconv"

	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/fatih/color"
	"github.com/jung-kurt/gofpdf"
	"github.com/jung-kurt/gofpdf/contrib/barcode"
)

var (
	red         = color.New(color.FgRed)
	boldRed     = red.Add(color.Bold)
	blue        = color.New(color.FgBlue)
	boldBlue    = blue.Add(color.Bold)
	orientation = "P"
	unit        = "mm"
	size        = "A4"
	Wd          = 0.0
	Ht          = 0.0
	ok          bool
	pdf         *gofpdf.Fpdf
)

func Report(jDocPdf []byte) {

	//boldBlue.Printf(string(jDocPdf[:]) + "\n")
	jDoc, err := gabs.ParseJSON(jDocPdf)
	if err != nil {

		boldRed.Printf("Failed to parse: %s\n", err.Error())
		message("Error JSon parser.!!!!")
		return
	}

	PageConfig(jDoc)

	children, err := jDoc.S("report", "pages").Children()
	if err != nil {
		message(err.Error())
		return
	}
	if len(children) <= 0 {
		message("There are not pages.!!!")
		return
	}
	for _, child := range children {
		Page(child)
	}
	children = nil
}

/***************************************************************************************/
/*
/*			FUNCTION PDF
/*
/***************************************************************************************/
func PageConfig(jDoc *gabs.Container) {
	if !jDoc.Exists("report", "pageconfig") {
		message("Not exist page config.!!!!")
		return
	}
	if jDoc.Exists("report", "pageconfig", "orientation") {
		orientation, ok = jDoc.Path("report.pageconfig.orientation").Data().(string)
	}
	if jDoc.Exists("report", "pageconfig", "unit") {
		unit, ok = jDoc.Path("report.pageconfig.unit").Data().(string)
	}
	if jDoc.Exists("report", "pageconfig", "size") {
		size, ok = jDoc.Path("report.pageconfig.size").Data().(string)
	}
	if jDoc.Exists("report", "pageconfig", "Wd") {
		Wd, ok = jDoc.Path("report.pageconfig.Wd").Data().(float64)
	}
	if jDoc.Exists("report", "pageconfig", "Ht") {
		Ht, ok = jDoc.Path("report.pageconfig.Ht").Data().(float64)
	}

	//pdf = gofpdf.New(orientation, unit, size, "")
	pdf = gofpdf.NewCustom(&gofpdf.InitType{
		OrientationStr: orientation,
		UnitStr:        unit,
		Size:           gofpdf.SizeType{Wd: Wd, Ht: Ht},
	})
}

func Page(pag *gabs.Container) {
	if !pag.Exists("content") {
		log.Fatal("Failed Config Page Not Exist.")
		return
	}

	//pdf.AddPageFormat(orientation, gofpdf.SizeType{Wd: Wd, Ht: Ht})
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)

	children, err := pag.S("content").Children()
	if err != nil {
		message(err.Error())
		return
	}
	for child := range children {
		var item *gabs.Container
		item = children[child]
		//boldRed.Printf(item.String())
		strItem := string(item.String())
		key := []rune(strItem)
		i := strings.Index(strItem, ":")
		strKey := string(key[2 : i-1])
		switch strKey {
		case "barcode128":
			var (
				x    float64
				y    float64
				w    float64
				h    float64
				code string
			)
			x, ok = children[child].Path(fmt.Sprintf("%s.x", strKey)).Data().(float64)
			y, ok = children[child].Path(fmt.Sprintf("%s.y", strKey)).Data().(float64)
			w, ok = children[child].Path(fmt.Sprintf("%s.w", strKey)).Data().(float64)
			h, ok = children[child].Path(fmt.Sprintf("%s.h", strKey)).Data().(float64)
			code, ok = children[child].Path(fmt.Sprintf("%s.code", strKey)).Data().(string)
			Code128(x, y, w, h, code)
		case "write":
			var (
				x    float64
				y    float64
				text string
			)
			x, ok = children[child].Path(fmt.Sprintf("%s.x", strKey)).Data().(float64)
			y, ok = children[child].Path(fmt.Sprintf("%s.y", strKey)).Data().(float64)
			text, ok = children[child].Path(fmt.Sprintf("%s.text", strKey)).Data().(string)
			pdf.Text(x, y, text)
		case "cell":
			var (
				x    float64
				y    float64
				w    float64
				h    float64
				text string
				//borderStr string
				//ln int
				alignStr string
			)
			x, ok = children[child].Path(fmt.Sprintf("%s.x", strKey)).Data().(float64)
			y, ok = children[child].Path(fmt.Sprintf("%s.y", strKey)).Data().(float64)
			w, ok = children[child].Path(fmt.Sprintf("%s.w", strKey)).Data().(float64)
			h, ok = children[child].Path(fmt.Sprintf("%s.h", strKey)).Data().(float64)
			text, ok = children[child].Path(fmt.Sprintf("%s.text", strKey)).Data().(string)
			alignStr, ok = children[child].Path(fmt.Sprintf("%s.aligned", strKey)).Data().(string)
			pdf.SetXY(x,y)
			pdf.CellFormat(w, h, text,"",0,alignStr,false,0,"")

		case "font":

			name := children[child].Path(fmt.Sprintf("%s.name", strKey)).Data().(string)
			style := children[child].Path(fmt.Sprintf("%s.style", strKey)).Data().(string)
			size := children[child].Path(fmt.Sprintf("%s.size", strKey)).Data().(float64)
			pdf.SetFont(name, style, size)

		default:

		}

		//pdf.CellFormat(100, 10, strconv.Itoa(child) , "RLTB", 0, "C", false, 0, "")
	}

	children = nil
	pag = nil

}

// BARCODES
func Code128(x float64, y float64, w float64, h float64, text string) {
	key := barcode.RegisterCode128(pdf, text)
	width := w
	height := h
	barcode.BarcodeUnscalable(pdf, key, x, y, &width, &height, false)
}

/***************************************************************************************/
/*
/*			FUNCTION SAVE UTILS
/*
/***************************************************************************************/

func message(text string) {
	pdf = gofpdf.New(orientation, unit, size, "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 20)
	pdf.CellFormat(100, 10, text, "RLTB", 0, "C", false, 0, "")

}

/***************************************************************************************/
/*
/*			FUNCTION SAVE FILE
/*
/***************************************************************************************/

func Save() {
	err := pdf.OutputFileAndClose("report.pdf")
	if err != nil {
		boldRed.Printf("Failed to Save: %s\n", err.Error())
		//log.Fatal("Failed to parse.")
		return
	}

}
func Write(w http.ResponseWriter) {

	X := new(bytes.Buffer)

	pdf.Output(X)
	pdf.Close()
	_, err := X.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Pragma", "public")
	w.Header().Set("Expires", "0")
	w.Header().Set("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Cache-Control", "private") // required for certain browsers
	w.Header().Set("Content-Type", "application/pdf; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=\"report.pdf\";")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(len(X.Bytes())))

}
