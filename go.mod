module github.com/ihatiko/chef

go 1.23.4

replace (
	github.com/ihatiko/chef/code-gen/file-builder => ./components/code-gen/file-builder
	github.com/ihatiko/chef/code-gen/package-update => ./components/code-gen/package-update
)
