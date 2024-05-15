go install go.uber.org/mock/mockgen@latest
mockgen -version
export PATH=$PATH:$(go env GOPATH)/bin

mockgen -help
# mockgen has two modes of operation: source and reflect.

# Source mode generates mock interfaces from a source file.
# It is enabled by using the -source flag. Other flags that
# may be useful in this mode are -imports and -aux_files.
# Example:
# 	mockgen -source=foo.go [other options]

# Reflect mode generates mock interfaces by building a program
# that uses reflection to understand interfaces. It is enabled
# by passing two non-flag arguments: an import path, and a
# comma-separated list of symbols.
# Example:
# 	mockgen database/sql/driver Conn,Driver


mockgen -destination db/mock/store.go -package mock_db github.com/saalikmubeen/backend-masterclass-go/db/sqlc Store

# mockgen -source ./db/sqlc/store.go  -destination db/mock/store.go -package mock_db Store