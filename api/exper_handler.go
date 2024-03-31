package api

import (
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jung-kurt/gofpdf"
)

type ExperHandler struct {
}

func NewExperHandler() *ExperHandler {
	return &ExperHandler{}
}

type Detalization struct {
	Msisdn      string
	Date        time.Time
	DurationSec int
	DurationMin int
	Direction   string
	Price       float64
}

type DetalizationData struct {
	Data []Detalization
}

func generateDetalization() DetalizationData {
	res := DetalizationData{}
	for i := 0; i <= 50; i++ {
		d := Detalization{
			Msisdn:      "111911910",
			Date:        time.Now(),
			DurationSec: i + 5,
			DurationMin: i,
			Direction:   "yo",
			Price:       float64(i) / 2,
		}
		res.Data = append(res.Data, d)
	}
	return res
}

func (h *ExperHandler) HandlePeriod(c *fiber.Ctx) error {
	detalizationPeriod := []string{}
	for i := 1; i <= 6; i++ {
		fmt.Println(i)
		previousMonth := time.Now().AddDate(0, -i, 0).Format("01.2006")
		detalizationPeriod = append(detalizationPeriod, previousMonth)
	}
	return c.JSON(detalizationPeriod)
}

func (h *ExperHandler) GeneratePDF(c *fiber.Ctx) error {

	data := generateDetalization()
	// Then we create a new PDF document and write the title and the current date.
	pdf := newReport()
    // And we should take the opportunity and beef up our report with a nice logo.
	pdf = image(pdf)
	// After that, we create the table header and fill the table.
	pdf = header(pdf, []string{"Msisdn", "Date", "DurationSec", "DurationMin", "Direction", "Price"})
	for _, d := range data.Data {
		row := []string{d.Msisdn, d.Date.Format("2006-01-02"), fmt.Sprintf("%d", d.DurationSec), fmt.Sprintf("%d", d.DurationMin), d.Direction, fmt.Sprintf("%.2f", d.Price)}
		pdf = table(pdf, [][]string{row})
	}


	// if pdf.Err() {
	// 	log.Fatalf("Failed creating PDF report: %s\n", pdf.Error())
	// }

	// And finally, we write out our finished record to a file.
	err := savePDF(pdf)

	if err != nil {
		log.Fatalf("Cannot save PDF: %s|n", err)
	}

	// fmt.Println(pdf)
	return c.JSON("33.pdf")
}

// Next, we create a new PDF document.
func newReport() *gofpdf.Fpdf {
	// The package provides a function named `New()` to create a PDF document with
	//
	// * landscape ("L") or portrait ("P") orientation,
	// * the unit used for expressing lengths and sizes ("mm"),
	// * the paper format ("Letter"), and
	// * the path to a font directory.
	//
	// All of these can remain empty, in which case `New()` provides suitable defaults.
	//
	// Function `New()` returns an object of type `*gofpdf.Fpdf` that
	// provides a number of methods for filling the document.
	pdf := gofpdf.New("L", "mm", "Letter", "")

	// We start by adding a new page to the document.
	pdf.AddPage()

	// Now we set the font to "Times", the style to "bold", and the size to 28 points.
	pdf.SetFont("Times", "B", 28)

	// Then we write a text cell of length 40 and height 10. There are no
	// starting coordinates used here; instead, the `Cell()` method moves
	// the current position to the end of the cell so that the next call
	// to `Cell()` continues after the previous cell.
	pdf.Cell(40, 10, "Daily Report")

	// The `Ln()` function moves the current position to a new line, with
	// an optional line height parameter.
	pdf.Ln(12)

	pdf.SetFont("Times", "", 20)
	pdf.Cell(40, 10, time.Now().Format("Mon Jan 2, 2006"))
	pdf.Ln(20)

	return pdf
}

/* ### How Cell() and Ln() advance the output position

As mentioned in the comments, the `Cell()` method takes no coordinates. Instead, the PDF document maintains the current output position internally, and advances it to the right by the length of the cell being written.

Method `Ln()` moves the output position back to the left border and down by the provided value. (Passing `-1` uses the height of the recently written cell.)

HYPE[pdf](pdf.html)
*/

// ## The Table Header: Formatted Cells

// Having created the initial document, we can now create the table header.
// This time, we generate a formatted cell with a light grey as the
// background color.
func header(pdf *gofpdf.Fpdf, hdr []string) *gofpdf.Fpdf {
	pdf.SetFont("Times", "B", 16)
	pdf.SetFillColor(240, 240, 240)
	for _, str := range hdr {
		// The `CellFormat()` method takes a couple of parameters to format
		// the cell. We make use of this to create a visible border around
		// the cell, and to enable the background fill.
		pdf.CellFormat(40, 7, str, "1", 0, "", true, 0, "")
	}

	// Passing `-1` to `Ln()` uses the height of the last printed cell as
	// the line height.
	pdf.Ln(-1)
	return pdf
}

// ## The Table Body

// In the same fashion, we can create the table body.
func table(pdf *gofpdf.Fpdf, tbl [][]string) *gofpdf.Fpdf {
	// Reset font and fill color.
	pdf.SetFont("Times", "", 16)
	pdf.SetFillColor(255, 255, 255)

	// Every column gets aligned according to its contents.
	align := []string{"L", "C", "L", "R", "R", "R"}
	for _, line := range tbl {
		for i, str := range line {
			// Again, we need the `CellFormat()` method to create a visible
			// border around the cell. We also use the `alignStr` parameter
			// here to print the cell content either left-aligned or
			// right-aligned.
			pdf.CellFormat(40, 7, str, "1", 0, align[i], false, 0, "")
		}
		pdf.Ln(-1)
	}

	return pdf
}

// Next, let's not forget to impress our boss by adding a fancy image.
func image(pdf *gofpdf.Fpdf) *gofpdf.Fpdf {
	// The `ImageOptions` method takes a file path, x, y, width, and height
	// parameters, and an `ImageOptions` struct to specify a couple of options.
	pdf.ImageOptions("ee.jpg", 225, 10, 25, 25, false, gofpdf.ImageOptions{ImageType: "JPG", ReadDpi: true}, 0, "")
	return pdf
}

// ## Saving The Document
//
// Finally, the convenience method `OutputFileAndClose()` lets us save the
// finished document.
func savePDF(pdf *gofpdf.Fpdf) error {
	return pdf.OutputFileAndClose("report.pdf")
}
