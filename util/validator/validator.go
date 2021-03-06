package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const tagName = "validate"

//邮箱验证正则
var mailRe = regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`)

//验证接口
type Validator interface {
	Validate(interface{}) (bool, error)
}

type DefaultValidator struct {
}

func (v DefaultValidator) Validate(val interface{}) (bool, error) {
	return true, nil
}

type StringValidator struct {
	Min int
	Max int
}

func (v StringValidator) Validate(val interface{}) (bool, error) {
	l := len(val.(string))

	if l == 0 {
		return false, fmt.Errorf("cannot be blank")
	}

	if l < v.Min {
		return false, fmt.Errorf("should be at least %v chars long", v.Min)
	}

	if v.Max >= v.Min && l > v.Max {
		return false, fmt.Errorf("should be less than %v chars long", v.Max)
	}

	return true, nil
}


type NumberValidator struct {
	Min int
	Max int
}

func (v NumberValidator) Validate(val interface{}) (bool, error) {
	num := val.(int)

	if num < v.Min {
		return false, fmt.Errorf("should be greater than %v", v.Min)
	}

	if v.Max >= v.Min && num > v.Max {
		return false, fmt.Errorf("should be less than %v", v.Max)
	}

	return true, nil
}

type EmailValidator struct {
}

func (v EmailValidator) Validate(val interface{}) (bool, error) {
	if !mailRe.MatchString(val.(string)) {
		return false, fmt.Errorf("is not a valid email address")
	}
	return true, nil
}

func getValidatorFromTag(tag string) Validator {
	args := strings.Split(tag, ",")

	switch args[0] {
	case "number":
		validator := NumberValidator{}
		//将structTag中的min和max解析到结构体中
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "email":
		return EmailValidator{}
	}

	return DefaultValidator{}
}

func validateStruct(s interface{}) []error {
	errs := []error{}

	v := reflect.ValueOf(s)

	for i := 0; i < v.NumField(); i++ {
		//利用反射获取structTag
		tag := v.Type().Field(i).Tag.Get(tagName)

		if tag == "" || tag == "-" {
			continue
		}

		validator := getValidatorFromTag(tag)

		valid, err := validator.Validate(v.Field(i).Interface())
		if !valid && err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}

	return errs
}

type User struct {
	Id    int    `validate:"number,min=1,max=1000"`
	Name  string `validate:"string,min=2,max=10"`
	Bio   string `validate:"string"`
	Email string `validate:"email"`
}

func main() {
	user1 := User{
		Id:    0,
		Name:  "superlongstring",
		Bio:   "",
		Email: "foobar",
	}

	user2 := User{
		Id:    1001,
		Name:  "superboy",
		Bio:   "",
		Email: "foobar@qq.com",
	}

	fmt.Println("Errors:")
	for i, err := range validateStruct(user1) {
		fmt.Printf("\t%d. %s\n", i+1, err.Error())
	}

	fmt.Println("Errors:")
	for i, err := range validateStruct(user2) {
		fmt.Printf("\t%d. %s\n", i+1, err.Error())
	}
}
