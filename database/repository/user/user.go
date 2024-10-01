package user

import (
	model "articles-api/models/user"
	"articles-api/utils"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

type IUserRepository interface {
	CreateUser(data *model.RegisterUserModel) (map[string]interface{}, error)
	GetUserByEmail(email string) (*model.FetchUserModel, error)
	GetAllUsers(current_userId *uuid.UUID, page string) ([]model.FetchUserModel, error)
	DeleteMe(current_userId uuid.UUID) (string, error)
	GetUserById(current_userId *uuid.UUID, user_id uuid.UUID) (*model.FetchUserModel, error)
	UpdateUserById(current_userId uuid.UUID, data map[string]interface{}) error
	GetMyBlockList(current_userId uuid.UUID, page string) ([]map[string]interface{}, error)
	GetUserFollowingListById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]map[string]interface{}, error)
	GetUserFollowersListById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]map[string]interface{}, error)
	FollowUserById(current_userId uuid.UUID, user_id uuid.UUID) (string, error)
	BlockUserById(current_userId uuid.UUID, user_id uuid.UUID) (string, error)
}

func NewUserRepo(pool *pgxpool.Pool) IUserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (userRepo *UserRepository) BlockUserById(current_userId uuid.UUID, id uuid.UUID) (string, error) {
	sql := `
	DO $$
	BEGIN
	    IF EXISTS (SELECT 1 FROM blocking WHERE user_id = '${value1}$' AND blocking_id = '${value2}$') THEN
	        DELETE FROM blocking WHERE user_id = '${value1}$' AND blocking_id = '${value2}$';
		ELSE
			INSERT INTO blocking(user_id ,blocking_id) VALUES('${value1}$' , '${value2}$');
		END IF;
	END $$;`

	sql = strings.ReplaceAll(sql, "${value1}$", current_userId.String())
	sql = strings.ReplaceAll(sql, "${value2}$", id.String())

	ctx := context.Background()

	command, err := userRepo.pool.Exec(ctx, sql)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return command.String(), nil

}

func (userRepo *UserRepository) GetMyBlockList(current_userId uuid.UUID, page string) ([]map[string]interface{}, error) {

	ctx := context.Background()

	sql := `SELECT id, username, profile_img, email, password FROM blocking WHERE id = $1 `

	pagination := `
		LIMIT 15 
        OFFSET $2 * 15
	`

	sql += pagination

	rows, err := userRepo.pool.Query(ctx, sql, current_userId, page)

	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf(err.Error())
	}

	defer rows.Close()

	users := []map[string]interface{}{}

	for rows.Next() {

		var (
			id          uuid.UUID
			username    string
			profile_img *string
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := rows.Scan(
			&id,
			&username,
			&profile_img,
			&createdAt,
			&updatedAt,
		); err != nil {
			return []map[string]interface{}{}, fmt.Errorf(err.Error())
		}

		user := map[string]interface{}{
			"id":          id,
			"username":    username,
			"profile_img": profile_img,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		}

		users = append(users, user)
	}

	return users, nil

}

func (userRepo *UserRepository) GetUserFollowingListById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]map[string]interface{}, error) {

	sql := `
		SELECT u.username, u.profile_img FROM following AS f
		LEFT JOIN users AS u
		ON f.following_id = u.id
	`
	var filter string

	condition := "WHERE f.user_id = $1 "

	pagination := `
		LIMIT 15 
        OFFSET $2 * 15
	`

	parameters := []interface{}{user_id, page}

	if current_userId != nil {
		filter = `
			LEFT JOIN (SELECT 
                CASE
                    WHEN user_id = $3 THEN blocking_id
                    WHEN blocking_id = $3 THEN user_id
                    ELSE NULL
                END as related
            FROM blocking) AS r
		`
		condition += "AND r.related == NULL"

		parameters = append(parameters, current_userId)
	}

	sql += filter + condition + pagination

	ctx := context.Background()

	rows, err := userRepo.pool.Query(ctx, sql, parameters...)

	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf(err.Error())
	}

	defer rows.Close()

	users := []map[string]interface{}{}

	for rows.Next() {

		var (
			id          uuid.UUID
			username    string
			profile_img *string
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := rows.Scan(
			&id,
			&username,
			&profile_img,
			&createdAt,
			&updatedAt,
		); err != nil {
			return []map[string]interface{}{}, fmt.Errorf(err.Error())
		}
		user := map[string]interface{}{
			"id":          id,
			"username":    username,
			"profile_img": profile_img,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		}
		users = append(users, user)
	}

	return users, nil
}

func (userRepo *UserRepository) GetUserFollowersListById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]map[string]interface{}, error) {

	sql := `
		SELECT u.username, u.profile_img FROM following AS f
		LEFT JOIN users AS u
		ON f.user_id = u.id
	`
	var filter string

	condition := "WHERE f.following_id = $1 "

	parameters := []interface{}{user_id, page}

	if current_userId != nil {
		filter = `
			LEFT JOIN (SELECT 
                CASE
                    WHEN user_id = $3 THEN blocking_id
                    WHEN blocking_id = $3 THEN user_id
                    ELSE NULL
                END as related
            FROM blocking) AS r
		`
		condition += "AND r.related IS NULL"

		parameters = append(parameters, current_userId)
	}

	sql += filter + condition

	ctx := context.Background()

	rows, err := userRepo.pool.Query(ctx, sql, parameters...)

	if err != nil {
		return []map[string]interface{}{}, fmt.Errorf(err.Error())
	}

	defer rows.Close()

	users := []map[string]interface{}{}

	for rows.Next() {

		var (
			id          uuid.UUID
			username    string
			profile_img *string
			createdAt   time.Time
			updatedAt   time.Time
		)

		if err := rows.Scan(
			&id,
			&username,
			&profile_img,
			&createdAt,
			&updatedAt,
		); err != nil {
			return []map[string]interface{}{}, fmt.Errorf(err.Error())
		}
		user := map[string]interface{}{
			"id":          id,
			"username":    username,
			"profile_img": profile_img,
			"created_at":  createdAt,
			"updated_at":  updatedAt,
		}
		users = append(users, user)
	}

	return users, nil
}

func (userRepo *UserRepository) FollowUserById(current_userId uuid.UUID, user_id uuid.UUID) (string, error) {

	sql := `
	DO $$
	BEGIN
		IF EXISTS (
			SELECT 1
			FROM blocking
			WHERE (user_id = '${value1}$' AND blocking_id = '${value2}$') OR user_id = '${value2}$' AND blocking_id = '${value1}$'
		) THEN
			RAISE NOTICE 'user cannot access this user' ;
		ELSE
			IF EXISTS (SELECT 1 FROM following WHERE user_id = '${value1}$' AND following_id = '${value2}$') THEN
				DELETE FROM following WHERE user_id = '${value1}$' AND following_id = '${value2}$';
				UPDATE users SET follower_count = GREATEST(follower_count - 1 ,0) WHERE id = '${value1}$';
				UPDATE users SET following_count = GREATEST(following_count - 1 ,0) WHERE id = '${value2}$';
			ELSE
				INSERT INTO following(user_id ,following_id) VALUES('${value1}$' , '${value2}$');
				UPDATE users SET follower_count = follower_count + 1 WHERE id = '${value1}$';
				UPDATE users SET following_count = following_count + 1 WHERE id = '${value2}$';
			END IF;
		END IF;
	END $$;`

	sql = strings.ReplaceAll(sql, "${value1}$", current_userId.String())
	sql = strings.ReplaceAll(sql, "${value2}$", user_id.String())

	ctx := context.Background()

	command, err := userRepo.pool.Exec(ctx, sql)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return command.String(), nil
}

func (userRepo *UserRepository) UpdateUserById(current_userId uuid.UUID, data map[string]interface{}) error {

	ctx := context.Background()

	sql, parameters := utils.SqlUpdateQuery(data, "users", map[string]interface{}{
		"id": current_userId,
	})

	_, err := userRepo.pool.Exec(ctx, sql, parameters...)

	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}

func (userRepo *UserRepository) GetAllUsers(current_userId *uuid.UUID, page string) ([]model.FetchUserModel, error) {

	ctx := context.Background()

	sql := `
		SELECT id, username, profile_img , email , followers_cnt, following_cnt ,created_at, updated_at FROM users AS u
	`
	var (
		filter    string
		condition string
	)

	pagination := `
		LIMIT 15 
        OFFSET $1 * 15
	`

	parameters := []interface{}{page}

	if current_userId != nil {
		filter = `
			LEFT JOIN (SELECT 
                CASE
                    WHEN user_id = $2 THEN blocking_id
                    WHEN blocking_id = $2 THEN user_id
                    ELSE NULL
                END as related
            FROM blocking) AS r
			ON r.related = u.id
		`
		condition += "WHERE r.related IS NULL"

		parameters = append(parameters, current_userId)
	}

	sql += filter + condition + pagination

	rows, err := userRepo.pool.Query(ctx, sql, parameters...)

	if err != nil {
		return []model.FetchUserModel{}, fmt.Errorf(err.Error())
	}

	defer rows.Close()

	users := utils.ExtractUsersFromRows(rows)

	return users, nil

}

func (userRepo *UserRepository) GetUserById(current_userId *uuid.UUID, user_id uuid.UUID) (*model.FetchUserModel, error) {

	if current_userId != nil {
		checkQuery := `
			SELECT 1
			FROM blocking
			WHERE (user_id = $1 AND blocking_id = $2) OR user_id = $2 AND blocking_id = $1
		`
		var check bool

		userRepo.pool.QueryRow(context.Background(), checkQuery, current_userId, user_id).Scan(&check)

		if check {
			return &model.FetchUserModel{}, fmt.Errorf("not accessible")
		}
	}

	ctx := context.Background()

	sql := `SELECT id, username, profile_img, email, followers_cnt, following_cnt, adress, age, gender FROM users WHERE id = $1`

	rows := userRepo.pool.QueryRow(ctx, sql, user_id)

	var user model.FetchUserModel

	if err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.Profile_img,
		&user.Email,
		&user.Followers_cnt,
		&user.Following_cnt,
		&user.Adress,
		&user.Age,
		&user.Gender,
	); err != nil {
		return &model.FetchUserModel{}, fmt.Errorf(err.Error())
	}

	return &user, nil
}

func (userRepo *UserRepository) CreateUser(data *model.RegisterUserModel) (map[string]interface{}, error) {
	ctx := context.Background()

	sql := `INSERT INTO users(username,email,password) VALUES($1 ,$2, $3) RETURNING id, username ,email`

	rows := userRepo.pool.QueryRow(ctx, sql, data.Username, data.Email, data.Password)

	var (
		id       uuid.UUID
		username string
		email    string
	)

	if err := rows.Scan(
		&id,
		&username,
		&email,
	); err != nil {
		return map[string]interface{}{}, fmt.Errorf(err.Error())
	}

	user := map[string]interface{}{
		"id":       id,
		"username": username,
		"email":    email,
	}

	return user, nil
}

func (userRepo *UserRepository) GetUserByEmail(user_email string) (*model.FetchUserModel, error) {

	ctx := context.Background()

	sql := `SELECT id, username, email, password FROM users WHERE email = $1`

	rows := userRepo.pool.QueryRow(ctx, sql, user_email)

	user := &model.FetchUserModel{}

	if err := rows.Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Password,
	); err != nil {
		return &model.FetchUserModel{}, fmt.Errorf(err.Error())
	}

	return user, nil
}

func (userRepo *UserRepository) DeleteMe(current_userId uuid.UUID) (string, error) {
	ctx := context.Background()

	sql := `
		DELETE FROM users
		WHERE id = $1
	`
	command, err := userRepo.pool.Exec(ctx, sql, current_userId)

	if err != nil {
		return "", fmt.Errorf(err.Error())
	}

	return command.String(), nil
}
