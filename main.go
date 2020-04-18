package main

//go:generate statik -src=templates/

import (
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/rakyll/statik/fs"
	_ "github.com/ryota0624/nullable_struct/statik"
)

type Param struct {
	Package string
	Type    string
}

var (
	emptyString = "___empty___"
	typeName    = flag.String("type", emptyString, ``)
	packageName = flag.String("package", emptyString, ``)
	destName    = flag.String("dest", emptyString, ``)
)

func main() {
	flag.Parse()

	if *typeName == emptyString {
		_, _ = os.Stderr.WriteString(errors.New("flag type is empty").Error())
		os.Exit(1)
	}

	if *packageName == emptyString {
		_, _ = os.Stderr.WriteString(errors.New("flag package is empty").Error())
		os.Exit(1)
	}
	param := Param{
		Package: *packageName,
		Type:    *typeName,
	}
	statikFS, err := fs.New()
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	tmplF, err := statikFS.Open("/null_struct.go.tpl")
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	tmplBytes, err := ioutil.ReadAll(tmplF)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
	tmpl, err := template.New("null_struct").Parse(string(tmplBytes))
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}

	var dest io.Writer
	if *destName == emptyString {
		dest = os.Stdout
	} else {
		writeFile, err := os.OpenFile(*destName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			_, _ = os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
		defer writeFile.Close()
		dest = writeFile
	}

	err = tmpl.Execute(dest, param)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}
