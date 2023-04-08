package postgresql

import (
	"app/api/models"
	"app/pkg/helper"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, req *models.Register) (string, error) {
	var (
		query string
	)
	id := uuid.New().String()

	query = `
		INSERT INTO users(
			user_id, 
			first_name,
			last_name,
			login,
			password,
			phone_number,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, now()) RETURNING user_id
	`
	err := r.db.QueryRow(ctx, query,
		id,
		req.FirstName,
		req.LastName,
		req.Login,
		req.Password,
		helper.NewNullString(req.PhoneNumber),
	).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *userRepo) GetByID(ctx context.Context, req *models.UserPKey) (*models.User, error) {

	var (
		query string
		user  models.User
	)

	query = `
		SELECT
			user_id, 
			first_name,
			last_name,
			login,
			password,
			COALESCE(phone_number, ''),
			TO_CHAR(created_at, 'YYYY-MM-DD HH24-MI-SS'),
            TO_CHAR(updated_at, 'YYYY-MM-DD HH24-MI-SS')
		FROM users
		WHERE user_id = $1
	`

	err := r.db.QueryRow(ctx, query, req.UserId).Scan(
		&user.UserId,
		&user.FirstName,
		&user.LastName,
		&user.Login,
		&user.Password,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetList(ctx context.Context, req *models.GetListUserRequest) (resp *models.GetListUserResponse, err error) {

	resp = &models.GetListUserResponse{}

	var (
		query  string
		filter = " WHERE TRUE "
		offset = " OFFSET 0"
		limit  = " LIMIT 10"
	)

	query = `
		SELECT
	    	COUNT(*) OVER(),
			user_id, 
			first_name,
			last_name,
			login,
			password,
			COALESCE(phone_number, ''),
			TO_CHAR(created_at, 'YYYY-MM-DD HH24-MI-SS'),
			TO_CHAR(updated_at, 'YYYY-MM-DD HH24-MI-SS')
		FROM users
	`

	if len(req.Search) > 0 {
		filter += " AND name ILIKE '%' || '" + req.Search + "' || '%' "
	}

	if req.Offset > 0 {
		offset = fmt.Sprintf(" OFFSET %d", req.Offset)
	}

	if req.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT %d", req.Limit)
	}

	query += filter + offset + limit

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(
			&resp.Count,
			&user.UserId,
			&user.FirstName,
			&user.LastName,
			&user.Login,
			&user.Password,
			&user.PhoneNumber,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		resp.Users = append(resp.Users, &user)
	}

	return resp, nil
}
