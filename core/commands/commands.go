package commands

import (
	"fmt"
	"github.com/ihatiko/olymp/components/observability/tech"
	"github.com/ihatiko/olymp/core/iface"
	_ "github.com/ihatiko/olymp/core/store"
	"github.com/ihatiko/olymp/core/utils"
	tC "github.com/ihatiko/olymp/infrastucture/components/utils/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go/ast"
	"go/parser"
	"go/token"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime/debug"
)

func WithDeployment[Deployment iface.IDeployment]() func() (*cobra.Command, error) {
	return func() (*cobra.Command, error) {
		return &cobra.Command{
			Use: utils.ParseTypeName[Deployment](),
			Run: func(cmd *cobra.Command, args []string) {
				d := new(Deployment)

				defer func() {
					if r := recover(); r != nil {
						stack := string(debug.Stack())
						name := reflect.TypeOf(*d).String()
						slog.Error(fmt.Sprintf("Recovered in core (Run) [%s] \n error: %s", name, stack))
					}
				}()
				err := tC.ToConfig(d)
				if err != nil {
					slog.Error("Error in config", slog.Any("error", err))
					os.Exit(1)
				}
				commandName := utils.ParseTypeName[Deployment]()
				err = os.Setenv("TECH_SERVICE_COMMAND", commandName)
				if err != nil {
					slog.Error("Error in setting environment variable TECH_SERVICE_COMMAND", slog.Any("error", err))
					os.Exit(1)
				}
				app := (*d).Dep()
				rApp := reflect.ValueOf(app)

				p, err := os.Getwd()
				if err != nil {
					slog.Error("Error in getting current working directory", slog.Any("error", err))
					os.Exit(1)
				}

				var collectErrors []struct {
					Type reflect.Type
					Name string
				}
				for i := 0; i < rApp.NumField(); i++ {
					if rApp.Field(i).IsZero() {
						collectErrors = append(collectErrors, struct {
							Type reflect.Type
							Name string
						}{Type: rApp.Type().Field(i).Type, Name: rApp.Type().Field(i).Name})
					}
				}
				if len(collectErrors) != 0 {
					rAppType := reflect.TypeOf(app)
					baseDir := filepath.Dir(p)
					fPath := path.Join(baseDir, rAppType.PkgPath())
					convertedPath := filepath.ToSlash(fPath)
					fSet := token.NewFileSet()
					nodes, err := parser.ParseDir(fSet, convertedPath, nil, parser.ParseComments)
					if err != nil {
						slog.Error("Error in parsing dir", slog.Any("error", err))
						os.Exit(1)
					}
					deploymentPosition := getPosition(nodes, rAppType, fSet)
					name := reflect.TypeOf(*d).String()
					fmt.Println(fmt.Sprintf("⛔️ Error construct deployment [%s] %s", name, deploymentPosition))
					for _, errTypes := range collectErrors {
						position := getFieldPosition(nodes, rAppType, fSet, errTypes.Name, errTypes.Type.String())
						fmt.Println(fmt.Sprintf("⛔️ [%s %s] %s", errTypes.Name, errTypes.Type, position))
					}
					os.Exit(1)
				}
				app.Run()
			},
		}, nil
	}
}
func getPosition(nodes map[string]*ast.Package, rAppType reflect.Type, fSet *token.FileSet) token.Position {
	var filePosition token.Position
	for _, v := range nodes {
		for _, f := range v.Files {
			for _, decl := range f.Decls {
				if fDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range fDecl.Specs {
						if tSpec, ok := spec.(*ast.TypeSpec); ok {
							if tSpec.Name.String() == rAppType.Name() {
								return fSet.Position(fDecl.Pos())
							}
						}
					}
				}
			}
		}
	}
	return filePosition
}
func getFieldPosition(nodes map[string]*ast.Package, rAppType reflect.Type, fSet *token.FileSet, innerField string, innerFieldType string) token.Position {
	var filePosition token.Position
	for _, v := range nodes {
		for _, f := range v.Files {
			for _, decl := range f.Decls {
				if fDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range fDecl.Specs {
						if tSpec, ok := spec.(*ast.TypeSpec); ok {
							if tSpec.Name.String() == rAppType.Name() {
								if structSpec, ok := tSpec.Type.(*ast.StructType); ok {
									for _, field := range structSpec.Fields.List {
										if len(field.Names) > 0 {
											firstElem := field.Names[0]
											if astField, ok := firstElem.Obj.Decl.(*ast.Field); ok {
												if selectorField, ok := astField.Type.(*ast.SelectorExpr); ok {
													if sIdent, ok := selectorField.X.(*ast.Ident); ok {
														t := fmt.Sprintf("%s.%s", sIdent.Name, selectorField.Sel.Name)
														if t == innerFieldType && innerField == firstElem.Name {
															return fSet.Position(field.Pos())
														}
													}
												}
											}
										}

									}
								}
							}
						}
					}
				}
			}
		}
	}
	return filePosition
}
func WithApp(operators ...func() (*cobra.Command, error)) {
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
		slog.Error("error compile command", zap.Any("err", err))
		fmt.Println(err)
		os.Exit(1)
	}
	if len(os.Args) > 1 {
		arg := os.Args[1]
		if os.Args[1] == "-test.v" {
			arg = os.Getenv("TEST_COMMAND")
		}
		rootCommand.SetArgs([]string{arg})
		err := tech.Use(arg)
		if err != nil {
			os.Exit(1)
		}
	}

	err = rootCommand.Execute()
	if err != nil {
		slog.Error("error execute command", slog.Any("error", err))
		os.Exit(1)
	}
}
