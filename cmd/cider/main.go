package main

import (
	"flag"
	"fmt"
	"github.com/walkert/cider"
	"github.com/walkert/pager"
	"os"
)

func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
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
	var err error
	output := pager.New()
	if *s > 1 {
		err = cider.PrintSubnetsFor(cidr, *s, output)
	} else if *a {
		err = cider.PrintAllSubnets(cidr, output)
	} else {
		err = cider.PrintNetwork(cidr, output)
	}
	exit(err)
	output.Page()
}
