package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
	"crypto-bot/pkg/utils/csvUtils"
	"crypto-bot/pkg/utils/pathUtils"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"time"
)

var _ repository.Position = (*positionRepo)(nil)

type positionRepo struct {
	positions          map[string]*repositoryModel.Position // store all positions
	csvFilePath        string
	googleSheetService *sheets.Service
}

func NewPositionRepository(googleSheetService *sheets.Service) (repo *positionRepo, err error) {
	rootPath, err := pathUtils.GetRootProjectPath()
	if err != nil {
		return
	}
	csvPath := fmt.Sprintf("%s/%s", rootPath, csvFileNamePosition)

	// empty file or creating new one
	err = csvUtils.ForceCreateFile(csvPath)
	if err != nil {
		return
	}

	// adding column names
	columnNames := []string{
		"ID",
		"PlatformName",
		"AskDate",
		"BidDate",
		"Pair",
		"Amount",
		"Ask",
		"Bid",
		"ExpectedBidPrice",
		"Result",
		"ResultInPercent",
		"Profit",
		"Closed",
	}
	err = csvUtils.AppendLines(csvPath, columnNames)
	if err != nil {
		return
	}

	repo = &positionRepo{
		positions:          make(map[string]*repositoryModel.Position),
		csvFilePath:        csvPath,
		googleSheetService: googleSheetService,
	}
	return
}

func (repo *positionRepo) Store(position *repositoryModel.Position) (err error) {
	if position != nil {
		repo.positions[position.ID] = position
	}

	// creating new csv line
	newLine := []string{
		position.ID,
		position.PlatformName,
		position.AskDate.Format(time.RFC3339),
		position.BidDate.Format(time.RFC3339),
		position.Pair,
		fmt.Sprintf("%0.2f", position.Amount),
		fmt.Sprintf("%0.2f", position.AskPrice),
		fmt.Sprintf("%0.2f", position.BidPrice),
		fmt.Sprintf("%0.2f", position.ExpectedBidPrice),
		fmt.Sprintf("%0.2f", position.Result),
		fmt.Sprintf("%0.2f", position.ResultInPercent),
		fmt.Sprintf("%0.2f", position.Profit),
		fmt.Sprintf("%v", position.Closed),
	}

	// update line or append it if not exists
	err = csvUtils.UpdateLineOrAppend(repo.csvFilePath, position.ID, newLine)
	if err != nil {
		return
	}

	// update spreadsheet google
	var values [][]string
	for _, pos := range repo.positions {
		values = append(values, []string{
			pos.ID,
			pos.PlatformName,
			pos.AskDate.Format(time.RFC3339),
			pos.BidDate.Format(time.RFC3339),
			pos.Pair,
			fmt.Sprintf("%0.2f", pos.Amount),
			fmt.Sprintf("%0.2f", pos.AskPrice),
			fmt.Sprintf("%0.2f", pos.BidPrice),
			fmt.Sprintf("%0.2f", pos.ExpectedBidPrice),
			fmt.Sprintf("%0.2f", pos.Result),
			fmt.Sprintf("%0.2f", pos.ResultInPercent),
			fmt.Sprintf("%0.2f", pos.Profit),
			fmt.Sprintf("%v", pos.Closed),
		})
	}

	err = csvUtils.UpdateSpreadsheet(googleSpreadsheetID, positionRangeUpdate, values)
	if err != nil {
		return
	}
	return
}

func (repo *positionRepo) Get(positionID string) (position *repositoryModel.Position, err error) {
	position, ok := repo.positions[positionID]
	if !ok {
		return nil, repository.ErrEntityNotFound
	}
	return
}

func (repo *positionRepo) GetOpenedPositions(platformName string, pairs ...string) (positions []*repositoryModel.Position, err error) {
	for _, position := range repo.positions {
		if position.PlatformName == platformName && !position.Closed {
			positionShouldBeReturned := false
			if len(pairs) == 0 {
				positionShouldBeReturned = true
			} else {
				for _, pair := range pairs {
					if position.Pair == pair {
						positionShouldBeReturned = true
					}
				}
			}
			if positionShouldBeReturned {
				positions = append(positions, position)
			}

			positions = append(positions, position)
		}
	}
	return
}
