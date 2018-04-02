package libs

import (
    "wio/cmd/wio/utils/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
    "errors"
    "strings"
    "gopkg.in/yaml.v2"
    "os"
    "regexp"
)

// Parses libs.yml file for the project and if that file does not exist, it throws an
// error explaining what is wrong
func ParseLibs(projectDirectory string) (*types.LibsConfig, error) {
    libsConfigPath := projectDirectory + io.Sep + "libs.yml"

    if !utils.PathExists(libsConfigPath) {
        return nil, errors.New("libs.yml does not exist, run wio create to update the project")
    }

    libsConfig := new(types.LibsConfig)
    if err := io.NormalIO.ParseYml(libsConfigPath, libsConfig); err != nil {
        return nil, err
    }

    return libsConfig, nil
}

// Writes libs configuration to libs.yml file. It includes help text and nicely formats the file
func WriteLibs(projectDirectory string, config *types.LibsConfig) (error) {
    libsConfigPath := projectDirectory + io.Sep + "libs.yml"
    libsHelpText, err := io.AssetIO.ReadFile("templates" + io.Sep + "config" + io.Sep + "libs-help.txt")
    if err != nil {
        return err
    }

    libsHelpTextSlice := strings.Split(string(libsHelpText), "\n\n")

    configString, err := yaml.Marshal(config)
    if err != nil {
        return err
    }
    configStringSlice := strings.Split(string(configString), "\n")
    finalString := libsHelpTextSlice[0]
    last := false

    for index := 0; index < len(configStringSlice); index++ {
        line := configStringSlice[index]

        if len(line) >= 2 {
            checker := strings.Trim(line, " ")
            length := len(checker)
            if !last && length >= 1 && checker[length-1:length] == ":" {
                last = true
                finalString += "\n"
            } else {
                last = false
            }

            if strings.Contains(checker,"dependencies:") {
                finalString += libsHelpTextSlice[1] + "\n"
            }
        }

        finalString += line + "\n"
    }

    finalString = strings.Trim(finalString, "\n") + "\n"
    return io.NormalIO.WriteFile(libsConfigPath, []byte(finalString))
}


func getDependecnyForLibrary(libraryName string, config *types.LibsConfig) {
    dependencies := make(map[string]*types.DependencyTag)

    for dependency := range config.DependenciesTag {
        if strings.Contains(dependency, libraryName + "/") {
            dependencies[dependency] = config.DependenciesTag[dependency]
        }
    }
}

type Library struct {
    Name string
    HashName string    /* Name to make sure this library is unique */
    Path string
    CompileFlags []string
    Libraries []*Library
}


func ChooseIncludePath(libraryPath string) (string) {
    if utils.PathExists(libraryPath + io.Sep + "include") {
        return libraryPath + "/include"
    } else if utils.PathExists(libraryPath + io.Sep + "src") {
        return libraryPath + "/src"
    } else {
        return libraryPath
    }
}


func CreateLibBuild(projectLibPath string, wioLibPath string, wioExecutablePath string,
    board string, hash string, library *Library) (*Library, error) {

    hashName := hash + library.Name
    if err := os.MkdirAll(wioLibPath + io.Sep + hashName, os.ModePerm); err != nil {
        return nil, err
    }
    library.HashName = hashName

    librariesTemplate := `# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "{{WIO_PATH}}/toolchain/cmake/CosaToolchain.cmake")
cmake_minimum_required(VERSION 3.0.0)
project({{library-name}} C CXX ASM)

file(GLOB_RECURSE {{glob-name}} {{glob-string}})
generate_arduino_library({{library-name}}
	SRCS ${{{glob-name}}}
	BOARD uno)
target_compile_definitions({{library-name}} PRIVATE __AVR_Cosa__ {{compile-flags}})
include_directories({{lib-path}})
`

    dependencies := make([]*Library, 0)
    for i := 0; i < len(library.Libraries); i++  {
        lastElem, err := CreateLibBuild(projectLibPath, wioLibPath, wioExecutablePath, board, hash + library.Name,
            library.Libraries[i])
        if err != nil {
            return nil, nil
        }
        dependencies = append(dependencies, lastElem)
    }

    var wioPathRe *regexp.Regexp
    var libraryNameRe *regexp.Regexp
    var globNameRe *regexp.Regexp
    var globStringRe *regexp.Regexp
    var compileFLagsRe *regexp.Regexp
    var libPathRe *regexp.Regexp
    var err error

    if wioPathRe, err = regexp.Compile(`{{WIO_PATH}}`); err != nil { return nil, err }
    if libraryNameRe, err = regexp.Compile(`{{library-name}}`); err != nil { return nil, err }
    if globNameRe, err = regexp.Compile(`{{glob-name}}`); err != nil { return nil, err }
    if globStringRe, err = regexp.Compile(`{{glob-string}}`); err != nil { return nil, err }
    if compileFLagsRe, err = regexp.Compile(`{{compile-flags}}`); err != nil { return nil, err }
    if libPathRe, err = regexp.Compile(`{{lib-path}}`); err != nil { return nil, err }

    librariesTemplate = wioPathRe.ReplaceAllString(librariesTemplate, wioExecutablePath)
    librariesTemplate =  libraryNameRe.ReplaceAllString(librariesTemplate, hashName)
    librariesTemplate = globNameRe.ReplaceAllString(librariesTemplate, strings.ToUpper(hashName + "_SRC_FILES"))

    srcPath := library.Path + io.Sep + "src"
    if utils.PathExists(srcPath) {
        librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, srcPath + "/*.cpp " + srcPath + "/*.cc " + srcPath + "/*.c")
    } else {
        librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, library.Path + "/*.cpp " + library.Path + "/*.cc " + library.Path + "/*.c")
    }

    librariesTemplate = compileFLagsRe.ReplaceAllString(librariesTemplate, strings.Join(library.CompileFlags, " "))
    librariesTemplate = libPathRe.ReplaceAllString(librariesTemplate, ChooseIncludePath(library.Path))

    for dep := range dependencies {
        librariesTemplate += "include_directories(" + ChooseIncludePath(dependencies[dep].Path) + ")\n"
        librariesTemplate += "target_link_libraries(" + library.HashName + " ${CMAKE_SOURCE_DIR}/../" +
            dependencies[dep].HashName + io.Sep + "lib" + dependencies[dep].HashName + ".a)\n"
    }

    io.NormalIO.WriteFile(wioLibPath + io.Sep + hashName + io.Sep + "CMakeLists.txt", []byte(librariesTemplate))

    f, err := os.OpenFile(wioLibPath + io.Sep + ".." + io.Sep + "libraries.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    f.Write([]byte(library.HashName + "\n"))
    f.Close()

    return library, nil
}




func CreateCMakeBuild(projectPath string, target string, board string) (error) {
    projectLibPath := projectPath + io.Sep + "lib"
    wioLibPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + "libraries"
    wioExecutablePath, err := io.NormalIO.GetRoot()
    if err != nil {
        return err
    }

    twoOne := &Library{Name:"Lib", Libraries:nil, Path: "/Users/deep/Development/gowork/src/wio/hello/lib/Brother/lib/Lib"}
    twoTwo := &Library{Name:"Lizard", Libraries:nil, Path: "/Users/deep/Development/gowork/src/wio/hello/lib/Brother/lib/Lizard"}
    one := &Library{Name:"Brother", Libraries:[]*Library{twoOne, twoTwo}, Path: "/Users/deep/Development/gowork/src/wio/hello/lib/Brother"}

    CreateLibBuild(projectLibPath, wioLibPath, wioExecutablePath, board, "", one)

    return nil
}






/*

func PostOrder(root *Node, hash string, board string, libPath string, wioPath string) (*Node, error) {
    if root.Nodes == nil {
        hashName := hash + root.Name
        root.HashName = hashName
        os.Mkdir(wioPath + io.Sep + hashName, os.ModePerm)

        librariesTemplate := `# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "{{WIO_PATH}}/toolchain/cmake/CosaToolchain.cmake")
cmake_minimum_required(VERSION 3.0.0)
project({{library-name}} C CXX ASM)

file(GLOB_RECURSE {{glob-name}} {{glob-string}})
generate_arduino_library({{library-name}}
	SRCS ${{{glob-name}}}
	BOARD uno)
target_compile_definitions({{library-name}} PRIVATE __AVR_Cosa__ )
include_directories({{lib-path}}/include)
`

        var wioPathRe *regexp.Regexp
        var libraryNameRe *regexp.Regexp
        var globNameRe *regexp.Regexp
        var globStringRe *regexp.Regexp
        var libPathRe *regexp.Regexp
        var err error

        if wioPathRe, err = regexp.Compile(`{{WIO_PATH}}`); err != nil { return nil, err }
        if libraryNameRe, err = regexp.Compile(`{{library-name}}`); err != nil { return nil, err }
        if globNameRe, err = regexp.Compile(`{{glob-name}}`); err != nil { return nil, err }
        if globStringRe, err = regexp.Compile(`{{glob-string}}`); err != nil { return nil, err }
        if libPathRe, err = regexp.Compile(`{{lib-path}}`); err != nil { return nil, err }


        librariesTemplate = wioPathRe.ReplaceAllString(librariesTemplate, wioPath)
        librariesTemplate =  libraryNameRe.ReplaceAllString(librariesTemplate, hashName)
        librariesTemplate = globNameRe.ReplaceAllString(librariesTemplate, strings.ToUpper(hashName + "_SRC_FILES"))

        srcPath := root.Path + io.Sep + "src"
        librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, srcPath + "/*.cpp " + srcPath + "/*.cc " + srcPath + "/*.c")
        librariesTemplate = libPathRe.ReplaceAllString(librariesTemplate, root.Path )

        io.NormalIO.WriteFile(wioPath + io.Sep + hashName + io.Sep + "CMakeLists.txt", []byte(librariesTemplate))

        //fmt.Println(wioPath + io.Sep + hashName + io.Sep + "CMakeLists.txt")
        return root, nil
    }

    deps := make([]*Node, 0)
    for i := 0; i < len(root.Nodes); i++  {
        lastElem, err := PostOrder(root.Nodes[i], hash + root.Name, board, libPath, wioPath)
        if err != nil {
            return nil, nil
        }
        deps = append(deps, lastElem)
    }

    os.Mkdir(libPath + io.Sep + hash + root.Name, os.ModePerm)

    hashName := hash + root.Name
    root.HashName = hashName
    os.Mkdir(wioPath + io.Sep + hashName, os.ModePerm)

    librariesTemplate := `# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "{{WIO_PATH}}/toolchain/cmake/CosaToolchain.cmake")
cmake_minimum_required(VERSION 3.0.0)
project({{library-name}} C CXX ASM)

file(GLOB_RECURSE {{glob-name}} {{glob-string}})
generate_arduino_library({{library-name}}
	SRCS ${{{glob-name}}}
	BOARD uno)
target_compile_definitions({{library-name}} PRIVATE __AVR_Cosa__ )
include_directories({{lib-path}}/include)
`

    var wioPathRe *regexp.Regexp
    var libraryNameRe *regexp.Regexp
    var globNameRe *regexp.Regexp
    var globStringRe *regexp.Regexp
    var libPathRe *regexp.Regexp
    var err error

    if wioPathRe, err = regexp.Compile(`{{WIO_PATH}}`); err != nil { return nil, err }
    if libraryNameRe, err = regexp.Compile(`{{library-name}}`); err != nil { return nil, err }
    if globNameRe, err = regexp.Compile(`{{glob-name}}`); err != nil { return nil, err }
    if globStringRe, err = regexp.Compile(`{{glob-string}}`); err != nil { return nil, err }
    if libPathRe, err = regexp.Compile(`{{lib-path}}`); err != nil { return nil, err }


    librariesTemplate = wioPathRe.ReplaceAllString(librariesTemplate, wioPath)
    librariesTemplate =  libraryNameRe.ReplaceAllString(librariesTemplate, hashName)
    librariesTemplate = globNameRe.ReplaceAllString(librariesTemplate, strings.ToUpper(hashName + "_SRC_FILES"))

    srcPath := root.Path + io.Sep + "src"
    librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, srcPath + "/*.cpp " + srcPath + "/*.cc " + srcPath + "/*.c")
    librariesTemplate = libPathRe.ReplaceAllString(librariesTemplate, root.Path)


    for dep := range deps {
        librariesTemplate += "\ninclude_directories(" + deps[dep].Path + "/include)"
    }

    librariesTemplate += "\n"

    for dep := range deps {
        librariesTemplate += "target_link_libraries(" + root.HashName + " {CMAKE_SOURCE_DIR}/../" + deps[dep].HashName + io.Sep + deps[dep].HashName + ".a)"
    }


    io.NormalIO.WriteFile(wioPath + io.Sep + hashName + io.Sep + "CMakeLists.txt", []byte(librariesTemplate))


    return root, nil
}*/

func CreateCMakeLibraries(projectPath string, target string, board string /*onfig *types.LibsConfig*/) (error) {
    /*
    twoOne := &Node{Name:"Lib", Nodes:nil, Path: "/Users/deep/Development/gowork/src/wio/hello/lib/Brother/lib/Lib"}
    twoTwo := &Node{Name:"Lizard", Nodes:nil, Path: "/Users/deep/Development/gowork/src/wio/hello/lib/Brother/lib/Lizard"}
    one := &Node{Name:"Brother", Nodes:[]*Node{twoOne, twoTwo}, Path: "/Users/deep/Development/gowork/src/wio/hello/lib/Brother"}

    targetPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target
    librariesPath := targetPath + io.Sep + "libraries"

    // create directory structure
    if err := os.MkdirAll(librariesPath, os.ModePerm); err != nil {
        return err
    }

    PostOrder(one, "", board, projectPath + io.Sep + "lib", librariesPath)
    */

    CreateCMakeBuild(projectPath, target, board)



    /*
    targetPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target
    librariesPath := targetPath + io.Sep + "libraries"
    cmakePath := librariesPath + io.Sep + "CMakeLists.txt"

    librariesTemplate := `# Cosa Toolchain
set(CMAKE_TOOLCHAIN_FILE "${WIO_PATH}/toolchain/cmake/CosaToolchain.cmake")
cmake_minimum_required(VERSION 3.0.0)
project({{library-name}} C CXX ASM)

{{libraries-generation}}
`
    // create directory structure
    if err := os.MkdirAll(librariesPath, os.ModePerm); err != nil {
        return err
    }

    var wioPath *regexp.Regexp
    var libraryName *regexp.Regexp
    var librariesGeneration *regexp.Regexp
    var err error

    if wioPath, err = regexp.Compile(`{{WIO_PATH}}`); err != nil { return err }
    if libraryName, err = regexp.Compile(`{{library-name}}`); err != nil { return err }
    if librariesGeneration, err = regexp.Compile(`{{libraries-generation}}`); err != nil { return err }


    // generate libraries cmake text
    for library := range config.LibrariesTag {
        srcPath := config.LibrariesTag[library].Path + io.Sep + "src"
        str := "(GLOB_RECURSE" + strings.ToUpper(library + "SRC_FILES ") + srcPath + "/*.cpp " + srcPath + "/*.cc " + srcPath + "/*.c)"
        str += "include_directories(" + config.LibrariesTag[library].Path + io.Sep + "include\")\n"
        str += "generate_arduino_library(" + library + "\n"
        str += "\tSRCS ${" + strings.ToUpper(library + "SRC_FILES ") + "}"

        if len(library.libs) > 0 {
            str += "\tLIBS "
            for i := 0; i < len(library.libs); i++ {
                str += library.libs[i].Name + " "
            }
            str += "\n"
        }
    }




    for lib := range config.LibrariesTag {
        fmt.Println("SRC: " + config.LibrariesTag[lib].Path + "/src")
        fmt.Println("INCLUDE: " + config.LibrariesTag[lib].Path + "/include")
    }
    */


    /*
        Project Structure inside .wio
        targets:
            main:
                 CMakeLists.txt -> builds the firmware and links the libraries
                 Libraries:
                    CMakeLists.txt -> builds the library that user included
                                   -> generate arduino library for all the dependencies

        Make CMakeLists.txt -> allows you to switch from one target to another
     */


    // src files from src folder
    // include_directory files from include folder

    return nil
}


/*
    * Check if libs.yml exist and if yes:
        * parse the yml and get the structure
        * based on the structure, download these packages and their dependencies
        * if these packages exist, remove them and use the downloaded ones
        * get all the packages
        * parse all the libs in lib folder
        * merge that with
    *



    packages will only be updated and gathered using wio packager
    wio create assume all the packages are valid and libs.yml is valid and if not valid throw error


 */
