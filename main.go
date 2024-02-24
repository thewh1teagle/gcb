package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.design/x/clipboard"
)

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

func main() {
	args := os.Args
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
	if len(args) > 1 {
		// read file / image
		path := args[1]
		ext := filepath.Ext(path)
		imgs_ext := []string{".jpg", ".png"}

		// Check if arg path exists
		if _, err := os.Stat(path); err != nil {
			if os.IsNotExist(err) {
				// file does not exist
				// write args as is
				args := os.Args[1:]                                    // Exclude the first element which is the program name
				joinedArgs := strings.Join(args, " ")                  // Join the args with spaces
				clipboard.Write(clipboard.FmtText, []byte(joinedArgs)) // Write the joined args to the clipboard
				os.Exit(0)
			} else {

				log.Fatalf("failed to read file: %v", err)
			}
		}

		// Check if image or text file
		if Contains(imgs_ext, ext) {
			// Detected image file
			content, err := ioutil.ReadFile(path)
			if err != nil {
				// UTF-8 error?
				log.Fatalf("failed to read file: %v", err)
			}
			clipboard.Write(clipboard.FmtImage, content)
		} else {
			// try to decode into UTF-8 buffer
			content, err := ioutil.ReadFile(path)
			if err != nil {
				// UTF-8 error?
				log.Fatalf("failed to read file: %v", err)
			}
			clipboard.Write(clipboard.FmtText, content)
		}
	} else {
		buf := new(bytes.Buffer)
		// Copy from os.Stdin to the buffer until the delimiter is encountered
		_, err = io.Copy(buf, os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error reading standard input:", err)
			return
		}
		clipboard.Write(clipboard.FmtText, buf.Bytes())
	}
}
