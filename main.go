package main

import (
	"fmt"
	"errors"
	"github.com/chandanbsd/gator/internal/config"
	"os"
	"github.com/chandanbsd/gator/internal/database"
	"time"
	"database/sql"
	"context"
	"github.com/google/uuid"
)

import _ "github.com/lib/pq"


type command struct {
	Name string
	Arguments []string
}

type commands struct {
	options map[string]func(*state, command)error
}

type state struct {
	db *database.Queries
	conf *config.Config
}

func (c *commands) run(s *state, cmd command) error {
	for name, handler  := range c.options {
		if name  == cmd.Name {
			err := handler(s, cmd )
			return err
		}
	}
	return errors.New("The command you have entered is not valid")
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.options[name] = f
}

func loginHandler(s *state, cmd command) error{
	if s == nil {
		fmt.Println("Please enter a command")
		return errors.New("command is empty")
	}

	if cmd.Arguments == nil || len(cmd.Arguments) == 0 {
		fmt.Println("Please provide a username")
		return errors.New("Failed to provide a username")
	}

	s.conf.SetUser(cmd.Arguments[0])

	fmt.Printf("Your username has been set to %s\n", s.conf.CurrentUserName)

	return nil
}

func registerHandlers(coms commands) {
	coms.register("login", loginHandler)
	coms.register("register", registerHandler)
}

func registerHandler(s *state, cmd command) error {
	if len(cmd.Arguments) == 0 {
		fmt.Println("Missing name argument")
		return errors.New("You missed name argument")
	}

	user, err := state.db.db.GetUser(context.Background(), cmd.Arguments[0])
	if err != nil || user != nil{
		fmt.Prinln("Cannot verify if the user exists, please try again later")
		return err
	}
	newUserParams := state.db.db.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: cmd.Arguments[0],
	}
		user, err = state.db.db.CreateUser(context.Background(), newUserParams)
	if err == nil {
		fmt.Println("Failed to create the user.")
		return err
	}
	return nil
}

func main() {
	arguments := os.Args
	s := state {}



	if len(arguments) < 2 {
		fmt.Println("You have not provided arguments")
		os.Exit(1)
	}

	com := command{
		Name: arguments[1],
		Arguments: arguments[2:],
	}

	c, err := config.Read()
	if err != nil {
		fmt.Println("Failed to read the config file")
		os.Exit(1)
	}

	s.conf = &c

	coms := commands{
		options: make(map[string]func(*state, command)error),
	}
	registerHandlers(coms)

	err = coms.run(&s, com)
	if err != nil {
		fmt.Println("Error occured")
		os.Exit(1)
	}

	db, err := sql.Open("postgres", s.conf.DbURL)

	state.db = database.New(db)

	c, err = config.Read()
	if err != nil {
		fmt.Println("Failed to update the config")
		os.Exit(1)
	}

	fmt.Print(c)
}
