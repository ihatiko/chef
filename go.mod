module github.com/ihatiko/chef/cli/go-chef

go 1.23.4

replace (
	github.com/ihatiko/chef/cli/go-chef => ./cli/go-chef
	github.com/ihatiko/chef/cli/go-chef-core => ./cli/go-chef-core
	github.com/ihatiko/components/code-gen/file-builder => ./components/file-builder
	github.com/ihatiko/components/code-gen/package-updater => ./components/package-updater
)
