package entity

import "errors"

var (
	ErrMetricNotFound            = errors.New("metric not set")
	ErrEnvVarNotFound            = errors.New("env var not set")
	ErrInvalidURLFormat          = errors.New("invalid URL format")
	ErrMethodNotAllowed          = errors.New("method not allowed")
	ErrInternalServerError       = errors.New("internal server error")
	ErrInputVarIsWrongType       = errors.New("metric value is wrong type")
	ErrStatusBadRequest          = errors.New("bad request")
	ErrInvalidGzipData           = errors.New("invalid gzip data")
	ErrReadingRequestBody        = errors.New("error reading request body")
	ErrInputMetricNotFound       = errors.New("metric not found")
	ErrNotImplementedServerError = errors.New("not implemented server error")
	ErrStorageInstance           = errors.New("data is not an instance of storage")
	ErrConfigFileNotFound        = errors.New("config file not found")
	ErrForbidden                 = errors.New("forbidden")
)
