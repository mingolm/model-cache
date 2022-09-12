package main

import (
	"context"
	"fmt"
)

func main() {
	userRepo := NewRepo()
	row, err := userRepo.Get(context.Background(), 1001)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("user: %+v\n", row)
}
