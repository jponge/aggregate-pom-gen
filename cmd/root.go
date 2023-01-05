package cmd

import (
	"html/template"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.Flags().String("groupId", "org.acme", "The groupId")
	rootCmd.Flags().String("artifactId", "aggregator", "The artifactId")
	rootCmd.Flags().String("version", "999-SNAPSHOT", "The version")
}

var rootCmd = &cobra.Command{
	Use:   "aggregate-pom-gen",
	Short: "Generate an aggregate POM file for direct sub-folders",
	Run:   doRun,
}

const pomTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" 
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" 
         xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/maven-v4_0_0.xsd">

	<groupId>{{.GroupId}}</groupId>
  	<artifactId>{{.ArtifactId}}</artifactId>
  	<version>{{.Version}}</version>
  	<packaging>pom</packaging>
  	<modelVersion>4.0.0</modelVersion>

	<modules>
		{{range .Modules}}
	  		<module>{{.}}</module>
		{{end}}
  	</modules>
</project>
`

type templateParams struct {
	GroupId    string
	ArtifactId string
	Version    string
	Modules    []string
}

func doRun(cmd *cobra.Command, args []string) {
	current, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	defer current.Close()

	var submodules []string

	sub, err := current.ReadDir(0)
	if err != nil {
		log.Fatal(err)
	}
	for _, folder := range sub {
		if folder.IsDir() {
			dir, err := os.Open(folder.Name())
			if err != nil {
				log.Fatal(err)
			}
			files, err := dir.Readdir(0)
			if err != nil {
				log.Fatal(err)
			}
			for _, file := range files {
				if !file.IsDir() && file.Name() == "pom.xml" {
					submodules = append(submodules, folder.Name())
				}
			}
		}
	}

	params := templateParams{
		GroupId:    cmd.Flags().Lookup("groupId").Value.String(),
		ArtifactId: cmd.Flags().Lookup("artifactId").Value.String(),
		Version:    cmd.Flags().Lookup("version").Value.String(),
		Modules:    submodules,
	}

	pomFile, err := os.OpenFile("pom.xml", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := template.New("pom").Parse(pomTemplate)
	if err != nil {
		log.Fatal(err)
	}
	err = tmpl.Execute(pomFile, params)
	if err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
