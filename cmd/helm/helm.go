package main

import (
	"fmt"
	"os"

	"github.com/StefanNyman/kubectl/lib"
)

func main() {
	ctx, err := lib.NewCtx(lib.Helm)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	err = ctx.Run()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
