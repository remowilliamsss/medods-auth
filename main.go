package main

import (
	"context"
	"medods-auth/app"
)

func main() {
	if err := app.Run(context.Background()); err != nil {
		panic(err)
	}
}
