package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
	"google.golang.org/api/sheets/v4"
	"time"
)

var _ repository.Position = (*positionRepo)(nil)

type positionRepo struct {
	positions          map[string]*repositoryModel.Position // store all positions
	googleSheetService *sheets.Service
}

func NewPositionRepository(googleSheetService *sheets.Service) (repo *positionRepo, err error) {
	repo = &positionRepo{
		positions:          make(map[string]*repositoryModel.Position),
		googleSheetService: googleSheetService,
	}
	return
}

func (repo *positionRepo) Store(position *repositoryModel.Position) (err error) {
	if position != nil {
		repo.positions[position.ID] = position
	}

	// append new line in google spreadsheet
	valueRange := sheets.ValueRange{
		Values: [][]interface{}{
			{
				position.ID,
				position.PlatformName,
				position.AskDate.Format(time.RFC3339),
				position.BidDate.Format(time.RFC3339),
				position.Pair,
				position.Amount,
				position.AskPrice,
				position.BidPrice,
				position.ExpectedBidPrice,
				position.Result,
				position.ResultInPercent,
				position.Profit,
				position.Closed,
			},
		},
	}
	request := repo.googleSheetService.Spreadsheets.Values.
		Append(googleSpreadsheetID, positionRangeUpdate, &valueRange).
		ValueInputOption("RAW")
	_, err = request.Do()
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
