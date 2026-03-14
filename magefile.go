//go:build mage

package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Dev runs the development server.
func Dev() error {
	return sh.RunV("go", "run", "./cmd/server")
}

// Build compiles server and seed binaries to bin/.
func Build() error {
	if err := sh.RunV("go", "build", "-o", "bin/server", "./cmd/server"); err != nil {
		return err
	}
	return sh.RunV("go", "build", "-o", "bin/seed", "./cmd/seed")
}

// Test runs domain tests.
func Test() error {
	return sh.RunV("go", "test", "./internal/domain/...", "-v")
}

// Sqlc generates repository code from SQL queries.
func Sqlc() error {
	return sh.RunV("sqlc", "generate")
}

// Clean removes built binaries.
func Clean() error {
	for _, f := range []string{"bin/server", "bin/seed"} {
		os.Remove(f)
	}
	return nil
}

// Deps downloads and tidies Go module dependencies.
func Deps() error {
	if err := sh.RunV("go", "mod", "download"); err != nil {
		return err
	}
	return sh.RunV("go", "mod", "tidy")
}

type Seed mg.Namespace

// Create creates a new user. Pass EMAIL, PASSWORD, and NAME env vars.
func (Seed) Create() error {
	email := os.Getenv("EMAIL")
	password := os.Getenv("PASSWORD")
	name := os.Getenv("NAME")
	if email == "" || password == "" || name == "" {
		return fmt.Errorf("usage: EMAIL=x PASSWORD=x NAME=x mage seed:create")
	}
	return sh.RunV("go", "run", "./cmd/seed", "create", "--email", email, "--password", password, "--name", name)
}

// List lists all users.
func (Seed) List() error {
	return sh.RunV("go", "run", "./cmd/seed", "list")
}

type UI mg.Namespace

// Install installs npm dependencies for the Svelte estimate builder.
func (UI) Install() error {
	cmd := exec.Command("npm", "install")
	cmd.Dir = "ui"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Build builds the Svelte estimate builder bundle.
func (UI) Build() error {
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = "ui"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Watch runs the Svelte estimate builder in watch mode.
func (UI) Watch() error {
	cmd := exec.Command("npm", "run", "watch")
	cmd.Dir = "ui"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
