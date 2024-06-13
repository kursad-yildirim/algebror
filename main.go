package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type config struct {
	difficultyLevel int
	apiPort         int
	rangeMin        int
	rangeMax        int
	opCount         int
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

func main() {
	var c config
	var q questions
	c.populate()
	ops := operations{
		{name: "addition", sign: "+", numType: "decimal"},
		{name: "subtraction", sign: "-", numType: "decimal"},
		{name: "multiplication", sign: "x", numType: "decimal"},
		{name: "division", sign: "/", numType: "integer"},
		{name: "percentage", sign: "%", numType: "integer"},
	}
	q.generateOps(c, ops)
	q.alignOps()
	for _, r := range q {
		for _, x := range r.opString {
			fmt.Printf("%v", x)
		}
		fmt.Println()
	}
	fmt.Println("===============================================================")
	for _, r := range q {
		for _, x := range r.opStringMasked {
			fmt.Printf("%v", x)
		}
		fmt.Println()
	}
	// http.HandleFunc("/generate", generateOps())
	// fmt.Printf("Server started at :%v", c.apiPort)
	// http.ListenAndServe(":"+fmt.Sprintf("%s", c.apiPort), nil)
}

func (q *questions) generateOps(c config, ops operations) {
	for i := 0; i < c.opCount; i++ {
		*q = append(*q, generateSingleOp(c, ops[rand.Intn(len(ops))]))
	}
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
		numI1 = getRandomInteger(c)
		numI2 = getRandomInteger(c)
		resultI = numI1 * numI2
		q.opString = []string{fmt.Sprintf("%v", numI1), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", resultI)}
	case "division":
		numI1 = getRandomInteger(c)
		numI2 = getRandomInteger(c)
		resultI = numI1 * numI2
		q.opString = []string{fmt.Sprintf("%v", resultI), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", numI1)}
	case "percentage":
		numI1 = getRandomInteger(c)
		numI2 = getRandomInteger(c)
		resultF = float64(numI1*numI2) / 100
		q.opString = []string{fmt.Sprintf("%v", numI1), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", resultF)}
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
	n := math.Round(rand.Float64()*float64(c.rangeMax*10*c.difficultyLevel)) / 10
	if n == 0 {
		return n + 0.1
	} else {
		return n
	}
}

func formatNumber(num float64, p int) string {
	return strconv.FormatFloat(num, 'f', p, 64)
}

func (c *config) populate() {
	c.difficultyLevel = 5
	c.apiPort = 8080
	c.rangeMin = 0
	c.rangeMax = 10
	c.opCount = 40
}

/*

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {

	index := generateRandomIndex()

	filename := fmt.Sprintf("operations_%s.pdf", index)
	title := fmt.Sprintf("Operations Practice - Level %d - %s", difficultyLevel, index)
	createPDF(filename, alignedQuestions, alignedAnswers, title)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, filename)
}


func createPDF(filename string, questions [][]string, answers [][]string, title string) {
	const (
		pageWidth    = 210.0
		pageHeight   = 297.0
		margin       = 10.0
		columnWidth  = (pageWidth - 2*margin) / 2
		lineSpacing  = 11.0
		rowsPerPage  = 20
		headerHeight = 20.0
	)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)

	addPage := func(ops [][]string, pageTitle string, header string) {
		pdf.AddPage()
		xLeft := margin
		xRight := margin + columnWidth
		y := margin

		pdf.SetFont("Arial", "B", 12)
		pdf.SetXY(margin, y)
		pdf.CellFormat(pageWidth-2*margin, 10, pageTitle, "0", 0, "C", false, 0, "")
		y += 10
		pdf.SetFont("Arial", "", 12)
		pdf.SetXY(margin, y)
		pdf.CellFormat(pageWidth-2*margin, 10, header, "0", 0, "L", false, 0, "")
		y += headerHeight

		for i := 0; i < rowsPerPage; i++ {
			op1 := strings.Join(ops[i*2], " ")
			op2 := ""
			if i*2+1 < len(ops) {
				op2 = strings.Join(ops[i*2+1], " ")
			}

			pdf.SetXY(xLeft, y)
			pdf.CellFormat(columnWidth, lineSpacing, op1, "0", 0, "C", false, 0, "")
			pdf.SetXY(xRight, y)
			pdf.CellFormat(columnWidth, lineSpacing, op2, "0", 0, "C", false, 0, "")
			y += lineSpacing
		}
	}

	addPage(questions, title, "Name: _______________________________          Time: _____          Score: _____ / 40")

	addPage(answers, "----- Answers -----", "")

	err := pdf.OutputFileAndClose(filename)
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
*/
