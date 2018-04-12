package parsers

import "wio/cmd/wio/utils/types"

// Type for Dependency tree used for parsing libraries
type DependencyTree struct {
    Config types.LibraryLockTag
    Child    []*DependencyTree
}
