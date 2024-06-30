// This file is part of Algebror.
//
// Algebror is free software: you can redistribute it and/or modify it under
//   the terms of the GNU General Public License as published by the Free
//   Software Foundation, either version 3 of the License, or (at your option)
//   any later version.
//
// Algebror is distributed in the hope that it will be useful, but WITHOUT ANY
//   WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
//   FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License along with
//   Algebror. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"middle.earth/algebror/functions"
)

func main() {
	if err := functions.C.Populate(); err != nil {
		functions.Clog(err.Error(), "error")
		os.Exit(1)
	}
	functions.Clog(fmt.Sprintf("%#v\n", functions.C), "debug")

	http.HandleFunc(functions.C.AppPath, generateTest)
	functions.Clog(fmt.Sprintf("Server %s has been initiated at %s:%s%s\n", functions.C.AppName, functions.C.AppServer, functions.C.AppPort, functions.C.AppPath), "info")
	if err := http.ListenAndServe(":"+fmt.Sprintf("%v", functions.C.AppPort), nil); err != nil {
		functions.Clog(err.Error(), "error")
		os.Exit(1)
	}

}

func generateTest(w http.ResponseWriter, r *http.Request) {
	functions.C.Populate()
	dExists := false
	params := r.URL.Query()
	for key, values := range params {
		for _, value := range values {
			if key == "d" {
				dExists = true
				d, err := strconv.Atoi(value)
				if (err == nil) && (d >= 1) && (d <= 5) {
					functions.C.DifficultyLevel = d
				} else {
					functions.Clog(fmt.Sprintf("Difficulty must be an integer in range [1..5], setting default as level: %v\n", functions.C.DifficultyLevel), "info")
					http.Error(w, fmt.Sprintf("Difficulty must be an integer in range [1..5], setting default as level: %v\n", functions.C.DifficultyLevel), http.StatusMethodNotAllowed)

				}
			}
		}
	}
	if !dExists {
		functions.Clog(fmt.Sprintf("Difficulty is not provided setting default as level: %v\n", functions.C.DifficultyLevel), "info")
		http.Error(w, fmt.Sprintf("Difficulty is not provided setting default as level: %v\n", functions.C.DifficultyLevel), http.StatusMethodNotAllowed)
	}
	functions.Q = nil
	functions.Q.GenerateOps(functions.C)
	functions.Q.AlignOps()
	functions.CreatePDF(functions.Q, &functions.C)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", strings.Split(functions.C.FilePath, "/")[len(strings.Split(functions.C.FilePath, "/"))-1]))
	w.Header().Set("Content-Type", "application/pdf")
	http.ServeFile(w, r, functions.C.FilePath)
}
