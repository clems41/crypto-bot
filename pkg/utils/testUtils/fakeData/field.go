package fakeData

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/icrowley/fake"
	"strings"
)

// UuidWithOnlyAlphaNumeric return Uuid from uuid.New() without any char like '-'
func UuidWithOnlyAlphaNumeric() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

// Bool return a bool randomly true or false
func Bool() bool {
	return FakeIntBetween(0, 1) == 1
}

// Password return string complex password
func Password() string {
	return fake.Password(passwordMinCharacter, passwordMaxCharacter, passwordAllowUpper, passwordAllowNumber, passwordAllowSpecial)
}

// UniqueEmail return a string as email format will be unique (using uuid)
func UniqueEmail() string {
	return fmt.Sprintf("%s@%s.%s", UuidWithOnlyAlphaNumeric(), fake.Word(), fake.Word())
}

// UniqueUsername return unique username (using uuid)
func UniqueUsername() string {
	return fmt.Sprintf("%s%s%s", fake.FirstName(), fake.LastName(), UuidWithOnlyAlphaNumeric())
}
