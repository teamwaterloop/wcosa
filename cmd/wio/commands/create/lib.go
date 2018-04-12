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

    if err = copyTemplates(lib.args.Directory, lib.args.AppType, lib.args.Ide, "config"+Sep+"create_paths.json"); err != nil {
        return err
    }
    inf, err := lib.FillConfig()
    if err != nil {
        return err
    }

    config = inf.(*types.LibConfig)

    return HandleCMakeCreation(lib.args.Directory, lib.args.Framework, config.TargetsTag, config.LibrariesTag)
}

// Prints all the commands relevant to library type
func (lib Lib) printNextCommands() {
    Norm.Cyan("`wio build -h`\n")
    Norm.Cyan("`wio run -h`\n")
    Norm.Cyan("`wio upload -h`\n")
    Norm.Cyan("`wio test -h`\n")
}

// Handles config file for lib
func (lib Lib) FillConfig() (interface{}, error) {
    Verb.Verbose("* Loaded wio.yml file template\n")

    libConfig := types.LibConfig{}
    if err := NormalIO.ParseYml(lib.args.Directory+Sep+"wio.yml", &libConfig);
        err != nil {
        return nil, err
    }

    // make modifications to the data
    libConfig.MainTag.Ide = lib.args.Ide
    libConfig.MainTag.Platform = lib.args.Platform
    libConfig.MainTag.Framework = AppendIfMissing(libConfig.MainTag.Framework, lib.args.Framework)
    libConfig.MainTag.Name = filepath.Base(lib.args.Directory)

    if libConfig.TargetsTag.Default_target == "" {
        libConfig.TargetsTag.Default_target = "test"
    }

    if target, ok := libConfig.TargetsTag.Targets[libConfig.TargetsTag.Default_target]; ok {
        defaultTarget := &types.TargetTag{}
        defaultTarget.Board = lib.args.Board
        defaultTarget.Compile_flags = target.Compile_flags
        libConfig.TargetsTag.Targets[libConfig.TargetsTag.Default_target] = defaultTarget
    } else {
        libConfig.TargetsTag.Targets[libConfig.TargetsTag.Default_target] = &types.TargetTag{}
        libConfig.TargetsTag.Targets[libConfig.TargetsTag.Default_target].Board = lib.args.Board
    }

    Verb.Verbose("* Modified information in the configuration\n")

    if err := PrettyPrintConfig(&libConfig, lib.args.Directory+Sep+"wio.yml");
        err != nil {
        return nil, err
    }
    Verb.Verbose("* Filled/Updated template written back to the file\n")

    return &libConfig, nil
}

func (lib Lib) update() (error) {
    return genericUpdate(lib, lib.args)
}
