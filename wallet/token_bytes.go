package wallet

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"time"

	password "github.com/1800alex/go-utilities-password"
	bls "github.com/chuwt/chia-bls-go"
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

// hash output uint32
func Hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// add mac and time.now() as seed
func Hashseed() int64 {
	seed, _ := GetRandomSHA256Seed()
	hostname, _ := os.Hostname()
	t := time.Now().UnixNano() // int64
	return int64(Hash(fmt.Sprintf("%d %s %s", t, seed, hostname)))
}

func GetRandomSHA256Seed() (result string, err error) {
	return password.Generate(16, true, false, false, true)
}

func CalculatePlotIdPk(poolPk, plotPK []byte) []byte {
	hash := sha256.New()
	hash.Write(poolPk)
	hash.Write(plotPK)
	return hash.Sum(nil)
}

func PublicKeyFromHexString(key string) (bls.PublicKey, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return bls.PublicKey{}, err
	}
	return bls.NewPublicKey(keyBytes)
}
