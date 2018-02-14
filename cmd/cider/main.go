package main

import (
	"flag"
	"fmt"
	"github.com/walkert/cider"
	"os"
)

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func main() {
	var s = flag.Int("s", 0, "")
	var a = flag.Bool("a", false, "")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-s SIZE] [-a] CIDR\n", os.Args[0])
		fmt.Printf("Options:\n  -s SIZE\tprint all subnets of this block size that can fit within CIDR\n")
		fmt.Printf("  -a\t\tprint all subnets that can fit within this CIDR\n")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	cidr := flag.Arg(0)
	if *s > 1 {
		err := cider.PrintSubnetsFor(cidr, *s, os.Stdout)
		exit(err)
	}
	if *a {
		err := cider.PrintAllSubnets(cidr, os.Stdout)
		exit(err)
	}
	err := cider.PrintNetwork(cidr, os.Stdout)
	exit(err)
}
