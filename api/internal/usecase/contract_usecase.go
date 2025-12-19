package usecase

import (
	"context"
	"nexus/internal/models"
	"nexus/internal/repository"
)

type ContractUsecase interface {
	CreateContract(ctx context.Context, contract models.Contract) (models.Contract, error)
}

type contractUsecase struct {
	contractRepo repository.ContractRepository
	companyRepo  repository.CompanyRepository
}
