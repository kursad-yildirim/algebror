package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"middle.earth/algebror/functions"
)

func main() {
	functions.C.Populate()
	http.HandleFunc("/generate-test", generateTest)
	fmt.Printf("Server started at :%v\n", functions.C.ApiPort)
	log.Fatal(http.ListenAndServe(":"+fmt.Sprintf("%v", functions.C.ApiPort), nil))
}

func generateTest(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	for key, values := range params {
		for _, value := range values {
			if key == "d" {
				d, err := strconv.Atoi(value)
				if err == nil {
					functions.C.DifficultyLevel = d
				}
			}
		}
	}
	functions.Q = nil
	functions.Q.GenerateOps(functions.C)
	functions.Q.AlignOps()
	functions.CreatePDF(functions.Q, &functions.C)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strings.Split(functions.C.FilePath, "/")[len(strings.Split(functions.C.FilePath, "/"))-1]))
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, functions.C.FilePath)
}
