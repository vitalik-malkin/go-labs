package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {

	wd, err := os.Getwd()
	tmpl, err := template.ParseFiles(wd + "/text.tmpl")
	if err != nil {
		fmt.Printf("error while parsing template: %v\n", err)
	} else {
		fmt.Printf("template parsed\n")
	}
	fs, err := os.Create(wd + "text.txt")
	err = tmpl.Execute(fs, struct{ NotApply bool }{false})
	if err != nil {
		fmt.Printf("error while template execution: %v\n", err)
	}
	fmt.Printf("done\n")

}