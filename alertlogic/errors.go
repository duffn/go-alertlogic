package alertlogic

// Error messages
const (
	errEmptyApiToken               = "API token must not be empty"
	errEmptyUsernameOrPassword     = "username or password must not be empty"
	errEmptyAccessKeyIdOrSecretKey = "accessKeyId or secretKey must not be empty"
	errEmptyAccountId              = "account ID must not be empty"
	errMakeRequestError            = "error from makeRequest"
	errUnmarshalError              = "error unmarshalling the JSON response"
)
