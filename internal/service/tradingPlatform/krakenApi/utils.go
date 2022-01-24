package krakenApiMock

import (
	"fmt"
)

func GetKrakenPair(projectPair string) (krakenPair string, err error) {
	krakenPair, ok := pairConverter[projectPair]
	if !ok {
		return krakenPair, fmt.Errorf("cannot find kraken pair for %s", projectPair)
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
	return projectPair, fmt.Errorf("cannot find project pair for %s", krakenPair)
}
