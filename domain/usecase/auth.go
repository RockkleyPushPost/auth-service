package usecase

import (
	"fmt"
	"log"
	"pushpost/pkg/kafka"
)

const VerificationRequestTopic = "verification-requests-topic"

type AuthUseCase struct {
	KafkaBroker string
}

func (u *AuthUseCase) CreateOTPVerificationRequest(email string) {
	producer := kafka.NewProducer(u.KafkaBroker, VerificationRequestTopic)
	defer producer.Close()

	err := producer.SendMessage([]byte(email))

	if err != nil {
		log.Println("error sending verification request:", err)
	} else {
		fmt.Println("verification request sent:", email)
	}
}
