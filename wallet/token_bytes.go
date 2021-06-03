package wallet

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"math/rand"
	"net"
	"os"
	"time"
)

func init() {
	// rand.Seed(time.Now().UnixNano())
	rand.Seed(Hashseed())
}

func TokenBytes(n int) []byte {
	if n <= 0 {
		return []byte{}
	}
	var need int
	if n&1 == 0 { // even
		need = n
	} else { // odd
		need = n + 1
	}
	size := need / 2
	dst := make([]byte, need)
	src := dst[size:]
	if _, err := rand.Read(src[:]); err != nil {
		return []byte{}
	}
	hex.Encode(dst, src)
	return dst[:n]
}

func GetMacAddrs() (string, error) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}
		return macAddr, nil
	}
	return "", err
}

// hash output uint32
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// add mac and time.now() as seed
func Hashseed() int64 {
	mac_adr, _ := GetMacAddrs()
	hostname, _ := os.Hostname()
	t := time.Now().UnixNano() // int64
	return int64(Hash(fmt.Sprintf("%d %s %s", t, mac_adr, hostname)))
}
