package domain

import (
	"github.com/RockkleyPushPost/auth-service/domain/dto"
)

type AuthUseCase interface {
	RegisterUser(dto *dto.RegisterUserDTO) error
	Login(dto dto.UserLoginDTO) (string, error)
	IsEmailVerified(email string) (bool, error)
	SendNewOTP(email string) error
	VerifyEmailOTP(otp, email string) (bool, error)
}
