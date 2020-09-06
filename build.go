// +build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	// mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	fmt.Println("Building...")

	files := getListOfProjectDirs()

	for _, file := range files {
		err := func() error {
			os.Chdir(file)
			defer os.Chdir("..")

			return sh.RunV("docker", "build", "-t", "example/"+file, ".")
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

func Load() error {
	mg.Deps(Build)
	fmt.Println("Loading images...")

	files := getListOfProjectDirs()

	for _, file := range files {
		func() {
			imageName := "example/" + file
			fmt.Println("loading...", imageName)
			sh.RunV("kind", "load", "docker-image", imageName)
		}()
	}

	return nil
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Load)
	fmt.Println("Installing...")

	files := getYamlFiles()
	for _, file := range files {
		fmt.Println("installing:", file)
		sh.RunV("kubectl", "apply", "-f", file)
	}
	return nil
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("MyApp")
}

func getListOfProjectDirs() []string {
	files := make([]string, 0)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == "." {
				return nil
			}
			if strings.HasPrefix(info.Name(), "cmd-") {
				files = append(files, info.Name())
			}
			return filepath.SkipDir
		}
		return nil
	})
	return files
}

func getYamlFiles() []string {
	files := make([]string, 0)
	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if strings.HasSuffix(info.Name(), ".yaml") {
				files = append(files, info.Name())
			}
		}
		return nil
	})
	return files
}
