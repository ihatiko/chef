module github.com/ihatiko/chef/cli/go-chef

go 1.23.4

replace (
	github.com/ihatiko/chef/code-gen/file-builder v0.0.0-00010101000000-000000000000 => ../../components/code-gen/file-builder
	github.com/ihatiko/chef/code-gen/package-update v0.0.0-00010101000000-000000000000 => ../../components/code-gen/package-update
)

require (
	github.com/ihatiko/chef/code-gen/file-builder v0.0.0-00010101000000-000000000000
	github.com/ihatiko/chef/code-gen/package-update v0.0.0-00010101000000-000000000000
)

require golang.org/x/mod v0.22.0 // indirect
