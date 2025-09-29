package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"nexus/api/internal/models"
)

var DB *sql.DB

func ConnectDB() {
	dsn := "postgres://nexususer:postgres@localhost:5432/nexusdb?sslmode=disable"

	var err error
	DB, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Houve um erro com a conexão com o banco de dados:", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Houve um erro ao tentar pingar o banco de dados:", err)
	}

	fmt.Println("Conexão com o banco de dados estabelecida com êxito")
}

func CreateEmpresa(empresa models.Empresa) (models.Empresa, error) {
	empresas := []models.Empresa{empresa}
	empresasSalvas, err := CreateEmpresasBatch(empresas)
	if err != nil {
		return models.Empresa{}, err
	}
	return empresasSalvas[0], nil
}

func CreateEmpresasBatch(empresas []models.Empresa) ([]models.Empresa, error) {
	tx, err := DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(context.Background(), `INSERT INTO empresas (nome, cnpj, email_contato) VALUES ($1, $2, $3) RETURNING id`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var empresasSalvas []models.Empresa

	for _, empresa := range empresas {
		var newID int64
		err := stmt.QueryRowContext(context.Background(), empresa.Nome, empresa.CNPJ, empresa.EmailContato).Scan(&newID)
		if err != nil {
			return nil, err
		}
		empresa.ID = int(newID)
		empresasSalvas = append(empresasSalvas, empresa)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return empresasSalvas, nil
}

func GetEmpresas(id *int64) ([]models.Empresa, error) {
	query := "SELECT id, nome, cnpj, email_contato FROM empresas"
	args := []interface{}{}
	if id != nil {
		query += " WHERE id = $1"
		args = append(args, *id)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var empresas []models.Empresa
	for rows.Next() {
		var emp models.Empresa
		if err := rows.Scan(&emp.ID, &emp.Nome, &emp.CNPJ, &emp.EmailContato); err != nil {
			return nil, err
		}
		empresas = append(empresas, emp)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(empresas) == 0 {
		return nil, fmt.Errorf("nenhuma empresa encontrada")
	}
	return empresas, nil
}

func UpdateEmpresa(empresa models.Empresa) (int64, error) {
	query := `UPDATE empresas SET nome = $1, cnpj = $2, email_contato_faturamento = $3 WHERE id = $4`

	// ExecContext é usado para queries que não retornam linhas (como UPDATE, DELETE).
	res, err := DB.ExecContext(context.Background(), query, empresa.Nome, empresa.CNPJ, empresa.EmailContato, empresa.ID)
	if err != nil {
		return 0, err
	}

	// RowsAffected retorna o número de linhas que foram afetadas pela query.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	// Se nenhuma linha foi afetada, significa que a empresa com aquele ID não foi encontrada.
	return rowsAffected, nil
}

func DeleteEmpresa(id int64) (int64, error) {
	query := `DELETE FROM empresas WHERE id = $1`

	res, err := DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffected, nil
}

func CreateUsuario(usuario models.Usuario) (models.Usuario, error) {
	var exists bool
	err := DB.QueryRowContext(context.Background(), "SELECT EXISTS(SELECT 1 FROM usuarios WHERE email = $1)", usuario.Email).Scan(&exists)
	if err != nil {
		return models.Usuario{}, fmt.Errorf("erro ao checar e-mail: %w", err)
	}
	if exists {
		return models.Usuario{}, fmt.Errorf("e-mail já cadastrado")
	}

	query := `INSERT INTO usuarios (nome, email, perfil) VALUES ($1, $2, $3) RETURNING id`
	err = DB.QueryRowContext(context.Background(), query, usuario.Nome, usuario.Email, usuario.Perfil).Scan(&usuario.ID)
	if err != nil {
		return models.Usuario{}, err
	}
	return usuario, nil
}

func GetUsuarios(id *int64) ([]models.Usuario, error) {
	query := "SELECT id, nome, email, perfil FROM usuarios"
	args := []interface{}{}
	if id != nil {
		query += " WHERE id = $1"
		args = append(args, *id)
	}

	rows, err := DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var usuarios []models.Usuario
	for rows.Next() {
		var usuario models.Usuario
		if err := rows.Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Perfil); err != nil {
			return nil, err
		}
		usuarios = append(usuarios, usuario)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(usuarios) == 0 {
		return nil, fmt.Errorf("nenhum usuário encontrado")
	}
	return usuarios, nil
}

func UpdateUsuario(usuario models.Usuario) (int64, error) {
	query := `UPDATE usuarios SET nome = $1, email = $2, perfil = $3 WHERE id = $4`
	res, err := DB.ExecContext(context.Background(), query, usuario.Nome, usuario.Email, usuario.Perfil, usuario.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteUsuario(id int64) (int64, error) {
	query := `DELETE FROM usuarios WHERE id = $1`
	res, err := DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func CreateContrato(contrato models.Contrato) (models.Contrato, error) {
	query := `INSERT INTO contratos (empresa_id, tipo_contrato, horas_contratadas, data_inicio, data_fim) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	if contrato.DataFim.Before(contrato.DataInicio) {
		return models.Contrato{}, fmt.Errorf("data de fim não pode ser anterior à data de início")
	}

	err := DB.QueryRowContext(context.Background(), query, contrato.EmpresaID, contrato.TipoContrato, contrato.HorasContratadas, contrato.DataInicio, contrato.DataFim).Scan(&contrato.ID)
	if err != nil {
		return models.Contrato{}, err
	}
	return contrato, nil
}

func GetContratos(id *int64) ([]models.Contrato, error) {
	query := "SELECT id, empresa_id, tipo_contrato, horas_contratadas, data_inicio, data_fim FROM contratos"
	args := []interface{}{}
	if id != nil {
		query += " WHERE id = $1"
		args = append(args, *id)
	}

	rows, err := DB.QueryContext(context.Background(), query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contratos []models.Contrato
	for rows.Next() {
		var contrato models.Contrato
		if err := rows.Scan(&contrato.ID, &contrato.EmpresaID, &contrato.TipoContrato, &contrato.HorasContratadas, &contrato.DataInicio, &contrato.DataFim); err != nil {
			return nil, err
		}
		contratos = append(contratos, contrato)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(contratos) == 0 {
		return nil, fmt.Errorf("nenhum contrato encontrado")
	}
	return contratos, nil
}

func UpdateContrato(contrato models.Contrato) (int64, error) {
	query := `UPDATE contratos SET empresa_id = $1, tipo_contrato = $2, horas_contratadas = $3, data_inicio = $4, data_fim = $5 WHERE id = $6`

	if contrato.DataFim.Before(contrato.DataInicio) {
		return 0, fmt.Errorf("data de fim não pode ser anterior à data de início")
	}

	err := DB.QueryRowContext(context.Background(), query, contrato.EmpresaID, contrato.TipoContrato, contrato.HorasContratadas, contrato.DataInicio, contrato.DataFim, contrato.ID).Scan(&contrato.ID)
	if err != nil {
		return 0, err
	}
	return contrato.ID, nil
}

// TODO: Contrato não deve ser deletado, mas gerar um campo "ativo" falso
func DeleteContrato(id int64) (int64, error) {
	query := `UPDATE contratos SET ativo = false WHERE id = $1`
	res, err := DB.ExecContext(context.Background(), query, id)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func GetContratosPorEmpresaID(empresaID int64) ([]models.Contrato, error) {
	query := "SELECT id, empresa_id, tipo_contrato, horas_contratadas, data_inicio, data_fim FROM contratos WHERE empresa_id = $1"
	rows, err := DB.QueryContext(context.Background(), query, empresaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contratos []models.Contrato
	for rows.Next() {
		var contrato models.Contrato
		if err := rows.Scan(&contrato.ID, &contrato.EmpresaID, &contrato.TipoContrato, &contrato.HorasContratadas, &contrato.DataInicio, &contrato.DataFim); err != nil {
			return nil, err
		}
		contratos = append(contratos, contrato)
	}
	return contratos, rows.Err()
}
