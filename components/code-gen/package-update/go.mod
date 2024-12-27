module github.com/ihatiko/chef/code-gen/package-update

go 1.23.4

require (
	github.com/ihatiko/chef/code-gen/file-builder v0.0.0-00010101000000-000000000000
	golang.org/x/mod v0.22.0
)

replace github.com/ihatiko/chef/code-gen/file-builder => ../file-builder
