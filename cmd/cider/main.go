package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/walkert/cider"
	"github.com/walkert/pager"
)

// exit prints errors to stderr and exits 1 if an error occurred
func exit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

// matchingCIDRs returns a list of all cidirs which contain ip
func matchingCIDRs(ip string, cidrs []string) (matches []string, err error) {
	for _, cidr := range cidrs {
		isIn, err := cider.IsInNet(ip, cidr)
		if err != nil {
			return matches, err
		}
		if isIn {
			matches = append(matches, cidr)
		}
		sort.Strings(matches)
	}
	return matches, nil
}

// getCidrs returns a list of all lines in fname
func getCidrs(fname string) (cidrs []string, err error) {
	fh, err := os.Open(fname)
	if err != nil {
		return cidrs, err
	}
	scanner := bufio.NewScanner(fh)
	for scanner.Scan() {
		cidrs = append(cidrs, scanner.Text())
	}
	return cidrs, nil
}

func main() {
	var a = flag.Bool("a", false, "")
	var f = flag.String("f", "", "")
	var i = flag.String("i", "", "")
	var s = flag.Int("s", 0, "")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-a] [-f FILE] [-i IP] [-s SIZE] CIDR\n", os.Args[0])
		fmt.Printf("Options:\n  -a\t\tprint all subnets that can fit within CIDR\n")
		fmt.Printf("  -f FILE\tRead CIDR blocks from this file\n")
		fmt.Printf("  -i IP\t\tcheck whether IP is a part of CIDR\n")
		fmt.Printf("  -s SIZE\tprint all subnets of this block size that can fit within CIDR\n")
	}
	flag.Parse()
	if len(flag.Args()) != 1 {
		if !(*i != "" && *f != "") {
			flag.Usage()
			os.Exit(1)
		}
	}
	cidr := flag.Arg(0)
	var err error
	output := pager.New()
	if *a {
		err = cider.PrintAllSubnets(cidr, output)
	} else if *s > 1 {
		err = cider.PrintSubnetsFor(cidr, *s, output)
	} else if *i != "" {
		var cidrs []string
		if *f != "" {
			cidrs, err = getCidrs(*f)
			exit(err)
		} else {
			cidrs = []string{cidr}
		}
		matches, err := matchingCIDRs(*i, cidrs)
		exit(err)
		l := len(matches)
		switch {
		case l == 0:
			fmt.Println("\u2716")
		case l == 1:
			if *f != "" {
				fmt.Println(matches[0])
			} else {
				fmt.Println("\u2714")
			}
		case l > 1:
			fmt.Printf("%s\n", strings.Join(matches, "\n"))
			os.Exit(0)
		}
	} else {
		err = cider.PrintNetwork(cidr, output)
	}
	exit(err)
	output.Page()
}
