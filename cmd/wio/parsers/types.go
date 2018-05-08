package parsers

import "wio/cmd/wio/types"

// Type for Dependency tree used for parsing libraries
type DependencyTree struct {
    Config types.PackageLockTag
    Child    []*DependencyTree
}
