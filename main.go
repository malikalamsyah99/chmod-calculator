package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type PermissionData struct {
	Numeric     string
	Symbolic    string
	OwnerRead   bool
	OwnerWrite  bool
	OwnerExec   bool
	GroupRead   bool
	GroupWrite  bool
	GroupExec   bool
	PublicRead  bool
	PublicWrite bool
	PublicExec  bool
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/convert", convertHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // Handler untuk CSS

	fmt.Println("Server started at http://localhost:8011")
	http.ListenAndServe(":8011", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, PermissionData{})
}

func convertHandler(w http.ResponseWriter, r *http.Request) {
	var data PermissionData
	r.ParseForm()

	// Jika input numerik diberikan
	if numeric, ok := r.Form["numeric"]; ok && numeric[0] != "" {
		data.Numeric = numeric[0]
		data.Symbolic = numericToSymbolic(numeric[0])
	}

	// Jika input simbolik diberikan
	if symbolic, ok := r.Form["symbolic"]; ok && symbolic[0] != "" {
		data.Symbolic = symbolic[0]
		data.Numeric = symbolicToNumeric(symbolic[0])
	}

	// Update checkbox berdasarkan izin numerik
	if data.Numeric != "" {
		setCheckboxesFromNumeric(&data)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}

// Fungsi untuk mengonversi izin numerik ke simbolik
func numericToSymbolic(numeric string) string {
	num, err := strconv.Atoi(numeric)
	if err != nil || num < 0 || num > 777 {
		return ""
	}

	perms := []string{"---", "---", "---"}
	for i := 2; i >= 0; i-- {
		perm := num % 10
		perms[i] = fmt.Sprintf("%c%c%c",
			ifElse(perm&4 != 0, 'r', '-'),
			ifElse(perm&2 != 0, 'w', '-'),
			ifElse(perm&1 != 0, 'x', '-'))
		num /= 10
	}
	return "-" + strings.Join(perms, "")
}

// Fungsi untuk mengonversi izin simbolik ke numerik
func symbolicToNumeric(symbolic string) string {
	if len(symbolic) != 10 || symbolic[0] != '-' {
		return ""
	}

	numeric := 0
	for i := 1; i <= 9; i += 3 {
		perm := 0
		if symbolic[i] == 'r' {
			perm += 4
		}
		if symbolic[i+1] == 'w' {
			perm += 2
		}
		if symbolic[i+2] == 'x' {
			perm += 1
		}
		numeric = numeric*10 + perm
	}
	return fmt.Sprintf("%03d", numeric)
}

// Fungsi untuk mengatur checkboxes berdasarkan izin numerik
func setCheckboxesFromNumeric(data *PermissionData) {
	if len(data.Numeric) != 3 {
		return
	}

	owner, _ := strconv.Atoi(string(data.Numeric[0]))
	group, _ := strconv.Atoi(string(data.Numeric[1]))
	public, _ := strconv.Atoi(string(data.Numeric[2]))

	data.OwnerRead = owner&4 != 0
	data.OwnerWrite = owner&2 != 0
	data.OwnerExec = owner&1 != 0

	data.GroupRead = group&4 != 0
	data.GroupWrite = group&2 != 0
	data.GroupExec = group&1 != 0

	data.PublicRead = public&4 != 0
	data.PublicWrite = public&2 != 0
	data.PublicExec = public&1 != 0
}

// Fungsi sederhana untuk ternary
func ifElse(cond bool, a, b rune) rune {
	if cond {
		return a
	}
	return b
}
