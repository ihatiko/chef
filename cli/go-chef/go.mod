module github.com/ihatiko/chef/cli/go-chef

go 1.23.4

replace (
	github.com/ihatiko/chef/components/code-gen-utils/command-executor v0.0.0-00010101000000-000000000000 => ./../../components/code-gen-utils/command-executor
	github.com/ihatiko/chef/components/code-gen-utils/package-updater v0.0.0-00010101000000-000000000000 => ./../../components/code-gen-utils/package-updater
)


require golang.org/x/mod v0.22.0 // indirect
