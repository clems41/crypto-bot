package krakenApiMock

import "crypto-bot/internal/service/tradingPlatform"

func GetKrakenPair(projectPair string) (krakenPair string, err error) {
	krakenPair, ok := pairConverter[projectPair]
	if !ok {
		return krakenPair, tradingPlatform.ErrPairNotFound
	}
	return
}

func GetProjectPair(krakenPair string) (projectPair string, err error) {
	for project, kraken := range pairConverter {
		if kraken == krakenPair {
			projectPair = project
			return
		}
	}
	return projectPair, tradingPlatform.ErrPairNotFound
}
