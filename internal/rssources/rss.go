package rssources

import (
	"context"
	"fmt"
	"log"

	"github.com/cmczk/rss-feed-bot/internal/models"

	"github.com/SlyMarbo/rss"
)

type RSSourceClient struct {
	URL          string
	RSSourceID   int64
	RSSourceName string
}

func NewFromModel(m models.RSSource) RSSourceClient {
	return RSSourceClient{
		URL:          m.FeedURL,
		RSSourceID:   m.ID,
		RSSourceName: m.Name,
	}
}

func (r RSSourceClient) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := r.loadFeed(ctx, r.URL)
	if err != nil {
		log.Println("failed to fetch RSS feed", err)
		return nil, fmt.Errorf("failed to fetch RSS feed: %s", err.Error())
	}

	items := make([]models.Item, 0, len(feed.Items))
	for _, it := range feed.Items {
		items = append(items, models.Item{
			Title:        it.Title,
			Categories:   it.Categories,
			Link:         it.Link,
			Date:         it.Date,
			Summary:      it.Summary,
			RSSourceName: r.RSSourceName,
		})
	}

	return items, nil
}

func (r RSSourceClient) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var (
		feedCh = make(chan *rss.Feed)
		errCh  = make(chan error)
	)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errCh <- err
			return
		}

		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}
