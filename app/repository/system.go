package repository

import (
	"suasor/repository"
)

type SystemRepositories interface {
	ConfigRepo() repository.ConfigRepository
}
