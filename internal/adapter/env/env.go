package env

import (
	"os"

	"github.com/korovindenis/go-pc-metrics/internal/domain/entity"
)

// functions
type IEnv interface {
	GetEnvVariable(varName string) (string, error)
}

type envAdapter struct {
}

func New() (IEnv, error) {
	return &envAdapter{}, nil
}

func (e *envAdapter) GetEnvVariable(varName string) (string, error) {
	if envVarValue, exists := os.LookupEnv(varName); exists {
		return envVarValue, nil
	}
	return "", entity.ErrEnvVarNotFound
}
