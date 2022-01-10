package postgres

import (
	"database/sql"
	"fmt"

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

func (r *orderRepo) CreateOrder(order pb.Order) (pb.Order, error) {
	var id string
	fmt.Println("OK")

	err := r.db.QueryRow(`
		INSERT INTO orders(order_id, book_id, description, updated_at)
		VALUES ($1, $2, $3, current_timestamp) RETURNING order_id`,
		order.Id,
		order.BookId,
		order.Description,
	).Scan(&id)
	if err != nil {
		return pb.Order{}, err
	}

	order, err = r.GetOrder(id)

	if err != nil {
		return pb.Order{}, nil
	}

	return order, nil
}

func (r *orderRepo) GetOrder(id string) (pb.Order, error) {
	var order pb.Order

	err := r.db.QueryRow(`
		SELECT order_id, book_id, description, created_at, updated_at FROM orders
		WHERE order_id=$1 AND deleted_at IS NULL`, id).Scan(
		&order.Id,
		&order.BookId,
		&order.Description,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return pb.Order{}, err
	}

	return order, nil
}

func (r *orderRepo) ListOrders(page, limit int64) ([]*pb.Order, int64, error) {
	offset := (page - 1) * limit

	rows, err := r.db.Queryx(`
		SELECT order_id, book_id, description, created_at FROM orders
		WHERE deleted_at IS NULL
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
			&order.Description,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		orders = append(orders, &order)
	}

	err = r.db.QueryRow(`SELECT count(*) FROM orders WHERE deleted_at IS NULL`).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

func (r *orderRepo) UpdateOrder(order pb.Order) (pb.Order, error) {
	result, err := r.db.Exec(`
		UPDATE orders SET book_id=$1, description=$2, updated_at=current_timestamp
		WHERE order_id=$3 AND deleted_at IS NULL`,
		order.BookId,
		order.Description,
		order.Id,
	)
	if err != nil {
		return pb.Order{}, err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return pb.Order{}, sql.ErrNoRows
	}

	order, err = r.GetOrder(order.Id)
	if err != nil {
		return pb.Order{}, err
	}

	return order, nil
}

func (r *orderRepo) DeleteOrder(id string) error {
	result, err := r.db.Exec(`
		UPDATE orders SET deleted_at=current_timestamp WHERE order_id=$1 AND deleted_at IS NULL`, id)
	if err != nil {
		return err
	}
	if i, _ := result.RowsAffected(); i == 0 {
		return sql.ErrNoRows
	}

	return nil
}
