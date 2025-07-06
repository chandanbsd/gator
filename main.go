package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/chandanbsd/gator/internal/config"
	"github.com/chandanbsd/gator/internal/database"
	"github.com/chandanbsd/gator/internal/feed"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type command struct {
	Name      string
	Arguments []string
}

type commands struct {
	options map[string]func(*state, command) error
}

type state struct {
	db   *database.Queries
	conf *config.Config
}

func (c *commands) run(s *state, cmd command) error {
	for name, handler := range c.options {
		if name == cmd.Name {
			err := handler(s, cmd)
			return err
		}
	}
	return errors.New("The command you have entered is not valid")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.options[name] = f
}

func loginHandler(s *state, cmd command) error {
	if s == nil {
		fmt.Println("Please enter a command")
		os.Exit(1)
	}

	if cmd.Arguments == nil || len(cmd.Arguments) == 0 {
		fmt.Println("Please provide a username")
		os.Exit(1)
	}
	user, err := s.db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil {
		os.Exit(1)
	}
	s.conf.SetUser(user.Name)


	fmt.Printf("Your username has been set to %s\n", s.conf.CurrentUserName)

	return nil
}

func registerHandlers(coms commands) {
	coms.register("login", loginHandler)
	coms.register("register", registerHandler)
	coms.register("users", usersHandler)
	coms.register("reset", deleteHandler)
	coms.register("agg", aggHandler)
	coms.register("feeds", feedsHandler)
	coms.register("addfeed", middlewareLoggedIn(addFeedHandler))
	coms.register("follow", middlewareLoggedIn(followHandler))
	coms.register("following", middlewareLoggedIn(followingHandler))
	coms.register("unfollow", middlewareLoggedIn(deleteFeedFollowHandler))
	coms.register("browse", middlewareLoggedIn(browseHandler))
}

func deleteHandler(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	s.conf.SetUser("")

	return nil
}

func registerHandler(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		fmt.Println("Missing name argument")
		os.Exit(1)
	}

	_, err := s.db.GetUser(context.Background(), cmd.Arguments[0])
	if err == nil {
		os.Exit(1)
	}

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	currentNullTime := sql.NullTime{
		Time:  time.Now(),
		Valid: false,
	}

	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentNullTime,
		Name:      cmd.Arguments[0],
	}
	_, err = s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		fmt.Println("Failed to create the user.")
		os.Exit(1)
	}
	s.conf.SetUser(cmd.Arguments[0])

	fmt.Printf("Created user %s\n", s.conf.CurrentUserName)
	return nil
}

func usersHandler(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		os.Exit(1)
	}

	for _, user := range users {

		if user.Name == s.conf.CurrentUserName {
			fmt.Println(user.Name + " (current)")
		} else {
			fmt.Println(user.Name)
		}
	}

	return nil
}

func aggHandler(s *state, cmd command) error {

	if len(cmd.Arguments) != 1 {
		fmt.Println("Incorrect arguments list")
		os.Exit(1)
	}

	duration, err := time.ParseDuration(cmd.Arguments[0])
	if err != nil {
		fmt.Println("Unable to parse duration")
		return err	
	}

	ticker := time.NewTicker(duration)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func browseHandler(s *state, cmd command, u database.User) error {
	limit32 := int32(2)
	if len(cmd.Arguments) == 1 {
		limit64, err := strconv.ParseInt(cmd.Arguments[0], 10, 32)
		if err != nil {
			return err
		}
		limit32 = int32(limit64)
	}

	posts, err := s.db.GetPostsForUser(
		context.Background(), 
		database.GetPostsForUserParams{
			UserID: u.ID,
			Limit: int32(limit32),
		})
	if err != nil {
		fmt.Println("Unable to get posts for user")
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %v\tCreatedAt: %v\tUpdatedAt: %v\tTitle:%v\tUrl:%v\tDescription:%v\tPublishedAt:%v",
			post.CreatedAt,
			post.UpdatedAt,
			post.Title,
			post.Url,
			post.Description,
			post.PublishedAt,
		)
	}
	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func (s *state, cmd command) error{
		user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
		if err != nil {
			fmt.Println("Failed to get the current user")
			return err
		}
		return handler(s, cmd, user)
	}

}

func addFeedHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) < 2 {
		fmt.Println("Missing name argument")
		os.Exit(1)
	}
	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	currentNullTime := sql.NullTime{
		Time:  time.Now(),
		Valid: false,
	}



	newFeed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentNullTime,
		Name:      cmd.Arguments[0],
		Url: cmd.Arguments[1],
		UserID: user.ID,
	}

	_, err := s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		fmt.Println("Failed to create the feed")
		os.Exit(1)
	} else {
		fmt.Println("User created successfully")
	}

	followHandler(s, command{
		Name: "follow",
		Arguments: []string{cmd.Arguments[1]},
	}, user)
	return nil
}

func feedsHandler(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		fmt.Println("Unable to fetch feeds")
		os.Exit(1)
	}

	for _, feed := range feeds {
		fmt.Printf("name: %v url: %v, user_name: %v\n", feed.Name, feed.Url, feed.UserName)
	}
	return nil
}

func followHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		fmt.Println("Missing url or wrong arguments provided");
		os.Exit(1)
	}

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	currentNullTime := sql.NullTime{
		Time:  time.Now(),
		Valid: false,
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Arguments[0])
	if err != nil {
		fmt.Println("Feed does not exist")
		os.Exit(1)
	}

	createdFeedFollowParam := database.CreateFeedFollowParams {
		ID: uuid.New(),
		CreatedAt: currentTime,
		UpdatedAt: currentNullTime,
		FeedID: feed.ID,
		UserID: user.ID,
	}

	createdFeedFollow, err := s.db.CreateFeedFollow(context.Background(), createdFeedFollowParam)
	if err != nil {
		fmt.Println("Failed to folow the feed")
		os.Exit(1)
	}

	fmt.Printf("Feed name: %s and current user: %s\n", createdFeedFollow.FeedName, createdFeedFollow.UserName)

	return nil
}

func followingHandler(s *state, cmd command,  user database.User) error {

	feedsForUser, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		fmt.Println("Failed to fetch feeds followed by user");
		os.Exit(1)
	}

	for _, feed := range feedsForUser {
		fmt.Println(feed)
	}
	return nil
}

func deleteFeedFollowHandler(s *state, cmd command, user database.User) error {
	if len(cmd.Arguments) != 1 {
		fmt.Println("Missing url argument");
	}
	
	feedFollowParams := database.DeleteFeedFollowParams {
		UserID: user.ID,
		Url: cmd.Arguments[0],
	}


	err := s.db.DeleteFeedFollow(context.Background(), feedFollowParams)
	if err != nil {
		fmt.Println("Failed to delete the feed follow")
		os.Exit(1)
	}
	return nil
}

func scrapeFeeds(s *state) error{
	fmt.Println("*************Fetching Feed**************")
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		fmt.Println("Failed to fetch the next feed")
		os.Exit(1)
	}

	markFeedFetchedParams := database.MarkFeedFetchedParams {
		LastFetchedAt: sql.NullTime{
			Time: time.Now(),
		},
		ID: nextFeed.ID,
	}

	rssFeed, err := feed.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		fmt.Println("Failed to fetch feed")
		return err
	}
	
	for _, item := range rssFeed.Channel.Item {

	
		currentNullTime := sql.NullTime{
			Time:  time.Now(),
			Valid: false,
		}

		pubAt, err :=time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			fmt.Print("Failed to parse time %v", err)
			return err
		}

		newPost := database.CreatePostParams {
			ID: uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: currentNullTime,
			Title: item.Title,
			Url: item.Link,
			Description: sql.NullString{
				String: item.Description ,
			},
			PublishedAt: sql.NullTime{
				Time: pubAt,
			},
			FeedID: nextFeed.ID,
		}

	err =	s.db.CreatePost(context.Background(), newPost)
	fmt.Println("Failed to create post")
	}

	err = s.db.MarkFeedFetched(context.Background(), markFeedFetchedParams)
	if err != nil {
		fmt.Println("Failed to mark the feed as fetched")
	}

	fmt.Println(rssFeed)
	
	return nil
}


func main() {
	arguments := os.Args
	s := state{}

	if len(arguments) < 2 {
		fmt.Println("You have not provided arguments")
		os.Exit(1)
	}

	com := command{
		Name:      arguments[1],
		Arguments: arguments[2:],
	}

	c, err := config.Read()
	if err != nil {
		fmt.Println("Failed to read the config file")
		os.Exit(1)
	}

	s.conf = &c

	coms := commands{
		options: make(map[string]func(*state, command) error),
	}
	registerHandlers(coms)

	db, err := sql.Open("postgres", s.conf.DbUrl)

	s.db = database.New(db)

	err = coms.run(&s, com)
	if err != nil {
		fmt.Println("Error occured")
		os.Exit(1)
	}

	c, err = config.Read()
	if err != nil {
		fmt.Println("Failed to update the config")
		os.Exit(1)
	}
}
