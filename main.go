package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
)

// Define the operation struct
type operation struct {
	name string
	sign string
}

// Define operations as a slice of operation
type operations []operation

// Define the difficulty level constant (1 to 5, 5 being the hardest)
const difficultyLevel = 1

// Define number ranges based on the difficulty level
var (
	addMin, addMax   = getRange(1.0, 10.0)
	subMin, subMax   = getRange(1.0, 10.0)
	mulMin, mulMax   = getRange(1.0, 10.0)
	divMin, divMax   = getRange(1, 10)
	percMin, percMax = getRange(0, 100)
	baseMin, baseMax = getRange(1, 10)
	opCount          = 40
)

func getRange(baseMin, baseMax float64) (float64, float64) {
	multiplier := float64(difficultyLevel)
	return baseMin * multiplier, baseMax * multiplier
}

func main() {
	http.HandleFunc("/generate", generatePDFHandler)
	fmt.Println("Server started at :3389")
	http.ListenAndServe(":3389", nil)
}

func generatePDFHandler(w http.ResponseWriter, r *http.Request) {
	// Seed the random number generator to get different results each time
	rand.Seed(time.Now().UnixNano())

	// Initialize the operations slice directly
	ops := operations{
		{name: "addition", sign: "+"},
		{name: "subtraction", sign: "-"},
		{name: "multiplication", sign: "x"},
		{name: "division", sign: "/"},
		{name: "percentage", sign: "%"},
	}

	// Prepare opCount random operations
	var randomOps [][]string
	var randomOpsWithAnswers [][]string

	for i := 0; i < opCount; i++ {
		randomIndex := rand.Intn(len(ops))
		op := ops[randomIndex]

		var num1, num2, result float64
		var opString []string
		var opWithAnswerString []string

		switch op.name {
		case "addition":
			num1 = getRandomDecimal(addMin, addMax)
			num2 = getRandomDecimal(addMin, addMax)
			result = num1 + num2
		case "subtraction":
			num1 = getRandomDecimal(subMin, subMax)
			num2 = getRandomDecimal(subMin, subMax)
			result = num1 - num2
		case "multiplication":
			num1 = getRandomDecimal(mulMin, mulMax)
			num2 = getRandomDecimal(mulMin, mulMax)
			result = num1 * num2
		case "division":
			num1 = float64(getRandomInteger(int(divMin), int(divMax)))
			num2 = float64(getRandomInteger(int(divMin), int(divMax)))
			if num2 == 0 {
				num2 = 1
			}
			result = num1 / num2
		case "percentage":
			num1 = float64(getRandomInteger(int(percMin), int(percMax))) // percentage
			num2 = float64(getRandomInteger(int(baseMin), int(baseMax))) // base number
			result = (num1 / 100) * num2
		}

		opString = []string{formatNumber(num1), op.sign, formatNumber(num2), "=", formatNumber(result)}
		opWithAnswerString = make([]string, len(opString))
		copy(opWithAnswerString, opString)

		// Randomly mask one operand or the result
		maskIndex := rand.Intn(3) // 0 for num1, 1 for num2, 2 for result
		if maskIndex == 0 {
			opString[0] = "_____"
		} else if maskIndex == 1 {
			opString[2] = "_____"
		} else {
			opString[4] = "_____"
		}

		randomOps = append(randomOps, opString)
		randomOpsWithAnswers = append(randomOpsWithAnswers, opWithAnswerString)
	}

	// Align operations to ensure equal signs are vertically aligned
	alignedQuestions := alignOperations(randomOps)
	alignedAnswers := alignOperations(randomOpsWithAnswers)

	// Generate a random 5-character index for the filename
	index := generateRandomIndex()

	// Generate PDF file with two pages and the index in the filename and title
	filename := fmt.Sprintf("operations_%s.pdf", index)
	title := fmt.Sprintf("Operations Practice - Level %d - %s", difficultyLevel, index)
	createPDF(filename, alignedQuestions, alignedAnswers, title)

	// Set headers to prompt the download with the suggested filename
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, filename)
}

func getRandomDecimal(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func getRandomInteger(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func formatNumber(num float64) string {
	return strconv.FormatFloat(num, 'f', 1, 64)
}

func alignOperations(ops [][]string) [][]string {
	maxLen := [5]int{}

	// Calculate maximum lengths for each part of the operations
	for _, op := range ops {
		for i, part := range op {
			if len(part) > maxLen[i] {
				maxLen[i] = len(part)
			}
		}
	}

	// Align operations
	alignedOps := make([][]string, len(ops))
	for i, op := range ops {
		alignedOp := make([]string, len(op))
		for j, part := range op {
			alignedOp[j] = fmt.Sprintf("%-"+strconv.Itoa(maxLen[j])+"s", part)
		}
		alignedOps[i] = alignedOp
	}

	return alignedOps
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

	// Function to add a page of operations
	addPage := func(ops [][]string, pageTitle string, header string) {
		pdf.AddPage()
		xLeft := margin
		xRight := margin + columnWidth
		y := margin

		// Add header
		pdf.SetFont("Arial", "B", 12)
		pdf.SetXY(margin, y)
		pdf.CellFormat(pageWidth-2*margin, 10, pageTitle, "0", 0, "C", false, 0, "")
		y += 10
		pdf.SetFont("Arial", "", 12)
		pdf.SetXY(margin, y)
		pdf.CellFormat(pageWidth-2*margin, 10, header, "0", 0, "L", false, 0, "")
		y += headerHeight

		// Add operations
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

	// Add questions page with full header
	addPage(questions, title, "Name: _______________________________          Time: _____          Score: _____ / 40")

	// Add answers page with a simpler header
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
