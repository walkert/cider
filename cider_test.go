package cider

import (
	"bytes"
	"testing"
)

var PrintNetworkTests = []struct {
	cidr    string
	want    string
	wantErr bool
}{
	{
		"192.168.1.0/28",
		`Network:    192.168.1.0
First:      192.168.1.1
Last:       192.168.1.14
Broadcast:  192.168.1.15
Netmask:    255.255.255.240
Usable:     14
Next:       192.168.1.16
`,
		false},
	{
		"192.168.",
		"invalid CIDR address: 192.168.",
		true,
	},
}

var PrintSubnetsForTests = []struct {
	cidr    string
	seg     int
	want    string
	wantErr bool
}{
	{
		"192.168.1.0/27",
		28,
		`192.168.1.0/28
192.168.1.16/28
`,
		false,
	},
	{
		"192.168.1.0/27",
		26,
		"block size (26) must be less than original CIDR (27)",
		true,
	},
	{
		"192.168.",
		28,
		"invalid CIDR address: 192.168.",
		true,
	},
}

var PrintAllSubnetsTests = []struct {
	cidr    string
	want    string
	wantErr bool
}{
	{
		"192.168.1.0/27",
		`/28:
192.168.1.0/28
192.168.1.16/28

/29:
192.168.1.0/29
192.168.1.8/29
192.168.1.16/29
192.168.1.24/29
`,
		false,
	},
	{
		"192.168.",
		"invalid CIDR address: 192.168.",
		true,
	},
}

var IsInNetTests = []struct {
	cidr string
	ip   string
	want bool
}{
	{
		"192.168.1.0/27",
		"10.1.134.124",
		false,
	},
	{
		"192.168.1.0/27",
		"192.168.1.1",
		true,
	},
}

func TestPrintNetwork(t *testing.T) {
	for _, tt := range PrintNetworkTests {
		buf := &bytes.Buffer{}
		err := PrintNetwork(tt.cidr, buf)
		if err != nil {
			if !tt.wantErr {
				t.Fatalf("Unexpected error: %s\n", err.Error())
			}
			if err.Error() != tt.want {
				t.Fatalf("Wanted error: %s\n, Got error: %s\n", tt.want, err.Error())
			}
		} else {
			if buf.String() != tt.want {
				t.Fatalf("Wanted:\n%s\nGot:\n%s\n", tt.want, buf.String())
			}
		}
	}
}

func TestPrintSubnetsFor(t *testing.T) {
	for _, tt := range PrintSubnetsForTests {
		buf := &bytes.Buffer{}
		err := PrintSubnetsFor(tt.cidr, tt.seg, buf)
		if err != nil {
			if !tt.wantErr {
				t.Fatalf("Unexpected error: %s\n", err.Error())
			}
			if err.Error() != tt.want {
				t.Fatalf("Wanted error: %s\n, Got error: %s\n", tt.want, err.Error())
			}
		} else {
			if buf.String() != tt.want {
				t.Fatalf("Wanted:\n%s\nGot:\n%s\n", tt.want, buf.String())
			}
		}
	}
}

func TestPrintAllSubnets(t *testing.T) {
	for _, tt := range PrintAllSubnetsTests {
		buf := &bytes.Buffer{}
		err := PrintAllSubnets(tt.cidr, buf)
		if err != nil {
			if !tt.wantErr {
				t.Fatalf("Unexpected error: %s\n", err.Error())
			}
			if err.Error() != tt.want {
				t.Fatalf("Wanted error: %s\n, Got error: %s\n", tt.want, err.Error())
			}
		} else {
			if buf.String() != tt.want {
				t.Fatalf("Wanted:\n%s\nGot:\n%s\n", tt.want, buf.String())
			}
		}
	}
}

func TestIsInNet(t *testing.T) {
	for _, tt := range IsInNetTests {
		got, _ := IsInNet(tt.ip, tt.cidr)
		if got != tt.want {
			t.Fatalf("Wanted: %s\nGot:%s\n", tt.want, got)
		}
	}
}
