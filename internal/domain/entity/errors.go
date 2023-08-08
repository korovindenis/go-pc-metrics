package entity

import "errors"

var ErrMetricNotFound = errors.New("metric not set")
var ErrEnvVarNotFound = errors.New("env var not set")
