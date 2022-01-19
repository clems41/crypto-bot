package localRepository

import (
	"crypto-bot/internal/repository"
	"crypto-bot/internal/repository/repositoryModel"
)

var _ repository.Position = (*positionRepo)(nil)

type positionRepo struct {
	positions map[string]*repositoryModel.Position // store all positions
}

func NewPositionRepository() (*positionRepo, error) {
	return &positionRepo{
		positions: make(map[string]*repositoryModel.Position),
	}, nil
}

func (repo *positionRepo) Store(position *repositoryModel.Position) (err error) {
	if position != nil {
		repo.positions[position.ID] = position
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

func (repo *positionRepo) GetOpenedPositions(platformName string) (positions []*repositoryModel.Position, err error) {
	for _, position := range repo.positions {
		if position.PlatformName == platformName && !position.Closed {
			positions = append(positions, position)
		}
	}
	return
}
