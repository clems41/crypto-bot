package retry

import (
	"crypto-bot/pkg/logger"
	"github.com/pkg/errors"
	"time"
)

func DoWithRetry(functionToExecute func() error, nbRetries int, delayBetweenRetries time.Duration) (err error) {
	var count int
	err = functionToExecute()
	for err != nil && count < nbRetries {
		count++
		logger.Errorf("Got error : %s, retry %d/%d after %0.0f seconds", err.Error(), count,
			nbRetries, delayBetweenRetries)
		time.Sleep(delayBetweenRetries)
		err = functionToExecute()
	}
	return errors.WithStack(err)
}
