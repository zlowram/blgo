package main

import (
	"io"
	"log"
	"os"
)

func copyFile(src string, dest string) error {
	sourcefile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destfile.Close()

	if _, err = io.Copy(destfile, sourcefile); err == nil {
		sourceinfo, err := os.Stat(src)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}
	}
	return err
}

func copyDir(src string, dest string) error {
	sourceinfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(dest, sourceinfo.Mode()); err != nil {
		return err
	}

	directory, _ := os.Open(src)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {
		sourcefilepointer := src + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			if err = copyDir(sourcefilepointer, destinationfilepointer); err != nil {
				log.Fatal(err)
			}
		} else {
			if err = copyFile(sourcefilepointer, destinationfilepointer); err != nil {
				log.Fatal(err)
			}
		}

	}
	return err
}
