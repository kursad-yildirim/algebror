package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

type config struct {
	difficultyLevel int
	apiPort         int
	rangeMin        int
	rangeMax        int
	opCount         int
	filePath        string
	fileDst         string
}

type operation struct {
	name    string
	sign    string
	numType string
}

type operations []operation

type question struct {
	opString       []string
	opStringMasked []string
}

type questions []question

var c config
var q questions
var ops operations

func main() {
	c.populate()
	ops = operations{
		{name: "addition", sign: "+", numType: "decimal"},
		{name: "subtraction", sign: "-", numType: "decimal"},
		{name: "multiplication", sign: "x", numType: "decimal"},
		{name: "division", sign: "/", numType: "integer"},
		{name: "percentage", sign: "%", numType: "integer"},
	}
	http.HandleFunc("/generate-test", generateTest)
	fmt.Printf("Server started at :%v\n", c.apiPort)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%v", c.apiPort), nil))
}

func (c *config) populate() {
	c.difficultyLevel = 5
	c.apiPort = 18080
	c.rangeMax = 10
	c.opCount = 40
	c.fileDst = "./out/"
}

func generateTest(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	for key, values := range params {
		for _, value := range values {
			if key == "d" {
				d, err := strconv.Atoi(value)
				if err == nil {
					c.difficultyLevel = d
				}
			}
		}
	}
	q = nil
	q.generateOps(c, ops)
	q.alignOps()
	q.createPDF(&c)
	fmt.Println(">>>>>>>>>>>>>>>", c.filePath)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strings.Split(c.filePath, "/")[len(strings.Split(c.filePath, "/"))-1]))
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, c.filePath)
}

func (q *questions) generateOps(c config, ops operations) {
	for i := 0; i < c.opCount; i++ {
		*q = append(*q, generateSingleOp(c, ops[rand.Intn(len(ops))]))
	}
}

func generateSingleOp(c config, op operation) question {
	var numF1, numF2, resultF float64
	var numI1, numI2, resultI int
	var q question
	switch op.name {
	case "addition":
		numF1 = getRandomDecimal(c)
		numF2 = getRandomDecimal(c)
		resultF = numF1 + numF2
		q.opString = []string{formatNumber(numF1, 1), " " + op.sign + " ", formatNumber(numF2, 1), " = ", formatNumber(resultF, 2)}
	case "subtraction":
		numF1 = getRandomDecimal(c)
		numF2 = getRandomDecimal(c)
		resultF = numF1 - numF2
		q.opString = []string{formatNumber(numF1, 1), " " + op.sign + " ", formatNumber(numF2, 1), " = ", formatNumber(resultF, 2)}
	case "multiplication":
		if c.difficultyLevel > 2 {
			numF1 = getRandomDecimal(c)
			numF2 = getRandomDecimal(c)
			resultF = math.Round(numF1*numF2*100) / 100
			q.opString = []string{fmt.Sprintf("%v", numF1), " " + op.sign + " ", fmt.Sprintf("%v", numF2), " = ", fmt.Sprintf("%v", resultF)}
		} else {
			numI1 = getRandomInteger(c)
			numI2 = getRandomInteger(c)
			resultI = numI1 * numI2
			q.opString = []string{fmt.Sprintf("%v", numI1), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", resultI)}
		}
	case "division":
		if c.difficultyLevel > 3 {
			numF1 = getRandomDecimal(c)
			numF2 = getRandomDecimal(c)
			resultF = math.Round(numF1*numF2*100) / 100
			q.opString = []string{fmt.Sprintf("%v", resultF), " " + op.sign + " ", fmt.Sprintf("%v", numF2), " = ", fmt.Sprintf("%v", numF1)}
		} else {
			numI1 = getRandomInteger(c)
			numI2 = getRandomInteger(c)
			resultI = numI1 * numI2
			q.opString = []string{fmt.Sprintf("%v", resultI), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", numI1)}
		}
	case "percentage":
		if c.difficultyLevel > 4 {
			numF1 = getRandomDecimal(c)
			numF2 = getRandomDecimal(c)
			resultF = math.Round(numF1*numF2/100*100) / 100
			q.opString = []string{fmt.Sprintf("%v", numF1), " " + op.sign + " ", fmt.Sprintf("%v", numF2), " = ", fmt.Sprintf("%v", resultF)}
		} else {
			numI1 = getRandomInteger(c)
			numI2 = getRandomInteger(c)
			resultF = math.Round(float64(numI1*numI2)/100*100) / 100
			q.opString = []string{fmt.Sprintf("%v", numI1), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", resultF)}
		}
	}
	q.opStringMasked = make([]string, len(q.opString))
	copy(q.opStringMasked, q.opString)
	maskIndex := rand.Intn(3)
	if maskIndex == 0 {
		q.opStringMasked[0] = "____"
	} else if maskIndex == 1 {
		q.opStringMasked[2] = "____"
	} else {
		q.opStringMasked[4] = "____"
	}
	return q
}

func getRandomInteger(c config) int {
	n := rand.Intn(c.rangeMax * c.difficultyLevel)
	if n == 0 {
		return n + 1
	} else {
		return n
	}
}

func getRandomDecimal(c config) float64 {
	n := math.Round(rand.Float64()*10000) / 10000 * float64(c.rangeMax*c.difficultyLevel)
	if c.difficultyLevel > 4 {
		n = math.Round(n*100) / 100
	} else {
		n = math.Round(n*10) / 10
	}
	if n == 0 {
		return n + 0.1
	} else {
		return n
	}
}

func formatNumber(num float64, p int) string {
	return strconv.FormatFloat(num, 'f', p, 64)
}

func (q *questions) alignOps() {
	col1 := (*q)[:len(*q)/2]
	lenBefore, lenBeforeMasked := 0, 0
	for _, r := range col1 {
		if len(r.opString[0])+len(r.opString[2]) > lenBefore {
			lenBefore = len(r.opString[0]) + len(r.opString[2])
		}
		if len(r.opStringMasked[0])+len(r.opStringMasked[2]) > lenBeforeMasked {
			lenBeforeMasked = len(r.opStringMasked[0]) + len(r.opStringMasked[2])
		}
	}
	for _, r := range col1 {
		r.opString[0] = strings.Repeat(" ", lenBefore-len(r.opString[0])-len(r.opString[2])) + r.opString[0]
		r.opStringMasked[0] = strings.Repeat(" ", lenBeforeMasked-len(r.opStringMasked[0])-len(r.opStringMasked[2])) + r.opStringMasked[0]
	}
	col2 := (*q)[len(*q)/2:]
	lenBefore, lenBeforeMasked = 0, 0
	for _, r := range col2 {
		if len(r.opString[0])+len(r.opString[2]) > lenBefore {
			lenBefore = len(r.opString[0]) + len(r.opString[2])
		}
		if len(r.opStringMasked[0])+len(r.opStringMasked[2]) > lenBeforeMasked {
			lenBeforeMasked = len(r.opStringMasked[0]) + len(r.opStringMasked[2])
		}
	}
	for _, r := range col2 {
		r.opString[0] = strings.Repeat(" ", lenBefore-len(r.opString[0])-len(r.opString[2])) + r.opString[0]
		r.opStringMasked[0] = strings.Repeat(" ", lenBeforeMasked-len(r.opStringMasked[0])-len(r.opStringMasked[2])) + r.opStringMasked[0]
	}
}

func (q questions) createPDF(c *config) {
	const (
		pageWidth    = 210.0
		pageHeight   = 297.0
		margin       = 10.0
		columnWidth  = (pageWidth - 2*margin) / 2
		lineSpacing  = 11.0
		rowsPerPage  = 20
		headerHeight = 20.0
		fontSize     = 12
		fontFace     = "Arial"
		header       = "Name: _______________________________          Time: _____          Score: _____ / 40"
	)
	index := generateRandomIndex()

	filename := fmt.Sprintf("operations_%v_%s.pdf", c.difficultyLevel, index)
	c.filePath = c.fileDst + filename
	pageTitle := fmt.Sprintf("Operations Practice - Level %d - %s", c.difficultyLevel, index)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont(fontFace, "", fontSize)

	pdf.AddPage()
	xLeft := margin
	xRight := margin + columnWidth
	y := margin
	pdf.SetFont(fontFace, "B", fontSize)
	pdf.SetXY(margin, y)
	pdf.CellFormat(pageWidth-2*margin, lineSpacing, pageTitle, "0", 0, "C", false, 0, "")
	y += 10
	pdf.SetFont(fontFace, "", fontSize)
	pdf.SetXY(margin, y)
	pdf.CellFormat(pageWidth-2*margin, lineSpacing, header, "0", 0, "C", false, 0, "")
	y += headerHeight
	for i := 0; i < rowsPerPage; i++ {
		pdf.SetXY(xLeft, y)
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i].opStringMasked, ""), "0", 0, "C", false, 0, "")
		pdf.SetXY(xRight, y)
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i+rowsPerPage].opStringMasked, ""), "0", 0, "C", false, 0, "")
		y += lineSpacing
	}
	pdf.AddPage()
	xLeft = margin
	xRight = margin + columnWidth
	y = margin
	pdf.SetFont(fontFace, "B", fontSize)
	pdf.SetXY(margin, y)
	pdf.CellFormat(pageWidth-2*margin, lineSpacing, "---- Answers ----", "0", 0, "C", false, 0, "")
	pdf.SetFont(fontFace, "", fontSize)
	y += headerHeight
	for i := 0; i < rowsPerPage; i++ {
		pdf.SetXY(xLeft, y)
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i].opString, ""), "0", 0, "C", false, 0, "")
		pdf.SetXY(xRight, y)
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i+rowsPerPage].opString, ""), "0", 0, "C", false, 0, "")
		y += lineSpacing
	}

	err := pdf.OutputFileAndClose(c.filePath)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

}

func generateRandomIndex() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 5)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}
