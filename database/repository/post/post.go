package post

import (
	"articles-api/utils"
	model "articles-api/models/post"
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PostRepository struct {
	pool *pgxpool.Pool
}

type IPostRepository interface {
	CreatePost(current_userId uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error)
	CreateComment(current_userId uuid.UUID, parent_id uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error)
	UpdatePostById(current_id uuid.UUID, post_id uuid.UUID, data map[string]interface{}) error
	GetPostCommentsById(current_userId *uuid.UUID, id uuid.UUID, page string) ([]model.FetchPostModel, error)
	DeletePostById(user_id uuid.UUID, id uuid.UUID) error
	GetSinglePostById(current_userId *uuid.UUID, id uuid.UUID) (*model.FetchPostModel, error)
	GetAllPost(current_userId *uuid.UUID, page string, category *string) ([]model.FetchPostModel, error)
	GetUserPostsById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]model.FetchPostModel, error)
	LikePostById(current_userId uuid.UUID, post_id uuid.UUID) (string, error)
	GetCategories(page string) ([]map[string]interface{}, error)
}

func NewUserRepo(pool *pgxpool.Pool) IPostRepository {
	return &PostRepository{
		pool: pool,
	}
}

func (postRepo *PostRepository) GetCategories(page string) ([]map[string]interface{}, error) {

	ctx := context.Background()

	sql := `
		SELECT category, COUNT(*) as post_count FROM posts
		GROUP BY category
		ORDER BY post_count DESC
		LIMIT 15 
        OFFSET $1 * 15
	`
	rows, err := postRepo.pool.Query(ctx, sql, page)

	if err != nil {
		return []map[string]interface{}{}, err
	}

	defer rows.Close()

	posts := []map[string]interface{}{}

	for rows.Next() {

		var (
			category   string
			post_count int
		)

		if err := rows.Scan(
			&category,
			&post_count,
		); err != nil {
			return []map[string]interface{}{}, err
		}

		post := map[string]interface{}{
			"category":   category,
			"post_count": post_count,
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (postRepo *PostRepository) LikePostById(current_userId uuid.UUID, post_id uuid.UUID) (string, error) {

	ctx := context.Background()

	sql := `
	DO $$
	DECLARE
    	post_user_id UUID;
	BEGIN
		SELECT user_id INTO post_user_id FROM posts WHERE id = $2;
		IF EXISTS (
			SELECT 1
			FROM blocking
			WHERE (user_id = '${value1}$' AND blocking_id = post_user_id) OR user_id = post_user_id AND blocking_id = '${value1}$'
		) THEN
			RAISE NOTICE 'user cannot access this post' ;
		ELSE 
			IF EXISTS (SELECT 1 FROM likes WHERE user_id = '${value1}$' AND post_id = '${value2}$' ) THEN
				DELETE FROM likes WHERE user_id = '${value1}$' AND post_id = '${value2}$';
				UPDATE posts SET like_count = GREATEST(like_count - 1 ,0) WHERE user_id = '${value1}$' AND id = '${value2}$';
			ELSE
				INSERT INTO likes(user_id ,post_id) VALUES('${value1}$' , '${value2}$');
				UPDATE posts SET like_count = like_count + 1 WHERE user_id = '${value1}$' AND id = '${value2}$';
			END IF;
		END IF;
	END $$;`

	sql = strings.ReplaceAll(sql, "${value1}$", current_userId.String())
	sql = strings.ReplaceAll(sql, "${value2}$", post_id.String())

	command, err := postRepo.pool.Exec(ctx, sql)

	if err != nil {
		return "", err
	}

	return command.String(), nil
}

func (postRepo *PostRepository) GetUserPostsById(current_userId *uuid.UUID, user_id uuid.UUID, page string) ([]model.FetchPostModel, error) {

	ctx := context.Background()

	sql := `
		SELECT p.id, p.parent_id, u.username, u.profile_img, p.title, p.content, p.like_count , p.created_at, p.updated_at FROM posts AS p
		LEFT JOIN users AS u
		ON u.id = p.user_id
	`

	pagination := `
		ORDER BY p.created_at DESC
		LIMIT 15 
        OFFSET $2 * 15
	`
	condition := "WHERE user_id = $1"

	var filter string

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
			ON r.related = p.user_id
		`
		condition += " AND r.related IS NULL"

		parameters = append(parameters, current_userId)
	}

	sql += filter + condition + pagination

	rows, err := postRepo.pool.Query(ctx, sql, parameters...)

	if err != nil {
		return []model.FetchPostModel{}, err
	}

	defer rows.Close()

	posts := utils.ExtractPostsFromRows(rows)

	return posts, nil

}
func (postRepo *PostRepository) CreateComment(current_userId uuid.UUID, parent_id uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error) {

	checkQuery := `
			DO $$
			DECLARE
				post_user_id UUID;
			BEGIN
				SELECT user_id INTO post_user_id FROM posts WHERE id = ${value2}$;
				IF EXISTS (
					SELECT 1
					FROM blocking
					WHERE (user_id = '${value1}$' AND blocking_id = post_user_id) OR user_id = post_user_id AND blocking_id = '${value1}$'
				) THEN
					RAISE NOTICE 'user cannot access this post' ;
				END IF;
			END $$;`

	checkQuery = strings.ReplaceAll(checkQuery, "${value1}$", current_userId.String())
	checkQuery = strings.ReplaceAll(checkQuery, "${value2}$", parent_id.String())

	_, err := postRepo.pool.Exec(context.Background(), checkQuery)

	if err != nil {
		return &model.FetchPostModel{}, err
	}

	ctx := context.Background()

	sql := `
		INSERT INTO posts(user_id, parent_id ,title, content) VALUES($1 ,$2 ,$3 ,$4) RETURNING id, parent_id, title, content
	`
	var post model.FetchPostModel

	err = postRepo.pool.QueryRow(ctx, sql, current_userId, parent_id, data.Title, data.Content).Scan(&post.Id, &post.Parent_id, &post.Title, &post.Content)

	if err != nil {
		return &post, err
	}

	return &post, nil
}

func (postRepo *PostRepository) CreatePost(current_userId uuid.UUID, data *model.CreatePostModel) (*model.FetchPostModel, error) {

	ctx := context.Background()

	sql := `
		INSERT INTO posts(user_id ,title,content) VALUES($1 ,$2 ,$3) RETURNING id, title ,content
	`
	var post model.FetchPostModel

	err := postRepo.pool.QueryRow(ctx, sql, current_userId, data.Title, data.Content).Scan(&post.Id, &post.Title, &post.Content)

	if err != nil {
		return &post, err
	}

	return &post, nil
}

func (postRepo *PostRepository) UpdatePostById(current_userId uuid.UUID, post_id uuid.UUID, data map[string]interface{}) error {

	ctx := context.Background()

	sql, parameters := utils.SqlUpdateQuery(data, "posts", map[string]interface{}{
		"id":      post_id,
		"user_id": current_userId,
	})

	_, err := postRepo.pool.Exec(ctx, sql, parameters...)

	if err != nil {
		return err
	}

	return nil
}

func (postRepo *PostRepository) GetPostCommentsById(current_userId *uuid.UUID, post_id uuid.UUID, page string) ([]model.FetchPostModel, error) {

	ctx := context.Background()

	sql := `
		SELECT p.id, p.parent_id, u.username, u.profile_img, p.title, p.content,p.like_count,  p.created_at, p.updated_at FROM posts AS p
		LEFT JOIN users AS u
		ON u.id = p.user_id
	`
	var filter string

	condition := "WHERE p.parent_id = $2 "

	pagination := `
		ORDER BY p.created_at DESC
		LIMIT 15 
        OFFSET $1 * 15
	`
	parameters := []interface{}{page, post_id}

	if current_userId != nil {
		filter = `
			LEFT JOIN (SELECT 
                CASE
                    WHEN user_id = $3 THEN blocking_id
                    WHEN blocking_id = $3 THEN user_id
                    ELSE NULL
                END as related
            FROM blocking) AS r
			ON r.related = p.user_id
		`
		condition += " AND r.related IS NULL"

		parameters = append(parameters, current_userId)
	}

	sql += filter + condition + pagination

	rows, err := postRepo.pool.Query(ctx, sql, parameters...)

	if err != nil {
		return []model.FetchPostModel{}, err
	}

	defer rows.Close()

	posts := utils.ExtractPostsFromRows(rows)

	return posts, nil
}

func (postRepo *PostRepository) DeletePostById(current_userId uuid.UUID, id uuid.UUID) error {

	ctx := context.Background()

	sql := `
		DELETE FROM posts p
		WHERE p.user_id = $1 AND p.id = $2
	`
	_, err := postRepo.pool.Exec(ctx, sql, current_userId, id)

	if err != nil {
		return err
	}

	return nil
}

func (postRepo *PostRepository) GetSinglePostById(current_userId *uuid.UUID, post_id uuid.UUID) (*model.FetchPostModel, error) {

	if current_userId != nil {
		checkQuery := `
			DO $$
			DECLARE
				post_user_id UUID;
			BEGIN
				SELECT user_id INTO post_user_id FROM posts WHERE id = ${value2}$;
				IF EXISTS (
					SELECT 1
					FROM blocking
					WHERE (user_id = '${value1}$' AND blocking_id = post_user_id) OR user_id = post_user_id AND blocking_id = '${value1}$'
				) THEN
					RAISE NOTICE 'user cannot access this post' ;
				END IF;
			END $$;`

		checkQuery = strings.ReplaceAll(checkQuery, "${value1}$", current_userId.String())
		checkQuery = strings.ReplaceAll(checkQuery, "${value2}$", post_id.String())

		_, err := postRepo.pool.Exec(context.Background(), checkQuery)

		if err != nil {
			return &model.FetchPostModel{}, err
		}
	}

	ctx := context.Background()

	sql := `
		SELECT p.id, p.parent_id, u.username, u.profile_img, p.title, p.content, p.like_count , p.created_at, p.updated_at FROM posts AS p
		LEFT JOIN users AS u
		ON u.id = p.user_id
		WHERE p.id = $1
	`
	var post model.FetchPostModel

	err := postRepo.pool.QueryRow(ctx, sql, post_id).Scan(
		&post.Id,
		&post.Parent_id,
		&post.Username,
		&post.Profile_img,
		&post.Title,
		&post.Content,
		&post.Like_count,
		&post.Created_at,
		&post.Updated_at,
	)

	if err != nil {
		return &post, err
	}

	return &post, nil
}

func (postRepo *PostRepository) GetAllPost(current_userId *uuid.UUID, page string, category *string) ([]model.FetchPostModel, error) {

	ctx := context.Background()

	sql := `
		SELECT p.id, p.parent_id, u.username, u.profile_img, p.title, p.content, p.like_count ,p.created_at, p.updated_at FROM posts AS p
		LEFT JOIN users AS u
		ON u.id = p.user_id
	`
	var filter string

	condition := "WHERE "

	parameters := []interface{}{page}

	pagination := `
		ORDER BY p.created_at DESC
		LIMIT 15 
        OFFSET $1 * 15
	`
	if category != nil {
		condition += "p.category = $2"
		parameters = append(parameters, category)
	}

	if current_userId != nil {
		filter = `
			LEFT JOIN (SELECT 
                CASE
                    WHEN user_id = $2 THEN blocking_id
                    WHEN blocking_id = $2 THEN user_id
                    ELSE NULL
                END as related
            FROM blocking) AS r
			ON r.related = p.user_id
		`

		if len(condition) > 6 {
			condition += "AND r.related IS NULL"
		} else {
			condition += "r.related IS NULL"
		}

		parameters = append(parameters, current_userId)

	}

	sql += filter + condition + pagination

	rows, err := postRepo.pool.Query(ctx, sql, parameters...)

	if err != nil {
		return []model.FetchPostModel{}, err
	}

	defer rows.Close()

	posts := utils.ExtractPostsFromRows(rows)

	return posts, nil
}
