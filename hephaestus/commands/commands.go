package commands

import (
	"fmt"
	"os"

	"github.com/ihatiko/olymp/hephaestus/iface"
	_ "github.com/ihatiko/olymp/hephaestus/store"
	"github.com/ihatiko/olymp/hephaestus/utils"
	tC "github.com/ihatiko/tech-config"
	"github.com/spf13/cobra"
)

// deployment описание внутренний функции деплоймента
type deployment = func() (*cobra.Command, error)

// WithDeployment Минимальная единица приложения
// Описание как запускать ваш код в k8s - работает исключительно на строгих типах
// Deployment any -> Структура которая распологается по пути internal/server/deployments
func WithDeployment[Deployment iface.IDeployment]() deployment {
	return func() (*cobra.Command, error) {
		d := new(Deployment)
		err := tC.ToConfig(d)
		if err != nil {
			return nil, err
		}
		return &cobra.Command{
			Use: utils.ParseTypeName[Deployment](),
			Run: func(cmd *cobra.Command, args []string) {
				app := (*d).Dep()
				app.Run()
			},
		}, nil
	}
}

func WithApp(operators ...deployment) {
	cmd := new(cobra.Command)
	var (
		err error
		c   *cobra.Command
	)
	for _, d := range operators {
		c, err = d()
		if err != nil {
			fmt.Printf("error: %v", err)
			continue
		}
		cmd.AddCommand(c)
	}
	Compile(cmd, err)
}

func Compile(rootCommand *cobra.Command, err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		arg := []string{os.Args[1]}
		if os.Args[1] == "-test.v" {
			arg = []string{os.Getenv("TEST_COMMAND")}
		}
		rootCommand.SetArgs(arg)
		//err := tech_components.Configure(os.Args[1])
		//if err != nil {
		//	fmt.Println(err)
		//}
	}

	err = rootCommand.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
