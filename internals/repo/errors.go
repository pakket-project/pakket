package repo

import (
	"fmt"
)

// Repository already exists
type RepositoryAlreadyExistsError struct {
	Repository string
}

func (repo RepositoryAlreadyExistsError) Error() string {
	return fmt.Sprintf("repository %s already exists", repo.Repository)
}

// Repository already exists, not defined in config
type UndefinedRepositoryAlreadyExistsError struct {
	Repository string
}

func (repo UndefinedRepositoryAlreadyExistsError) Error() string {
	return fmt.Sprintf("repository %s already exists, but is not defined in the config", repo.Repository)
}

// Repository is not a valid Stew repository
type UnvalidRepository struct {
	GitLink string
}

func (repo UnvalidRepository) Error() string {
	return fmt.Sprintf("%s is not a valid Stew repository", repo.GitLink)
}
