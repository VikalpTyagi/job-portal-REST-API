// * User model
package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email"`
	DOB    string `json:"dateOfBirth"`
	PassHash string `json:"-"`
}
type NewUser struct {
	gorm.Model
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" gorm:"unique;not null" validate:"required,email"`
	DOB    string `json:"dateOfBirth" gorm:"not null" validate:"required"`
	Password string `json:"password" validate:"required"`
}

//*====================== for forget password =========================================
type ForgetPass struct {
	DOB   string `json:"dateOfBirth" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type OTPcont struct{
	OTP int `json:"otp" validate:"required"`
	Email string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
	ConfirmPassword string `json:"confirmPass" validate:"required"`
}