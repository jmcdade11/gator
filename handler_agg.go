package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmcdade11/gator/internal/database"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <time_between_reqs>", cmd.Name)
	}

	timeBetweenReqsArg := cmd.Args[0]
	timeBetweenRequests, err := time.ParseDuration(timeBetweenReqsArg)
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) error {
	nextToFetch, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	lastFetchedAt := sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		ID:            nextToFetch.ID,
		LastFetchedAt: lastFetchedAt,
		UpdatedAt:     time.Now().UTC(),
	})

	if err != nil {
		return err
	}

	rssFeed, err := fetchFeed(nextToFetch.Url)
	if err != nil {
		return err
	}
	for _, item := range rssFeed.Channel.Item {
		nullDescription := sql.NullString{
			String: item.Description,
		}
		publishedAt := sql.NullTime{}
		if t, err := time.Parse(time.RFC1123Z, item.PubDate); err == nil {
			publishedAt = sql.NullTime{
				Time:  t,
				Valid: true,
			}
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: nullDescription,
			Url:         item.Link,
			PublishedAt: publishedAt,
			FeedID:      nextToFetch.ID,
		})

		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}

	return nil
}
