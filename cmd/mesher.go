package main

import (
	"fmt"
	"github.com/spf13/cobra"
	f3ds "github.com/vaguilera/mesher/pkg/3ds"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "Mesher",
	Short: "Mesher allows you to convert some 3d format to web-usable ones",
	Long:  "Mesher allows you to convert some 3d format to web-usable ones",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args[0]) < 4 {
			log.Fatalf("File source extension should be .3ds or .obj")
		}

		switch args[0][len(args[0])-3:] {
		case "3ds":
			process3dsFile(args)
		default:
			log.Fatalf("File source extension should be .3ds or .obj")
		}

		fmt.Println("fsdgdfgdf")
	},
}

func main() {

	if err := rootCmd.Execute(); err != nil {

		os.Exit(1)
	}
}

func process3dsFile(args []string) {
	f := f3ds.F3DS{}
	// "/Users/aguilera_sal/Downloads/box.3ds"
	err := f.LoadFile(args[0])
	if err != nil {
		log.Fatalf("%s", err)
	}
}
