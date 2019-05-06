package main

import "github.com/jung-kurt/gofpdf"
import (
	"bytes"
	"strings"
	"io/ioutil"
)
import (
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

import (
    "bufio"
    "fmt"
    "os"
)

const bg_image = "background.jpg"


func add_background(pdf *gofpdf.Fpdf) {
	cert_image, _ := assetsCertificateJpeg()
	r_image := bytes.NewReader(cert_image.bytes)
	pdf.RegisterImageReader(bg_image, pdf.ImageTypeFromMime("image/jpeg"), r_image)
	pagew, pageh := pdf.GetPageSize()
	pdf.Image(bg_image, 0, 0, pagew, pageh, false, "", 0, "")
}

func load_fonts(pdf *gofpdf.Fpdf) {
	font_sans_json, _ := assetsPtsansJsonBytes()
	font_sans_z, _ := assetsPtsansZBytes()
	pdf.AddFontFromBytes("PTSans", "", font_sans_json, font_sans_z)
	
	font_serif_json, _ := assetsPtserifJsonBytes()
	font_serif_z, _ := assetsPtserifZBytes()
	pdf.AddFontFromBytes("PTSerif", "", font_serif_json, font_serif_z)
}

func to_cp1251(text string) string {
	sr := strings.NewReader(text)
	tr := transform.NewReader(sr, charmap.Windows1251.NewEncoder())
	buf, err := ioutil.ReadAll(tr)
	if err != nil {
		return "???"
	}
	return string(buf)
}

func write_main(pdf *gofpdf.Fpdf, text string, row int) {
	page_width, _ := pdf.GetPageSize()
	pdf.SetFont("PTSerif", "", main_font_size)
	width := pdf.GetStringWidth(to_cp1251(text))

	x := (page_width - width) / 2.0
	y := first_line + float64(row) * line_delta
	
	pdf.MoveTo(x, y)
	pdf.Cell(width, 0, to_cp1251(text))
}


func Readln(r *bufio.Reader) (string, error) {
  var (isPrefix bool = true
       err error = nil
       line, ln []byte
      )
  for isPrefix && err == nil {
      line, isPrefix, err = r.ReadLine()
      ln = append(ln, line...)
  }
  return string(ln),err
}

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	add_background(pdf)
	load_fonts(pdf)

    f, err := os.Open("cert.txt")
    if err != nil {
        fmt.Printf("error opening file: %v\n",err)
        os.Exit(1)
    }
    r := bufio.NewReader(f)

	for row := 0; row < 5; row++ {
	    text, _ := Readln(r)
		write_main(pdf, text, row)
	}
    defer f.Close()
	
	_ = pdf.OutputFileAndClose("cert.pdf")
}
