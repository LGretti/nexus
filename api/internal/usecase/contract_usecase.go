package usecase

import (
	"context"
	"nexus/api/internal/models"
	"nexus/api/internal/repository"
)

type ContractUsecase interface {
	CreateContract(ctx context.Context, contract models.Contract) (models.Contract, error)
}

type contractUsecase struct {
	contractRepo repository.ContractRepository
	companyRepo  repository.CompanyRepository
}
