package csvUtils

import (
	"crypto-bot/pkg/utils/testUtils/fakeData"
	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestForceCreateFile(t *testing.T) {
	filePath := "test.csv"

	// try to create, should work and get empty file
	err := ForceCreateFile(filePath)
	assert.NoError(t, err)
	file, err := os.Open(filePath)
	assert.NoError(t, err, "file should be created")
	stat, err := file.Stat()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), stat.Size())

	// adding lines in file to check that it will be removed when trying to create file again
	lines := []string{
		strings.Join([]string{fake.Word(), fake.Word()}, ","),
		strings.Join([]string{fake.Word(), fake.Word()}, ","),
		strings.Join([]string{fake.Word(), fake.Word()}, ","),
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	assert.NoError(t, err)

	// try to create, should work and get empty file
	err = ForceCreateFile(filePath)
	assert.NoError(t, err)
	file, err = os.Open(filePath)
	assert.NoError(t, err, "file should be created")
	stat, err = file.Stat()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), stat.Size())

	// delete file after test
	err = os.Remove(filePath)
	assert.NoError(t, err)
}

func TestCreateFileIfNotExists(t *testing.T) {
	filePath := "test.csv"

	// try to create, should work and get empty file
	err := CreateFileIfNotExists(filePath)
	assert.NoError(t, err)
	file, err := os.Open(filePath)
	assert.NoError(t, err, "file should be created")
	stat, err := file.Stat()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), stat.Size())

	// adding lines in file to check that it will be removed when trying to create file again
	lines := []string{
		strings.Join([]string{fake.Word(), fake.Word()}, ","),
		strings.Join([]string{fake.Word(), fake.Word()}, ","),
		strings.Join([]string{fake.Word(), fake.Word()}, ","),
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	assert.NoError(t, err)

	// try to create, should do nothing because file already exists
	err = CreateFileIfNotExists(filePath)
	assert.NoError(t, err)
	file, err = os.Open(filePath)
	assert.NoError(t, err, "file should be created")
	stat, err = file.Stat()
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), stat.Size())

	// delete file after test
	err = os.Remove(filePath)
	assert.NoError(t, err)
}

func TestAppendLines(t *testing.T) {
	filePath := "test.csv"

	// create file
	file, err := os.Create(filePath)
	assert.NoError(t, err)

	// appending lines
	lines := [][]string{
		{fake.Word(), fake.Word()},
		{fake.Word(), fake.Word()},
	}
	err = AppendLines(filePath, lines...)
	assert.NoError(t, err)

	// check that file is not empty and contains data
	stat, err := file.Stat()
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), stat.Size())
	input, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)
	for _, line := range lines {
		for _, word := range line {
			assert.True(t, strings.Contains(string(input), word))
		}
	}

	// delete file after test
	err = os.Remove(filePath)
	assert.NoError(t, err)
}

func TestUpdateLineOrAppend(t *testing.T) {
	filePath := "test.csv"

	// create file
	file, err := os.Create(filePath)
	assert.NoError(t, err)

	// appending existing lines
	matchingString := fakeData.UuidWithOnlyAlphaNumeric()
	previousLine := []string{matchingString, fakeData.UniqueEmail()}
	existingLine := []string{fake.Word(), fake.Word()}
	newLine := []string{matchingString, fakeData.UniqueEmail()}
	lines := [][]string{
		previousLine,
		existingLine,
	}
	var input string
	var linesStr []string
	for _, line := range lines {
		linesStr = append(linesStr, strings.Join(line, ","))
	}
	input = strings.Join(linesStr, "\n")
	err = ioutil.WriteFile(filePath, []byte(input), 0644)
	assert.NoError(t, err)

	// try to update specific line
	err = UpdateLineOrAppend(filePath, matchingString, newLine)
	assert.NoError(t, err)

	// check that file is not empty and contains data
	stat, err := file.Stat()
	assert.NoError(t, err)
	assert.NotEqual(t, int64(0), stat.Size())
	fileUpdated, err := ioutil.ReadFile(filePath)
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(fileUpdated), strings.Join(newLine, ",")),
		"file should contain %s", strings.Join(newLine, ","))
	assert.True(t, strings.Contains(string(fileUpdated), strings.Join(existingLine, ",")),
		"file should contain %s", strings.Join(existingLine, ","))
	assert.False(t, strings.Contains(string(fileUpdated), strings.Join(previousLine, ",")),
		"file should not contain %s", strings.Join(previousLine, ","))

	// delete file after test
	err = os.Remove(filePath)
	assert.NoError(t, err)
}
