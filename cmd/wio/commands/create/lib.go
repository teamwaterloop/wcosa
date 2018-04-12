// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates a library to be published
package create

import (
    "path/filepath"

    . "wio/cmd/wio/utils/io"
    . "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/types"
    . "wio/cmd/wio/parsers/cmake"
)

// Creates project structure for library type
func (lib Lib) createStructure() (error) {
    Verb.Verbose("\n")
    if err := createStructure(lib.args.Directory, "src", "lib", "test", ".wio/targets"); err != nil {
        return err
    }

    return nil
}

// Prints the project structure for library type
func (lib Lib) printProjectStructure() {
    Norm.Cyan("src    - put your source files here.\n")
    Norm.Cyan("lib    - libraries for the project go here.\n")
    Norm.Cyan("test   - put your files for unit testing here.\n")
}

// Creates a template project that is ready to build and upload for library type
func (lib Lib) createTemplateProject() (error) {
    config := &types.LibConfig{}
    var err error

    if err = copyTemplates(lib.args); err != nil { return err }
    if config, err = lib.FillConfig(); err != nil { return err }

    for target := range config.TargetsTag.Targets {
        // create cmake files for each target libraries
        if err = PopulateCMakeFilesForLibs(lib.args.Directory, lib.args.Board, target, config.LibrariesTag); err != nil {
            return err
        }
    }

    CreateMainCMakeListsFile(lib.args.Directory, config.TargetsTag.Targets[config.TargetsTag.Default_target].Board,
        lib.args.Framework, config.TargetsTag.Default_target,
            config.TargetsTag.Targets[config.TargetsTag.Default_target].Compile_flags)

    return nil
}

// Prints all the commands relevant to library type
func (lib Lib) printNextCommands() {
    Norm.Cyan("`wio build -h`\n")
    Norm.Cyan("`wio run -h`\n")
    Norm.Cyan("`wio upload -h`\n")
    Norm.Cyan("`wio test -h`\n")
}

// Handles config file for lib
func (lib Lib) FillConfig() (*types.LibConfig, error) {
    Verb.Verbose("* Loaded wio.yml file template\n")

    libConfig := types.LibConfig{}
    if err := NormalIO.ParseYml(lib.args.Directory + Sep + "wio.yml", &libConfig);
    err != nil { return nil, err }

    // make modifications to the data
    libConfig.MainTag.Ide = lib.args.Ide
    libConfig.MainTag.Platform = lib.args.Platform
    libConfig.MainTag.Framework = AppendIfMissing(libConfig.MainTag.Framework, lib.args.Framework)
    libConfig.MainTag.Name = filepath.Base(lib.args.Directory)

    if libConfig.TargetsTag.Default_target == "" {
        libConfig.TargetsTag.Default_target = "test"
    }

    if libConfig.TargetsTag.Targets[libConfig.TargetsTag.Default_target].Board == "" {
        libConfig.TargetsTag.Targets[libConfig.TargetsTag.Default_target] = &types.TargetTag{}
        libConfig.TargetsTag.Targets[libConfig.TargetsTag.Default_target].Board = lib.args.Board
    }

    Verb.Verbose("* Modified information in the configuration\n")

    if err := PrettyPrintConfig(&libConfig, lib.args.Directory + Sep + "wio.yml");
    err != nil { return nil, err }
    Verb.Verbose("* Filled/Updated template written back to the file\n")

    return &libConfig, nil
}

func (lib Lib) FillCMake(paths map[string]string) (error) {
    return nil
}
