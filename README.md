# Maven Aggregator POM Generator 

A simple tool to generate an aggregate Maven POM file from independent projects in sub-folders.

## Overview

```
$ aggregate-pom-gen --help
Generate an aggregate POM file for direct sub-folders

Usage:
  aggregate-pom-gen [flags]

Flags:
      --artifactId string   The artifactId (default "aggregator")
      --groupId string      The groupId (default "org.acme")
  -h, --help                help for aggregate-pom-gen
      --version string      The version (default "999-SNAPSHOT")
```

It will write a `pom.xml` file in the current folder, where each referenced module corresponds to a direct sub-folder with a `pom.xml` file.

## Installation

```
go install github.com/jponge/aggregate-pom-gen
```
