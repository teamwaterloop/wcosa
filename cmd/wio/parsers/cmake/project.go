// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/cmake package, which contains parser to create cmake files
// This file parses dependencies and creates CMakeLists.txt file for the whole project
package cmake

import (
    "path/filepath"
    . "wio/cmd/wio/utils/io"
    "strings"
    "wio/cmd/wio/utils/types"
)

// This creates the main cmake file based on the target. This method is used for creating the main cmake for project
// type of "lib"
func CreateMainCMakeListsFileLib(projectPath string, board string, framework string, target string, targetFlags []string, libFlags []string) (error) {
    projectName := filepath.Base(projectPath)
    executablePath, err := NormalIO.GetRoot()

    if err != nil {
        return err
    }

    lockFilePath := projectPath + Sep + ".wio" + Sep + lockFileName
    targetPath := projectPath + Sep + ".wio" + Sep + "targets" + Sep + target

    toolChainPath := executablePath + "/toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := AssetIO.ReadFile("templates/cmake/CMakeListsLib.txt.tpl")

    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{toolchain-path}}", toolChainPath, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{project-name}}", projectName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-name}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{board}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{framework}}", strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-flags}}", strings.Join(targetFlags, " "), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{lib-flags}}", strings.Join(libFlags, " "), -1)

    lockConfig := types.LibrariesLockConfig{}

    if err = NormalIO.ParseYml(lockFilePath, &lockConfig); err != nil {
        return err
    }

    linkLibraryString := ""

    for lib := range lockConfig.Libraries {
        if !strings.Contains(lib, "__") {
            linkLibraryString += "include_directories(\"" + lockConfig.Libraries[lib].Path + "/include" + "\")\n"
        }
    }

    linkLibraryString += "\n"

    for lib := range lockConfig.Libraries {
        linkLibraryString += "target_link_libraries(" + target + " \"" + targetPath + Sep + "libraries" + Sep + lib + "/" +
            "lib" + lockConfig.Libraries[lib].Hash + ".a" + "\")\n"
    }

    if linkLibraryString == "\n" {
        linkLibraryString = ""
    }

    templateDataStr = strings.Replace(templateDataStr, "{{link-library}}", linkLibraryString, -1)

    return NormalIO.WriteFile(projectPath + Sep + ".wio" + Sep + "CMakeLists.txt", []byte(templateDataStr))
}

// This creates the main cmake file based on the target. This method is used for creating the main cmake for project
// type of "app"
func CreateMainCMakeListsFileApp(projectPath string, board string, framework string, target string, targetFlags []string) (error) {
    projectName := filepath.Base(projectPath)
    executablePath, err := NormalIO.GetRoot()

    if err != nil {
        return err
    }

    lockFilePath := projectPath + Sep + ".wio" + Sep + lockFileName
    targetPath := projectPath + Sep + ".wio" + Sep + "targets" + Sep + target

    toolChainPath := executablePath + "/toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := AssetIO.ReadFile("templates/cmake/CMakeLists.txt.tpl")

    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{toolchain-path}}", toolChainPath, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{project-name}}", projectName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{project-path}}", projectPath, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-name}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{board}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{framework}}", strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-flags}}", strings.Join(targetFlags, " "), -1)

    lockConfig := types.LibrariesLockConfig{}

    if err = NormalIO.ParseYml(lockFilePath, &lockConfig); err != nil {
        return err
    }

    for lib := range lockConfig.Libraries {
        if !strings.Contains(lib, "__") {
            templateDataStr += "include_directories(\"" + lockConfig.Libraries[lib].Path + "/include" + "\")\n"
        }
    }

    templateDataStr += "\n"

    for lib := range lockConfig.Libraries {
        templateDataStr += "target_link_libraries(" + target + " \"" + targetPath + Sep + "libraries" + Sep + lib + "/" +
            "lib" + lockConfig.Libraries[lib].Hash + ".a" + "\")\n"
    }

    return NormalIO.WriteFile(projectPath + Sep + ".wio" + Sep + "CMakeLists.txt", []byte(templateDataStr))
}

// Handles creation and update of all the cmake files. It parses libraries and creates cmake files for each target
// and based on project type
func HandleCMakeCreation(projectPath string, framework string, targetsTag types.TargetsTag, librariesTag types.LibrariesTag, isLib bool, libFlags []string) (error) {
    defaultTarget := targetsTag.Targets[targetsTag.Default_target]

    for target := range targetsTag.Targets {
        currTarget := targetsTag.Targets[target]

        // create cmake files for each target libraries
        if err := PopulateCMakeFilesForLibs(projectPath, currTarget.Board, framework, target, librariesTag); err != nil {
            return err
        }
    }

    if isLib {
        return CreateMainCMakeListsFileLib(projectPath, defaultTarget.Board, framework, targetsTag.Default_target, defaultTarget.Compile_flags, libFlags)
    } else {
        return CreateMainCMakeListsFileApp(projectPath, defaultTarget.Board, framework, targetsTag.Default_target, defaultTarget.Compile_flags)
    }
}
