package model

import (
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	customerNameMinLength     = 3
	customerNameMaxLength     = 255
	customerEmailMaxLength    = 255
	customerPhoneMaxLength    = 30
	customerPhoneMinDigits    = 10
	customerPhoneMaxDigits    = 15
	customerPasswordMinLength = 8
	customerPasswordMaxBytes  = 64
)

var (
	customerEmailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	customerPhoneFormatRegex = regexp.MustCompile(`^\+?[0-9\s\-()]+$`)
)

type Customer struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func NewCustomerNoPassword(name, email, phone string) (*Customer, error) {
	name = NormalizeString(name)
	email = NormalizeString(strings.ToLower(email))
	phone = NormalizeString(phone)

	if err := ValidateCustomerName(name); err != nil {
		return nil, err
	}
	if err := ValidateCustomerEmail(email); err != nil {
		return nil, err
	}
	if err := ValidateCustomerPhone(phone); err != nil {
		return nil, err
	}

	return &Customer{
		Name:  name,
		Email: email,
		Phone: phone,
	}, nil
}

func NormalizeString(s string) string {
	return strings.TrimSpace(s)
}

func ValidateCustomerName(name string) error {
	switch length := utf8.RuneCountInString(name); {
	case length == 0:
		return ErrCustomerNameRequired
	case length < customerNameMinLength:
		return ErrCustomerNameTooShort
	case length > customerNameMaxLength:
		return ErrCustomerNameTooLong
	default:
		return nil
	}
}

func ValidateCustomerEmail(email string) error {
	switch {
	case email == "":
		return ErrCustomerEmailRequired
	case len(email) > customerEmailMaxLength:
		return ErrCustomerEmailTooLong
	case !customerEmailRegex.MatchString(email):
		return ErrCustomerEmailInvalid
	default:
		return nil
	}
}

func ValidateCustomerPhone(phone string) error {
	switch {
	case phone == "":
		return ErrCustomerPhoneRequired
	case len(phone) > customerPhoneMaxLength:
		return ErrCustomerPhoneTooLong
	case !customerPhoneFormatRegex.MatchString(phone):
		return ErrCustomerPhoneInvalid
	default:
		digitsCount := countDigits(phone)
		if digitsCount < customerPhoneMinDigits || digitsCount > customerPhoneMaxDigits {
			return ErrCustomerPhoneInvalid
		}
		return nil
	}
}

func ValidateCustomerPassword(password string) error {
	normalized := NormalizeString(password)
	switch length := utf8.RuneCountInString(normalized); {
	case normalized == "":
		return ErrCustomerPasswordRequired
	case length < customerPasswordMinLength:
		return ErrCustomerPasswordTooShort
	case len(password) > customerPasswordMaxBytes:
		return ErrCustomerPasswordTooLong
	default:
		return nil
	}
}

func countDigits(s string) int {
	count := 0
	for _, r := range s {
		if r >= '0' && r <= '9' {
			count++
		}
	}
	return count
}
