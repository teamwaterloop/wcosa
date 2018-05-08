// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/cmake package, which contains parser to create cmake files
// This file parses dependencies and creates CMakeLists.txt file for the whole project
package cmake

import (
    . "wio/cmd/wio/utils/io"
    "strings"
    "wio/cmd/wio/types"
    "wio/cmd/wio/utils"
)

// This creates the main cmake file based on the target provided. This method is used for creating the main cmake for
// project type of "pkg". This CMake file links the package with the target provided so that it can tested and run
// before getting shipped
func CreatePkgMainCMakeLists(pkgName string, pkgPath string, board string, framework string, target string,
    targetFlags []string, pkgFlags []string) (error) {

    executablePath, err := NormalIO.GetRoot()
    if err != nil {
        return err
    }

    lockFilePath := pkgPath + Sep + ".wio" + Sep + lockFileName
    targetPath := pkgPath + Sep + ".wio" + Sep + "targets" + Sep + target

    toolChainPath := executablePath + "/toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := AssetIO.ReadFile("templates/cmake/CMakeListsPkg.txt.tpl")

    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{toolchain-path}}", toolChainPath, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{project-name}}", pkgName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-name}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{board}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{framework}}", strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-flags}}", strings.Join(targetFlags, " "), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{pkg-flags}}", strings.Join(pkgFlags, " "), -1)

    linkPackagesString := ""

    if utils.PathExists(lockFilePath) {
        lockConfig := types.PackagesLockConfig{}

        if err = NormalIO.ParseYml(lockFilePath, &lockConfig); err != nil {
            return err
        }

        for pkg := range lockConfig.Packages {
            if !strings.Contains(pkg, "__") {
                linkPackagesString += "include_directories(\"" + lockConfig.Packages[pkg].Path + "/include" + "\")\n"
            }
        }

        linkPackagesString += "\n"

        for pkg := range lockConfig.Packages {
            linkPackagesString += "target_link_libraries(" + target + " \"" + targetPath + Sep + "libraries" + Sep +
                pkg + "/" + "lib" + lockConfig.Packages[pkg].Hash + ".a" + "\")\n"
        }

        if linkPackagesString == "\n" {
            linkPackagesString = ""
        }
    }

    templateDataStr = strings.Replace(templateDataStr, "{{link-library}}", linkPackagesString, -1)

    return NormalIO.WriteFile(pkgPath + Sep + ".wio" + Sep + "CMakeLists.txt", []byte(templateDataStr))
}

// This creates the main cmake file based on the target. This method is used for creating the main cmake for project
// type of "app". In this it does not link any library but rather just populates a target that can be uploaded
func CreateAppMainCMakeLists(appName string, appPath string, board string, framework string, target string,
    targetFlags []string) (error) {

    executablePath, err := NormalIO.GetRoot()
    if err != nil {
        return err
    }

    lockFilePath := appPath + Sep + ".wio" + Sep + lockFileName
    targetPath := appPath + Sep + ".wio" + Sep + "targets" + Sep + target

    toolChainPath := executablePath + "/toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := AssetIO.ReadFile("templates/cmake/CMakeListsApp.txt.tpl")
    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{toolchain-path}}", toolChainPath, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{project-name}}", appName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-name}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{board}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{framework}}", strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-flags}}", strings.Join(targetFlags, " "), -1)

    if utils.PathExists(lockFilePath) {
        lockConfig := types.PackagesLockConfig{}

        if err = NormalIO.ParseYml(lockFilePath, &lockConfig); err != nil {
            return err
        }

        for lib := range lockConfig.Packages {
            if !strings.Contains(lib, "__") {
                templateDataStr += "include_directories(\"" + lockConfig.Packages[lib].Path + "/include" + "\")\n"
            }
        }

        templateDataStr += "\n"

        for pkg := range lockConfig.Packages {
            templateDataStr += "target_link_libraries(" + target + " \"" + targetPath + Sep + "packages" + Sep + pkg + "/" +
                "lib" + lockConfig.Packages[pkg].Hash + ".a" + "\")\n"
        }
    }

    return NormalIO.WriteFile(appPath + Sep + ".wio" + Sep + "CMakeLists.txt", []byte(templateDataStr))
}

// This creates all the cmake files necessary for the project to compile. It also generates a dependency tree based
// on the project structure and links all the packages together. It generates all this for all the targets that are
// defined in the wio.yml file
func FullCMakeCreationWithDepsParsing(projectName string, projectPath string, framework string, targets types.TargetsTag,
    dependencies types.DependenciesTag, isPkg bool, pkgFlags []string) (error) {

    // create a separate set of CMake files for each target
    for target := range targets.Targets {
        currTarget := targets.Targets[target]

        // Parse dependencies and then create CMakeLists file for each dependency
        if err := ParseDepsAndCreateCMakes(projectPath, currTarget.Board, framework, target, dependencies);
        err != nil {
            return err
       }
    }

    defaultTarget := targets.Targets[targets.Default_target]

    if isPkg {
        return CreatePkgMainCMakeLists(projectName, projectPath, defaultTarget.Board, framework,
            targets.Default_target, defaultTarget.Compile_flags, pkgFlags)
    } else {
        return CreateAppMainCMakeLists(projectName, projectPath, defaultTarget.Board, framework,
            targets.Default_target, defaultTarget.Compile_flags)
    }
}
