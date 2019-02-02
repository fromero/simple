package main

import (
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
)

func init() {
	files, err := ioutil.ReadDir(sessionsDirectory)
	if err != nil {
		log.Println("Can't read saved sessions: " + err.Error())
	}
	for _, f := range files {
		sess := new(Session)
		idSession := f.Name()
		err = readGob(sessionsDirectory+f.Name(), idSession)
		if err != nil {
			log.Printf("Ca't read session %s", idSession)
		}
		sessions[f.Name()] = sess
	}

}

func writeGob(filePath string, object interface{}) error {
	file, err := os.Create(filePath)
	if err == nil {
		encoder := gob.NewEncoder(file)
		err = encoder.Encode(object)
	}
	err = file.Close()
	return err
}

func readGob(filePath string, object interface{}) error {
	file, err := os.Open(filePath)
	if err == nil {
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(object)
	}
	err = file.Close()
	return err
}
