package utils

import (
  "net"
  "strings"
  danmtypes "github.com/nokia/danm/crd/apis/danm/v1"
  "github.com/nokia/danm/pkg/bitarray"
  "github.com/nokia/danm/pkg/ipam"
  "github.com/nokia/danm/pkg/admit"
  stubs "github.com/nokia/danm/test/stubs/danm"
)

const (
  AllocFor5k = "gAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
                "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
                "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
                "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
                "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
                "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
                "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA" +
                "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=="
)

func SetupAllocationPools(nets []danmtypes.DanmNet) error {
  for index, dnet := range nets {
    if dnet.Spec.Options.Cidr != "" {
      err := admit.CreateAllocationArray(&dnet)
      if err != nil {
        return err
      }
      _, ipnet, err := net.ParseCIDR(dnet.Spec.Options.Cidr)
      if err != nil {
        return err
      }
      if dnet.Spec.Options.Pool.Start == "" {
        dnet.Spec.Options.Pool.Start = (ipam.Int2ip(ipam.Ip2int(ipnet.IP) + 1)).String()
      }
      if dnet.Spec.Options.Pool.End == "" {
        dnet.Spec.Options.Pool.End = (ipam.Int2ip(ipam.Ip2int(admit.GetBroadcastAddress(ipnet)) - 1)).String()
      }
      if strings.HasPrefix(dnet.ObjectMeta.Name, "full") {
        exhaustNetwork(&dnet)
      }
      nets[index].Spec = dnet.Spec
    }
  }
  return nil
}

func GetTestNet(netId string, testNets []danmtypes.DanmNet) *danmtypes.DanmNet {
  for _, net := range testNets {
    if net.ObjectMeta.Name == netId {
      return &net
    }
  }
  return nil
}

func CreateExpectedAllocationsList(ip string, isExpectedToBeSet bool, networkId string) []stubs.ReservedIpsList {
  var ips []stubs.ReservedIpsList
  if ip != "" {
    reservation := stubs.Reservation {Ip: ip, Set: isExpectedToBeSet,}
    expectedAllocation := stubs.ReservedIpsList{NetworkId: networkId, Reservations: []stubs.Reservation {reservation,},}
    ips = append(ips, expectedAllocation)
  }
  return ips
}

func exhaustNetwork(netInfo *danmtypes.DanmNet) {
    ba := bitarray.NewBitArrayFromBase64(netInfo.Spec.Options.Alloc)
    _, ipnet, _ := net.ParseCIDR(netInfo.Spec.Options.Cidr)
    ipnetNum := ipam.Ip2int(ipnet.IP)
    begin := ipam.Ip2int(net.ParseIP(netInfo.Spec.Options.Pool.Start)) - ipnetNum
    end := ipam.Ip2int(net.ParseIP(netInfo.Spec.Options.Pool.End)) - ipnetNum
    for i:=begin;i<=end;i++ {
        ba.Set(uint32(i))
    }
    netInfo.Spec.Options.Alloc = ba.Encode()
}

func GetTconf(tconfName string, tconfSet []danmtypes.TenantConfig) *danmtypes.TenantConfig {
  for _, tconf := range tconfSet {
    if tconf.ObjectMeta.Name == tconfName {
      return &tconf
    }
  }
  return nil
}
