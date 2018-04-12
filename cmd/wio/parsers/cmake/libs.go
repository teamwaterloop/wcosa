// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/libs package, which contains parser to parse libraries and their dependencies
// This file parses dependencies and creates CMakeLists.txt file for each of them
package cmake

import (
    "wio/cmd/wio/utils/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
    "path/filepath"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
    . "wio/cmd/wio/parsers"
)

var configFileName = "wio.yml"
var localLibFolder = "local"
var remoteLibFolder = "remote"
var lockFileName = "libs.lock"
var executionWayFile = "executionWay.txt"

// Parses Library based on the path provides and its dependencies and it modifies structures to store the data
func parseLibrary(libPath string, depTree *DependencyTree, depConfig types.LibrariesLockConfig, hash string) (error) {
    configFile := libPath + io.Sep + configFileName
    libraryLibPath := libPath + io.Sep + "lib"

    var config types.LibConfig
    var libName string

    depLock := types.LibraryLockTag{}

    // check if wio.yml file exists
    if utils.PathExists(configFile) {
        if err := io.NormalIO.ParseYml(configFile, &config); err != nil {
            return err
        }

        libName = config.MainTag.Name
        depLock.Compile_flags = config.MainTag.Compile_flags
    } else {
        libName = filepath.Base(libPath)
    }

    if len(hash) <= 0 {
        hash += libName
    } else {
        hash += "__" + libName
    }

    depLock.Name = libName
    depLock.Hash = hash
    depLock.Path = libPath

    depConfig.Libraries[hash] = &depLock
    depTree.Config = depLock
    depTree.Child = make([]*DependencyTree, 0)

    if !utils.PathExists(libraryLibPath) {
        return nil
    }

    if dirs, err := ioutil.ReadDir(libraryLibPath); err != nil {
        return err
    } else if len(dirs) > 0 {
        for dir := range dirs {
            if dirs[dir].Name()[0] == '.' {
                continue
            }

            currTree := DependencyTree{}

            if err := parseLibrary(libraryLibPath+io.Sep+dirs[dir].Name(), &currTree, depConfig, hash);
                err != nil {
                return err
            }
            depTree.Child = append(depTree.Child, &currTree)
        }
    }

    return nil
}

// Parses all the libraries recursively with the help of ParseLibrary function
func parseLibs(librariesPath string, depTrees []*DependencyTree, depLock types.LibrariesLockConfig,
    libTag types.LibrariesTag, hash string) ([]*DependencyTree, error) {
    if !utils.PathExists(librariesPath) {
        return depTrees, nil
    }

    if dirs, err := ioutil.ReadDir(librariesPath); err != nil {
        return depTrees, err
    } else if len(dirs) > 0 {
        for dir := range dirs {
            currTree := DependencyTree{}

            // ignore hidden files
            if dirs[dir].Name()[0] == '.' {
                continue
            }

            if err := parseLibrary(librariesPath+io.Sep+dirs[dir].Name(), &currTree, depLock, hash);
                err != nil {
                return depTrees, err
            }

            if val, ok := libTag[currTree.Config.Name]; ok {
                currTree.Config.Compile_flags = val.Compile_flags
                depLock.Libraries[currTree.Config.Name].Compile_flags = val.Compile_flags
            }

            depTrees = append(depTrees, &currTree)
        }
    }

    return depTrees, nil
}

// Parses all the libraries and their dependencies and creates libs.lock file
func createLibsLockFile(projectPath string, libTag types.LibrariesTag) ([]*DependencyTree, error) {
    librariesLocalPath := projectPath + io.Sep + "lib" + io.Sep + localLibFolder
    librariesRemotePath := projectPath + io.Sep + "lib" + io.Sep + remoteLibFolder
    wioPath := projectPath + io.Sep + ".wio"

    dependencyTrees := make([]*DependencyTree, 0)
    dependencyLock := types.LibrariesLockConfig{}
    dependencyLock.Libraries = types.LibrariesLockTag{}

    // parse all the libraries and their dependencies in lib/local folder
    dependencyTrees, err := parseLibs(librariesLocalPath, dependencyTrees, dependencyLock, libTag, "")
    if err != nil {
        return nil, err
    }
    // parse all the libraries and their dependencies in lib/remote folder
    dependencyTrees, err = parseLibs(librariesRemotePath, dependencyTrees, dependencyLock, libTag, "")
    if err != nil {
        return nil, err
    }

    // write lock file
    if err = io.NormalIO.WriteYml(wioPath+io.Sep+lockFileName, &dependencyLock); err != nil {
        return nil, err
    }

    // return Dependency Tree
    return dependencyTrees, nil
}

// Chooses include path based on the the template we support
func chooseIncludePath(libraryPath string) (string) {
    return libraryPath + "/include"
}

// Create CMakeLists file for each target
func createCMakes(exePath string, projectPath string, board string, target string, depTree *DependencyTree) (error) {
    wioLibPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + "libraries"

    if err := os.MkdirAll(wioLibPath+io.Sep+depTree.Config.Hash, os.ModePerm); err != nil {
        return err
    }

    librariesTemplate := `# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "{{WIO_PATH}}/toolchain/cmake/CosaToolchain.cmake")
cmake_minimum_required(VERSION 3.0.0)
project({{library-name}} C CXX ASM)

file(GLOB_RECURSE {{glob-name}} {{glob-string}})
generate_arduino_library({{library-name}}
	SRCS ${{{glob-name}}}
	BOARD uno)
target_compile_definitions({{library-name}} PRIVATE __AVR_Cosa__ {{compile-flags}})
include_directories({{lib-path}})
`

    for i := 0; i < len(depTree.Child); i++ {
        err := createCMakes(exePath, projectPath, board, target, depTree.Child[i])
        if err != nil {
            return nil
        }
    }

    var wioPathRe *regexp.Regexp
    var libraryNameRe *regexp.Regexp
    var globNameRe *regexp.Regexp
    var globStringRe *regexp.Regexp
    var compileFLagsRe *regexp.Regexp
    var libPathRe *regexp.Regexp
    var err error

    if wioPathRe, err = regexp.Compile(`{{WIO_PATH}}`); err != nil {
        return err
    }
    if libraryNameRe, err = regexp.Compile(`{{library-name}}`); err != nil {
        return err
    }
    if globNameRe, err = regexp.Compile(`{{glob-name}}`); err != nil {
        return err
    }
    if globStringRe, err = regexp.Compile(`{{glob-string}}`); err != nil {
        return err
    }
    if compileFLagsRe, err = regexp.Compile(`{{compile-flags}}`); err != nil {
        return err
    }
    if libPathRe, err = regexp.Compile(`{{lib-path}}`); err != nil {
        return err
    }

    librariesTemplate = wioPathRe.ReplaceAllString(librariesTemplate, exePath)
    librariesTemplate = libraryNameRe.ReplaceAllString(librariesTemplate, depTree.Config.Hash)
    librariesTemplate = globNameRe.ReplaceAllString(librariesTemplate,
        strings.ToUpper(depTree.Config.Hash+"_SRC_FILES"))

    srcPath := depTree.Config.Path + io.Sep + "src"
    if utils.PathExists(srcPath) {
        librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, srcPath+"/*.cpp "+srcPath+"/*.cc "+srcPath+"/*.c")
    } else {
        librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, depTree.Config.Path+"/*.cpp "+depTree.Config.Path+"/*.cc "+depTree.Config.Path+"/*.c")
    }

    librariesTemplate = compileFLagsRe.ReplaceAllString(librariesTemplate, strings.Join(depTree.Config.Compile_flags, " "))
    librariesTemplate = libPathRe.ReplaceAllString(librariesTemplate, chooseIncludePath(depTree.Config.Path))
    librariesTemplate += "\n"

    for dep := range depTree.Child {
        librariesTemplate += "include_directories(" + chooseIncludePath(depTree.Child[dep].Config.Path) + ")\n"
        librariesTemplate += "target_link_libraries(" + depTree.Config.Hash + " ${CMAKE_SOURCE_DIR}/../" +
            depTree.Child[dep].Config.Hash + io.Sep + "lib" + depTree.Child[dep].Config.Hash + ".a)\n\n"
    }

    io.NormalIO.WriteFile(wioLibPath+io.Sep+depTree.Config.Hash+io.Sep+"CMakeLists.txt", []byte(librariesTemplate))

    executionWayFileName := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + executionWayFile

    f, err := os.OpenFile(executionWayFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil
    }

    LibPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + depTree.Config.Hash

    if _, err = f.WriteString(LibPath + "\n"); err != nil {
        return err
    }

    return nil
}

// Parsers each library and its dependency and then populate libraries CMakes for each target
func PopulateCMakeFilesForLibs(projectPath string, board string, target string, libTag types.LibrariesTag) (error) {
    dependencyTree, err := createLibsLockFile(projectPath, libTag)

    exePath, err := io.NormalIO.GetRoot()

    if err != nil {
        return err
    }
    if err != nil {
        return err
    }

    executionWayFilePath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + executionWayFile

    // remove execution way text file
    if utils.PathExists(executionWayFilePath) {
        err = os.Remove(executionWayFilePath)
    }

    if err != nil {
        return err
    }

    for tree := range dependencyTree {
        if err = createCMakes(exePath, projectPath, board, target, dependencyTree[tree]); err != nil {
            return err
        }
    }

    return nil
}
