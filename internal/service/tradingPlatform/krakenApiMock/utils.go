package krakenApiMock

import "crypto-bot/internal/service/tradingPlatform"

func GetKrakenPair(projectPair string) (krakenPair string, err error) {
	krakenPair, ok := pairConversion[projectPair]
	if !ok {
		return krakenPair, tradingPlatform.ErrPairNotFound
	}
	return
}

func GetProjectPair(krakenPair string) (projectPair string, err error) {
	for project, kraken := range pairConversion {
		if kraken == krakenPair {
			projectPair = project
			return
		}
	}
	return projectPair, tradingPlatform.ErrPairNotFound
}
