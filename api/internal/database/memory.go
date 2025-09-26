package database

import (
	"errors"
	"nexus/api/internal/models"
)

// Empresa
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

// Usuario
var usuarios []models.Usuario
var nextIDUsuario int = 1

func CadUsuario(usuario models.Usuario) (models.Usuario, error) {
	usuario.ID = nextIDUsuario
	nextIDUsuario++

	usuarios = append(usuarios, usuario)
	return usuario, nil
}

func GetUsuarios() ([]models.Usuario, error) {
	if len(usuarios) == 0 {
		return nil, errors.New("Nenhum usuário encontrado")
	}
	return usuarios, nil
}

func GetUsuarioByID(id int) (*models.Usuario, error) {
	for i, usuario := range usuarios {
		if usuario.ID == id {
			return &usuarios[i], nil
		}
	}
	return nil, errors.New("Usuário não encontrado")
}

func UpdateUsuario(id int, usuarioAtualizado models.Usuario) (*models.Usuario, error) {
	for i, usuario := range usuarios {
		if usuario.ID == id {
			usuarios[i].Nome = usuarioAtualizado.Nome
			usuarios[i].Email = usuarioAtualizado.Email
			usuarios[i].Perfil = usuarioAtualizado.Perfil
			return &usuarios[i], nil
		}
	}
	return nil, errors.New("Usuário não encontrado")
}

func DeleteUsuario(id int) error {
	for i, usuario := range usuarios {
		if usuario.ID == id {
			usuarios = append(usuarios[:i], usuarios[i+1:]...)
			return nil
		}
	}
	return errors.New("Usuário não encontrado")
}
