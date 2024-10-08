package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	f3ds "github.com/vaguilera/mesher/pkg/3ds"
	"github.com/vaguilera/mesher/pkg/obj"
	"github.com/vaguilera/mesher/pkg/w3d"
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
		case "obj":
			processObjFile(args)
		default:
			log.Fatalf("File source extension should be .3ds or .obj")
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func process3dsFile(args []string) {
	f := f3ds.F3DS{}

	err := f.LoadFile(args[0])
	if err != nil {
		log.Fatalf("%s", err)
	}

	w3dfile := w3d.New3WDFrom3DS(&f)

	jsonW3d, err := json.Marshal(w3dfile)
	if err != nil {
		fmt.Printf("error marshalling json: %s\n", err.Error())
	}

	outFile, err := os.Create(args[1])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	n, err := outFile.WriteString("export const w3d = '" + string(jsonW3d) + "';")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("%s created. %d bytes\n", args[1], n)
}

func processObjFile(args []string) {
	f := obj.ObjFile{}

	err := f.LoadFile(args[0])
	if err != nil {
		log.Fatalf("%s", err)
	}

	w3dfile := w3d.New3WDFromOBJ(&f)

	jsonW3d, err := json.Marshal(w3dfile)
	if err != nil {
		fmt.Printf("error marshalling json: %s\n", err.Error())
	}

	outFile, err := os.Create(args[1])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	n, err := outFile.WriteString("export const w3d = '" + string(jsonW3d) + "';")
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("%s created. %d bytes\n", args[1], n)
}
