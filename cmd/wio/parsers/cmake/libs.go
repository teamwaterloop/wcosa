// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/cmake package, which contains parser to create cmake files
// This file parses dependencies and creates CMakeLists.txt file for each of them
package cmake

import (
    "wio/cmd/wio/utils/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
    "path/filepath"
    "io/ioutil"
    "os"
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
func createCMakes(exePath string, projectPath string, board string, framework string, target string, depTree *DependencyTree) (error) {
    wioLibPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + "libraries"

    if err := os.MkdirAll(wioLibPath+io.Sep+depTree.Config.Hash, os.ModePerm); err != nil {
        return err
    }

    librariesTemplate, err := io.AssetIO.ReadFile("templates" + io.Sep + "cmake" + io.Sep + "CMakeListsTarget.txt")
    if err != nil {
        return err
    }
    librariesTemplateStr := string(librariesTemplate)

    for i := 0; i < len(depTree.Child); i++ {
        err := createCMakes(exePath, projectPath, board, framework, target, depTree.Child[i])
        if err != nil {
            return nil
        }
    }

    librariesTemplateStr = strings.Replace(librariesTemplateStr, "{{WIO_PATH}}", exePath, -1)
    librariesTemplateStr = strings.Replace(librariesTemplateStr, "{{library-name}}", depTree.Config.Hash, -1)
    librariesTemplateStr = strings.Replace(librariesTemplateStr, "{{lib-path}}", depTree.Config.Path, -1)
    librariesTemplateStr = strings.Replace(librariesTemplateStr, "{{board}}", board, -1)
    librariesTemplateStr = strings.Replace(librariesTemplateStr, "{{framework}}", strings.ToUpper(framework), -1)
    librariesTemplateStr = strings.Replace(librariesTemplateStr, "{{compile-flags}}",
        strings.Join(depTree.Config.Compile_flags, " "), -1)
    librariesTemplateStr += "\n"

    for dep := range depTree.Child {
        librariesTemplateStr += "include_directories(\"" + chooseIncludePath(depTree.Child[dep].Config.Path) + "\")\n"
        librariesTemplateStr += "target_link_libraries(" + depTree.Config.Hash + " \"${CMAKE_SOURCE_DIR}/../" +
            depTree.Child[dep].Config.Hash + io.Sep + "lib" + depTree.Child[dep].Config.Hash + ".a\")\n\n"
    }

    io.NormalIO.WriteFile(wioLibPath+io.Sep+depTree.Config.Hash+io.Sep+"CMakeLists.txt", []byte(librariesTemplateStr))

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
func PopulateCMakeFilesForLibs(projectPath string, board string, framework string, target string, libTag types.LibrariesTag) (error) {
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
        if err = createCMakes(exePath, projectPath, board, framework, target, dependencyTree[tree]); err != nil {
            return err
        }
    }

    return nil
}
