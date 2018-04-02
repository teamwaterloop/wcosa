// Copyright 2018 Waterloop. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.


// Package type contains types for use by other packages
// This file contains all the types that are used throughout the application

package types

type CliArgs struct {
    AppType     string
    Directory   string
    Board       string
    Framework   string
    Platform    string
    Ide         string
    Tests       bool
}

// type for the targets tag in the configuration file
type TargetsTag map[string]*TargetTag

// Structure to handle individual target inside targets
type TargetTag struct {
    Default string
    Board string
    Compile_flags []string
}

// Structure to hold information about project type: app
type AppTag struct {
    Name string
    Platform string
    Framework string
    Ide string
}

// Structure to hold information about project type: lib
type LibTag struct {
    Name string
    Version string
    Authors []string
    License []string
    Platform string
    Framework []string
    Board []string
    Compile_flags []string
    Ide string
}

type AppConfig struct {
    MainTag AppTag              `yaml:"app"`
    Targets TargetsTag          `yaml:"targets"`
}

type LibConfig struct {
    MainTag LibTag              `yaml:"lib"`
    Targets TargetsTag          `yaml:"targets"`
}

// Structure to handle individual library inside libraries
type LibraryTag struct {
    Url           string
    Version       string
    Path          string
    Compile_flags []string
}

// Structure to handle individual dependency inside dependencies
type DependencyTag struct {
    Name          string
    Url           string
    Version       string
    Compile_flags []string
}

// type for the libraries tag in the libs.yml file
type LibrariesTag map[string]*LibraryTag

// type for the dependencies tag in the libs.yml file
type DependenciesTag map[string]*DependencyTag

// type for whole libs.yml file
type LibsConfig struct {
    LibrariesTag    LibrariesTag    `yaml:"libraries"`
    DependenciesTag DependenciesTag `yaml:"dependencies"`
}
