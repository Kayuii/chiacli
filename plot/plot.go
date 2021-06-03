package plot

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
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
	LogPath string
	PlotID  string
	LogFile string
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
	LogPath    string `yaml:"LogPath"`
	FarmerKey  string `yaml:"FarmerKey"`
	PoolKey    string `yaml:"PoolKey"`
}

func (p *Plot) Chia(config *Config) error {
	if !IsDir(config.TempPath) {
		fmt.Println("获取缓存目录失败")
		os.Exit(0)
	}
	if !IsDir(config.Temp2Path) {
		fmt.Println("获取缓存目录2失败")
		os.Exit(0)
	}
	if !IsDir(config.FinalPath) {
		fmt.Println("获取最终目录失败")
		os.Exit(0)
	}
	if !IsDir(config.LogPath) {
		err := os.MkdirAll(config.LogPath, fs.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
	logPath, err := filepath.Abs(config.LogPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	p.LogPath = logPath

	var (
		ChiaExec string = "chia"
		args     []string
	)

	fmt.Printf("chia utils %s by %s <%s> %s \n", chiacli.Version, chiacli.Author, chiacli.Email, chiacli.Github)

	for i := 1; i <= config.NumPlots; i++ {
		log.SetFlags(log.LstdFlags)
		log.Printf("Plotting %d file \n", i)
		args = p.MakeChiaPlots(*config)
		res, err := p.RunExec(ChiaExec, strconv.Itoa(i), args...)
		if err != nil {
			fmt.Println(err)
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
		fmt.Println("获取缓存目录失败")
		os.Exit(0)
	}
	if !IsDir(config.Temp2Path) {
		fmt.Println("获取缓存目录2失败")
		os.Exit(0)
	}
	if !IsDir(config.FinalPath) {
		fmt.Println("获取最终目录失败")
		os.Exit(0)
	}
	if !IsDir(config.LogPath) {
		err := os.MkdirAll(config.LogPath, fs.ModePerm)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
	logPath, err := filepath.Abs(config.LogPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	p.LogPath = logPath

	var (
		ChiaExec string = "ProofOfSpace"
		args     []string
	)

	fmt.Printf("chia utils %s by %s <%s> %s \n", chiacli.Version, chiacli.Author, chiacli.Email, chiacli.Github)

	for i := 1; i <= config.NumPlots; i++ {
		fmt.Printf("Plotting %d file \n", i)
		args = p.MakeChiaPos(*config)
		res, err := p.RunExec(ChiaExec, strconv.Itoa(i), args...)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if res {
			return nil
		}
	}

	log.SetFlags(log.LstdFlags)
	return nil
}

func (p *Plot) RunExec(ChiaExec, plotnum string, args ...string) (b bool, e error) {

	// cmd := cmd.NewCmd(ChiaExec, args...)
	cmd := cmd.NewCmdOptions(cmd.Options{Streaming: true}, ChiaExec, args...)

	fmt.Println("commandline: ", ChiaExec, strings.Join(args, " "))

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-sigs
		fmt.Println("signal", sig, "called", ". Terminating...")
		cmd.Stop()
		cancel()
	}()

	statusChan := cmd.Start()

	fmt.Printf("Process ID: #%d \n", cmd.Status().PID)
	f, _ := os.OpenFile(fmt.Sprintf("%s/%s", p.LogPath, p.LogFile), os.O_WRONLY|os.O_CREATE|os.O_SYNC|os.O_APPEND, 0644)
	defer f.Close()
	logger := log.New(io.MultiWriter(f, os.Stdout), "", log.Lmsgprefix)

	go func() {
		for cmd.Stdout != nil || cmd.Stderr != nil {
			select {
			case line, open := <-cmd.Stdout:
				if !open {
					cmd.Stdout = nil
					continue
				}
				logger.Println(line)
			case line, open := <-cmd.Stderr:
				if !open {
					cmd.Stderr = nil
					continue
				}
				fmt.Fprintln(os.Stderr, line)
			}
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("context done, cli exiting...")
		return true, nil
	default:
	}

	finalStatus := <-statusChan
	if finalStatus.Error != nil {
		return true, finalStatus.Error
	}

	logger.Printf("CommandLine Use %s", time.Duration(finalStatus.StopTs-finalStatus.StartTs).String())

	return false, nil
}

func (p *Plot) MakeChiaPlots(confYaml Config) []string {
	ChiaCmd := []string{
		"plots",
		"create",
	}

	sk := wallet.KeyGen(wallet.TokenBytes(32))
	farmerPk, err := wallet.DecodePointG1(confYaml.FarmerKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	poolPk, err := wallet.DecodePointG1(confYaml.PoolKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	plot_public_key := wallet.GeneratePlotPublicKey(wallet.MaterSkToLocalSk(sk).GetG1(), farmerPk)
	plotID := wallet.CalculatePlotIdPk(poolPk, plot_public_key)
	p.PlotID = hex.EncodeToString(plotID)[:12]

	fmt.Println("plot id: " + hex.EncodeToString(plotID))

	dt_string := time.Now().Format("2006-01-02-15-04")
	p.LogFile = strings.Join([]string{
		"plot",
		"k" + strconv.Itoa(confYaml.KSize),
		dt_string,
		p.PlotID + ".log",
	}, "-")

	ChiaCmd = append(ChiaCmd,
		"-i", hex.EncodeToString(plotID),
		"-f", confYaml.FarmerKey,
		"-p", confYaml.PoolKey,
		"-k", strconv.Itoa(confYaml.KSize),
		"-r", strconv.Itoa(confYaml.Threads),
		"-u", strconv.Itoa(confYaml.Buckets),
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

func (p *Plot) MakeChiaPos(confYaml Config) []string {
	ChiaCmd := []string{
		"create",
	}

	sk := wallet.KeyGen(wallet.TokenBytes(32))
	farmerPk, err := wallet.DecodePointG1(confYaml.FarmerKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	poolPk, err := wallet.DecodePointG1(confYaml.PoolKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	// log.Printf("sk:" + hex.EncodeToString(sk.Bytes()))
	plot_public_key := wallet.GeneratePlotPublicKey(wallet.MaterSkToLocalSk(sk).GetG1(), farmerPk)
	plotID := wallet.CalculatePlotIdPk(poolPk, plot_public_key)
	p.PlotID = hex.EncodeToString(plotID)[:12]

	fmt.Println("plot id: " + hex.EncodeToString(plotID))

	g1 := bls12381.NewG1()
	plotMemo := make([]byte, 0, 128)
	plotMemo = append(plotMemo, g1.ToCompressed(poolPk)...)   // Len 48
	plotMemo = append(plotMemo, g1.ToCompressed(farmerPk)...) // Len 48
	plotMemo = append(plotMemo, sk.Bytes()...)                // Len 32

	// fmt.Printf("memo: " + hex.EncodeToString(plotID))

	//  "plot-k{args.size}-{dt_string}-{plot_id}.plot"
	dt_string := time.Now().Format("2006-01-02-15-04")
	filename := strings.Join([]string{
		"plot",
		"k" + strconv.Itoa(confYaml.KSize),
		dt_string,
		hex.EncodeToString(plotID) + ".plot",
	}, "-")
	p.LogFile = strings.Join([]string{
		"plot",
		"k" + strconv.Itoa(confYaml.KSize),
		dt_string,
		p.PlotID + ".log",
	}, "-")

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
	fmt.Println(cmd.Args)
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
