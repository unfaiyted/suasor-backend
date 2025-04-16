package repository

import (
	"suasor/repository"
)

type JobRepositories interface {
	JobRepo() repository.JobRepository
}
