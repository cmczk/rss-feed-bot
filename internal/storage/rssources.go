package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/cmczk/rss-feed-bot/internal/models"
)

type RSSourcesPostgresStorage struct {
	db *sql.DB
}

type dbRSSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}

func (r *RSSourcesPostgresStorage) RSSources(ctx context.Context) ([]models.RSSource, error) {
	const op = "storage.rssources.RSSources"

	conn, err := r.db.Conn(ctx)
	if err != nil {
		log.Printf("%s: cannot connect to db: %s", op, err.Error())
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}
	defer conn.Close()

	var rssources []models.RSSource
	res, err := conn.QueryContext(ctx, "SELECT * FROM rssources")
	if err != nil {
		return nil, fmt.Errorf("query to select all the rssources failed: %w", err)
	}
	defer res.Close()

	for res.Next() {
		var rssource dbRSSource
		if err := res.Scan(&rssource.ID, &rssource.Name, &rssource.FeedURL, &rssource.CreatedAt); err != nil {
			log.Printf("%s: failed to scan row: %s", op, err.Error())
			return nil, fmt.Errorf("%s: failed to scan row: %w", op, err)
		}
		rssources = append(rssources, models.RSSource(rssource))
	}

	return rssources, nil
}

func (r *RSSourcesPostgresStorage) RSSourceByID(ctx context.Context, id int64) (*models.RSSource, error) {
	const op = "storage.rssources.RSSourceByID"

	conn, err := r.db.Conn(ctx)
	if err != nil {
		log.Printf("%s: cannot connect to db: %s", op, err.Error())
		return nil, fmt.Errorf("cannot connect to db: %w", err)
	}
	defer conn.Close()

	var rssource dbRSSource
	if err := conn.QueryRowContext(
		ctx,
		"SELECT * FROM rssources WHERE id = $1 LIMIT 1",
		id,
	).Scan(
		&rssource.ID,
		&rssource.Name,
		&rssource.FeedURL,
		&rssource.CreatedAt,
	); err != nil {
		log.Printf("%s: cannot scan rssource: %s", op, err.Error())
		return nil, fmt.Errorf("%s: cannot scan rssource: %w", op, err)
	}

	rssourceModel := models.RSSource(rssource)

	return &rssourceModel, nil
}

func (r *RSSourcesPostgresStorage) Add(ctx context.Context, rssource models.RSSource) (int64, error) {
	const op = "storage.rssources.Add"

	conn, err := r.db.Conn(ctx)
	if err != nil {
		log.Printf("%s: cannot connect to db: %s", op, err.Error())
		return 0, fmt.Errorf("cannot connect to db: %w", err)
	}
	defer conn.Close()

	res, err := conn.ExecContext(
		ctx,
		"INSERT INTO rssources (name, feed_url) VALUES ($1, $2)",
		rssource.Name, rssource.FeedURL,
	)
	if err != nil {
		log.Printf("%s: cannot insert new rssource: %s", op, err.Error())
		return 0, fmt.Errorf("%s: cannot insert new rssource: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("%s: cannot get last inserted id: %s", op, err.Error())
		return 0, fmt.Errorf("%s: cannot get last inserted id: %w", op, err)
	}

	return id, nil
}

func (r *RSSourcesPostgresStorage) Delete(ctx context.Context, id int64) error {
	const op = "storage.rssources.Delete"

	conn, err := r.db.Conn(ctx)
	if err != nil {
		log.Printf("%s: cannot connect to db: %s", op, err.Error())
		return fmt.Errorf("cannot connect to db: %w", err)
	}
	defer conn.Close()

	_, err = conn.ExecContext(
		ctx,
		"DELETE FROM rssources WHERE id = $1",
		id,
	)
	if err != nil {
		log.Printf("%s: cannot delete rssource: %s", op, err.Error())
		return fmt.Errorf("cannot delete rssource: %w", err)
	}

	return nil
}
