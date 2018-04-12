package cmake

import (
    "path/filepath"
    . "wio/cmd/wio/utils/io"
    "strings"
    "wio/cmd/wio/utils/types"
)

func CreateMainCMakeListsFile(projectPath string, board string, framework string, target string, flags []string) (error) {
    projectName := filepath.Base(projectPath)
    executablePath, err := NormalIO.GetRoot()

    if err != nil {
        return err
    }

    lockFilePath := projectPath + Sep + ".wio" + Sep + lockFileName
    targetPath := projectPath + Sep + ".wio" + Sep + "targets" + Sep + target

    toolChainPath := executablePath + "/toolchain/cmake/CosaToolchain.cmake"

    // read the CMakeLists.txt file template
    templateData, err := AssetIO.ReadFile("templates/cmake/CMakeLists.txt.tpl")

    if err != nil {
        return err
    }

    templateDataStr := strings.Replace(string(templateData), "{{toolchain-path}}", toolChainPath, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{project-name}}", projectName, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{project-path}}", projectPath, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{target-name}}", target, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{board}}", board, -1)
    templateDataStr = strings.Replace(templateDataStr, "{{framework}}", strings.ToUpper(framework), -1)
    templateDataStr = strings.Replace(templateDataStr, "{{flags}}", strings.Join(flags, " "), -1)

    lockConfig := types.LibrariesLockConfig{}

    if err = NormalIO.ParseYml(lockFilePath, &lockConfig); err != nil {
        return err
    }

    for lib := range lockConfig.Libraries {
        if !strings.Contains(lib, "__") {
            templateDataStr += "include_directories(" + lockConfig.Libraries[lib].Path + "/include" + ")\n"
        }
    }

    templateDataStr += "\n"

    for lib := range lockConfig.Libraries {
        templateDataStr += "target_link_libraries(" + target + " " + targetPath + Sep + "libraries" + Sep + lib + "/" +
            "lib" + lockConfig.Libraries[lib].Hash + ".a" + ")\n"
    }

    return NormalIO.WriteFile(projectPath + Sep + ".wio" + Sep + "CMakeLists.txt", []byte(templateDataStr))
}
