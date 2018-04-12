package utils

import (
    "os"
    "path/filepath"
)

// Checks if path exists and returns true and false based on that
func PathExists(path string) (bool) {
    if _, err := os.Stat(path); err != nil {
        return false
    }
    return true
}

// Checks if the give path is a director and based on the returns
// true or false. If path does not exist, it throws an error
func IsDir(path string) (bool, error) {
    fi, err := os.Stat(path)
    if err != nil {
        return false, err
    }

    return fi.IsDir(), nil
}

// Checks if the path contains the extensions provided and it returns true and false
// based on that.If path does not exist, it throws an error
func HasExtension(path string, extensions ...string) (bool, error) {
    if !PathExists(path) { return false, nil }
    for extension := 0; extension < len(extensions) ; extension++ {
        var ext = filepath.Ext(path)
        if extensions[extension] == ext {
            return true, nil
        }
    }

    return false, nil
}
