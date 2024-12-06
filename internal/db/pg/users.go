package pg

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	gostreamv1 "github.com/lordvidex/gostream/pkg/api/gostream/v1"
)

// CreateUser ...
func (r *Repository) CreateUser(ctx context.Context, p *gostreamv1.User) (uint64, error) {
	q := sq.Insert("stream_users").
		Columns("name", "age", "nationality").
		Values(p.GetName(), p.GetAge(), p.GetNationality()).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	query, params, err := q.ToSql()
	if err != nil {
		fmt.Println("sq error: ", err)
		return 0, err
	}
	if err = r.pool.QueryRow(ctx, query, params...).Scan(&p.Id); err != nil {
		return 0, err
	}

	return p.Id, nil
}

// ListUsers ...
func (r *Repository) ListUsers(ctx context.Context) ([]*gostreamv1.User, error) {
	q := sq.Select("id", "name", "age", "nationality").From("stream_users").
		PlaceholderFormat(sq.Dollar)

	query, params, err := q.ToSql()
	if err != nil {
		fmt.Println("sq error: ", err)
	}
	rows, err := r.pool.Query(ctx, query, params...)
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()
	result := make([]*gostreamv1.User, 0)
	for rows.Next() {
		var u gostreamv1.User
		if err = rows.Scan(&u.Id, &u.Name, &u.Age, &u.Nationality); err != nil {
			fmt.Println(err)
		}

		result = append(result, &u)
	}
	return result, err
}


// UpdateUser ...
func (r *Repository) UpdateUser(ctx context.Context, p *gostreamv1.User) (error) {
	id := p.GetId()
	q := sq.Update("stream_users").Where("id = ?", id).
		PlaceholderFormat(sq.Dollar)
	q = q.Set("name", p.GetName())
	q = q.Set("age", p.GetAge())
	q = q.Set("nationality", p.GetNationality())

	query, params, err := q.ToSql()
	if err != nil {
		fmt.Println("sq error: ", err)
		return err
	}
	if _, err := r.pool.Exec(ctx, query, params...); err != nil {
		fmt.Println("Exec error: ", err)
		return err
	}
	return nil
}