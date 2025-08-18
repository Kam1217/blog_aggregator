package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Kam1217/blog_aggregator/internal/config"
	"github.com/Kam1217/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	cfg        *config.Config
	cfgManager *config.ConfigManager
	db         *database.Queries
}

type command struct {
	name string
	args []string
}

type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to fetch the next feed: %w", err)
	}
	_, err = s.db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		return fmt.Errorf("failed to mark the fetched feed: %w", err)
	}
	data, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to make a HTTP request: %w", err)
	}
	for _, post := range data.Channel.Item {
		t, err := time.Parse(time.RFC1123Z, post.PubDate)
		if err != nil {
			fmt.Printf("could not parse pubDate: %v\n", err)
			continue
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Title:     post.Title,
			Url:       post.Link,
			Description: sql.NullString{
				String: post.Description,
				Valid:  true,
			},
			PublishedAt: t,
			FeedID:      feed.ID,
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

func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login handler expects a single argument but got an empty slice")
	}
	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			os.Exit(1)
		} else {
			return fmt.Errorf("database error getting user: %w", err)
		}
	}
	if err := s.cfgManager.SetUser(s.cfg, cmd.args[0]); err != nil {
		return fmt.Errorf("error setting the username to config: %w", err)
	}
	fmt.Printf("Username has been set to: %s\n", s.cfg.CurrentUserName)
	return nil
}

func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.registeredCommands[cmd.name]
	if !exists {
		return fmt.Errorf("command does not exist: %v", cmd.name)
	}
	if err := handler(s, cmd); err != nil {
		return fmt.Errorf("error calling the command: %w", err)
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("login handler expects a single argument but got an empty slice")
	}
	user, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if user.Name != "" {
		os.Exit(1)
	}
	newUser, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	if err := s.cfgManager.SetUser(s.cfg, newUser.Name); err != nil {
		return fmt.Errorf("error setting the config username: %w", err)
	}

	fmt.Printf("New user has been created %v:\n", newUser)
	return nil
}

func handlerReset(s *state, _ command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleteing users: %w", err)
	}
	fmt.Println("succesfully deleted users")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error getting users: %w", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}

	}
	return nil
}

func handlerAgg(s *state, cmd command) error {
	timeBetweenRequests, err := time.ParseDuration(cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to set time between requests: %w", err)
	}
	fmt.Printf("Collecting feeds every %v\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)
	for ; ; <-ticker.C {
		err = scrapeFeeds(s)
		if err != nil {
			fmt.Printf("failed to scrape the feed: %v", err)
		}
	}
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 2 {
		return fmt.Errorf("add feed needs 2 arguments")
	}
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}
	_, err = s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to auto-follow new feed: %w", err)
	}
	fmt.Println(feed)
	return nil
}

func handlerListFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %w", err)
	}
	if len(feeds) == 0 {
		fmt.Println("There is no feeds in the database, try adding a feed")
	}
	for _, feed := range feeds {
		fmt.Printf("%s\n", feed.Name)
		fmt.Printf("%s\n", feed.Url)
		fmt.Printf("%s\n", feed.Name_2)
	}
	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.args) != 1 {
		return fmt.Errorf("follow command requires 1 argument")
	}
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed by url")
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to follow feed: %w", err)
	}
	fmt.Printf("%s is now following %s\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get follows: %w", err)
	}
	if len(follows) == 0 {
		fmt.Println("There are currently no follows")
		return nil
	}
	for _, follow := range follows {
		fmt.Println(follow.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	feed, err := s.db.GetFeedByURL(context.Background(), cmd.args[0])
	if err != nil {
		return fmt.Errorf("failed to get feed by url")
	}
	err = s.db.DeleteFollows(context.Background(), database.DeleteFollowsParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}
	fmt.Println("succesfully unfollowed feed")
	return nil
}

func handlerBrowse(s *state, cmd command) error {
	limit := 2 // default
	if len(cmd.args) > 0 {
		if l, err := strconv.Atoi(cmd.args[0]); err == nil {
			limit = l
		}
	}

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	posts, err := s.db.GetPostForUser(context.Background(), database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("failed to get posts for user: %w", err)
	}
	return nil
}
