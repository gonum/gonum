//+build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	cmd := exec.Command("git", "clone", "https://github.com/xianyi/OpenBLAS")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	fmt.Println("done cloning")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir("OpenBLAS")
	if err != nil {
		log.Fatal(err)
	}
	return

	cmd = exec.Command("make")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
