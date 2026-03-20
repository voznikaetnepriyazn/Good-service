package postgresql

import (
	"database/sql"
	"fmt"

	"github.com/voznikaetnepriyazn/Good-service/internal/models/good"

	"github.com/google/uuid"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgresql.New"

	db, err := sql.Open("pgx", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: failed to ping db: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) AddURL(good good.Good) (uuid.UUID, error) {
	const op = "storage.postgresql.addURL"

	newID := uuid.New()

	stmt, err := s.db.Prepare(
		`INSERT INTO Order ("Id", "idOfCustomer") 
		VALUES ($1, $2)
		`)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var insertedID uuid.UUID
	err = stmt.QueryRow(newID, good.Id).Scan(&insertedID)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return insertedID, nil
}

func (s *Storage) DeleteURL(id uuid.UUID) error {
	const op = "storage.postgresql.deleteURL"

	stmt, err := s.db.Prepare(
		`DELETE FROM Order WHERE id=Id`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetAllURL() ([]good.Good, error) {
	const op = "storage.postgresql.getAllURL"

	stmt, err := s.db.Prepare(`
		SELECT Order.Id
		FROM Order 
		INNER JOIN dbo.GoodInOrder ON Order.IdOfClient = GoodInOrder.IdOfClient 
		INNER JOIN Good ON GoodInOrder.Id = Good.Id
		`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var goods []good.Good
	for row.Next() {
		var good good.Good
		err := row.Scan(&good)
		if err != nil {
			return nil, fmt.Errorf("%s: scann failed: %w", op, err)
		}
		goods = append(goods, good)
	}

	return goods, nil
}

func (s *Storage) GetByIdURL(id uuid.UUID) (good.Good, error) {
	const op = "storage.postgresql.getByIdURL"

	stmt, err := s.db.Prepare(`
	SELECT * FROM dbo.Order WHERE id=Id'
	`)
	if err != nil {
		return good.Good{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var g good.Good
	err = stmt.QueryRow(id).Scan(&g)
	if err != nil {
		if err == sql.ErrNoRows {
			return good.Good{}, fmt.Errorf("%s: good not found", op)
		}
		return good.Good{}, fmt.Errorf("%s: %w", op, err)
	}

	return g, nil
}

func (s *Storage) UpdateURL(good good.Good) error {
	const op = "storage.postgresql.updateURL"

	newID := uuid.New()

	stmt, err := s.db.Prepare(
		`INSERT INTO Order ("Id", "idOfCustomer") 
		VALUES ($1, $2)
		`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var insertedID uuid.UUID
	err = stmt.QueryRow(newID, good.Id).Scan(&insertedID)
	if err != nil {
		return fmt.Errorf("%s: execute query: %w", op, err)
	}

	return nil
}

func (s *Storage) GetListOfGoodsByBrand(id uuid.UUID) ([]good.Good, error) {
	const op = "storage.postgresql.getlistofgoodsbybrand"

	res, err := s.GetAllURL()
	if err != nil {
		return []good.Good{}, fmt.Errorf("%s, %w", op, err)
	}

	var list []good.Good
	for _, l := range res {
		if l.BrandId == id {
			list = append(list, l)
		}
	}

	return list, nil
}

func (s *Storage) GetListOfGoodsByType(id uuid.UUID) ([]good.Good, error) {
	const op = "storage.postgresql.getlistofgoodsbytype"

	res, err := s.GetAllURL()
	if err != nil {
		return []good.Good{}, fmt.Errorf("%s, %w", op, err)
	}

	var list []good.Good
	for _, l := range res {
		if l.TypeId == id {
			list = append(list, l)
		}
	}

	return list, nil
}

func (s *Storage) IsAvaliableForOrder(id uuid.UUID) (bool, error) {
	const op = "storage.postgresql.isavaliablefororder"

	resGood, err := s.GetByIdURL(id)
	if err != nil {
		return false, fmt.Errorf("%s, %w", op, err)
	}

	if resGood.Rest == 0 {
		return false, nil
	}
	return true, nil
}

func (s *Storage) RestOfGood(id uuid.UUID) (int16, error) {
	const op = "storage.postgresql.restofgood"

	resGood, err := s.GetByIdURL(id)
	if err != nil {
		return 0, fmt.Errorf("%s, %w", op, err)
	}

	return resGood.Rest, nil
}
