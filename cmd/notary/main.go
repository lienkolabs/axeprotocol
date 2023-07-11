package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const usage = `usage: notary secure-vault-file <command> [<args>]

These are the common notary commands used in various situations:

synchronize local data to network
    sync      Synchronizes to default node
    pending   Show actions pending of inforporation by the network

create a new handle
    all       Show all available users managed by notary
    new       New ID on existing breeze network
    update    Update existing ID 

power of attorney
    show      Show all live power of attorney
	grant     Grant new power of attorney
	revoke    Revoke existing power of attorney
`

func askConfirm(text string) bool {
	var s string
	fmt.Printf("%v (Y/n): ", text)
	_, err := fmt.Scan(&s)
	if err != nil {
		panic(err)
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	if s == "y" || s == "yes" {
		return true
	}
	return false
}

func defaultNotaryDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	if home := os.Getenv(env); home != "" {
		def := filepath.Join(home, ".notary")
		return def
	}
	log.Fatalf("could not find home directory")
	return ""
}

func main() {
	dirname := defaultNotaryDir()
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		if askConfirm(fmt.Sprintf("directory %v does not exists. Do you want to create one?", dirname)) {
			os.Mkdir(dirname, os.ModePerm)
		} else {
			fmt.Println("exiting notary.")
			return
		}
	}
	fmt.Println(dirname)
	if len(os.Args) > 20 {
		fmt.Printf("%v\n", usage)
		return
	}
}
