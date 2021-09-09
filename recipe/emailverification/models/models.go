package models

type TypeInput struct {
	GetEmailForUserID        func(userID string) (string, error)
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
	Override                 *OverrideStruct
}

type TypeNormalisedInput struct {
	GetEmailForUserID        func(userID string) (string, error)
	GetEmailVerificationURL  func(user User) (string, error)
	CreateAndSendCustomEmail func(user User, emailVerificationURLWithToken string)
	Override                 OverrideStruct
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}
