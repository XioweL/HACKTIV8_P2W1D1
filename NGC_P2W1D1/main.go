package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("Field %s: %s", e.Field, e.Err)
}

func Validate(input interface{}) error {
	v := reflect.ValueOf(input)
	if v.Kind() != reflect.Struct {
		return errors.New("input must be a struct")
	}

	var validationErrors []ValidationError

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)
		tags := field.Tag.Get("validate")

		if tags == "" {
			continue
		}

		tagMap := parseTags(tags)

		if err := validateField(field.Name, value, tagMap); err != nil {
			validationErrors = append(validationErrors, ValidationError{Field: field.Name, Err: err.Error()})
		}
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors: %v", validationErrors)
	}

	return nil
}

func parseTags(tag string) map[string]string {
	tagMap := make(map[string]string)
	tags := strings.Split(tag, ",")

	for _, t := range tags {
		parts := strings.SplitN(t, "=", 2)
		if len(parts) == 2 {
			tagMap[parts[0]] = parts[1]
		} else {
			tagMap[parts[0]] = ""
		}
	}

	return tagMap
}

func validateField(name string, value reflect.Value, tagMap map[string]string) error {
	if tagMap["required"] != "" {
		if isZero(value) {
			return errors.New("is required")
		}
	}

	if minVal, ok := tagMap["min"]; ok && value.Kind() == reflect.Int {
		if value.Int() < int64(parseInt(minVal)) {
			return fmt.Errorf("must be >= %s", minVal)
		}
	}

	if maxVal, ok := tagMap["max"]; ok && value.Kind() == reflect.Int {
		if value.Int() > int64(parseInt(maxVal)) {
			return fmt.Errorf("must be <= %s", maxVal)
		}
	}

	if minLen, ok := tagMap["minLen"]; ok && value.Kind() == reflect.String {
		if len(value.String()) < parseInt(minLen) {
			return fmt.Errorf("length must be >= %s", minLen)
		}
	}

	if maxLen, ok := tagMap["maxLen"]; ok && value.Kind() == reflect.String {
		if len(value.String()) > parseInt(maxLen) {
			return fmt.Errorf("length must be <= %s", maxLen)
		}
	}

	if _, ok := tagMap["email"]; ok && value.Kind() == reflect.String {
		if !isValidEmail(value.String()) {
			return errors.New("must be a valid email")
		}
	}

	return nil
}

func parseInt(str string) int {
	num, _ := strconv.Atoi(str)
	return num
}

func isZero(v reflect.Value) bool {
	return v.IsZero()
}

func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

type Avenger struct {
	Name     string `validate:"required,minLen=3,maxLen=50"`
	Age      int    `validate:"required,min=18,max=100"`
	Email    string `validate:"required,email"`
	Rank     string `validate:"maxLen=20"`
	Missions int    `validate:"min=1"`
}

func main() {
	avenger := Avenger{
		Name:     "Steve Rogers",
		Age:      105,
		Email:    "steve.rogers@avengers.com",
		Rank:     "Captain",
		Missions: 10,
	}

	if err := Validate(avenger); err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation successful!")
	}
}
