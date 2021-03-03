package main

import (
	"fmt"
	"io/ioutil"
	"os"
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

//
//for each requirement:
//	if requirement is already present at exact version,
//		continue
//	if requirement already has newer version:
//		continue

func main2() error {
	fmt.Println(os.Args)
	wd, _ := os.Getwd()
	fmt.Println("cwd: ", wd)
	if len(os.Args) != 4 {
		return fmt.Errorf("usage: gomerge current-path(%%A) other-path(%%B) final-path(%%P)")
	}
	// Note: the contract for custom merge drivers is defined here:
	// https://git-scm.com/docs/gitattributes#_defining_a_custom_merge_driver
	// The driver should write the resolved file by overwriting the current path.
	current, remote, finalPath := os.Args[1], os.Args[2], os.Args[3]
	if filepath.Base(finalPath) != "go.mod" {
		// TODO look at Go files too.
		return fmt.Errorf("cannot fix conflicts in non-go.mod files (base %q)", finalPath)
	}

	currentMod, err := parseGoMod(current)
	if err != nil {
		return fmt.Errorf("cannot parse local go.mod: %v", err)
	}
	remoteMod, err := parseGoMod(remote)
	if err != nil {
		return fmt.Errorf("cannot parse remote go.mod: %v", err)
	}
	for _, req := range remoteMod.Require {
		currentMod.AddNewRequire(req.Mod.Path, req.Mod.Version, req.Indirect)
	}
	// TODO support exclude and replace too
	currentMod.Cleanup()
	data, err := currentMod.Format()
	if err != nil {
		return fmt.Errorf("cannot format merged go.mod file: %v", err)
	}
	if err := ioutil.WriteFile(current, data, 0666); err != nil {
		return fmt.Errorf("cannot write merged file: %v", err)
	}
	fmt.Printf("overwrite %v with %q\n", current, data)
	//c := exec.Command("go", "mod", "tidy")
	//c.Stdout = os.Stdout
	//c.Stderr = os.Stderr
	//if err := c.Run(); err != nil {
	//	return fmt.Errorf("go mod tidy failed: %v", err)
	//}
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
