package libs

import (
    "wio/cmd/wio/utils/types"
    "wio/cmd/wio/utils/io"
    "wio/cmd/wio/utils"
    "path/filepath"
    "io/ioutil"
    "os"
    "regexp"
    "strings"
)

type DependencyTree struct {
    Config types.LibraryLockTag
    Child    []*DependencyTree
}

// Parses Library based on the path provides and its dependencies and it modifies structures to store the data
func ParseLibrary(libPath string, depTree *DependencyTree, depConfig types.LibrariesLockConfig, hash string) (error) {
    configFile := libPath + io.Sep + "wio.yml"
    libraryLibPath :=  libPath + io.Sep + "lib"

    var config types.LibConfig
    var libName string

    depLock := types.LibraryLockTag{}

    // check if config.yml file exists
    if utils.PathExists(configFile) {
        if err := io.NormalIO.ParseYml(configFile, &config); err != nil {
            return err
        }

        libName = config.MainTag.Name
        depLock.Compile_flags = config.MainTag.Compile_flags
    } else {
        libName = filepath.Base(libPath)
    }
    hash += libName
    depLock.Hash = hash
    depLock.Path = libPath

    depConfig.Libraries[libName] = &depLock
    depTree.Config = depLock
    depTree.Child = make([]*DependencyTree, 0)

    if !utils.PathExists(libraryLibPath) {
        return nil
    }

    if dirs, err := ioutil.ReadDir(libraryLibPath); err != nil {
        return err
    } else if len(dirs) > 0 {
        for dir := range dirs {
            if dirs[dir].Name()[0] == '.' {
                continue
            }

            currTree := DependencyTree{}

            if err := ParseLibrary(libraryLibPath + io.Sep + dirs[dir].Name(), &currTree, depConfig, hash);
            err != nil {
                return err
            }
            depTree.Child = append(depTree.Child, &currTree)
        }
    }

    return nil
}

// Parses all the libraries recursively with the help of ParseLibrary function
func ParseLibs(librariesPath string, depTrees []*DependencyTree, depLock types.LibrariesLockConfig) ([]*DependencyTree, error){
    if !utils.PathExists(librariesPath) {
        return depTrees, nil
    }

    if dirs, err := ioutil.ReadDir(librariesPath); err != nil {
        return depTrees, err
    } else if len(dirs) > 0 {
        for dir := range dirs {
            currTree := DependencyTree{}

            // ignore hidden files
            if dirs[dir].Name()[0] == '.' {
                continue
            }

            if err := ParseLibrary(librariesPath + io.Sep + dirs[dir].Name(), &currTree, depLock, "");
            err != nil {
                return depTrees, err
            }

            depTrees = append(depTrees, &currTree)
        }
    }

    return depTrees, nil
}

// Parses all the libraries and their dependencies and creates libs.lock file
func CreateLibsLockFile(projectPath string) ([]*DependencyTree, error) {
    librariesPath := projectPath + io.Sep + "lib"
    wioPath := projectPath + io.Sep + ".wio"

    dependencyTrees := make([]*DependencyTree, 0)
    dependencyLock := types.LibrariesLockConfig{}
    dependencyLock.Libraries = types.LibrariesLockTag{}

    // parse all the libraries and their dependencies in lib folder
    dependencyTrees, err := ParseLibs(librariesPath, dependencyTrees, dependencyLock)
    if err != nil {
        return nil, err
    }

    // write lock file
    if err = io.NormalIO.WriteYml(wioPath + io.Sep + "libraries.lock", &dependencyLock); err != nil {
        return nil, err
    }

    // return Dependency Tree
    return dependencyTrees, nil
}


// Chooses include path based on the the template we support
func ChooseIncludePath(libraryPath string) (string) {
    return libraryPath + "/include"
}

// Create CMakeLists file for each target
func CreateCMakes(exePath string, projectPath string, board string, target string, depTree *DependencyTree) (error) {
    wioLibPath := projectPath + io.Sep + ".wio" + io.Sep + "targets" + io.Sep + target + io.Sep + "libraries"

    if err := os.MkdirAll(wioLibPath + io.Sep + depTree.Config.Hash, os.ModePerm); err != nil {
        return err
    }

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

    for i := 0; i < len(depTree.Child); i++ {
        err := CreateCMakes(exePath, projectPath, board, target, depTree.Child[i])
        if err != nil {
            return nil
        }
    }

    var wioPathRe *regexp.Regexp
    var libraryNameRe *regexp.Regexp
    var globNameRe *regexp.Regexp
    var globStringRe *regexp.Regexp
    var compileFLagsRe *regexp.Regexp
    var libPathRe *regexp.Regexp
    var err error

    if wioPathRe, err = regexp.Compile(`{{WIO_PATH}}`); err != nil {
        return err
    }
    if libraryNameRe, err = regexp.Compile(`{{library-name}}`); err != nil {
        return err
    }
    if globNameRe, err = regexp.Compile(`{{glob-name}}`); err != nil {
        return err
    }
    if globStringRe, err = regexp.Compile(`{{glob-string}}`); err != nil {
        return err
    }
    if compileFLagsRe, err = regexp.Compile(`{{compile-flags}}`); err != nil {
        return err
    }
    if libPathRe, err = regexp.Compile(`{{lib-path}}`); err != nil {
        return err
    }

    librariesTemplate = wioPathRe.ReplaceAllString(librariesTemplate, exePath)
    librariesTemplate = libraryNameRe.ReplaceAllString(librariesTemplate, depTree.Config.Hash)
    librariesTemplate = globNameRe.ReplaceAllString(librariesTemplate,
        strings.ToUpper(depTree.Config.Hash+"_SRC_FILES"))

    srcPath := depTree.Config.Path + io.Sep + "src"
    if utils.PathExists(srcPath) {
        librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, srcPath+"/*.cpp "+srcPath+"/*.cc "+srcPath+"/*.c")
    } else {
        librariesTemplate = globStringRe.ReplaceAllString(librariesTemplate, depTree.Config.Path+"/*.cpp "+depTree.Config.Path+"/*.cc "+depTree.Config.Path+"/*.c")
    }

    librariesTemplate = compileFLagsRe.ReplaceAllString(librariesTemplate, strings.Join(depTree.Config.Compile_flags, " "))
    librariesTemplate = libPathRe.ReplaceAllString(librariesTemplate, ChooseIncludePath(depTree.Config.Path))

    for dep := range depTree.Child {
        librariesTemplate += "include_directories(" + ChooseIncludePath(depTree.Child[dep].Config.Path) + ")\n"
        librariesTemplate += "target_link_libraries(" + depTree.Config.Hash + " ${CMAKE_SOURCE_DIR}/../" +
            depTree.Child[dep].Config.Hash + io.Sep + "lib" + depTree.Child[dep].Config.Hash + ".a)\n"
    }

    io.NormalIO.WriteFile(wioLibPath+io.Sep+depTree.Config.Hash+io.Sep+"CMakeLists.txt", []byte(librariesTemplate))

    return nil
}

func PopulateCMakeFilesforLibs(projectPath string, board string, target string) (error) {
   dependencyTree, err := CreateLibsLockFile(projectPath)
   exePath, err := io.NormalIO.GetRoot()

    if err != nil {
        return err
    }
   if err != nil {
       return err
   }

   for tree := range dependencyTree {
       if err = CreateCMakes(exePath, projectPath, board, target, dependencyTree[tree]); err != nil {
           return err
       }
   }

   return nil
}
