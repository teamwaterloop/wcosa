// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// This contains helper function for generic template parsing and creating the project
package create

import (
    "regexp"
    "path/filepath"

    . "wio/cmd/wio/utils/io"
    "os"
    "wio/cmd/wio/utils/types"
    "wio/cmd/wio/utils"
    "wio/cmd/wio/parsers/cmake"
    "strings"
    "errors"
)

// Parses the paths.json file and uses that to get paths to copy files to and from
// It also stores all the paths as a map to be used later on
func parsePathsAndCopy(jsonPath string, projectPath string, tags []string) (error) {
    var paths = Paths{}
    if err := AssetIO.ParseJson(jsonPath, &paths); err != nil {
        return err
    }

    var re, e = regexp.Compile(`{{.+}}`)
    if e != nil {
        return e
    }

    var sources []string
    var destinations []string
    var overrides []bool

    for i := 0; i < len(paths.Paths); i++ {
        for t := 0; t < len(tags); t++ {
            if paths.Paths[i].Id == tags[t] {
                sources = append(sources, paths.Paths[i].Src)

                destination, e := filepath.Abs(re.ReplaceAllString(paths.Paths[i].Des, projectPath))
                if e != nil {
                    return e
                }

                destinations = append(destinations, destination)
                overrides = append(overrides, paths.Paths[i].Override)
            }
        }
    }

    return AssetIO.CopyMultipleFiles(sources, destinations, overrides)
}

// generic method to create structure based on a list of paths provided
func createStructure(projectPath string, relPaths ...string) (error) {
    for path := 0; path < len(relPaths); path++ {
        fullPath := projectPath + Sep + relPaths[path]

        if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
            return err
        }
        Verb.Verbose(`* Created "` + relPaths[path] + `" folder` + "\n")
    }

    return nil
}

// generic method to copy templates based on command line arguments
func copyTemplates(projectPath string, appType string, ide string, jsonPath string) (error) {
    strArray := make([]string, 1)
    strArray[0] = appType + "-gen"

    Verb.Verbose("\n")
    if ide == "clion" {
        Verb.Verbose("* Clion Ide available so ide template set up will be used\n")
        strArray = append(strArray, appType+"-clion")
    } else {
        Verb.Verbose("* General template setup will be used\n")
    }

    if err := parsePathsAndCopy(jsonPath, projectPath, strArray); err != nil {
        return err
    }
    Verb.Verbose("* All Template files created in their right position\n")

    return nil
}

// Appends a string to string slice only if it is missing
func AppendIfMissing(slice []string, i string) []string {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}

// Generic update updates the project. This is a generic method used by both app and lib
func genericUpdate(projectType ProjectTypes, cliArgs *types.CliArgs) (error) {
    wioFile := cliArgs.Directory + Sep + "wio.yml"
    configApp := &types.AppConfig{}
    configLib := &types.LibConfig{}
    var targetsTag types.TargetsTag
    var librariesTag types.LibrariesTag
    var config interface{}
    var err error

    fillConfig := false
    if utils.PathExists(wioFile) {
        // check if the config file is valid yml
        if cliArgs.AppType == "app" {
            err = NormalIO.ParseYml(wioFile, configApp)
        } else if cliArgs.AppType == "lib" {
            err = NormalIO.ParseYml(wioFile, configLib)
        }

        if err != nil {
            Verb.Verbose("* Invalid yml file (wio.yml), deleting the file to create a new one\n")
            fillConfig = true
            os.Remove(wioFile)
        } else {
            // check if config is of right project type
            configStr, err := NormalIO.ReadFile(wioFile)
            if err != nil {
                return err
            }

            if !strings.Contains(string(configStr), cliArgs.AppType + ":") {
                return errors.New("current project is of a different type than the one specified in the update\n")
            }
        }
    } else {
        fillConfig = true
    }

    projectType.createStructure()

    // copy templates
    if err := copyTemplates(cliArgs.Directory, cliArgs.AppType, cliArgs.Ide, "config"+Sep+"update_paths.json"); err != nil {
        return err
    }

    if fillConfig {
        Verb.Verbose("* Filling the new wio.yml with details\n")
        config, err = projectType.FillConfig()
        if err != nil {
            return err
        }
    } else {
        Verb.Verbose("* Using configurations from the current wio.yml file\n")
    }

    Verb.Verbose("* Creating CMake files for all the targets and libraries\n")
    if cliArgs.AppType == "app" {
        if fillConfig {
            configApp = config.(*types.AppConfig)
        }
        targetsTag = configApp.TargetsTag
        librariesTag = configApp.LibrariesTag

        return cmake.HandleCMakeCreation(cliArgs.Directory, cliArgs.Framework, targetsTag, librariesTag, false, nil)
    } else if cliArgs.AppType == "lib" {
        if fillConfig {
            configLib = config.(*types.LibConfig)
        }
        targetsTag = configLib.TargetsTag
        librariesTag = configLib.LibrariesTag

        return cmake.HandleCMakeCreation(cliArgs.Directory, cliArgs.Framework, targetsTag, librariesTag, true, configLib.MainTag.Compile_flags)
    }

    return nil
}
