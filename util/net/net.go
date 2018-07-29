package net_util

import (
	"net"
)

func GetLocalIp(name string) (string, error){
	i, iErr := net.InterfaceByName(name)
	if iErr != nil {
		return "", iErr
	}

	addr, getAddrError := i.Addrs()
	if getAddrError != nil {
		return "", getAddrError
	}

	for _, v := range addr {
		tmp := v.(*net.IPNet).IP.To4()
		if tmp == nil {
			continue
		}
		return tmp.String(), nil
	}

	return "", nil
}
