// +build mage

package main

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	cmd := exec.Command("go", "build", "-o", "kpkg", ".")
	return cmd.Run()
}

// A install step
func Install() error {
	mg.Deps(InstallDeps)
	fmt.Println("Installing...")
	cmd := exec.Command("go", "install", "./...")
	return cmd.Run()
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "list", "./...")
	if err := cmd.Run(); err != nil {
		return err
	}
	cmd = exec.Command("go", "mod", "tidy")
	return cmd.Run()
}

// Clean up after yourself
func Clean() error {
	fmt.Println("Cleaning...")
	cmd := exec.Command("rm", "-rf", "./kpkg")
	return cmd.Run()
}
