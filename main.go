package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func main() {
	os.Exit(main1())
}

func main1() int {
	if err := main2(); err != nil {
		fmt.Fprintf(os.Stderr, "gomerge: %v\n", err)
		return 1
	}
	return 0
}

func main2() error {
	base, local, remote, merged := os.Args[1], os.Args[2], os.Args[3], os.Args[4]
	_ = base
	if filepath.Base(merged) != "go.mod" {
		// TODO look at Go files too.
		return fmt.Errorf("cannot fix conflicts in non-go.mod files (base %q)", base)
	}
	localMod, err := parseGoMod(local)
	if err != nil {
		return fmt.Errorf("cannot parse local go.mod: %v", err)
	}
	remoteMod, err := parseGoMod(remote)
	if err != nil {
		return fmt.Errorf("cannot parse remote go.mod: %v", err)
	}
	for _, req := range remoteMod.Require {
		localMod.Require = append(localMod.Require, req)
	}
	// TODO support exclude and replace too
	localMod.Cleanup()
	data, err := localMod.Format()
	if err != nil {
		return fmt.Errorf("cannot format merged go.mod file: %v", err)
	}
	if err := ioutil.WriteFile(merged, data, 0666); err != nil {
		return fmt.Errorf("cannot write merged file: %v", err)
	}
	c := exec.Command("go", "mod", "tidy")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %v", err)
	}
	return nil
}

func parseGoMod(name string) (*modfile.File, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, err
	}
	mf, err := modfile.Parse(name, data, nil)
	if err != nil {
		return nil, fmt.Errorf("%v (contents %q)", err, data)
	}
	return mf, nil
}

func fileContents(name string) string {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return fmt.Sprintf("%s not found", name)
	}
	return fmt.Sprintf("%s=%q", name, data)
}
