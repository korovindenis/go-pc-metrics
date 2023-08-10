package env

import (
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

type Env interface {
	GetEnvVariable(varName string) (string, error)
}

type EnvAdapter struct {
}

func New() (*EnvAdapter, error) {
	return &EnvAdapter{}, nil
}

func (e *EnvAdapter) GetEnvVariable(varName string) (string, error) {
	if envVarValue, exists := os.LookupEnv(varName); exists {
		return envVarValue, nil
	}
	return "", entity.ErrEnvVarNotFound
}
