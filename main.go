package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/netip"
	"os"
	"strings"
)

type Networks []netip.Prefix

func (p Networks) Contains(ip netip.Addr) bool {
	for _, prefix := range p {
		if prefix.Contains(ip) {
			return true
		}
	}
	return false
}

func scan(input string, fn func(string)) error {
	var err error
	var file *os.File
	if file, err = os.Open(input); err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			fn(line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func loadNetworks(input string) (Networks, error) {
	networks := Networks{}
	var err error

	err = scan(input, func(line string) {
		var network netip.Prefix
		if network, err = netip.ParsePrefix(line); err == nil {
			networks = append(networks, network)
		} else {
			fmt.Println(err.Error())
		}
	})
	return networks, err
}

var Reverse bool

func main() {
	var err error

	ptrAddressesPath := flag.String("i", "", "the file holding IP Addresses")
	ptrNetworksPath := flag.String("n", "", "the file holding Subnets in CIDR format")
	ptrReverse := flag.Bool("r", false, "asd")
	flag.Parse()

	addressesPath := *ptrAddressesPath
	networksPath := *ptrNetworksPath
	Reverse = *ptrReverse

	if len(addressesPath) == 0 || len(networksPath) == 0 {
		fmt.Println("Usage: ipmatch -i <addressfile> -n <subnetfile> [-r]")
		return
	}

	var networks Networks

	if networks, err = loadNetworks(networksPath); err != nil {
		fmt.Println(err.Error())
		return
	}
	if err = scan(addressesPath, func(line string) {
		ip, err := netip.ParseAddr(line)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			if networks.Contains(ip) != Reverse {
				fmt.Println(ip.String())
			}
		}
	}); err != nil {
		fmt.Println(err.Error())
	}
}
