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
	var a = flag.Bool("a", false, "")
	var i = flag.String("i", "", "")
	var s = flag.Int("s", 0, "")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-a] [-i IP] [-s SIZE] CIDR\n", os.Args[0])
		fmt.Printf("Options:\n  -a\t\tprint all subnets that can fit within CIDR\n")
		fmt.Printf("  -i IP\t\tcheck whether IP is a part of CIDR\n")
		fmt.Printf("  -s SIZE\tprint all subnets of this block size that can fit within CIDR\n")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}
	cidr := flag.Arg(0)
	var err error
	output := pager.New()
	if *a {
		err = cider.PrintAllSubnets(cidr, output)
	} else if *s > 1 {
		err = cider.PrintSubnetsFor(cidr, *s, output)
	} else if *i != "" {
		isIn, err := cider.IsInNet(*i, cidr)
		exit(err)
		var result string
		if isIn {
			result = " "
		} else {
			result = " not "
		}
		fmt.Printf("IP %s is%sin CIDR %s\n", *i, result, cidr)
		os.Exit(0)
	} else {
		err = cider.PrintNetwork(cidr, output)
	}
	exit(err)
	output.Page()
}
