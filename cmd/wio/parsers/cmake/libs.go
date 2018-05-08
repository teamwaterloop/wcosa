// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/cmake package, which contains parser to create cmake files
// This file parses dependencies and creates CMakeLists.txt file for each of them

package cmake

import (
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/parsers"
    "wio/cmd/wio/utils"
    "io/ioutil"
    "path/filepath"
    "os"
    "strings"
)

var configFileName = "wio.yml"
var vendorPkgFolder = "vendor"
var remotePkgFolder = "remote"
var lockFileName = "pkg.lock"
var executionWayFile = "executionWay.txt"

// Parses Package based on the path provides and its dependencies and it modifies dependency tree to store the data
func parsePackage(packagePath string, depTree *parsers.DependencyTree, depConfig types.PackagesLockConfig,
    dependencies types.DependenciesTag, hash string) (error) {

    configFile := packagePath + io.Sep + configFileName
    dependenciesPath := packagePath + io.Sep + "pkg"

    var config types.PkgConfig
    var pkgName string

    depLock := types.PackageLockTag{}

    // check if wio.yml file exists for the dependency and based on that gather the name
    // of the package
    if utils.PathExists(configFile) {
        if err := io.NormalIO.ParseYml(configFile, &config); err != nil {
            return err
        }

        pkgName = config.MainTag.Name
        depLock.Compile_flags = config.MainTag.Compile_flags
    } else {
        pkgName = filepath.Base(packagePath)
    }

    if len(hash) <= 0 {
        hash += pkgName
    } else {
        hash += "__" + pkgName
    }

    depLock.Name = pkgName
    depLock.Hash = hash
    depLock.Path = packagePath

    depConfig.Packages[hash] = &depLock
    depTree.Config = depLock
    depTree.Child = make([]*parsers.DependencyTree, 0)

    // parse vendor dependencies
    vendorPath := dependenciesPath + io.Sep + "vendor"
    if utils.PathExists(vendorPath) {
        if dirs, err := ioutil.ReadDir(vendorPath); err != nil {
            return err
        } else if len(dirs) > 0 {
            for dir := range dirs {
                if dirs[dir].Name()[0] == '.' {
                    continue
                }

                currTree := parsers.DependencyTree{}

                if err := parsePackage(vendorPath+io.Sep+dirs[dir].Name(), &currTree, depConfig,
                    config.DependenciesTag, hash);
                    err != nil {
                    return err
                }
                depTree.Child = append(depTree.Child, &currTree)
            }
        }
    }

    // parse remote dependencies
    remotePath := dependenciesPath + io.Sep + "remote"
    if utils.PathExists(remotePath) {
        if dirs, err := ioutil.ReadDir(remotePath); err != nil {
            return err
        } else if len(dirs) > 0 {
            for dir := range dirs {
                if dirs[dir].Name()[0] == '.' {
                    continue
                }

                currTree := parsers.DependencyTree{}

                if err := parsePackage(remotePath+io.Sep+dirs[dir].Name(), &currTree, depConfig,
                    config.DependenciesTag, hash);
                    err != nil {
                    return err
                }
                depTree.Child = append(depTree.Child, &currTree)
            }
        }
    }


    if val, ok := dependencies[depTree.Config.Name]; ok {
        allFlags := utils.AppendIfMissing(val.Compile_flags, depTree.Config.Compile_flags)

        depTree.Config.Compile_flags = allFlags
        depConfig.Packages[depTree.Config.Hash].Compile_flags = allFlags
    }


    return nil
}

// Parses all the packages recursively with the help of ParsePackage function
func parsePackages(packagesPath string, depTrees []*parsers.DependencyTree, depConfig types.PackagesLockConfig,
    dependencies types.DependenciesTag, hash string) ([]*parsers.DependencyTree, error) {

    if !utils.PathExists(packagesPath) {
        return depTrees, nil
    }

    if dirs, err := ioutil.ReadDir(packagesPath); err != nil {
        return depTrees, err
    } else if len(dirs) > 0 {
        for dir := range dirs {
            currTree := parsers.DependencyTree{}

            // ignore hidden files
            if dirs[dir].Name()[0] == '.' {
                continue
            }

            if err := parsePackage(packagesPath+io.Sep+dirs[dir].Name(), &currTree, depConfig, dependencies, hash);
                err != nil {
                return depTrees, err
            }

            depTrees = append(depTrees, &currTree)
        }
    }

    return depTrees, nil
}

// Parses all the packages and their dependencies and creates pkg.lock file
func createPkgLockFile(projectPath string, dependencies types.DependenciesTag) ([]*parsers.DependencyTree, error) {
    wioPath := projectPath + io.Sep + ".wio"
    librariesLocalPath := wioPath + io.Sep + "pkg" +  io.Sep + vendorPkgFolder
    librariesRemotePath := wioPath + io.Sep + "pkg" +  io.Sep + remotePkgFolder


    dependencyTrees := make([]*parsers.DependencyTree, 0)
    dependencyLock := types.PackagesLockConfig{}
    dependencyLock.Packages = types.PackagesLockTag{}

    // parse all the libraries and their dependencies in vendor folder
    dependencyTrees, err := parsePackages(librariesLocalPath, dependencyTrees, dependencyLock, dependencies, "")
    if err != nil {
        return nil, err
    }

    // parse all the libraries and their dependencies in remote folder
    dependencyTrees, err = parsePackages(librariesRemotePath, dependencyTrees, dependencyLock, dependencies, "")
    if err != nil {
        return nil, err
    }

    // write the lock file
    if err = io.NormalIO.WriteYml(wioPath+io.Sep+lockFileName, &dependencyLock); err != nil {
        return nil, err
    }

    // return Dependency Tree
    return dependencyTrees, nil
}

// Given all the dependency tree data, it will create cmake files for each dependency
func createCMakes(exePath string, projectPath string, board string, framework string, target string,
    depTree *parsers.DependencyTree) (error) {

    // this is where build files for each package specific to each target will be held
    wioBuildPkgPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + "packages"
    if err := os.MkdirAll(wioBuildPkgPath+io.Sep+depTree.Config.Hash, os.ModePerm); err != nil {
        return err
    }

    packagesTemplate, err := io.AssetIO.ReadFile("templates" + io.Sep + "cmake" + io.Sep + "CMakeListsTarget.txt")
    if err != nil {
        return err
    }

    packagesTemplateStr := string(packagesTemplate)
    for i := 0; i < len(depTree.Child); i++ {
        err := createCMakes(exePath, projectPath, board, framework, target, depTree.Child[i])
        if err != nil {
            return nil
        }
    }

    packagesTemplateStr = strings.Replace(packagesTemplateStr, "{{WIO_PATH}}", exePath, -1)
    packagesTemplateStr = strings.Replace(packagesTemplateStr, "{{library-name}}", depTree.Config.Hash, -1)
    packagesTemplateStr = strings.Replace(packagesTemplateStr, "{{lib-path}}", depTree.Config.Path, -1)
    packagesTemplateStr = strings.Replace(packagesTemplateStr, "{{board}}", board, -1)
    packagesTemplateStr = strings.Replace(packagesTemplateStr, "{{framework}}", strings.ToUpper(framework), -1)
    packagesTemplateStr = strings.Replace(packagesTemplateStr, "{{compile-flags}}",
        strings.Join(depTree.Config.Compile_flags, " "), -1)
    packagesTemplateStr += "\n"

    for dep := range depTree.Child {
        packagesTemplateStr += "include_directories(\"" + depTree.Child[dep].Config.Path + "/include\")\n"
        packagesTemplateStr += "target_link_libraries(" + depTree.Config.Hash + " \"${CMAKE_SOURCE_DIR}/../" +
            depTree.Child[dep].Config.Hash + io.Sep + "lib" + depTree.Child[dep].Config.Hash + ".a\")\n\n"
    }

    cMakefilePath := wioBuildPkgPath+io.Sep+depTree.Config.Hash+io.Sep+"CMakeLists.txt"
    io.NormalIO.WriteFile(cMakefilePath,
        []byte(packagesTemplateStr))

    // create execution way file
    executionWayFileName := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + executionWayFile

    f, err := os.OpenFile(executionWayFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil
    }

    BuildPkgPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + depTree.Config.Hash

    if _, err = f.WriteString(BuildPkgPath + "\n"); err != nil {
        return err
    }

    return nil
}

// This parses dependencies and creates cmake files for all of these. It does this for a target that is provided
// to it.
func ParseDepsAndCreateCMakes(projectPath string, board string, framework string, target string,
    dependencies types.DependenciesTag) (error) {

    // this will create a lock file for packages and will provide us with a dependency tree for us to use
    dependencyTree, err := createPkgLockFile(projectPath, dependencies)
    if err != nil {
        return err
    }
    exePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    // This file describe the way these dependencies are linked
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
