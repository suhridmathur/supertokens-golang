package smtp

import (
	"errors"
	"net/smtp"
	"strconv"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/emaildeliverymodels"
	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery/services/smtpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func makeServiceImplementation(smtpAuth smtp.Auth, host string, port int, fromEmail string) smtpmodels.ServiceInterface {
	sendRawEmail := func(input smtpmodels.GetContentResult, userContext supertokens.UserContext) error {
		msg := []byte("From: " + fromEmail + "\r\n" +
			"To: " + input.ToEmail + "\r\n" +
			"Subject: " + input.Subject + "\r\n\r\n" +
			input.Body + "\r\n")
		return smtp.SendMail(host+":"+strconv.Itoa(port), smtpAuth, fromEmail, []string{input.ToEmail}, msg)
	}

	getContent := func(input emaildeliverymodels.EmailType, userContext supertokens.UserContext) (smtpmodels.GetContentResult, error) {
		if input.EmailVerification != nil {
			return smtpmodels.GetContentResult{}, errors.New("TODO")
		} else {
			return smtpmodels.GetContentResult{}, errors.New("should never come here")
		}
	}

	return smtpmodels.ServiceInterface{
		SendRawEmail: &sendRawEmail,
		GetContent:   &getContent,
	}
}
