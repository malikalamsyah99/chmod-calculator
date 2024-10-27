package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/convert", convertHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Render the homepage
func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, nil)
}

// Handle conversion
func convertHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	symbolic := r.FormValue("symbolic")
	numeric := r.FormValue("numeric")

	var result string
	var err error

	if symbolic != "" {
		result, err = symbolicToNumeric(symbolic)
	} else if numeric != "" {
		result, err = numericToSymbolic(numeric)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, map[string]string{
		"Symbolic": symbolic,
		"Numeric":  numeric,
		"Result":   result,
	})
}

// Convert symbolic permissions (e.g., -rw-r--r--) to numeric (e.g., 644)
func symbolicToNumeric(perm string) (string, error) {
	if len(perm) != 10 {
		return "", fmt.Errorf("invalid permission format")
	}

	owner := perm[1:4]
	group := perm[4:7]
	public := perm[7:10]

	numeric := fmt.Sprintf("%d%d%d", toOctalSymbol(owner), toOctalSymbol(group), toOctalSymbol(public))
	return numeric, nil
}

// Convert each symbolic part (e.g., rw-) to its octal equivalent
func toOctalSymbol(symbols string) int {
	val := 0
	if symbols[0] == 'r' {
		val += 4
	}
	if symbols[1] == 'w' {
		val += 2
	}
	if symbols[2] == 'x' {
		val += 1
	}
	return val
}

// Convert numeric permissions (e.g., 755) to symbolic (e.g., -rwxr-xr-x)
func numericToSymbolic(numeric string) (string, error) {
	if len(numeric) != 3 {
		return "", fmt.Errorf("invalid numeric permission format")
	}

	owner, err := strconv.Atoi(string(numeric[0]))
	if err != nil {
		return "", err
	}
	group, err := strconv.Atoi(string(numeric[1]))
	if err != nil {
		return "", err
	}
	public, err := strconv.Atoi(string(numeric[2]))
	if err != nil {
		return "", err
	}

	symbolic := fmt.Sprintf("-%s%s%s", toSymbol(owner), toSymbol(group), toSymbol(public))
	return symbolic, nil
}

// Convert each octal value (e.g., 7) to its symbolic equivalent (e.g., rwx)
func toSymbol(octal int) string {
	var symbols strings.Builder
	if octal >= 4 {
		symbols.WriteString("r")
		octal -= 4
	} else {
		symbols.WriteString("-")
	}
	if octal >= 2 {
		symbols.WriteString("w")
		octal -= 2
	} else {
		symbols.WriteString("-")
	}
	if octal >= 1 {
		symbols.WriteString("x")
	} else {
		symbols.WriteString("-")
	}
	return symbols.String()
}
