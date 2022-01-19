package tradingPlatformMock

import (
	"crypto-bot/internal/constant/currencyConst"
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/pkg/utils/pathUtils"
	"encoding/csv"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

var (
	openDataset = map[uint]func() (dataset []repositoryModel.Price, err error){
		1: openDataset1,
		2: openDataset2,
	} // define which method to use to open dataset depending on its ID (first number in file name)
)

// openDataset1 will open data from btc_14012022_1s.csv file
func openDataset1() (dataset []repositoryModel.Price, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}

	datasetPath := rootPath + "/" + datasetRelativePath1

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

		timestampStr := line[0]
		var timestamp int64
		timestamp, err = strconv.ParseInt(timestampStr, 0, 64)
		if err != nil {
			return
		}

		valueStr := line[4]
		var convertedValue float64
		convertedValue, err = strconv.ParseFloat(valueStr, 64)
		if err != nil {
			return
		}
		askPrice := convertedValue * (1 + fakeFeesInPercent/2/100)
		bidPrice := convertedValue * (1 - fakeFeesInPercent/2/100)
		dataset = append(dataset, repositoryModel.Price{
			Date:     time.Unix(timestamp, 0),
			Pair:     currencyConst.BtcEurPair,
			AskPrice: askPrice,
			BidPrice: bidPrice,
		})
	}
	return
}

// openDataset2 will open data from btc_17012022_1min.json file
func openDataset2() (dataset []repositoryModel.Price, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}

	datasetPath := rootPath + "/" + datasetRelativePath2
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

		// apply fake fees of 0.15% for ask and bid
		askPrice := convertedValue * fakeFeesInPercent / 100
		bidPrice := convertedValue * (1 - fakeFeesInPercent/100)
		dataset = append(dataset, repositoryModel.Price{
			Date:     time.Unix(int64(data.Ct), 0),
			Pair:     currencyConst.BtcEurPair,
			AskPrice: askPrice,
			BidPrice: bidPrice,
		})
	}
	return
}
