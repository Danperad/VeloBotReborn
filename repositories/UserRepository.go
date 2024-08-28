package repositories

import (
	"context"
	"github.com/danperad/tgBotVeloBot/dbModels"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

func InsertUser(user *dbModels.User) error {
	var databaseUrl = os.Getenv("DATABASE_URL")

	db, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	var count int
	err = tx.QueryRow(context.Background(), "SELECT count(user_id) FROM users WHERE user_id LIKE $1", user.UserId).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	_, err = tx.Exec(context.Background(), "INSERT INTO users (user_id, user_name) VALUES ($1, $2)", user.UserId, user.UserName)
	if err != nil {
		return err
	}

	_ = tx.Commit(context.Background())

	return nil
}

func AddResult(userId string, speed float64, distance float64) (*dbModels.User, error) {
	var databaseUrl = os.Getenv("DATABASE_URL")

	db, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tx, err := db.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())
	user := dbModels.User{}

	err = tx.QueryRow(context.Background(), "SELECT user_id, user_name FROM users WHERE user_id = $1", userId).Scan(&user.UserId, &user.UserName)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(context.Background(), "INSERT INTO result (user_id, max_speed, distance) VALUES ($1, $2, $3)", userId, speed, distance)
	if err != nil {
		return nil, err
	}

	_ = tx.Commit(context.Background())

	return &user, nil
}

func GetResults(userId string) ([]dbModels.Result, error) {
	var databaseUrl = os.Getenv("DATABASE_URL")

	db, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	user := dbModels.User{}

	err = db.QueryRow(context.Background(), "SELECT user_id, user_name FROM users WHERE user_id = $1", userId).Scan(&user.UserId, &user.UserName)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(context.Background(), "SELECT result_id, max_speed, distance FROM result WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	results := make([]dbModels.Result, 0)
	for rows.Next() {
		result := dbModels.Result{}
		if err := rows.Scan(&result.ResultId, &result.MaxSpeed, &result.Distance); err != nil {
			return nil, err
		}
		result.User = user
		results = append(results, result)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func GetAllResults() ([]dbModels.Raiting, error) {
	var databaseUrl = os.Getenv("DATABASE_URL")

	db, err := pgxpool.New(context.Background(), databaseUrl)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query(context.Background(), "SELECT user_id, user_name, max_speed, max_distance, distance FROM rating")
	if err != nil {
		return nil, err
	}
	results := make([]dbModels.Raiting, 0)
	for rows.Next() {
		result := dbModels.Raiting{}
		if err := rows.Scan(&result.User.UserId, &result.User.UserName, &result.MaxSpeed, &result.Distance, &result.SumDistance); err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
