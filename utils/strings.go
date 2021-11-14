package utils

import (
	"regexp"
	"sync"
	"time"
	"unicode"
	"golang.org/x/crypto/bcrypt"

	"mini_pos/constants"
)

var (
	emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)
var lock = &sync.Mutex{}

// ValidateEmail return true if email address has a valid format
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidatePassword(s string) bool {
	lock.Lock()
	defer lock.Unlock()
	letters := 0
	var number, upper, special bool
	for _, c := range s {
		letters++
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		}
	}
	return letters >= 8 && number && upper && special
}

func HashPassword(password string) (string, error) {
	lock.Lock()
	defer lock.Unlock()
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	lock.Lock()
	defer lock.Unlock()
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// DateString return string representation of a date pointer
func DateString(dateTime *time.Time) string {
	lock.Lock()
	defer lock.Unlock()
	if dateTime == nil {
		return ""
	}

	return dateTime.Format(constants.DateFormatStd)
}


func NullStringScan(value *string) string {
	lock.Lock()
	defer lock.Unlock()
	if value == nil {
		return ""
	}

	return *value
}
