package cider

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// operators maps operator strings to functions
var operators = map[string]func(x, y uint32) uint32{
	"inc": func(x, y uint32) uint32 { return x + y },
	"dec": func(x, y uint32) uint32 { return x - y },
}

// Network object holds details about a given CIDR
type Network struct {
	Usable      int
	Net         net.IP
	First       net.IP
	Last        net.IP
	Next        net.IP
	Mask        net.IP
	MaskSize    int
	Broadcast   net.IP
	NextNetwork net.IP
}

// Summary hold information about Network segments
type Summary struct {
	Count       int
	Networks    []string
	TotalUsable int
	MaskSize    int
}

// update returns a net.IP object that has been
// incremented or decremented by the supplied value
func update(ip net.IP, by int, op string) net.IP {
	ival := binary.BigEndian.Uint32(ip)
	ival = operators[op](ival, uint32(by))
	bval := make([]byte, net.IPv4len)
	binary.BigEndian.PutUint32(bval, ival)
	return net.IP(bval)
}

// Details returns a fully populated Network object based on the cidr string
func Details(cidr string) (Network, error) {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return Network{}, err
	}
	count := net.IPv4Mask(0, 0, 0, 0)
	network := net.IPv4(0, 0, 0, 0).To4()
	netmask := net.IPv4(0, 0, 0, 0).To4()
	broadcast := net.IPv4(0, 0, 0, 0).To4()
	for i := 0; i < net.IPv4len; i++ {
		// To get the count, we need to discover the opposite of the mask.
		// to do this, we eXclusive OR the mask with another mask which is
		// all 1s. The shorthand for this is ^number which is the equivalent
		// to 0xff ^ number.
		count[i] = ^ipnet.Mask[i]
		// To get the network - we take the mask and AND it with the original IP.
		// This will leave the network bits of the IP set to 1
		network[i] = ipnet.Mask[i] & ip.To4()[i]
		// To get the broadcast, we need to take the opposite of the mask (ones
		// in the host bits) and then OR this with zeros
		broadcast[i] = ^ipnet.Mask[i] | ip.To4()[i]
		// For netmask we just take the bits from the mask and OR them
		netmask[i] = ipnet.Mask[i] | netmask[i]
	}
	maskSize, _ := ipnet.Mask.Size()
	return Network{
		Usable:      int(binary.BigEndian.Uint32(count) - 1),
		Net:         network,
		First:       update(network, 1, "inc"),
		Last:        update(broadcast, 1, "dec"),
		Next:        update(broadcast, 1, "inc"),
		Mask:        netmask,
		MaskSize:    maskSize,
		Broadcast:   broadcast,
		NextNetwork: update(broadcast, 1, "inc"),
	}, nil
}

// GetSubnetsFor returns a Summary object for the cidr/size combination
func GetSubnetsFor(cidr string, size int) (Summary, error) {
	summary := Summary{}
	outer, err := Details(cidr)
	if err != nil {
		return summary, err
	}
	if size < outer.MaskSize {
		return summary, fmt.Errorf("block size (%d) must be less than original CIDR (%d)", size, outer.MaskSize)
	}
	end := update(outer.Broadcast, 1, "inc")
	current, _ := Details(fmt.Sprintf("%s/%d", outer.Net.String(), size))
	for {
		summary.Networks = append(summary.Networks, current.Net.String())
		summary.TotalUsable += current.Usable
		summary.MaskSize = current.MaskSize
		next := current.NextNetwork.String()
		if next == end.String() {
			break
		}
		current, _ = Details(fmt.Sprintf("%s/%d", next, size))
	}
	summary.Count = len(summary.Networks)
	return summary, nil
}

// PrintNetwork prints the network detail for the cidr string to output
func PrintNetwork(cidr string, output io.Writer) error {
	netinfo, err := Details(cidr)
	if err != nil {
		return err
	}
	fmt.Fprintf(output, "%-11s %s\n", "Network:", netinfo.Net.String())
	fmt.Fprintf(output, "%-11s %s\n", "First:", netinfo.First.String())
	fmt.Fprintf(output, "%-11s %s\n", "Last:", netinfo.Last.String())
	fmt.Fprintf(output, "%-11s %s\n", "Broadcast:", netinfo.Broadcast.String())
	fmt.Fprintf(output, "%-11s %s\n", "Netmask:", netinfo.Mask.String())
	fmt.Fprintf(output, "%-11s %d\n", "Usable:", netinfo.Usable)
	fmt.Fprintf(output, "%-11s %s\n", "Next:", netinfo.Next.String())
	return nil
}

// PrintAllSubnets prints all networks for the cidr string to output
func PrintAllSubnets(cidr string, output io.Writer) error {
	details, err := Details(cidr)
	if err != nil {
		return err
	}
	var blocks []Summary
	for i := details.MaskSize + 1; i < 30; i++ {
		s, _ := GetSubnetsFor(cidr, i)
		blocks = append(blocks, s)
	}
	for idx, s := range blocks {
		fmt.Fprintf(output, "/%d:\n", s.MaskSize)
		for _, n := range s.Networks {
			fmt.Fprintf(output, "%s/%d\n", n, s.MaskSize)
		}
		if idx != len(blocks)-1 {
			fmt.Fprintln(output)
		}
	}
	return nil
}

// PrintSubnetsFor prints all networks of segment size for the cidr string to output
func PrintSubnetsFor(cidr string, segment int, output io.Writer) error {
	sum, err := GetSubnetsFor(cidr, segment)
	if err != nil {
		return err
	}
	for _, net := range sum.Networks {
		fmt.Fprintf(output, "%s/%d\n", net, segment)
	}
	return nil
}

// IsInNet returns a boolean to indicate whether an IP is in a CIDR block
func IsInNet(ip, cidr string) (bool, error) {
	inNet, err := Details(cidr)
	if err != nil {
		return false, err
	}
	ipNet, err := Details(fmt.Sprintf("%s/%d", ip, inNet.MaskSize))
	if err != nil {
		return false, err
	}
	return inNet.Net.Equal(ipNet.Net), nil
}
