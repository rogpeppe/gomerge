package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	os.Exit(main1())
}

func main1() int {
	fmt.Printf("mergetool running\n")
	base, local, remote, merged := os.Args[1], os.Args[2], os.Args[3], os.Args[4]
	fmt.Printf("	BASE=%s\n", fileContents(base))
	fmt.Printf("	LOCAL=%s\n", fileContents(local))
	fmt.Printf("	REMOTE=%s\n", fileContents(remote))
	fmt.Printf("	MERGED=%s\n", fileContents(merged))
	return 1
}

func fileContents(name string) string {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return fmt.Sprintf("%s not found", name)
	}
	return fmt.Sprintf("%s=%q", name, data)
}
