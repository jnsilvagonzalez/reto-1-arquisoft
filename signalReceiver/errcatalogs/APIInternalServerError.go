package errcatalogs

import "fmt"

const StatusCodeInternalServerError = 500
const StatusCodeGatewayTimeoutError = 504

type APIInternalServerError struct {
	HTTPCode        int
	HTTPMessage     string
	moreInformation string
}

func (apiInternalServerError *APIInternalServerError) Error() string {
	return fmt.Sprintf(
		"Status code: %v, HttpMessage: %v, MoreInformation: %v",
		apiInternalServerError.HTTPCode,
		apiInternalServerError.HTTPMessage,
		apiInternalServerError.moreInformation,
	)
}

func MakeResponseInternalServerError(httpMessage, stackError string) error {
	return &APIInternalServerError{
		HTTPCode:        StatusCodeInternalServerError,
		HTTPMessage:     httpMessage,
		moreInformation: stackError,
	}
}
func MakeResponseGatewayTimeoutError(httpMessage, stackError string) error {
	return &APIInternalServerError{
		HTTPCode:        StatusCodeGatewayTimeoutError,
		HTTPMessage:     httpMessage,
		moreInformation: stackError,
	}
}

func (apiInternalServerError *APIInternalServerError) GetHTTPCode() int {
	return apiInternalServerError.HTTPCode
}

func (apiInternalServerError *APIInternalServerError) GetHTTPMessage() string {
	return apiInternalServerError.HTTPMessage
}

func (apiInternalServerError *APIInternalServerError) GetMoreInformation() string {
	return apiInternalServerError.moreInformation
}
