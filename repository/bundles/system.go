package bundles

import (
	"suasor/repository"
)

type SystemRepositories interface {
	ConfigRepo() repository.ConfigRepository
}
