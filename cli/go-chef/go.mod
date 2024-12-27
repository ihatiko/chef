module github.com/ihatiko/chef/cli/go-chef

go 1.23.4

replace (
	github.com/ihatiko/chef/code-gen/file-builder => ../../components/code-gen/file-builder
	github.com/ihatiko/chef/code-gen/package-update => ../../components/code-gen/package-update
)

require (
	github.com/ihatiko/chef/code-gen/file-builder v0.0.18
	github.com/ihatiko/chef/code-gen/package-update v0.0.18
)

require golang.org/x/mod v0.22.0 // indirect
