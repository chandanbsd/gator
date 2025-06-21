package main

import (
	"fmt"
	"errors"
	"github.com/chandanbsd/gator/internal/config"
	"os"
)


type command struct {
	Name string
	Arguments []string
}

type commands struct {
	options map[string]func(*state, command)error
}

type state struct {
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
	coms.register("login", loginHandler)
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

	fmt.Print(c)
}
