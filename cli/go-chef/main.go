package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
	"io"
	"net/http"
	"os"
	"strings"
)

var rootCmd = &cobra.Command{
	Use:   "zero",
	Short: "zero is a cli tool for performing basic mathematical operations",
	Long:  "zero is a cli tool for performing basic mathematical operations - addition, multiplication, division and subtraction.",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Oops. An error while executing Zero '%s'\n", err)
		os.Exit(1)
	}
}

func main() {
	response, err := http.Get("https://proxy.golang.org/github.com/ihatiko/go-chef-sandbox-test/@v/list")
	if err != nil {
		return
	}

	reader := bufio.NewReader(response.Body)

	bytes, err := reader.ReadBytes(0)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return
	}

	versions := strings.Split(string(bytes), "\n")

	semver.Sort(versions)

	lastVersion := versions[len(versions)-1]

	fmt.Println(lastVersion)

	//fmt.Println(rootCmd.Execute())
	//Execute()
}
