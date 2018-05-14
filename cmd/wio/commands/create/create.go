// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create and update command and sub commands provided by the tool.
// Creates, updates and initializes a wio project.
package create

import (
    "github.com/urfave/cli"
    "path/filepath"
    "wio/cmd/wio/utils/io/log"
    "os"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/types"
    "regexp"
    "wio/cmd/wio/utils"
    "bufio"
    "strings"
    "wio/cmd/wio/commands"
)

const (
    APP = "app"
    PKG = "pkg"
)

/// This structure wraps all the important features needed for a create and update command
type PacketCreate struct {
    directory string
    name string
    board string
    framework string
    platform string
    ide string
    tests bool
}

type Create struct {
    Context *cli.Context
    Type    string
    Update  bool
    error error
}

// Executes the create command
func (create Create) Execute() {
    commands.RecordError(create.error, "")

    if !create.Update && len(create.Context.Args()) < 2 {
        // When we are creating a project we need both directory and a board from args
        log.Error(true, "Project directory or Board not specified")
    } else if !create.Update && len(create.Context.Args()) < 1 {
        // When we are updating a project we only need directory from args
        log.Error(true, "Project directory not specified")
    }

    createPacket := &PacketCreate{}

    // fetch directory based on the argument
    directory, err := filepath.Abs(create.Context.Args()[0])
    commands.RecordError(err, "")

    // based on the command line context, get all the important fields
    createPacket.directory = directory
    createPacket.name = filepath.Base(directory)
    createPacket.framework = create.Context.String("framework")
    createPacket.platform = create.Context.String("platform")
    createPacket.ide = create.Context.String("ide")
    createPacket.tests = create.Context.Bool("tests")

    if create.Update {
        // if update command is called
        if create.checkUpdate(directory) {
            // after check if an update can be performed
            createPacket.board = create.Context.String("board")
            create.updateProject(createPacket)
        } else {
            // project structure is invalid or too distorted for an update to be applied
            message := `Based on your current project structure and wio.yml file, update cannot be performed.
You can use: wio create <app type> DIRECTORY BOARD`
            log.Norm.Cyan(true, message)
        }
    } else if create.createCheck(directory) {
        // if create command is called and it's check is passed
        createPacket.board = create.Context.Args()[1]

        create.prePrint(createPacket)
        create.createStructure(createPacket.directory, true)
        create.initialProjectSetup(createPacket)
        create.postPrint()
    }
}

/// This method is a helper function that prints instructions to the console before "create"
/// or "update" command is executed. This has information about the project and it's ty[e
func (create Create) prePrint(createPacket *PacketCreate) {
    log.Norm.Yellow(false, "Project Name: ")
    log.Norm.Cyan(true, createPacket.name)
    log.Norm.Yellow(false, "Project Type: ")
    log.Norm.Cyan(true, create.Type)
    log.Norm.Yellow(false, "Project Path: ")
    log.Norm.Cyan(true, createPacket.directory)
}

/// This method is used to create the project. This erases everything that there is in the directory
/// and creates the project from scratch. This is the reason that this method is only called after
/// a check has been performed
func (create Create) createStructure(directory string, delete bool) {
    if delete {
        log.Norm.Yellow(false, "Creating project structure ... ")
    } else {
        log.Norm.Yellow(false, "Updating project structure ... ")
    }

    if delete && utils.PathExists(directory) {
        commands.RecordError(os.RemoveAll(directory), "failure")
    }

    commands.RecordError(os.MkdirAll(directory+io.Sep+"src", os.ModePerm), "failure")
    commands.RecordError(os.MkdirAll(directory+io.Sep+".wio"+io.Sep+"build", os.ModePerm), "failure")

    if create.Type == PKG {
        /// each package will have an include method
        commands.RecordError(os.MkdirAll(directory+io.Sep+"include", os.ModePerm),
            "failure")
    }

    if create.Context.Bool("tests") {
        commands.RecordError(os.MkdirAll(directory+io.Sep+"tests", os.ModePerm), "failure")
    }
}

// This is a method that updates the created project. It will fix the structure and apply
// updates. It also makes sure all the configurations are applied
func (create Create) updateProject(createPacket *PacketCreate) {
    create.prePrint(createPacket)
    create.createStructure(createPacket.directory, false)
    create.updateProjectSetup(createPacket)
    create.postPrint()
}


/// This is one of the most important step as this sets up the project when update command is used.
/// This also updates the wio.yml file so that it can be fixed and current configurations can be applied.
func (create Create) updateProjectSetup(createPacket *PacketCreate) {
    log.Norm.Green(true, "success")
    log.Norm.Yellow(false, "Updating the project ... ")

    if createPacket.ide == "clion" {
        // copy gitignore file
        io.AssetIO.CopyFile("templates/gitignore/.gitignore-clion", createPacket.directory + "/.gitignore",
            false)
    } else {
        // copy gitignore file
        io.AssetIO.CopyFile("templates/gitignore/.gitignore-general", createPacket.directory + "/.gitignore",
            false)
    }

    // get default configuration values
    defaults := types.DConfig{}
    commands.RecordError(io.AssetIO.ParseYml("config/defaults.yml", &defaults), "failure")


    var config interface{}

    if create.Type == APP {
        projectConfig := &types.AppConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.directory + io.Sep + "wio.yml", projectConfig),
            "failure")

        // update the name of the project
        projectConfig.MainTag.Name = createPacket.name
        // update the targets to make sure they are valid and there is a default target
        create.handleTargets(&projectConfig.TargetsTag, defaults.Board)
        // set the default board to be from the default target
        createPacket.board = projectConfig.TargetsTag.Targets[projectConfig.TargetsTag.Default_target].Board

        // check framework and platform
        checkFrameworkAndPlatform(&projectConfig.MainTag.Framework,
            &projectConfig.MainTag.Platform, &defaults)

        config = projectConfig
    } else {
        projectConfig := &types.PkgConfig{}

        commands.RecordError(io.NormalIO.ParseYml(createPacket.directory + io.Sep + "wio.yml", projectConfig),
            "failure")

        // update the name of the project
        projectConfig.MainTag.Name = createPacket.name
        // update the targets to make sure they are valid and there is a default target
        create.handleTargets(&projectConfig.TargetsTag, defaults.Board)
        // set the default board to be from the default target
        createPacket.board = projectConfig.TargetsTag.Targets[projectConfig.TargetsTag.Default_target].Board

        // make sure boards are updated in yml file
        if !utils.StringInSlice("ALL", projectConfig.MainTag.Board) {
            projectConfig.MainTag.Board = []string{"ALL"}
        } else if !utils.StringInSlice(createPacket.board, projectConfig.MainTag.Board) {
            projectConfig.MainTag.Board = append(projectConfig.MainTag.Board, createPacket.board)
        }

        // check frameworks and platform
        checkFrameworkArrayAndPlatform(&projectConfig.MainTag.Framework,
            &projectConfig.MainTag.Platform, &defaults)

        config = projectConfig
    }

    commands.RecordError(utils.PrettyPrintConfig(config, createPacket.directory+io.Sep+"wio.yml"),
        "failure")
}

// This function checks if framework and platform are not empty. It future we can in force in valud
// frameworks and platforms using this
func checkFrameworkAndPlatform(framework *string, platform *string, defaults *types.DConfig) {
    if *framework == "" {
        *framework = defaults.Framework
    }

    if *platform == "" {
        *platform = defaults.Platform
    }
}

// This function is similar to the above but in this case it checks if multiple frameworks are invalid
// and same goes for platform
func checkFrameworkArrayAndPlatform(framework *[]string, platform *string, defaults *types.DConfig) {
    if len(*framework) == 0 {
        *framework = append(*framework, defaults.Framework)
    }

    if *platform == "" {
        *platform = defaults.Platform
    }
}

/// This is one of the most important step as this sets up the project when create command is used.
/// This also fills up the wio.yml file so that default configuration along with user choices
/// are applied.
func (create Create) initialProjectSetup(createPacket *PacketCreate) {
    log.Norm.Green(true, "success")
    log.Norm.Yellow(false, "Creating template project ... ")

    defaultTarget := "main"

    commands.RecordError(copyTemplates(createPacket.directory, create.Type, createPacket.ide,
        "config"+io.Sep+"create_paths.json"), "failure")

    var config interface{}

    if create.Type == APP {
        projectConfig := &types.AppConfig{}

        // make modifications to the data
        projectConfig.MainTag.Ide = createPacket.ide
        projectConfig.MainTag.Platform = createPacket.platform
        projectConfig.MainTag.Framework = createPacket.framework
        projectConfig.MainTag.Name = createPacket.name
        projectConfig.TargetsTag.Default_target = defaultTarget
        targets := make(map[string]*types.TargetTag, 1)
        projectConfig.TargetsTag.Targets = targets

        targetsTag := projectConfig.TargetsTag
        create.handleTargets(&targetsTag, createPacket.board)

        config = projectConfig
    } else {
        projectConfig := &types.PkgConfig{}
        defaultTarget = "tests"

        // make modifications to the data
        projectConfig.MainTag.Ide = createPacket.ide
        projectConfig.MainTag.Platform = createPacket.platform
        projectConfig.MainTag.Framework = []string{createPacket.framework}
        projectConfig.MainTag.Name = createPacket.name
        projectConfig.TargetsTag.Default_target = defaultTarget
        targets := make(map[string]*types.TargetTag, 0)
        projectConfig.TargetsTag.Targets = targets
        projectConfig.MainTag.Board = []string{createPacket.board}

        targetsTag := projectConfig.TargetsTag
        create.handleTargets(&targetsTag, createPacket.board)

        config = projectConfig
    }

    commands.RecordError(utils.PrettyPrintConfig(config, createPacket.directory+io.Sep+"wio.yml"),
        "failure")
}

/// This method handles the targets that a user can create and what these targets are
/// in wio.yml file. It targets are not there, it will create a default target. Unless
/// it will keep the targets that are already there
func (create Create) handleTargets(targetsTag *types.TargetsTag, board string) {
    defaultTarget := &types.TargetTag{}

    if target, ok := targetsTag.Targets[targetsTag.Default_target]; ok {
        defaultTarget.Board = target.Board
        defaultTarget.Compile_flags = target.Compile_flags
        targetsTag.Targets[targetsTag.Default_target] = defaultTarget
    } else {
        defaultTarget.Board = board
        targetsTag.Targets[targetsTag.Default_target] = defaultTarget
    }
}

/// This method prints next steps for any type of create/update command. This will help user
/// decide what they can do next
func (create Create) postPrint() {
    log.Norm.Green(true, "success")
    log.Norm.Yellow(true, "Project has been successfully created/updated and initialized!!")
    log.Norm.Yellow(true, "Check following commands:")

    log.Norm.Cyan(true, "`wio build -h`")
    log.Norm.Cyan(true, "`wio run -h`")
    log.Norm.Cyan(true, "`wio upload -h`")

    if create.Context.Bool("tests") {
        log.Norm.Cyan(true, "`wio test -h`")
    }
}

/// This method is a crucial peace of check to make sure people do not lose their work. It makes
/// sure that if people are creating the project when there are files in the folder, they mean it
/// and not doing it by mistake. It will warn them to update instead if they want
func (create Create) createCheck(directory string) (bool) {
    if !utils.PathExists(directory) {
        return true
    }

    if status, err := utils.IsEmpty(directory); err != nil {
        commands.RecordError(err, "")
        return false
    } else if status {
        return true
    } else {
        message := `The directory is not empty!!
This action will erase everything and will create a new project.
An alternative is to do: wio update <app type> DIRECTORY
Please type y/yes to indicate creation and anything else to indicate abortion: `
        log.Norm.Cyan(false, message)
        reader := bufio.NewReader(os.Stdin)
        text, err := reader.ReadString('\n')
        commands.RecordError(err, "")

        text = strings.TrimSuffix(strings.ToLower(text), "\n")

        if text == "y" || text == "yes" {
            log.Norm.Write(true, "")
            return true
        } else {
            return false
        }
    }
}

/// This method checks if an update can be performed. This basically checks for compatibility issues
/// and configurations to make sure the update can be performed
func (create Create) checkUpdate(directory string) (bool) {
    wioPath := directory + io.Sep + "wio.yml"

    if status, err := utils.IsEmpty(directory); err != nil {
        // if there is an error
        commands.RecordError(err, "")
        return false
    } else if status {
        // if the folder is empty
        return false
    } else if !utils.PathExists(wioPath) {
        // if wio.yml file does not exist
        return false
    } else {
        if create.Type == APP {
            config := &types.AppConfig{}
            if err := io.NormalIO.ParseYml(wioPath, config); err != nil {
                // can't parse wio.yml file for app type
                return false
            } else if config.MainTag.Name == "" {
                // type of the project is wrong compared to the one specified
                return false
            } else {
                // all checks passed
                return true
            }
        } else {
            config := &types.PkgConfig{}
            if err := io.NormalIO.ParseYml(wioPath, config); err != nil {
                // can't parse wio.yml file for lib type
                return false
            } else if config.MainTag.Name == "" {
                // type of the project is wrong compared to the one specified
                return false
            } else {
                // all checks passed
                return true
            }
        }
    }
}

/// This function copies the templates needed to set up the project. It uses parsePathsAndCopy
/// function to parse the paths.json file and then get files based on the project type.
func copyTemplates(projectPath string, appType string, ide string, jsonPath string) (error) {
    strArray := make([]string, 0)
    strArray = append(strArray, appType + "-gen")

    if ide == "clion" {
        strArray = append(strArray, appType+"-clion")
    }

    if err := parsePathsAndCopy(jsonPath, projectPath, strArray); err != nil {
        return err
    }

    return nil
}

/// This function Parses the paths.json file and uses that to get the required files to set
/// up the wio project. This also decides the override based on the paths.json file
func parsePathsAndCopy(jsonPath string, projectPath string, tags []string) (error) {
    var paths = Paths{}
    if err := io.AssetIO.ParseJson(jsonPath, &paths); err != nil {
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

    return io.AssetIO.CopyMultipleFiles(sources, destinations, overrides)
}
