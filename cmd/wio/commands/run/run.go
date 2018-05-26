// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Builds, Tests, and Uploads the project to a device
package run

import (
    "github.com/urfave/cli"
    "wio/cmd/wio/commands/build"
    "path/filepath"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/utils/io/log"
    "os/exec"
    "bytes"
    "os"
    "strings"
    "wio/cmd/wio/utils/io"
)

type Run struct {
    Context *cli.Context
    error
}

// get context for the command
func (run Run) GetContext() (*cli.Context) {
    return run.Context
}

// Runs the build, upload command (acts as one in all command)
func (run Run) Execute() {
    directory, err := filepath.Abs(run.Context.String("dir"))
    commands.RecordError(err, "")

    port := run.Context.String("port")
    var targetName string

    if run.Context.IsSet("port") {
        log.Norm.Cyan(true, "Including port: \"" + port + "\" in the build process\n")

        // run the build so that port can be included in the build process
        targetName = build.RunBuild(directory, run.Context.String("target"), false, port)
    } else {
        // run the build normally
        targetName = build.RunBuild(directory, run.Context.String("target"), false, "")
    }

    // execute the upload command if port is valid
    if run.Context.IsSet("port") {
        targetsDirectory := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets"
        targetPath := targetsDirectory + io.Sep + targetName

        // execute make command
        makeCommand := exec.Command("make", "upload")
        makeCommand.Dir = targetPath

        // Stderr buffer
        cmdErrOutput := &bytes.Buffer{}
        makeCommand.Stderr = cmdErrOutput

        if log.Verb.IsVerbose() {
            makeCommand.Stdout = os.Stdout
        }

        log.Norm.Cyan(false, "\nuploading the target ... ")
        log.Verb.Write(true, "")

        err = makeCommand.Run()
        if err != nil {
            commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
        } else {
            log.Norm.Green(true, "success")
        }

        // print the ending message
        log.Norm.Yellow(true, "Upload successful for target: \"" + targetName + "\"")
    }
}
