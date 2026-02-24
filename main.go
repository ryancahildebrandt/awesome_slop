// -*- coding: utf-8 -*-

// Created on Tue Feb 17 09:53:08 PM EST 2026
// author: Ryan Hildebrandt, github.com/ryancahildebrandt

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/goccy/go-yaml"
)

type Article struct {
	Id       int64
	Category string `yaml:"category"`
	Title    string `yaml:"title"`
	Url      string `yaml:"url"`
	Abstract string `yaml:"abstract"`
}

func main() {
	var articles []Article
	var err error

	b, err := os.ReadFile("./articles.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(b, &articles)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for i := range articles {
		articles[i].Id = int64(i + 1)
	}

	funcMap := template.FuncMap{
		"categoryReferences": func(v string) string {
			var refs []string
			for _, a := range articles {
				if a.Category == v {
					refs = append(refs, fmt.Sprintf("[%v](#%v)", a.Id, a.Id))
				}
			}
			return fmt.Sprint("[", strings.Join(refs, ", "), "]")
		}}
	tmpl, err := template.New("main").Funcs(funcMap).ParseGlob("./templates/*.md.tmpl")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	res, err := os.OpenFile("./README.md", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	w := bufio.NewWriter(res)
	err = tmpl.ExecuteTemplate(w, "main", articles)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = w.Flush()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
