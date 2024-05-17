package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"unicode"
)

func ValidateBody(body any, w http.ResponseWriter, r *http.Request) error {
	if r.Body == http.NoBody {
		return fmt.Errorf("body is empty")
	}

	return json.NewDecoder(r.Body).Decode(&body)
}

func ValidateEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func ValidatePassword(pass string) []string {
	var errors []string

	var (
		upper  = false
		lower  = false
		number = false
		space  = true
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upper = true
			continue
		case unicode.IsLower(char):
			lower = true
			continue
		case unicode.IsNumber(char):
			number = true
			continue
		case unicode.IsSpace(char):
			space = false
			continue
		}
	}

	if !upper {
		errors = append(errors, "La contraseña debe incluir una mayuscula")
	}

	if !lower {
		errors = append(errors, "La contraseña debe incluir una minuscula")
	}

	if !number {
		errors = append(errors, "La contraseña debe incluir un número")
	}

	if !space {
		errors = append(errors, "La contraseña no debe incluir espacios")
	}

	if len(pass) < 8 {
		errors = append(errors, "La contraseña debe tener mas de 8 caracteres")
	}

	return errors
}
