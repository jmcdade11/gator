package main

import (
	"context"
	"fmt"
)

func handlerGetUsers(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}
	currUser := s.cfg.CurrentUserName
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}

	for _, user := range users {
		userName := user.Name
		if userName == currUser {
			fmt.Printf("* %s (current)\n", userName)
		} else {
			fmt.Printf("* %s\n", userName)
		}
	}
	return nil
}
