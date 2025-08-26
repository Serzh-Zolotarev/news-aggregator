package postgres

import (
	"context"
	"github.com/jackc/pgx/v4"
	"news-aggregator/pkg/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

const maxNewsOnPage = 10

type Store struct {
	db *pgxpool.Pool
}

func New(dbURL string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Posts(page int, filter string) (*storage.NewsList, error) {
	if page == 0 {
		page = 1
	}

	var (
		rows      pgx.Rows
		countRows pgx.Row
		err       error
	)

	offset := (page - 1) * maxNewsOnPage

	if len(filter) != 0 {
		rows, err = s.db.Query(context.Background(), `
	SELECT p.id, p.title, p.content, p.pub_time, p.link 
	FROM posts p
	WHERE p.title iLIKE $1
	ORDER BY id DESC
	OFFSET $2
	LIMIT $3
	`,
			"%"+filter+"%",
			offset,
			maxNewsOnPage,
		)

		countRows = s.db.QueryRow(context.Background(), `
	SELECT COUNT(*) 
	FROM posts p
	WHERE p.title iLIKE $1
	`,
			"%"+filter+"%",
		)
	} else {
		rows, err = s.db.Query(context.Background(), `
	SELECT p.id, p.title, p.content, p.pub_time, p.link 
	FROM posts p
	ORDER BY id DESC 
	OFFSET $1
	LIMIT $2
	`,
			offset,
			maxNewsOnPage,
		)

		countRows = s.db.QueryRow(context.Background(), `SELECT COUNT(*) FROM posts`)
	}

	if err != nil {
		return nil, err
	}

	var count int
	err = countRows.Scan(&count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []storage.NewsShortDetailed

	for rows.Next() {
		var p storage.NewsShortDetailed
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &storage.NewsList{
		Pagination: storage.Pagination{
			Page:       page,
			MaxPages:   count / maxNewsOnPage,
			NewsOnPage: maxNewsOnPage,
		},
		News: posts,
	}, rows.Err()
}

func (s *Store) Post(id int) (*storage.NewsShortDetailed, error) {
	row := s.db.QueryRow(context.Background(), `
	SELECT p.id, p.title, p.content, p.pub_time, p.link 
	FROM posts p
	WHERE id=$1
	`,
		id,
	)

	post := &storage.NewsShortDetailed{}
	err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.PubTime,
		&post.Link,
	)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *Store) AddPost(p storage.NewsShortDetailed) error {
	_, err := s.db.Exec(context.Background(), ` 
	INSERT INTO posts (title, content, pub_time, link)
	VALUES ($1, $2, $3, $4)
	`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
	)

	return err
}

func (s *Store) AddPosts(p []storage.NewsShortDetailed) error {
	for _, post := range p {
		err := s.AddPost(post)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) UpdatePost(p storage.NewsShortDetailed) error {
	_, err := s.db.Exec(context.Background(), ` 
	UPDATE posts
	SET  title = $1, content = $2, pub_time = $3, link = $4
	WHERE id = $5
	`,
		p.Title,
		p.Content,
		p.PubTime,
		p.Link,
		p.ID,
	)
	return err
}

func (s *Store) DeletePost(p storage.NewsShortDetailed) error {
	_, err := s.db.Exec(context.Background(), `
    DELETE FROM posts 
	WHERE id = $1
	`, p.ID)
	return err
}
