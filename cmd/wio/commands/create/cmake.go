// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Part of commands/create package, which contains create command and sub commands provided by the tool.
// This contains helper function for cmake template parsing
package create

import (
    . "wio/cmd/wio/utils/io"
    "os"
    "io/ioutil"
    "path/filepath"
    . "wio/cmd/wio/utils"
    "wio/cmd/wio/utils/types"
    "regexp"
    "strings"
)

func GetLibrariesSlice(projectPath string) ([]*CMakeLibrary, error) {
    libraries := make([]*CMakeLibrary, 0)

    libPath := projectPath + Sep + "lib"
    if !PathExists(projectPath) {
        return nil, nil
    }

    // read all the files in the directory
    files, err := ioutil.ReadDir(libPath)
    if err != nil {
        return nil, err
    }

    // we need to create a library for each file
    for _, file := range files {
        libraryPath := libPath + Sep + file.Name()

        isDirectory, err := IsDir(libraryPath)
        if err != nil {
            return nil, err
        }

        if isDirectory {
            library := &CMakeLibrary{}
            library, err = libSearch(libPath+Sep+file.Name(), library)
            if err != nil {
                return nil, err
            }

            libraries = append(libraries, library)
        }
    }

    return libraries, nil
}

func libSearch(libraryPath string, library *CMakeLibrary) (*CMakeLibrary, error) {
    library.Src = make([]string, 0)
    library.Hdr = make([]string, 0)
    library.libs = make([]*CMakeLibrary, 0)
    library.Name = filepath.Base(libraryPath)

    librarySrcPath := libraryPath + Sep + "src"
    libraryLibPath := libraryPath + Sep + "lib"

    var fileList []string
    if PathExists(librarySrcPath) {
        // read all the files in src directory
        err := filepath.Walk(librarySrcPath, func(path string, f os.FileInfo, err error) error {
            isDirectory, err := IsDir(path)
            if err != nil {
                return err
            }

            if !isDirectory {
                fileList = append(fileList, path)
            }
            return nil
        })

        if err != nil {
            return nil, err
        }
        library.SourcePath = librarySrcPath
    } else {
        // read all the files in the root directory
        err := filepath.Walk(libraryPath, func(path string, f os.FileInfo, err error) error {
            isDirectory, err := IsDir(path)
            if err != nil {
                return err
            }

            if !isDirectory {
                fileList = append(fileList, path)
            }
            return nil
        })

        if err != nil {
            return nil, err
        }
        library.SourcePath = libraryPath
    }

    for file := 0; file < len(fileList); file++ {
        if value, _ := HasExtension(fileList[file], ".cpp", ".cc", ".c"); value {
            library.Src = append(library.Src, fileList[file])
        } else if value, _ := HasExtension(fileList[file], ".hh", ".h"); value {
            library.Hdr = append(library.Hdr, fileList[file])
        }
    }

    if PathExists(libraryLibPath) {
        libraries, err := GetLibrariesSlice(libraryPath)
        if err != nil {
            return nil, err
        }
        library.libs = libraries
    }

    return library, nil
}

func makeLibString(library *CMakeLibrary, librariesTag types.LibrariesTag, board string) (string) {
    str := "include_directories(" + library.SourcePath + ")\n"
    str += "generate_arduino_library(" + library.Name + "\n"

    str += "\tSRCS "
    for i := 0; i < len(library.Src); i++ {
        str += library.Src[i] + " "
    }
    str += "\n"

    str += "\tHDRS "
    for i := 0; i < len(library.Hdr); i++ {
        str += library.Hdr[i] + " "
    }
    str += "\n"

    if len(library.libs) > 0 {
        str += "\tLIBS "
        for i := 0; i < len(library.libs); i++ {
            str += library.libs[i].Name + " "
        }
        str += "\n"
    }

    compileFlags := ""

    if val, ok := librariesTag[library.Name]; ok {
        compileFlags += strings.Join(val.Compile_flags, " ")
    }

    str += "\tBOARD " + board + ")\n"
    str += "target_compile_definitions(" + library.Name + " PRIVATE __AVR_Cosa__ " + compileFlags + ")\n"

    for i := 0; i < len(library.libs); i++ {
        //if !strings.Contains(str, "generate_arduino_library(" + library.libs[i].Name + "\n") {
        str2 := makeLibString(library.libs[i], librariesTag, board)
        //}
        if !(strings.Contains(str2, "generate_arduino_library(" + library.libs[i].Name + "\n") &&
            strings.Contains(str, "generate_arduino_library(" + library.libs[i].Name + "\n")){
            str += str2
        }
    }

    return str
}

func GetLibraries(board string, librariesTag types.LibrariesTag, projectPath string) (string, []*CMakeLibrary) {
    libs := make([]*CMakeLibrary, 0)
    libs, _ = GetLibrariesSlice(projectPath)

    str := ""
    for i := 0; i < len(libs); i++ {
        str += makeLibString(libs[i], librariesTag, board)
    }

    return str, libs
}

func WriteCMakeLibraries(libString string, targetName string, projectPath string) (error) {
    return NormalIO.WriteFile(projectPath + Sep + ".wio" + Sep + "targets" + Sep + targetName + "Libs.cmake",
        []byte(libString))
}

func WriteCMakeFramework(libraries []*CMakeLibrary, targetName string, target *types.TargetSubTags, args *types.CliArgs) (error) {
    template :=`file(GLOB_RECURSE SRC_FILES "{{src-path}}/*.cpp" "{{src-path}}/*.cc" "{{src-path}}/*.c")

# create the firmware
generate_arduino_firmware({{project-name}}
	SRCS ${SRC_FILES}
	LIBS {{libs}}
	PORT 
	BOARD {{board}})
target_compile_definitions({{project-name}} PRIVATE __{{platform}}_{{framework}}__ {{compile-flags}})
`
    var srcPathRe *regexp.Regexp
    var projectNameRe *regexp.Regexp
    var libsRe *regexp.Regexp
    var boardRe *regexp.Regexp
    var platformRe *regexp.Regexp
    var frameworkRe *regexp.Regexp
    var compileFlagsRe *regexp.Regexp
    var e error

    if srcPathRe, e = regexp.Compile(`{{src-path}}`); e != nil { return e }
    if projectNameRe, e = regexp.Compile(`{{project-name}}`); e != nil { return e}
    if libsRe, e = regexp.Compile(`{{libs}}`); e != nil { return e }
    if boardRe, e = regexp.Compile(`{{board}}`); e != nil { return e }
    if platformRe, e = regexp.Compile(`{{platform}}`); e != nil { return e }
    if frameworkRe, e = regexp.Compile(`{{framework}}`); e != nil { return e }
    if compileFlagsRe, e = regexp.Compile(`{{compile-flags}}`); e != nil { return e }

    template = srcPathRe.ReplaceAllString(template, args.Directory + "/src")
    template = projectNameRe.ReplaceAllString(template, filepath.Base(args.Directory))

    libs := ""
    for i := 0; i < len(libraries); i++ {
        libs += libraries[i].Name + " "
    }
    libs = strings.Trim(libs, " ")

    template = libsRe.ReplaceAllString(template, libs)
    template = boardRe.ReplaceAllString(template, target.Board)
    template = platformRe.ReplaceAllString(template, strings.ToUpper(args.Platform))
    template = frameworkRe.ReplaceAllString(template, strings.Title(args.Framework))
    template = compileFlagsRe.ReplaceAllString(template, strings.Join(target.Compile_flags, " "))

    return NormalIO.WriteFile(args.Directory + Sep + ".wio" + Sep + "targets" + Sep + targetName + "Framework.cmake",
        []byte(template))
}
