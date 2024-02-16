package errcatalogs

import "fmt"

const StatusCodeBadRequest = 400
const BadRequestMoreInformationMessage = "Error"

type APIBadRequestError struct {
	statusCode      int
	httpMessage     string
	moreInformation string
}

func (apiBadRequestError *APIBadRequestError) Error() string {
	return fmt.Sprintf(
		"Status code: %v, HttpMessage: %v, MoreInformation: %v",
		apiBadRequestError.statusCode,
		apiBadRequestError.httpMessage,
		apiBadRequestError.moreInformation,
	)
}

func MakeBadRequestResponseError(messageError string) error {
	return &APIBadRequestError{
		statusCode:      StatusCodeBadRequest,
		httpMessage:     messageError,
		moreInformation: BadRequestMoreInformationMessage,
	}
}

func (apiBadRequestError *APIBadRequestError) GetStatusCode() int {
	return apiBadRequestError.statusCode
}

func (apiBadRequestError *APIBadRequestError) GetHTTPMessage() string {
	return apiBadRequestError.httpMessage
}

func (apiBadRequestError *APIBadRequestError) GetMoreInformation() string {
	return apiBadRequestError.moreInformation
}
