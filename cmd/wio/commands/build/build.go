// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands package, which contains all the commands provided by the tool.
// Builds the project
package build

import (
    "wio/cmd/wio/commands"
    "os"
    "github.com/urfave/cli"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io/log"
    "path/filepath"
    "os/exec"
    "wio/cmd/wio/parsers/cmake"
)

type Build struct {
    Context *cli.Context
    error
}

// Compiles the project
func (build Build) Execute() {
    directory, err := filepath.Abs(build.Context.String("directory"))
    commands.RecordError(err, "")

    clean := build.Context.Bool("clean")

    targetToBuildName := build.Context.String("target")
    var name string
    var currTargets types.TargetsTag
    var targetToBuild *types.TargetTag
    var pkgCompileFlags []string
    var framework string
    var dependencies types.DependenciesTag

    configPath := directory + io.Sep + "wio.yml"

    // try parsing app type and if it comes empty this means our project is pkg type
    appConfig := types.AppConfig{}
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

    log.Norm.Red(true, name)

    // clean the build files if clean flag is true
    if clean {
        cleanBuildFiles(directory)
    }

    // create the target
    createTarget(name, directory, targetToBuild.Board, framework, targetToBuildName, targetToBuild.Compile_flags,
        pkgCompileFlags, dependencies)

    // build the target
    buildTarget(directory, targetToBuildName)

    // remove the target
    removeTarget(directory)
}

func removeTarget(directory string) {
    cmakePath := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "CMakeLists.txt"
    depPath := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "dependencies.cmake"

    os.Remove(cmakePath)
    os.Remove(depPath)
}

func createTarget(name string, directory string, board string, framework string, target string, targetFlags []string,
    pkgFlags []string, dependencies types.DependenciesTag) (error) {
    // parse dependencies and create a dependencies.cmake file
    dependencyTree, err := cmake.ParseDepsAndCreateCMake(directory, dependencies)
    if err != nil {
        return err
    }

    // create the main CMakeLists.txt file
    if pkgFlags != nil {
        // create cmake for package type
        if err := cmake.CreatePkgMainCMakeLists(name, directory, board, framework, target, targetFlags, pkgFlags,
            dependencyTree); err != nil {
            return err
        }
    } else {
        // create cmake for app type
        if err := cmake.CreateAppMainCMakeLists(name, directory, board, framework, target, targetFlags,
            dependencyTree); err != nil {
            return err
        }
    }

    return nil
}

func buildTarget(directory string, targetName string) (error) {
    targetsDirectory := directory + io.Sep + ".wio" + io.Sep + "build" + io.Sep + "targets"
    targetPath := targetsDirectory + io.Sep + targetName

    // create a folder for the target
    if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
        return err
    }

    cmakeCommand := exec.Command("cmake", "../../")
    makeCommand := exec.Command("make")
    cmakeCommand.Dir = targetPath
    makeCommand.Dir = targetPath
    cmakeCommand.Stdout = os.Stdout
    makeCommand.Stdout = os.Stdout
    cmakeCommand.Stderr = os.Stderr
    makeCommand.Stderr = os.Stderr

    if err := cmakeCommand.Run(); err != nil {
        return err
    }

    log.Norm.Write(true, "\n\n")

    if err := makeCommand.Run(); err != nil {
        return err
    }

    log.Norm.Write(true, "\n\n")

    return nil
}

func cleanBuildFiles(directory string) {
    // clean.Execute()
}
