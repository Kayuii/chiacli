package plot

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/go-cmd/cmd"
	"github.com/kayuii/chiacli"
	"github.com/kayuii/chiacli/wallet"
	bls12381 "github.com/kilic/bls12-381"
)

type Plot struct {
	dirName string
}

func New() *Plot {
	return &Plot{}
}

type Config struct {
	NumPlots   int    `yaml:"NumPlots"`
	KSize      int    `yaml:"KSize"`
	Stripes    int    `yaml:"Stripes"`
	Buffer     int    `yaml:"Buffer"`
	Threads    int    `yaml:"Threads"`
	Buckets    int    `yaml:"Buckets"`
	NoBitfield bool   `yaml:"NoBitfield"`
	Progress   bool   `yaml:"Progress"`
	TempPath   string `yaml:"TempPath"`
	Temp2Path  string `yaml:"Temp2Path"`
	FinalPath  string `yaml:"FinalPath"`
	Total      int    `yaml:"Total"`
	Sleep      int    `yaml:"Sleep"`
	RunPath    string `yaml:"RunPath"`
	FarmerKey  string `yaml:"FarmerKey"`
	PoolKey    string `yaml:"PoolKey"`
}

func (p *Plot) Chia(config *Config) error {
	if !IsDir(config.TempPath) {
		log.Println("获取缓存目录失败")
		os.Exit(0)
	}
	if !IsDir(config.Temp2Path) {
		log.Println("获取缓存目录2失败")
		os.Exit(0)
	}
	if !IsDir(config.FinalPath) {
		log.Println("获取最终目录失败")
		os.Exit(0)
	}
	var (
		ChiaExec string = "chia"
		args     []string
	)

	log.Printf("chia utils %s by %s <%s> %s \n", chiacli.Version, chiacli.Author, chiacli.Email, chiacli.Github)

	args = MakeChiaPlots(*config)

	for i := 1; i <= config.NumPlots; i++ {
		log.SetFlags(log.LstdFlags)
		log.Printf("Plotting %d file \n", i)
		res, err := RunExec(ChiaExec, strconv.Itoa(i), args...)
		if err != nil {
			log.Println(err)
			return nil
		}
		if res {
			return nil
		}
	}
	return nil
}

func (p *Plot) Pos(config *Config) error {
	if !IsDir(config.TempPath) {
		log.Println("获取缓存目录失败")
		os.Exit(0)
	}
	if !IsDir(config.Temp2Path) {
		log.Println("获取缓存目录2失败")
		os.Exit(0)
	}
	if !IsDir(config.FinalPath) {
		log.Println("获取最终目录失败")
		os.Exit(0)
	}

	var (
		ChiaExec string = "ProofOfSpace"
		args     []string
	)

	log.Printf("chia utils %s by %s <%s> %s \n", chiacli.Version, chiacli.Author, chiacli.Email, chiacli.Github)

	args = MakeChiaPos(*config)

	for i := 1; i <= config.NumPlots; i++ {
		log.SetFlags(log.LstdFlags)
		log.Printf("Plotting %d file \n", i)
		res, err := RunExec(ChiaExec, strconv.Itoa(i), args...)
		if err != nil {
			log.Println(err)
			return nil
		}
		if res {
			return nil
		}
	}

	log.SetFlags(log.LstdFlags)
	return nil
}

func RunExec(ChiaExec, plotnum string, args ...string) (b bool, e error) {

	cmd := cmd.NewCmd(ChiaExec, args...)

	log.Println("commandline: ", ChiaExec, strings.Join(args, " "))

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-sigs
		log.Println("signal", sig, "called", ". Terminating...")
		cmd.Stop()
		cancel()
	}()

	ticker := time.NewTicker(time.Second)
	ln := 0
	le := 0
	statusChan := cmd.Start()

	log.Printf("Process ID: #%d \n", cmd.Status().PID)
	log.SetPrefix(fmt.Sprintf("Plot-%s ", plotnum))
	log.SetFlags(log.Lmsgprefix)

	go func() {
		for range ticker.C {
			status := cmd.Status()
			ne := len(status.Stderr)
			for i := le; i < ne; i++ {
				log.Println(status.Stderr[i])
			}
			le = ne

			n := len(status.Stdout)
			for i := ln; i < n; i++ {
				log.Println(status.Stdout[i])
			}
			ln = n
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("context done, runner exiting...")
		return true, nil
	default:
	}
	finalStatus := <-statusChan
	if finalStatus.Error != nil {
		return true, finalStatus.Error
	}
	ne := len(finalStatus.Stderr)
	for i := le; i < ne; i++ {
		log.Println(finalStatus.Stderr[i])
	}
	n := len(finalStatus.Stdout)
	for i := ln; i < n; i++ {
		log.Println(finalStatus.Stdout[i])
	}
	log.SetPrefix("")
	log.Printf("CommandLine Use %s", time.Duration(finalStatus.StopTs-finalStatus.StartTs).String())

	return false, nil
}

func MakeChiaPlots(confYaml Config) []string {
	ChiaCmd := []string{
		"plots",
		"create",
		"-f", confYaml.FarmerKey,
		"-p", confYaml.PoolKey,
		"-k", strconv.Itoa(confYaml.KSize),
		"-r", strconv.Itoa(confYaml.Threads),
		"-u", strconv.Itoa(confYaml.Buckets),
		"-b", strconv.Itoa(confYaml.Buffer),
	}

	if strings.Compare(confYaml.TempPath, confYaml.Temp2Path) == 0 || strings.Compare(confYaml.Temp2Path, ".") == 0 {
		ChiaCmd = append(ChiaCmd,
			"-t", confYaml.TempPath,
		)
	} else {
		ChiaCmd = append(ChiaCmd,
			"-t", confYaml.TempPath,
			"-2", confYaml.Temp2Path,
		)
	}

	ChiaCmd = append(ChiaCmd,
		"-d", confYaml.FinalPath,
	)

	if confYaml.KSize < 32 {
		ChiaCmd = append(ChiaCmd,
			"--override-k",
		)
	}

	if confYaml.NoBitfield {
		ChiaCmd = append(ChiaCmd,
			"-e",
		)
	}
	return ChiaCmd
}

func MakeChiaPos(confYaml Config) []string {
	ChiaCmd := []string{
		"create",
	}

	sk := wallet.KeyGen(wallet.TokenBytes(32))
	farmerPk, err := wallet.DecodePointG1(confYaml.FarmerKey)
	if err != nil {
		log.Fatal(err)
	}
	poolPk, err := wallet.DecodePointG1(confYaml.PoolKey)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("sk:" + hex.EncodeToString(sk.Bytes()))
	plot_public_key := wallet.GeneratePlotPublicKey(wallet.MaterSkToLocalSk(sk).GetG1(), farmerPk)
	plotID := wallet.CalculatePlotIdPk(poolPk, plot_public_key)

	log.Printf("plot id: " + hex.EncodeToString(plotID))

	g1 := bls12381.NewG1()
	plotMemo := make([]byte, 0, 128)
	plotMemo = append(plotMemo, g1.ToCompressed(poolPk)...)   // Len 48
	plotMemo = append(plotMemo, g1.ToCompressed(farmerPk)...) // Len 48
	plotMemo = append(plotMemo, sk.Bytes()...)                // Len 32

	log.Printf("memo: " + hex.EncodeToString(plotID))

	//  "plot-k{args.size}-{dt_string}-{plot_id}.plot"
	dt_string := time.Now().Format("2006-01-02-15-04")
	filename := strings.Join([]string{
		"plot",
		"k" + strconv.Itoa(confYaml.KSize),
		dt_string,
		hex.EncodeToString(plotID) + ".plot",
	}, "-")

	log.Printf("filename: " + filename)

	ChiaCmd = append(ChiaCmd,
		"-i", "0x"+hex.EncodeToString(plotID),
		"-m", "0x"+hex.EncodeToString(plotMemo),
		"-f", filename,
		"-k", strconv.Itoa(confYaml.KSize),
		"-r", strconv.Itoa(confYaml.Threads),
		"-u", strconv.Itoa(confYaml.Buckets),
		"-s", strconv.Itoa(confYaml.Stripes),
		"-b", strconv.Itoa(confYaml.Buffer),
	)

	if strings.Compare(confYaml.TempPath, confYaml.Temp2Path) == 0 || strings.Compare(confYaml.Temp2Path, ".") == 0 {
		ChiaCmd = append(ChiaCmd,
			"-t", confYaml.TempPath,
		)
	} else {
		ChiaCmd = append(ChiaCmd,
			"-t", confYaml.TempPath,
			"-2", confYaml.Temp2Path,
		)
	}

	ChiaCmd = append(ChiaCmd,
		"-d", confYaml.FinalPath,
	)

	if confYaml.NoBitfield {
		ChiaCmd = append(ChiaCmd,
			"-e",
		)
	}

	return ChiaCmd
}

func GetChieExec(ChiaAppPath string) (ChiaExec string) {
	ChiaExe := "chia"
	LineString := `/`

	ChiaExec = strings.Join([]string{ChiaAppPath, ChiaExe}, LineString)
	return
}

func CmdAndChangeDirToFile(commandName string, params []string) {
	cmd := exec.Command(commandName, params...)
	log.Println(cmd.Args)
	cmd.Start()
	cmd.Wait()
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		err := errors.New("can't find file")
		return "", err
	}
	return string(path[0 : i+1]), nil
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
