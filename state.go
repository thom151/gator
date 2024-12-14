package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/thom151/gator/internal/config"
	"github.com/thom151/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	arg  []string
}

func handlerAgg(s *state, cmd command) error {

	timeBetween, err := time.ParseDuration(cmd.arg[0])
	if err != nil {
		return err
	}

	fmt.Printf("Collecting feeds every %s\n", timeBetween)

	ticker := time.NewTicker(timeBetween)
	fmt.Println(ticker)

	for ; ; <-ticker.C {
		fmt.Println("trying to print")
		err = scrapeFeeds(s)
		if err != nil {
			return err
		}
	}

}

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 2
	if len(cmd.arg) > 0 {
		limitArg, err := strconv.Atoi(cmd.arg[0])
		if err != nil {
			return err
		}
		limit = limitArg
	}

	postUserParams := database.GetPostForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}

	posts, err := s.db.GetPostForUser(context.Background(), postUserParams)
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title : %s\n", post.Title)
		fmt.Printf("Description: %s\n", post.Description)
		fmt.Println()
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	innerHandler := func(s *state, cmd command) error {
		currUser, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		err = handler(s, cmd, currUser)
		return nil

	}
	return innerHandler
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	url := cmd.arg[0]
	feed, err := s.db.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return err
	}

	feedByUserParams := database.UnfollowFeedByUserParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.UnfollowFeedByUser(context.Background(), feedByUserParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed Unfollowed by %s: \n", user.Name)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	for _, feed := range feedFollows {
		fmt.Printf("* %s\n", feed.FeedName)
	}

	return nil

}

func handlerFollow(s *state, cmd command, user database.User) error {
	url := cmd.arg[0]

	feed, err := s.db.GetFeedFromUrl(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollowParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feeFollorRow, err := s.db.CreateFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		return err
	}

	fmt.Printf("Feed: %s\n", feeFollorRow.FeedName)
	fmt.Printf("User: %s\n", feeFollorRow.UserName)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {

	name := cmd.arg[0]
	url := cmd.arg[1]

	feedParams := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	fmt.Printf("Feed Created!\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s",
		feed.ID,
		feed.CreatedAt,
		feed.UpdatedAt,
		feed.Name,
		feed.Url,
		feed.UserID,
	)

	_, err = s.db.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}

	fmt.Printf("Feed Follow Added!\n")
	return nil
}

func handlerFeeds(s *state, cmd command) error {
	creatorRows, err := s.db.GetFeedCreator(context.Background())
	if err != nil {
		return err
	}
	for _, row := range creatorRows {
		fmt.Printf("Feed: %s\nURL: %s\nCreator: %s\n", row.Name_2, row.Url, row.Name)
		fmt.Println()
	}

	return nil

}

func handlerReset(s *state, cmd command) error {

	if len(cmd.arg) > 0 {
		return fmt.Errorf("Too many arguments")
	}

	err := s.db.ResetUser(context.Background())
	if err != nil {
		fmt.Println("There was a problem with resetting")
		os.Exit(1)
	}

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		}
		fmt.Printf("* %s\n", user.Name)
	}
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arg) < 1 {
		return fmt.Errorf("No username, please add one")
	}

	_, err := s.db.GetUser(context.Background(), cmd.arg[0])
	if err != nil {
		fmt.Println("User not found")
		os.Exit(1)
	}

	err = s.cfg.SetUser(cmd.arg[0])
	if err != nil {
		return err
	}

	fmt.Println("User has been set!")
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arg) < 1 {
		return fmt.Errorf("Please provide a name")
	}
	ctxt := context.Background()

	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      cmd.arg[0],
	}

	//checking if name already exists
	_, err := s.db.GetUser(ctxt, userParams.Name)
	if err == nil {
		fmt.Println("User already exists")
		os.Exit(1)
	}

	user, err := s.db.CreateUser(ctxt, userParams)
	if err != nil {
		return err
	}

	s.cfg.CurrentUserName = user.Name

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User Created! Details: \n\t %s \n\t %s \n\t %s \n\t %s \n",
		user.ID, user.CreatedAt, user.UpdatedAt, user.Name)
	return nil
}

type commands struct {
	commandMap map[string]func(*state, command) error
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.commandMap[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	if s == nil {
		return fmt.Errorf("State doesn't exist")
	}

	function, ok := c.commandMap[cmd.name]
	if !ok {
		return fmt.Errorf("Invalid Command")
	}

	err := function(s, cmd)
	if err != nil {
		return err
	}

	return nil
}
