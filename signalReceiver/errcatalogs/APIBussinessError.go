package errcatalogs

import "fmt"

type APIBusinessError struct {
	statusCode       int
	serverStatusCode int
	httpMessage      string
}

func (apiBusinessError *APIBusinessError) Error() string {
	return fmt.Sprintf(
		"Status code: %v, HttpMessage: %v, ServerStatusCode: %v",
		apiBusinessError.statusCode,
		apiBusinessError.httpMessage,
		apiBusinessError.serverStatusCode,
	)
}

func MakeBusinessResponseError(statusCode int, httpMessage string) error {
	return &APIBusinessError{
		statusCode:       statusCode,
		httpMessage:      httpMessage,
		serverStatusCode: 206,
	}
}

func (apiBusinessError *APIBusinessError) GetStatusCode() int {
	return apiBusinessError.statusCode
}

func (apiBusinessError *APIBusinessError) GetHTTPMessage() string {
	return apiBusinessError.httpMessage
}

func (apiBusinessError *APIBusinessError) GetServerStatusCode() int {
	return apiBusinessError.serverStatusCode
}
