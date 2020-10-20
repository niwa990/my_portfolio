package app

import (
	"fmt"

	"github.com/go-playground/validator"
)

var errorMessages []string

func UserInsertValidete(user *User) []string {
	// validateの実行
	validate := validator.New()
	err := validate.Struct(user)

	errorMessages = nil
	if err != nil {
		// MailとPasswordがValidateの対象
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Mail" || err.Field() == "Password" || err.Field() == "Name1" || err.Field() == "Name2" {
				// エラーメッセージ = 発生箇所 + エラーの種類
				fieldName := err.Field()
				tagName := err.Tag()

				fmt.Println(fieldName, tagName)

				errorMessage := GetErrorField(fieldName) + GetErrorType(tagName)
				errorMessages = append(errorMessages, errorMessage)
			}
		}
	}

	return errorMessages
}

func UserLoginValidate(user *User) []string {
	// validate 実行
	validate := validator.New()
	err := validate.Struct(user)

	errorMessages = nil
	if err != nil {
		// MailとPasswordがValidateの対象
		for _, err := range err.(validator.ValidationErrors) {
			if err.Field() == "Mail" || err.Field() == "Password" {
				// エラーメッセージ = 発生箇所 + エラーの種類
				fieldName := err.Field()
				tagName := err.Tag()
				errorMessage := GetErrorField(fieldName) + GetErrorType(tagName)
				errorMessages = append(errorMessages, errorMessage)
			}
		}
	}

	return errorMessages
}

func SkillValidate(skill *Skill) []string {
	// validate 実行
	validate := validator.New()
	err := validate.Struct(skill)

	errorMessages = nil
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			// エラーメッセージ = 発生箇所 + エラーの種類
			fieldName := err.Field()
			tagName := err.Tag()
			errorMessage := GetErrorField(fieldName) + GetErrorType(tagName)
			errorMessages = append(errorMessages, errorMessage)
		}
	}

	return errorMessages
}

// エラーの対象を取得
func GetErrorField(field string) string {
	var fieldName string
	switch field {
	case "SkillName":
		fieldName = "スキル："
	case "Period":
		fieldName = "期間："
	case "Mail":
		fieldName = "メールアドレス："
	case "Password":
		fieldName = "パスワード："
	case "Name1":
		fieldName = "名前：姓"
	case "Name2":
		fieldName = "名前：名"
	default:
		fieldName = "エラー："
	}

	return fieldName
}

// エラーの種類を取得
func GetErrorType(tag string) string {
	var errType string
	switch tag {
	case "required":
		errType = "必須項目です。入力してください"
	case "email":
		errType = "フォーマットが正しくありません。再度入力してください。"
	case "gte":
		errType = "入力桁数が不足しています。"
	default:
		errType = "無効な値です。"
	}

	return errType
}
