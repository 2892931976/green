package config

import (
	"log"
	"os"
)

var Log = log.New(os.Stdout, "[ goweb ] ", log.Ltime|log.Lshortfile)
