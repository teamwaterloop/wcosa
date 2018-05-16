// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of parsers/cmake package, which contains parser to create cmake files
// This file parses dependencies and creates CMakeLists.txt file for the whole project
package cmake

import (
    "strings"
    "wio/cmd/wio/parsers"
    "wio/cmd/wio/utils/io"
    "path/filepath"
)

// This creates the main cmake file based on the target provided. This method is used for creating the main cmake for
// project type of "pkg". This CMake file links the package with the target provided so that it can tested and run
// before getting shipped
func CreatePkgMainCMakeLists(pkgName string, pkgPath string, board string, framework string, target string,
    targetFlags []string, pkgFlags []string, depTree []*parsers.DependencyTree) (error) {

    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    toolChainPath := executablePath + io.Sep + "toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := io.AssetIO.ReadFile("templates/cmake/CMakeListsPkg.txt.tpl")

    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{TOOLCHAIN_FILE}}",
        filepath.ToSlash(toolChainPath), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PROJECT_PATH}}", filepath.ToSlash(pkgPath), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PROJECT_NAME}}", pkgName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_NAME}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{BOARD}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{FRAMEWORK}}",
        strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_COMPILE_FLAGS}}",
        strings.Join(targetFlags, " "), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PKG_COMPILE_FLAGS}}",
        strings.Join(pkgFlags, " "), -1)
    templateDataStr += "\n"

    for dep := range depTree {
        currLinkString := "target_link_libraries({{TARGET_NAME}} PUBLIC {{DEP_NAME}})"
        currLinkString = strings.Replace(currLinkString, "{{TARGET_NAME}}", pkgName, -1)
        currLinkString = strings.Replace(currLinkString, "{{DEP_NAME}}", depTree[dep].Config.Hash, -1)

        templateDataStr += currLinkString + "\n"
    }

    return io.NormalIO.WriteFile(pkgPath+io.Sep+".wio"+io.Sep+"build"+io.Sep+"CMakeLists.txt",
        []byte(templateDataStr))
}

// This creates the main cmake file based on the target. This method is used for creating the main cmake for project
// type of "app". In this it does not link any library but rather just populates a target that can be uploaded
func CreateAppMainCMakeLists(appName string, appPath string, board string, framework string, target string,
    targetFlags []string, depTree []*parsers.DependencyTree) (error) {

    executablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    toolChainPath := executablePath + io.Sep + "toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := io.AssetIO.ReadFile("templates/cmake/CMakeListsApp.txt.tpl")
    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{TOOLCHAIN_FILE}}",
        filepath.ToSlash(toolChainPath), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PROJECT_PATH}}", filepath.ToSlash(appPath), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{PROJECT_NAME}}", appName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_NAME}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{BOARD}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{FRAMEWORK}}",
        strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{TARGET_COMPILE_FLAGS}}",
        strings.Join(targetFlags, " "), -1)
    templateDataStr += "\n"

    for dep := range depTree {
        currLinkString := "target_link_libraries(${TARGET_NAME} PUBLIC {{DEP_NAME}})"
        currLinkString = strings.Replace(currLinkString, "{{DEP_NAME}}", depTree[dep].Config.Hash, -1)

        templateDataStr += currLinkString + "\n"
    }

    return io.NormalIO.WriteFile(appPath+io.Sep+".wio"+io.Sep+"build"+io.Sep+"CMakeLists.txt",
        []byte(templateDataStr))
}
