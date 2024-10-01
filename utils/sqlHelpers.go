package utils

import (
	post_model "articles-api/models/post"
	user_model "articles-api/models/user"
	"fmt"
	"github.com/jackc/pgx/v4"
)

func SqlUpdateQuery(fields map[string]interface{}, table string, conditions map[string]interface{}) (string, []interface{}) {

	sql := fmt.Sprintf("UPDATE %s SET", table)

	condition := " WHERE "

	size := len(fields)

	count := 0

	parameters := make([]interface{}, 0)

	for key, value := range fields {
		s := fmt.Sprintf(" %s = $%v ", key, count+1)
		parameters = append(parameters, value)
		sql += s
		count++
		if count < size {
			sql += ","
		}
	}

	for key, value := range conditions {
		s := fmt.Sprintf(" %s = $%v ", key, count+1)
		parameters = append(parameters, value)
		condition += s
		count++
		if count < size+len(conditions) {
			condition += " AND "
		}
	}
	sql += condition

	return sql, parameters
}

func ExtractPostsFromRows(rows pgx.Rows) []post_model.FetchPostModel {

	posts := []post_model.FetchPostModel{}

	for rows.Next() {

		var post post_model.FetchPostModel

		if err := rows.Scan(
			&post.Id,
			&post.Parent_id,
			&post.Username,
			&post.Profile_img,
			&post.Title,
			&post.Content,
			&post.Like_count,
			&post.Created_at,
			&post.Updated_at,
		); err != nil {
			continue
		}

		posts = append(posts, post)
	}
	return posts
}


func ExtractUsersFromRows(rows pgx.Rows) []user_model.FetchUserModel {

	users := []user_model.FetchUserModel{}

	for rows.Next() {

		var user user_model.FetchUserModel

		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Profile_img,
			&user.Email,
			&user.Followers_cnt,
			&user.Following_cnt,
			&user.Created_at,
			&user.Updated_at,
		); err != nil {
			continue
		}
	
		users = append(users, user)
	}

	return users
}
