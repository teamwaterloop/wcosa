// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// Sub command of create which creates an executable application
package create

import (
    "path/filepath"

    . "wio/cmd/wio/utils/io"
    . "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/types"
    . "wio/cmd/wio/parsers/cmake"
)

// Creates project structure for application type
func (app App) createStructure() (error) {
    Verb.Verbose("\n")
    if err := createStructure(app.args.Directory, "src", "lib", ".wio/targets"); err != nil {
        return err
    }

    if app.args.Tests {
        if err := createStructure(app.args.Directory, "test"); err != nil {
            return err
        }
    }

    return nil
}

// Prints the project structure for application type
func (app App) printProjectStructure() {
    Norm.Cyan("src    - put your source files here.\n")
    Norm.Cyan("lib    - libraries for the project go here.\n")
    if app.args.Tests {
        Norm.Cyan("test   - put your files for unit testing here.\n")
    }
}

// Creates a template project that is ready to build and upload for application type
func (app App) createTemplateProject() (error) {
    config := &types.AppConfig{}
    var err error

    if err = copyTemplates(app.args.Directory, app.args.AppType, app.args.Ide, "config"+Sep+"create_paths.json"); err != nil {
        return err
    }
    inf, err := app.FillConfig()
    if err != nil {
        return err
    }

    config = inf.(*types.AppConfig)

    return HandleCMakeCreation(app.args.Directory, app.args.Framework, config.TargetsTag, config.LibrariesTag, false, nil)
}

// Prints all the commands relevant to application type
func (app App) printNextCommands() {
    Norm.Cyan("`wio build -h`\n")
    Norm.Cyan("`wio run -h`\n")
    Norm.Cyan("`wio upload -h`\n")

    if app.args.Tests {
        Norm.Cyan("`wio test -h`\n")
    }
}

// Fill config file for the app
func (app App) FillConfig() (interface{}, error) {
    Verb.Verbose("* Loaded wio.yml file template\n")

    appConfig := types.AppConfig{}
    if err := NormalIO.ParseYml(app.args.Directory+Sep+"wio.yml", &appConfig);
        err != nil {
        return nil, err
    }

    // make modifications to the data
    appConfig.MainTag.Ide = app.args.Ide
    appConfig.MainTag.Platform = app.args.Platform
    appConfig.MainTag.Framework = app.args.Framework
    appConfig.MainTag.Name = filepath.Base(app.args.Directory)

    if appConfig.TargetsTag.Default_target == "" {
        appConfig.TargetsTag.Default_target = "main"
    }

    if target, ok := appConfig.TargetsTag.Targets[appConfig.TargetsTag.Default_target]; ok {
        defaultTarget := &types.TargetTag{}
        defaultTarget.Board = app.args.Board
        defaultTarget.Compile_flags = target.Compile_flags
        appConfig.TargetsTag.Targets[appConfig.TargetsTag.Default_target] = defaultTarget
    } else {
        appConfig.TargetsTag.Targets[appConfig.TargetsTag.Default_target] = &types.TargetTag{}
        appConfig.TargetsTag.Targets[appConfig.TargetsTag.Default_target].Board = app.args.Board
    }

    Verb.Verbose("* Modified information in the configuration\n")

    if err := PrettyPrintConfig(&appConfig, app.args.Directory+Sep+"wio.yml");
        err != nil {
        return nil, err
    }
    Verb.Verbose("* Filled/Updated template written back to the file\n")

    return &appConfig, nil
}

// Handles the update of app
func (app App) update() (error) {
    return genericUpdate(app, app.args)
}
