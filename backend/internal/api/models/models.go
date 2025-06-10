package models

// Error is a struct that holds the error response
// @Schema Error
// @description Error response object
// @property code {string} "VALIDATION_ERROR" "The error code"
// @property field {string} "Field name" "The field name that caused the error"
// @property message {string} "Error message" "The error message"
type Error struct {
	Code    string `json:"code" example:"404"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// @Schema GlobalErrorHandlerResp
// @description Generic API error response
// @property status {string} "error" "Indicates the status of the response"
// @property message {string} "An error occurred" "A human-readable message"
// @property code {integer} 500 "The HTTP status code"
type GlobalErrorHandlerResp struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// @Schema SuccessResponse
// @description Generic API response
// @property status {string} "success" "Indicates the status of the response"
// @property data {object} "The data object"
// @property message {string} "A human-readable message"
// @property code {integer} 200 "The HTTP status code"
// @property meta {object} "The meta object"
type SuccessResponse struct {
	Status  string      `json:"status" example:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message" example:"Object updated successfully"`
	Code    int         `json:"code" example:"200"`
	Meta    Meta        `json:"meta,omitempty"`
}

// @Schema ErrorResponse
// @description Generic API error response
// @property status {string} "error" "Indicates the status of the response"
// @property message {string} "An error occurred" "A human-readable message"
// @property code {integer} 500 "The HTTP status code"
// @property errors {array} "An array of errors"
// @property meta {object} "The meta object"
type ErrorResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Code    int     `json:"code"`
	Errors  []Error `json:"errors,omitempty"`
	Meta    Meta    `json:"meta,omitempty"`
}

// @Schema Meta
// @description Meta object
// @property timestamp {string} "2021-01-01T00:00:00Z" "The timestamp of the response"
type Meta struct {
	Timestamp string `json:"timestamp" example:"2025-02-12T18:31:40Z"`
}
