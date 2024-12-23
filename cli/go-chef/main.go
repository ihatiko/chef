package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
	"io"
	"net/http"
	"os"
	"os/exec"
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
	packageName := "github.com/ihatiko/go-chef-sandbox-test"
	fullPathName := fmt.Sprintf("https://proxy.golang.org/%s/@v/list", packageName)
	response, err := http.Get(fullPathName)
	if err != nil {
		fmt.Println(err)
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
	installInstruction := fmt.Sprintf("%s@%s", packageName, lastVersion)
	command := fmt.Sprintf("go install %s", installInstruction)
	fmt.Println(command)
	cmd := exec.Command(command)
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
	//fmt.Println(rootCmd.Execute())
	//Execute()
}
