module github.com/ihatiko/chef

go 1.23.4

replace (
	github.com/ihatiko/chef/code-gen/file-builder => ./components/code-gen/file-builder
	github.com/ihatiko/chef/code-gen/package-update => ./components/code-gen/package-update
)

require (
	github.com/ihatiko/chef/code-gen/file-builder v0.0.0-00010101000000-000000000000
	github.com/ihatiko/chef/code-gen/package-update v0.0.0-00010101000000-000000000000
)
