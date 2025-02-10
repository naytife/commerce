package models

type ErrorResponse struct {
	Error       bool
	FailedField string
	Tag         string
	Value       interface{}
}

// @Schema GlobalErrorHandlerResp
// @description Generic API error response
// @property success {boolean} false "Indicates if the request was successful"
// @property message {string} "Invalid request body" "A human-readable message"
// @property statusCode {integer} 400 "The HTTP status code"
type GlobalErrorHandlerResp struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// @Schema Response
// @description Generic API response
// @property success {boolean} true "Indicates if the request was successful"
// @property message {string} "User fetched successfully" "A human-readable message"
// @property data {object} db.User "The data returned by the API"
// @property statusCode {integer} 200 "The HTTP status code"
type ResponseHTTP struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Message    string      `json:"message"`
	StatusCode int         `json:"statusCode"`
}
