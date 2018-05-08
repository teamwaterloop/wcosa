// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Package manager for wio
package libraries

import (
    "github.com/urfave/cli"
    "github.com/go-errors/errors"
    "fmt"
)

const (
    GET = "get"
    UPDATE = "update"
    COLLECT = "collect"
)

type Libraries struct {
    Context *cli.Context
    Type string
    error
}

// Executes the libraries command
func (libraries Libraries) Execute() {

    switch libraries.Type {
    case GET:
        libraries.error = errors.New("GG")
        libraries.handleGet(libraries.Context)
        break
    case UPDATE:
        break
    case COLLECT:
        break
    }
}
func (libraries Libraries) handleGet(context *cli.Context) {
    if libraries.error != nil {
        return
    }

    fmt.Println("GGGGG")
}
