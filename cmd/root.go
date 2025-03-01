package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var root = &cobra.Command{
	Use:   "automsg",
	Short: "automatic message sending system",
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
