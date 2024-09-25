package main

import (
	"context"
	"fmt"
	"os"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %s", cmd.Name)
	}

	err := s.db.ResetUsers(context.Background())
	if err != nil {
		fmt.Println("error while resetting users")
		os.Exit(1)
		return err
	}
	return nil
}
