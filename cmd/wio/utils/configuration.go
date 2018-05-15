// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package utils contains utilities/files useful throughout the app
// This file contains all the function to manipulate project configuration file

package utils

import (
    "bufio"
    "strings"
    "fmt"
    "gopkg.in/yaml.v2"
    "os"

)

// Adds spacing and other formatting for project configuration
func writeProjectConfig(lines []string, path string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    for _, line := range lines {
        tokens := strings.Split(line, "\n")
        for _, token := range tokens {
            if strings.Contains(token, "targets:") ||
                (strings.Contains(token, "libraries:") && !strings.Contains(token, "#   libraries:")) {
                fmt.Fprint(w, "\n")
            }
            fmt.Fprintln(w, token)
        }
    }

    return w.Flush()
}

// Write configuration for the project with information on top and nice spacing
func PrettyPrintConfig(projectConfig interface{}, filePath string) (error) {
    appInfoPath := "templates" + Sep + "config" + Sep + "app-helper.txt"
    pkgInfoPath := "templates" + Sep + "config" + Sep + "pkg-helper.txt"
    targetsInfoPath := "templates" + Sep + "config" + Sep + "targets-helper.txt"
    dependenciesInfoPath := "templates" + Sep + "config" + Sep + "dependencies-helper.txt"

    var ymlData []byte
    var appInfoData []byte
    var pkgInfoData []byte
    var targetsInfoData []byte
    var dependenciesInfoData []byte
    var err error

    // get data
    if ymlData, err = yaml.Marshal(projectConfig); err != nil {
        return err
    }
    if appInfoData, err = AssetIO.ReadFile(appInfoPath); err != nil {
        return err
    }
    if pkgInfoData, err = AssetIO.ReadFile(pkgInfoPath); err != nil {
        return err
    }
    if targetsInfoData, err = AssetIO.ReadFile(targetsInfoPath); err != nil {
        return err
    }
    if dependenciesInfoData, err = AssetIO.ReadFile(dependenciesInfoPath); err != nil {
        return err
    }

    finalString := ""
    currentString := strings.Split(string(ymlData), "\n")

    beautify := false
    first := false

    for line := range currentString {
        currLine := currentString[line]

        if len(currLine) <= 1 {
            continue
        }

        if strings.Contains(currLine, "app:") {
            finalString += string(appInfoData) + "\n"
        } else if strings.Contains(currLine, "pkg:") {
            finalString += string(pkgInfoData) + "\n"
        } else if strings.Contains(currLine, "targets:") {
            finalString += "\n" + string(targetsInfoData) + "\n"
        } else if strings.Contains(currLine, "create:") {
            beautify = true
        } else if strings.Contains(currLine, "dependencies:") {
            beautify = true
            first = false
            finalString += "\n" + string(dependenciesInfoData) + "\n"
        } else if beautify && !first {
            first = true
        } else if !strings.Contains(currLine, "compile_flags:") && beautify {
            simpleString := strings.Trim(currLine, " ")

            if simpleString[len(simpleString)-1] == ':' {
                finalString += "\n"
            }
        }

        finalString += currLine + "\n"
    }

    err = NormalIO.WriteFile(filePath, []byte(finalString))

    return err
}
