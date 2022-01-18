package tradingPlatformMock

import (
	"crypto-bot/pkg/utils/pathUtils"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)

// openDataset1 will open data from btc_17012022_1min.json file
func openDataset1() (dataset []float64, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}

	datasetPath := rootPath + "/" + datasetRelativePath1
	datasetFile, err := ioutil.ReadFile(datasetPath)
	if err != nil {
		return nil, err
	}
	var values Dataset1
	err = json.Unmarshal(datasetFile, &values)
	if err != nil {
		return nil, err
	}

	for _, data := range values.Values {
		valueStr := data.Hp
		var convertedValue float64
		convertedValue, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return
		}
		dataset = append(dataset, float64(convertedValue))
	}
	return
}

// openDataset2 will open data from btc_14012022_1s.csv file
func openDataset2() (dataset []float64, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}

	datasetPath := rootPath + "/" + datasetRelativePath2

	csvFile, err := os.Open(datasetPath)
	if err != nil {
		return
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return
	}
	for idx, line := range csvLines {
		if idx == 0 {
			continue // first line is column names
		}
		valueStr := line[4]
		var convertedValue float64
		convertedValue, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return
		}
		dataset = append(dataset, float64(convertedValue))
	}
	return
}
