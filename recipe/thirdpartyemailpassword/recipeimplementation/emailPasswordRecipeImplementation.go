package recipeimplementation

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func MakeEmailPasswordRecipeImplementation(recipeImplementation tpepmodels.RecipeInterface) epmodels.RecipeInterface {
	return epmodels.RecipeInterface{
		SignUp: func(email, password string) (epmodels.SignUpResponse, error) {
			response, err := recipeImplementation.SignUp(email, password)
			if err != nil {
				return epmodels.SignUpResponse{}, err
			}
			if response.EmailAlreadyExistsError != nil {
				return epmodels.SignUpResponse{
					EmailAlreadyExistsError: &struct{}{},
				}, nil
			}
			return epmodels.SignUpResponse{
				OK: &struct{ User epmodels.User }{
					User: epmodels.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
					},
				},
			}, nil
		},

		SignIn: func(email, password string) (epmodels.SignInResponse, error) {
			response, err := recipeImplementation.SignIn(email, password)
			if err != nil {
				return epmodels.SignInResponse{}, err
			}
			if response.WrongCredentialsError != nil {
				return epmodels.SignInResponse{
					WrongCredentialsError: &struct{}{},
				}, nil
			}
			return epmodels.SignInResponse{
				OK: &struct{ User epmodels.User }{
					User: epmodels.User{
						ID:         response.OK.User.ID,
						Email:      response.OK.User.Email,
						TimeJoined: response.OK.User.TimeJoined,
					},
				},
			}, nil
		},

		GetUserByID: func(userId string) (*epmodels.User, error) {
			user, err := recipeImplementation.GetUserByID(userId)
			if err != nil {
				return nil, err
			}
			if user == nil || user.ThirdParty != nil {
				return nil, nil
			}
			return &epmodels.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
			}, nil
		},

		GetUserByEmail: func(email string) (*epmodels.User, error) {
			users, err := recipeImplementation.GetUsersByEmail(email)
			if err != nil {
				return nil, err
			}

			for _, user := range users {
				if user.ThirdParty == nil {
					return &epmodels.User{
						ID:         user.ID,
						Email:      user.Email,
						TimeJoined: user.TimeJoined,
					}, nil
				}
			}
			return nil, nil
		},

		CreateResetPasswordToken: func(userID string) (epmodels.CreateResetPasswordTokenResponse, error) {
			return recipeImplementation.CreateResetPasswordToken(userID)
		},
		ResetPasswordUsingToken: func(token, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error) {
			return recipeImplementation.ResetPasswordUsingToken(token, newPassword)
		},
		UpdateEmailOrPassword: func(userId string, email, password *string) (epmodels.UpdateEmailOrPasswordResponse, error) {
			return recipeImplementation.UpdateEmailOrPassword(userId, email, password)
		},
	}
}
