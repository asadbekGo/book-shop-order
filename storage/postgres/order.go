package postgres

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	pb "github.com/asadbekGo/book-shop-order/genproto/order_service"
)

type orderRepo struct {
	db *sqlx.DB
}

// NewOrderRepo ...
func NewOrderRepo(db *sqlx.DB) *orderRepo {
	return &orderRepo{db: db}
}

func (r *orderRepo) Create(order pb.Order) (pb.Order, error) {
	var id string
	err := r.db.QueryRow(`
		INSERT INTO orders(order_id, book_id, quantity, description, updated_at)
		VALUES ($1, $2, $3, $4, current_timestamp) returning order_id`,
		order.Id,
		order.BookId,
		order.Quantity,
		order.Description,
	).Scan(&id)
	if err != nil {
		return pb.Order{}, err
	}

	order, err = r.Get(id)

	if err != nil {
		return pb.Order{}, nil
	}

	return order, nil
}

func (r *orderRepo) Get(id string) (pb.Order, error) {
	var order pb.Order

	err := r.db.QueryRow(`
		SELECT order_id, book_id, quantity, description, created_at, updated_at FROM orders 
		WHERE order_id=$1 and deleted_at is null`, id).Scan(
		&order.Id,
		&order.BookId,
		&order.Quantity,
		&order.Description,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return pb.Order{}, err
	}

	return order, nil
}

func (r *orderRepo) List(page, limit int64) ([]*pb.Order, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
		SELECT order_id, book_id, quantity, description, created_at FROM orders 
		WHERE deleted_at is null
		LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	if err = rows.Err(); err != nil {
		return nil, 0, err
	}
	defer rows.Close() // nolint:errcheck

	var (
		orders []*pb.Order
		count  int64
	)

	for rows.Next() {
		var order pb.Order
		err = rows.Scan(
			&order.Id,
			&order.BookId,
			&order.Quantity,
			&order.Description,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		orders = append(orders, &order)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM orders WHERE deleted_at is null`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (r *orderRepo) Update(order pb.Order) (pb.Order, error) {
	result, err := r.db.Exec(`
		UPDATE orders SET book_id=$1, quantity=$2, description=$3, updated_at=current_timestamp
		WHERE order_id=$4 and deleted_at is null`,
		order.BookId,
		order.Quantity,
		order.Description,
		order.Id,
	)
	if err != nil {
		return pb.Order{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Order{}, sql.ErrNoRows
	}

	order, err = r.Get(order.Id)
	if err != nil {
		return pb.Order{}, err
	}

	return order, nil
}

func (r *orderRepo) Delete(id string) error {
	result, err := r.db.Exec(`
		UPDATE orders SET deleted_at=current_timestamp WHERE order_id=$1`, id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
