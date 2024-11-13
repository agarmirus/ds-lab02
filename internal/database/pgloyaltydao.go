package database

import (
	"container/list"
	"database/sql"
	"errors"

	_ "github.com/jackc/pgx"

	"github.com/agarmirus/ds-lab02/internal/models"
	"github.com/agarmirus/ds-lab02/internal/serverrors"
)

type PostgresLoyaltyDAO struct {
	connStr string
}

func NewPostgresLoyaltyDAO(connStr string) IDAO[models.Loyalty] {
	return &PostgresLoyaltyDAO{connStr}
}

func (dao *PostgresLoyaltyDAO) SetConnectionString(connStr string) {
	dao.connStr = connStr
}

func (dao *PostgresLoyaltyDAO) Create(loyalty *models.Loyalty) (models.Loyalty, error) {
	return models.Loyalty{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) Get() (list.List, error) {
	return list.List{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) GetById(loyalty *models.Loyalty) (models.Loyalty, error) {
	return models.Loyalty{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) GetByAttribute(attrName string, attrValue string) (resLst list.List, err error) {
	db, localErr := sql.Open(`postgres`, dao.connStr)

	if localErr == nil {
		var rows *sql.Rows
		rows, localErr = db.Query(
			`select * from loyalty where $1 = '$2';`,
			attrName, attrValue,
		)
		for localErr == nil && rows.Next() {
			var loyalty models.Loyalty
			localErr = rows.Scan(&loyalty)

			if localErr != nil {
				err = errors.New(serverrors.ErrQueryResRead)
			}

			resLst.PushBack(loyalty)
		}
	} else {
		err = errors.New(serverrors.ErrDatabaseConnection)
	}

	return resLst, err
}

func (dao *PostgresLoyaltyDAO) Update(loyalty *models.Loyalty) (models.Loyalty, error) {
	return models.Loyalty{}, errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) Delete(loyalty *models.Loyalty) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}

func (dao *PostgresLoyaltyDAO) DeleteByAttr(attrName string, attrValue string) error {
	return errors.New(serverrors.ErrMethodIsNotImplemented)
}
