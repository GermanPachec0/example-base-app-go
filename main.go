package main

import (
	"context"
	"os"

	"github.com/GermanPachec0/app-go/cmd"
)

func main() {
	ctx := context.Background()
	ret := cmd.Execute(ctx)
	os.Exit(ret)
}
