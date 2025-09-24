package database

import (
	"errors"
	"nexus/api/internal/models"
)

var empresas []models.Empresa
var nextIDEmpresa int = 1

func CadEmpresa(empresa models.Empresa) (models.Empresa, error) {
	empresa.ID = nextIDEmpresa
	nextIDEmpresa++

	empresas = append(empresas, empresa)
	return empresa, nil
}

func CadEmpresasEmLote(novasEmpresas []models.Empresa) ([]models.Empresa, error) {
	var empresasSalvas []models.Empresa

	for _, empresa := range novasEmpresas {
		empresa.ID = nextIDEmpresa
		nextIDEmpresa++
		empresas = append(empresas, empresa)
		empresasSalvas = append(empresasSalvas, empresa)
	}

	return empresasSalvas, nil
}

func GetEmpresas() ([]models.Empresa, error) {
	if len(empresas) == 0 {
		return nil, errors.New("Nenhuma empresa encontrada")
	}
	return empresas, nil
}

func GetEmpresaByID(id int) (*models.Empresa, error) {
	for i, empresa := range empresas {
		if empresa.ID == id {
			return &empresas[i], nil
		}
	}
	return nil, errors.New("Empresa não encontrada")
}

func UpdateEmpresa(id int, empresaAtualizada models.Empresa) (*models.Empresa, error) {
	for i, empresa := range empresas {
		if empresa.ID == id {
			empresas[i].Nome = empresaAtualizada.Nome
			empresas[i].CNPJ = empresaAtualizada.CNPJ
			empresas[i].EmailContato = empresaAtualizada.EmailContato
			return &empresas[i], nil
		}
	}
	return nil, errors.New("Empresa não encontrada")
}

func DeleteEmpresa(id int) error {
	for i, empresa := range empresas {
		if empresa.ID == id {
			empresas = append(empresas[:i], empresas[i+1:]...)
			return nil
		}
	}
	return errors.New("Empresa não encontrada")
}