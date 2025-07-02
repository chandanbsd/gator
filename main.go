package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/chandanbsd/gator/internal/config"
	"github.com/chandanbsd/gator/internal/database"
	"github.com/chandanbsd/gator/internal/feed"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"os"
	"time"
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
	coms.register("addFeed", addFeedHandler)
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

	var feedURL string = "https://www.wagslane.dev/index.xml"
	rssFeed, err := feed.FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}

	fmt.Println(rssFeed)
	return nil
}

func addFeedHandler(s *state, cmd command) error {
	if len(cmd.Arguments) < 2 {
		fmt.Println("Missing name argument")
		os.Exit(1)
	}

	user, err := s.db.GetUser(context.Background(), s.conf.CurrentUserName)
	if err != nil {
		fmt.Println("Failed to get the current user")
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

	_, err = s.db.CreateFeed(context.Background(), newFeed)
	if err != nil {
		fmt.Println("Failed to create the feed")
		os.Exit(1)
	} else {
		fmt.Println("User created successfully")
	}

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

	fmt.Println(c)
}
