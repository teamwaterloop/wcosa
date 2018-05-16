package utils

import (
    "os"
    "io"
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

// This checks if the directory is empty or not
func IsEmpty(name string) (bool, error) {
    f, err := os.Open(name)
    if err != nil {
        return false, err
    }
    defer f.Close()

    _, err = f.Readdirnames(1) // Or f.Readdir(1)
    if err == io.EOF {
        return true, nil
    }
    return false, err // Either not empty or error, suits both cases
}

// This checks if a string is in the slice
func StringInSlice(a string, list []string) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}

// It takes in a slice and an element and then ut appends that element to the slice only
// if that element in not already in the slice
func AppendIfMissingElem(slice []string, i string) []string {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}

// It takes two slices and appends the second one onto the first one. It does
// not allow duplicates
func AppendIfMissing(slice []string, slice2 []string) []string {
    newSlice := make([]string, 0)

    for _, ele1 := range slice {
        newSlice = AppendIfMissingElem(newSlice, ele1)
    }

    for _, ele2 := range slice2 {
        newSlice = AppendIfMissingElem(newSlice, ele2)
    }

    return newSlice
}

func ToLinixPath(path string) {

}
