package logger

import (
	"fmt"
	"log"
	"os"
)

var CommonLog *log.Logger
var ErrorLog *log.Logger

func init() {
	openLogfile, err := os.OpenFile("/home/ayush/Desktop/project2/server/log.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	CommonLog = log.New(openLogfile, "Common Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(openLogfile, "Error Logger:\t", log.Ldate|log.Ltime|log.Lshortfile)

}
