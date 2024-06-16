/*This file is part of Algebror.

Algebror is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

Algebror is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with Algebror. If not, see <https://www.gnu.org/licenses/>.

*/

package functions

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"middle.earth/algebror/tools"
)

type Config struct {
	DifficultyLevel int
	ApiPort         int
	RangeMax        int
	OpCount         int
	FilePath        string
	FileDst         string
	Ops             operations
}

type operation struct {
	name    string
	sign    string
	numType string
}

type operations []operation

type question struct {
	OpString       []string
	OpStringMasked []string
}

type Questions []question

var C Config
var Q Questions

func (c *Config) Populate() {
	c.DifficultyLevel = 2
	c.ApiPort = 18080
	c.RangeMax = 10
	c.OpCount = 40
	c.FileDst = "./"
	c.Ops = operations{
		{name: "addition", sign: "+", numType: "decimal"},
		{name: "subtraction", sign: "-", numType: "decimal"},
		{name: "multiplication", sign: "x", numType: "decimal"},
		{name: "division", sign: "/", numType: "integer"},
		{name: "percentage", sign: "%", numType: "integer"},
	}
}

func (q *Questions) GenerateOps(c Config) {
	for i := 0; i < c.OpCount; i++ {
		*q = append(*q, generateSingleOp(c, c.Ops[rand.Intn(len(c.Ops))]))
	}
}

func generateSingleOp(c Config, op operation) question {
	var numF1, numF2, resultF float64
	var numI1, numI2, resultI int
	var q question
	switch op.name {
	case "addition":
		numF1 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
		numF2 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
		resultF = numF1 + numF2
		q.OpString = []string{tools.FormatNumber(numF1, 1), " " + op.sign + " ", tools.FormatNumber(numF2, 1), " = ", tools.FormatNumber(resultF, 2)}
	case "subtraction":
		numF1 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
		numF2 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
		resultF = numF1 - numF2
		q.OpString = []string{tools.FormatNumber(numF1, 1), " " + op.sign + " ", tools.FormatNumber(numF2, 1), " = ", tools.FormatNumber(resultF, 2)}
	case "multiplication":
		if c.DifficultyLevel > 2 {
			numF1 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
			numF2 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
			resultF = math.Round(numF1*numF2*100) / 100
			q.OpString = []string{fmt.Sprintf("%v", numF1), " " + op.sign + " ", fmt.Sprintf("%v", numF2), " = ", fmt.Sprintf("%v", resultF)}
		} else {
			numI1 = tools.GetRandomInteger(c.RangeMax, c.DifficultyLevel)
			numI2 = tools.GetRandomInteger(c.RangeMax, c.DifficultyLevel)
			resultI = numI1 * numI2
			q.OpString = []string{fmt.Sprintf("%v", numI1), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", resultI)}
		}
	case "division":
		if c.DifficultyLevel > 3 {
			numF1 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
			numF2 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
			resultF = math.Round(numF1*numF2*100) / 100
			q.OpString = []string{fmt.Sprintf("%v", resultF), " " + op.sign + " ", fmt.Sprintf("%v", numF2), " = ", fmt.Sprintf("%v", numF1)}
		} else {
			numI1 = tools.GetRandomInteger(c.RangeMax, c.DifficultyLevel)
			numI2 = tools.GetRandomInteger(c.RangeMax, c.DifficultyLevel)
			resultI = numI1 * numI2
			q.OpString = []string{fmt.Sprintf("%v", resultI), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", numI1)}
		}
	case "percentage":
		if c.DifficultyLevel > 4 {
			numF1 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
			numF2 = tools.GetRandomDecimal(c.RangeMax, c.DifficultyLevel)
			resultF = math.Round(numF1*numF2/100*100) / 100
			q.OpString = []string{fmt.Sprintf("%v", numF1), " " + op.sign + " ", fmt.Sprintf("%v", numF2), " = ", fmt.Sprintf("%v", resultF)}
		} else {
			numI1 = tools.GetRandomInteger(c.RangeMax, c.DifficultyLevel)
			numI2 = tools.GetRandomInteger(c.RangeMax, c.DifficultyLevel)
			resultF = math.Round(float64(numI1*numI2)/100*100) / 100
			q.OpString = []string{fmt.Sprintf("%v", numI1), " " + op.sign + " ", fmt.Sprintf("%v", numI2), " = ", fmt.Sprintf("%v", resultF)}
		}
	}
	q.OpStringMasked = make([]string, len(q.OpString))
	copy(q.OpStringMasked, q.OpString)
	maskIndex := rand.Intn(3)
	if maskIndex == 0 {
		q.OpStringMasked[0] = "____"
	} else if maskIndex == 1 {
		q.OpStringMasked[2] = "____"
	} else {
		q.OpStringMasked[4] = "____"
	}
	return q
}

func (q *Questions) AlignOps() {
	col1 := (*q)[:len(*q)/2]
	lenBefore, lenBeforeMasked := 0, 0
	for _, r := range col1 {
		if len(r.OpString[0])+len(r.OpString[2]) > lenBefore {
			lenBefore = len(r.OpString[0]) + len(r.OpString[2])
		}
		if len(r.OpStringMasked[0])+len(r.OpStringMasked[2]) > lenBeforeMasked {
			lenBeforeMasked = len(r.OpStringMasked[0]) + len(r.OpStringMasked[2])
		}
	}
	for _, r := range col1 {
		r.OpString[0] = strings.Repeat(" ", lenBefore-len(r.OpString[0])-len(r.OpString[2])) + r.OpString[0]
		r.OpStringMasked[0] = strings.Repeat(" ", lenBeforeMasked-len(r.OpStringMasked[0])-len(r.OpStringMasked[2])) + r.OpStringMasked[0]
	}
	col2 := (*q)[len(*q)/2:]
	lenBefore, lenBeforeMasked = 0, 0
	for _, r := range col2 {
		if len(r.OpString[0])+len(r.OpString[2]) > lenBefore {
			lenBefore = len(r.OpString[0]) + len(r.OpString[2])
		}
		if len(r.OpStringMasked[0])+len(r.OpStringMasked[2]) > lenBeforeMasked {
			lenBeforeMasked = len(r.OpStringMasked[0]) + len(r.OpStringMasked[2])
		}
	}
	for _, r := range col2 {
		r.OpString[0] = strings.Repeat(" ", lenBefore-len(r.OpString[0])-len(r.OpString[2])) + r.OpString[0]
		r.OpStringMasked[0] = strings.Repeat(" ", lenBeforeMasked-len(r.OpStringMasked[0])-len(r.OpStringMasked[2])) + r.OpStringMasked[0]
	}
}

func CreatePDF(q Questions, c *Config) {
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

	filename := fmt.Sprintf("operations_%v_%s.pdf", c.DifficultyLevel, index)
	c.FilePath = c.FileDst + filename
	pageTitle := fmt.Sprintf("Operations Practice - Level %d - %s", c.DifficultyLevel, index)

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
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i].OpStringMasked, ""), "0", 0, "C", false, 0, "")
		pdf.SetXY(xRight, y)
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i+rowsPerPage].OpStringMasked, ""), "0", 0, "C", false, 0, "")
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
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i].OpString, ""), "0", 0, "C", false, 0, "")
		pdf.SetXY(xRight, y)
		pdf.CellFormat(columnWidth, lineSpacing, strings.Join(q[i+rowsPerPage].OpString, ""), "0", 0, "C", false, 0, "")
		y += lineSpacing
	}

	err := pdf.OutputFileAndClose(c.FilePath)
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
