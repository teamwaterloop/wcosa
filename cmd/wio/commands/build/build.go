// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Builds the project
package build

import (
    "os"
    "github.com/urfave/cli"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "wio/cmd/wio/parsers/cmake"
    "os/exec"
    "wio/cmd/wio/commands"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/io/log"
    "path/filepath"
    "wio/cmd/wio/commands/clean"
    "bytes"
    "strings"
)

type Build struct {
    Context *cli.Context
    error
}

// get context for the command
func (build Build) GetContext() (*cli.Context) {
    return build.Context
}

// Runs the build command when cli build option is provided
func (build Build) Execute() {
    RunBuild(build.Context.String("dir"), build.Context.String("target"),
        build.Context.Bool("clean"), "")
}

// This function allows other packages to call build as well. This is also used when cli build is executed
func RunBuild(directoryCli string, targetCli string, cleanCli bool, port string) string {
    directory, err := filepath.Abs(directoryCli)
    commands.RecordError(err, "")

    cleanStr := cleanCli

    targetToBuildName := targetCli
    var name string
    var currTargets types.TargetsTag
    var targetToBuild *types.TargetTag
    var pkgCompileFlags []string
    var framework string
    var dependencies types.DependenciesTag

    configPath := directory + io.Sep + "wio.yml"

    // try parsing app type and if it comes empty this means our project is pkg type
    appConfig := types.AppConfig{}

    // check if wio.yml file exits
    if !utils.PathExists(directory + io.Sep + "wio.yml") {
        log.Error(true, "wio.yml file not found. Not a valid wio project!")
    }

    commands.RecordError(io.NormalIO.ParseYml(configPath, &appConfig), "")

    // this means we need to parse pkg type
    if appConfig.MainTag.Name == "" {
        pkgConfig := types.PkgConfig{}
        commands.RecordError(io.NormalIO.ParseYml(configPath, &pkgConfig), "")

        currTargets = pkgConfig.TargetsTag
        name = pkgConfig.MainTag.Name
        pkgCompileFlags = pkgConfig.MainTag.Compile_flags
        framework = pkgConfig.MainTag.Framework[0]
        dependencies = pkgConfig.DependenciesTag
    } else {
        currTargets = appConfig.TargetsTag
        name = appConfig.MainTag.Name
        framework = appConfig.MainTag.Framework
        pkgCompileFlags = nil
        dependencies = appConfig.DependenciesTag
    }

    // if toBuildTarget is "default" we need to build the default
    if targetToBuildName == "default" {
        targetToBuildName = currTargets.Default_target
        targetToBuild = currTargets.Targets[targetToBuildName]
    } else {
        if target, ok := currTargets.Targets[targetToBuildName]; !ok {
            panic("bad tag")
        } else {
            targetToBuild = target
        }
    }

    // clean the build files if clean flag is true
    if cleanStr {
        cleanBuildFiles(directory, targetToBuildName)
    }

    log.Norm.Yellow(true, "Building the project")

    // create the target
    createTarget(name, directory, targetToBuild.Board, port, framework, targetToBuildName,
        targetToBuild.Compile_flags, pkgCompileFlags, dependencies)

    // build the target
    buildTarget(directory, targetToBuildName)

    // print the ending message
    log.Norm.Yellow(true, "Build successful for \"" + name + "\" project!")

    return targetToBuildName
}

// Scans dependency tree and based on that creates CMake build files
func createTarget(name string, directory string, board string, port string, framework string, target string,
    targetFlags []string, pkgFlags []string, dependencies types.DependenciesTag) {

    log.Norm.Cyan(false, "Scanning dependency tree for changes ... ")

    // parse dependencies and create a dependencies.cmake file
    dependencyTree, err := cmake.ParseDepsAndCreateCMake(directory, dependencies)
    if err != nil {
        log.Norm.Red(true, "failure")
        commands.RecordError(err, "")
    } else {
        log.Norm.Green(true, "success")
    }

    log.Verb.Verbose(false, "Creating Build Tool files ... ")

    // create the main CMakeLists.txt file
    if pkgFlags != nil {
        // create cmake for package type
        if err := cmake.CreatePkgMainCMakeLists(name, directory, board, port, framework, target, targetFlags,
            pkgFlags, dependencyTree); err != nil {
            log.Verb.Verbose(true, "failure")
            commands.RecordError(err, "")
        } else {
            log.Verb.Verbose(true, "success")
        }
    } else {
        // create cmake for app type
        if err := cmake.CreateAppMainCMakeLists(name, directory, board, port, framework, target, targetFlags,
            dependencyTree); err != nil {
            log.Verb.Verbose(true, "failure")
            commands.RecordError(err, "")
        } else {
            log.Verb.Verbose(true, "success")
        }
    }
}

// Execute CMake and Make commands to build the project
func buildTarget(directory string, targetName string) {
    targetsDirectory := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets"
    targetPath := targetsDirectory + io.Sep + targetName

    log.Norm.Cyan(false, "creating build environment for target: \"" + targetName + "\" ... ")

    // create a folder for the target
    if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
        commands.RecordError(err, "failure")
    } else {
        log.Norm.Green(true, "success")
    }

    // execute cmake command
    cmakeCommand := exec.Command("cmake", "../../", "-G", "Unix Makefiles")
    cmakeCommand.Dir = targetPath

    // Stderr buffer
    cmdErrOutput := &bytes.Buffer{}
    cmakeCommand.Stderr = cmdErrOutput

    if log.Verb.IsVerbose() {
        cmakeCommand.Stdout = os.Stdout
    }

    log.Norm.Cyan(false, "generating build files ... ")
    log.Verb.Write(true, "")

    err := cmakeCommand.Run()
    if err != nil {
        commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
    } else {
        log.Norm.Green(true, "success")
    }

    // execute make command
    makeCommand := exec.Command("make")
    makeCommand.Dir = targetPath

    cmdErrOutput = &bytes.Buffer{}
    makeCommand.Stderr = cmdErrOutput

    if log.Verb.IsVerbose() {
        makeCommand.Stdout = os.Stdout
    }

    log.Verb.Write(true, "")
    log.Norm.Cyan(false, "running the build ... ")
    log.Verb.Write(true, "")

    err = makeCommand.Run()
    if err != nil {
        commands.RecordError(err, "failure", strings.Trim(cmdErrOutput.String(), "\n"))
    } else {
        log.Norm.Green(true, "success")
        log.Verb.Write(true, "")
    }
}

// cleans current build files
func cleanBuildFiles(directory string, target string) {
    clean.RunClean(directory, target, true)
    log.Norm.Write(true, "")
}
