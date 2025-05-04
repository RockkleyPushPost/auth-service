package dto

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

const (
	MinNameLength     = 2
	MaxNameLength     = 32
	MinAge            = 0
	MaxAge            = 150
	MinPasswordLength = 6
	MaxPasswordLength = 64
	EmailRegex        = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

// CheckList is used to construct informative messages if the data is not valid
type CheckList map[string]bool

type RegisterUserDTO struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Age      uint   `json:"age"`
}

func (dto *RegisterUserDTO) Validate() error {
	dto.Name = strings.TrimSpace(dto.Name)
	dto.Email = strings.TrimSpace(dto.Email)
	dto.Password = strings.TrimSpace(dto.Password)

	if isValid, userDataError := validateUserData(*dto); !isValid {

		return errors.New(userDataError)
	}

	if isValid, passwordError := validatePassword(dto.Password); !isValid {

		return errors.New(passwordError)
	}

	return nil
}

func validateUserData(data RegisterUserDTO) (bool, string) {
	check := checkUserData(data)

	dataMissingRequirements := make([]string, 0, len(*check))

	// collecting missing requirements
	for str, req := range *check {
		if !req {
			dataMissingRequirements = append(dataMissingRequirements, str)
		}
	}

	if len(dataMissingRequirements) > 0 {

		return false, strings.Join(dataMissingRequirements, ", ")
	}

	return true, ""
}

func validatePassword(password string) (bool, string) {
	check := checkPassword(password)

	passwordMissingRequirements := make([]string, 0, len(*check))

	// collecting missing requirements
	for str, req := range *check {
		if !req {
			passwordMissingRequirements = append(passwordMissingRequirements, str)
		}
	}

	if len(passwordMissingRequirements) > 0 {

		return false, "password must also contain " + strings.Join(passwordMissingRequirements, ", ")
	}

	return true, ""
}
func checkUserData(dto RegisterUserDTO) *CheckList {
	nameLen := utf8.RuneCountInString(dto.Name)
	nameLengthMsg := fmt.Sprintf("name must be from %d to %d characters long", MinNameLength, MaxNameLength)
	ageMsg := fmt.Sprintf("age must be from %d to %d", MinAge, MaxAge)

	return &CheckList{
		nameLengthMsg:   nameLen >= MinNameLength && nameLen <= MaxNameLength,
		ageMsg:          dto.Age > MinAge && dto.Age < MaxAge,
		"invalid email": isEmailValid(dto.Email),
	}
}

func checkPassword(password string) *CheckList {
	passwordLength := utf8.RuneCountInString(password)
	lengthMsg := fmt.Sprintf("from %d to %d symbols", MinPasswordLength, MaxPasswordLength)

	return &CheckList{
		"uppercase symbol(s)": regexp.MustCompile(`[A-Z]`).MatchString(password),
		"lowercase symbol(s)": regexp.MustCompile(`[a-z]`).MatchString(password),
		"digits":              regexp.MustCompile(`[0-9]`).MatchString(password),
		lengthMsg:             passwordLength > MinPasswordLength && passwordLength < MaxPasswordLength,
	}
}

func isEmailValid(e string) bool {

	emailRegex := regexp.MustCompile(EmailRegex)

	return emailRegex.MatchString(strings.TrimSpace(e))
}
