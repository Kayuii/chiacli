package wallet_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	bls "github.com/chuwt/chia-bls-go"
	"github.com/kayuii/chiacli/wallet"
	_ "github.com/kayuii/chiacli/wallet"
	"github.com/stretchr/testify/require"
)

func TestMemo(t *testing.T) {
	size := 32

	farmerPk, err := wallet.PublicKeyFromHexString("96160804d76ccb56d937536935da2f5ecd32b19d55b56c1ca6c9bc24044ef1d118a8d773ec146130354f19a43483bac0")
	require.NoError(t, err)

	poolPk, err := wallet.PublicKeyFromHexString("b6e26610006b42b33bbc458dc42e8a41bcf25403382dd0074d61679a792f3570e54c22bca6d9863f6c4b22a68355e614")
	require.NoError(t, err)

	masterSk, err := bls.KeyFromHexString("436a41e23b17762e48c93c563aae817b182ec9501a54eb31bb6ecfb2de6221e7")
	require.NoError(t, err)

	plotPk := farmerPk.Add(masterSk.LocalSk().GetPublicKey())
	plotID := wallet.CalculatePlotIdPk(poolPk.Bytes(), plotPk.Bytes())

	// g1 := bls12381.NewG1()
	plotMemo := make([]byte, 0, 128)
	plotMemo = append(plotMemo, poolPk.Bytes()...)   // Len 48
	plotMemo = append(plotMemo, farmerPk.Bytes()...) // Len 48
	plotMemo = append(plotMemo, masterSk.Bytes()...) // Len 32

	filename := fmt.Sprintf("plot-k%d-%s-%s.plot",
		size,
		time.Now().Format("06-01-02-15-04"),
		hex.EncodeToString(plotID),
	)

	// require.Equal(t,
	// 	"436a41e23b17762e48c93c563aae817b182ec9501a54eb31bb6ecfb2de6221e7",
	// 	hex.EncodeToString(masterSk.Bytes()),
	// )
	// require.Equal(t,
	// 	"8b77e652199d65d5c8b8462bc6e48b081bb43ce5fc82d0b2ee4570c040bafa04b4af9f294281127e3db18c49e218afe0",
	// 	hex.EncodeToString(plotPk.Bytes()),
	// )
	require.Equal(t,
		"b715a64d889f437b09866ddb45e56284f5cf9a183b809dff51139029fd5864d7",
		// "e6a21bc95f7ade749cedd5694a9447909c246df3f7b040b5c98911b23c2b91a9",
		hex.EncodeToString(plotID),
	)

	require.Equal(t, ""+
		"b6e26610006b42b33bbc458dc42e8a41"+
		"bcf25403382dd0074d61679a792f3570"+
		"e54c22bca6d9863f6c4b22a68355e614"+
		"96160804d76ccb56d937536935da2f5e"+
		"cd32b19d55b56c1ca6c9bc24044ef1d1"+
		"18a8d773ec146130354f19a43483bac0"+
		"499bc89e17efd8658f8bf22b00f50c4d"+
		"fe5ce23e9f5363e200e9331957dfe586",
		hex.EncodeToString(plotMemo),
	)

	t.Log(filename)
}
