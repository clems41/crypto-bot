package securityUtils

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSaltPassword hash and salt password using bcrypt with cost of authenticationConst.CostFunction (min 4 max 31)
func HashAndSaltPassword(password string) (hashedPassword string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), costFunction)
	if err != nil {
		return hashedPassword, errors.WithStack(err)
	}
	hashedPassword = string(bytes)
	return
}

// CompareHashAndPassword check if hashedPassword is corresponding to password
func CompareHashAndPassword(hashedPassword, password string) (err error) {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}