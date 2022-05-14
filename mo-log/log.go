package mlog

import (
	"io"
	"log"
	"os"
)

type logger struct {
	*log.Logger
}

func (l logger) Errorln(v ...any) {
	if Debug {
		l.Logger.Fatalln(v...)
	} else {
		l.Logger.Println("Error:")
		l.Logger.Println(v...)
	}
}

var Log logger

func init() {
	f, err := os.OpenFile("./output.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, os.FileMode(0666))
	if err != nil {
		panic(err.Error())
	}

	var w io.Writer
	if Debug {
		//w = io.MultiWriter(f, os.Stdout)
		w = os.Stdout
	} else {
		w = f
	}
	Log = logger{Logger: log.New(w, "", log.LstdFlags)}
}
