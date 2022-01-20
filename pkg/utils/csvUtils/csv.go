package csvUtils

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// CreateFileIfNotExists will create new file based on filePath, if file already exists, nothing will be done.
func CreateFileIfNotExists(filePath string) (err error) {
	// check if file already exists
	_, err = os.Stat(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return
		}
	} else {
		// do nothing if already exists
		return
	}

	// creating new one if not exists
	_, err = os.Create(filePath)
	if err != nil {
		return
	}
	return
}

// ForceCreateFile will create new file based on filePath, if file already exists, it will be removed before creating new one.
func ForceCreateFile(filePath string) (err error) {
	// check if file already exists
	_, err = os.Stat(filePath)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return
		}
	} else {
		// delete file if already exists
		err = os.Remove(filePath)
		if err != nil {
			return
		}
	}

	// creating new one
	_, err = os.Create(filePath)
	if err != nil {
		return
	}
	return
}

// AppendLines will append all given line in file. One line is a string array (columns)
func AppendLines(filePath string, lines ...[]string) (err error) {

	// opening file
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	// creating lineStr, with all cell separated with coma
	var linesStr []string
	for _, line := range lines {
		linesStr = append(linesStr, strings.Join(line, ","))
	}

	// update file
	output := string(input) + strings.Join(linesStr, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return
	}
	return
}

// UpdateLineOrAppend will check if line already exists in file using matchingString (can be ID for example).
// If exists, line will be updated, if not, line will be appended in file.
func UpdateLineOrAppend(filePath string, matchingString string, newLine []string) (err error) {

	// opening file
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	// creating new line
	newLineStr := strings.Join(newLine, ",")

	// get actual lines
	lines := strings.Split(string(input), "\n")

	var lineExists bool
	for i, line := range lines {
		// if line exists, replace it
		if strings.Contains(line, matchingString) {
			lines[i] = newLineStr
			lineExists = true
		}
	}

	// if line doesn't exist, append new one
	if !lineExists {
		lines = append(lines, newLineStr)
	}

	// push updated lines
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return
	}
	return
}
