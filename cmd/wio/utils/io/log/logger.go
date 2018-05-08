// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package io contains helper functions related to io
// This file contains an interface to print output to io in various colors and modes
package log

import (
    "fmt"
    "wio/cmd/wio/utils/io"
)

var Norm = writer{verbose:true, status:true}
var Verb = writer{verbose:false, status:true}


// user should not touch this
type writer struct {
    verbose bool
    status bool
}

// Turns verbose mode on. This is the mode when Verbose functions work
func SetVerbose() {
    Verb.verbose = true
}

// This is used to turn normal mode print on and off. This way a silent mode can be implemented
func SetStatus(status bool) {
    Norm.status = status
}

func chooseSpecial(special1 string, special2 string) (string) {
    operatingSystem = io.WINDOWS
    if operatingSystem == io.WINDOWS {
        return special2
    } else  {
        return special1
    }
}

func write(color string, newLine bool, text string, a... interface{}) {
    //fmt.Print(color)
    fmt.Print(text)
    //fmt.Print(noColor)

    if newLine {
        fmt.Println()
    }
}


var operatingSystem = io.GetOS()
var cyan = chooseSpecial(`\u001b[36m`,`\x1b[36m`)
var green = chooseSpecial(`\u001b[32m`, `\x1b[32m`)
var magenta = chooseSpecial(`\u001b[35m`, `\x1b[35m`)
var red = chooseSpecial(`\u001b[31m`, `\x1b[31m`)
var yellow = chooseSpecial(`\u001b[33m`, `\x1b[33m`)
var blue = chooseSpecial(`\u001b[34m`, `\x1b[34m`)
var white = chooseSpecial(`\u001b[37m`, `\x1b[37m`)
var noColor = chooseSpecial(`\u001b[39m`, `\x1b[39m`)


// Red is a convenient helper function to print with red foreground.
func (writer writer) Red(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(red, newLine, text, a)
}

// Green is a convenient helper function to print with green foreground.
func (writer writer) Green(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(green, newLine, text, a)
}

// Yellow is a convenient helper function to print with yellow foreground.
func (writer writer) Yellow(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(yellow, newLine, text, a)
}

// Blue is a convenient helper function to print with blue foreground.
func (writer writer) Blue(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(blue, newLine, text, a)
}

// Magenta is a convenient helper function to print with magenta foreground.
func (writer writer) Magenta(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(magenta, newLine, text, a)
}

// Cyan is a convenient helper function to print with cyan foreground.
func (writer writer) Cyan(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(cyan, newLine, text, a)
}

// White is a convenient helper function to print with white foreground.
func (writer writer) White(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(white, newLine, text, a)
}

// Normal is a convenient helper function to print with default/normal foreground.
func (writer writer) Write(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    write(noColor, newLine, text, a)
}

// Special function to be used when using Verbose mode.
// In this mode, color can be set and other verbose default things can be defined
func (writer writer) Verbose(newLine bool, text string, a ...interface{}) {
    if !writer.status || !writer.verbose { return }

    writer.Write(newLine, text, a)
}
