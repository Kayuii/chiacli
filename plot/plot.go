package plot

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	bls "github.com/chuwt/chia-bls-go"
	"github.com/go-cmd/cmd"
	"github.com/kayuii/chiacli"
	"github.com/kayuii/chiacli/wallet"
)

type Plot struct {
	dirName   string
	LogPath   string
	PlotID    string
	LogFile   string
	Phase     int
	Len       int
	Buckets   int
	Table     int
	EndPhase  int
	PhaseTime [4]int64
}

func New() *Plot {
	return &Plot{
		Phase:    0,
		Len:      0,
		Buckets:  0,
		Table:    0,
		EndPhase: 0,
	}
}

type Config struct {
	NumPlots            int    `yaml:"NumPlots"`
	KSize               int    `yaml:"KSize"`
	Stripes             int    `yaml:"Stripes"`
	Buffer              int    `yaml:"Buffer"`
	Threads             int    `yaml:"Threads"`
	Rmulti2             int    `yaml:"Rmulti2"`
	Buckets             int    `yaml:"Buckets"`
	NoBitfield          bool   `yaml:"NoBitfield"`
	Progress            bool   `yaml:"Progress"`
	TempPath            string `yaml:"TempPath"`
	Temp2Path           string `yaml:"Temp2Path"`
	FinalPath           string `yaml:"FinalPath"`
	Total               int    `yaml:"Total"`
	Sleep               int    `yaml:"Sleep"`
	LogPath             string `yaml:"LogPath"`
	FarmePublicKey      string `yaml:"FarmePublicKey"`
	PoolPublicKey       string `yaml:"PoolPublicKey"`
	PoolContractAddress string `yaml:"PoolContractAddress"`
	LocalSk             string `yaml:"LocalSk"`
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

func (p *Plot) FastPos(config *Config) error {
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
		ChiaExec string = "fastpos"
		args     []string
	)

	fmt.Printf("chia utils %s by %s <%s> %s \n", chiacli.Version, chiacli.Author, chiacli.Email, chiacli.Github)

	for i := 1; i <= config.NumPlots; i++ {
		fmt.Printf("Plotting %d file \n", i)
		args = p.MakeFastPos(*config)
		res, err := p.RunExec(ChiaExec, strconv.Itoa(i), args...)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if res {
			return nil
		}
		fmt.Printf("Sleep %d sec \n", config.Sleep)
		time.Sleep(time.Duration(config.Sleep) * time.Second)
	}

	log.SetFlags(log.LstdFlags)
	return nil
}

func (p *Plot) Bladebit(config *Config) error {
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
		ChiaExec string = "bladebit"
		args     []string
	)

	fmt.Printf("chia utils %s by %s <%s> %s \n", chiacli.Version, chiacli.Author, chiacli.Email, chiacli.Github)

	for i := 1; i <= config.NumPlots; i++ {
		fmt.Printf("Plotting %d file \n", i)
		args = p.MakeBladeBit(*config)
		res, err := p.RunExec(ChiaExec, strconv.Itoa(i), args...)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if res {
			return nil
		}
		fmt.Printf("Sleep %d sec \n", config.Sleep)
		if i < config.NumPlots-1 {
			time.Sleep(time.Duration(config.Sleep) * time.Second)
		}
	}

	log.SetFlags(log.LstdFlags)
	return nil
}

func (p *Plot) RunExec(ChiaExec, plotnum string, args ...string) (b bool, e error) {

	p.Len = 0
	p.Phase = 0
	p.Buckets = 0
	p.Table = 0
	p.EndPhase = 1
	p.PhaseTime = [4]int64{0, 0, 0, 0}

	// cmd := cmd.NewCmd(ChiaExec, args...)
	cmd := cmd.NewCmdOptions(cmd.Options{Streaming: true}, ChiaExec, args...)

	fmt.Println("commandline: \"", ChiaExec, strings.Join(args, " "), "\"")

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
	// logger := log.New(io.MultiWriter(f, os.Stdout), "", log.Lmsgprefix)

	log.SetOutput(f)
	log.SetPrefix("")
	log.SetFlags(log.Lmsgprefix)

	go func() {
		for cmd.Stdout != nil || cmd.Stderr != nil {
			select {
			case line, open := <-cmd.Stdout:
				if !open {
					cmd.Stdout = nil
					continue
				}
				p.Len++
				log.Println(line)
				// p.FormatProgressShow(line)
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
	usedtime := time.Duration(finalStatus.StopTs - finalStatus.StartTs).String()

	log.Printf("CommandLine Use %s", usedtime)
	fmt.Printf("CommandLine Use %s \n", usedtime)

	return false, nil
}

func (p *Plot) MakeChiaPlots(confYaml Config) []string {
	ChiaCmd := []string{
		"plots",
		"create",
	}

	sk := bls.KeyGen(wallet.TokenBytes(32))
	farmerPk, err := wallet.PublicKeyFromHexString(confYaml.FarmePublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	poolPk, err := wallet.PublicKeyFromHexString(confYaml.PoolPublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	plotPk := farmerPk.Add(sk.LocalSk().GetPublicKey())
	plotID := wallet.CalculatePlotIdPk(poolPk.Bytes(), plotPk.Bytes())

	p.PlotID = hex.EncodeToString(plotID)[:12]

	fmt.Println("plot id: " + hex.EncodeToString(plotID))

	dt_string := time.Now().Format("2006-01-02-15-04")
	p.LogFile = strings.Join([]string{
		"plot",
		"k" + strconv.Itoa(confYaml.KSize),
		dt_string,
		p.PlotID + ".log",
	}, "-")

	// "-i", hex.EncodeToString(plotID),

	ChiaCmd = append(ChiaCmd,
		"-f", confYaml.FarmePublicKey,
		"-p", confYaml.PoolPublicKey,
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

	sk := bls.KeyGen(wallet.TokenBytes(32))
	farmerPk, err := wallet.PublicKeyFromHexString(confYaml.FarmePublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	poolPk, err := wallet.PublicKeyFromHexString(confYaml.PoolPublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	plotPk := farmerPk.Add(sk.LocalSk().GetPublicKey())
	plotID := wallet.CalculatePlotIdPk(poolPk.Bytes(), plotPk.Bytes())

	p.PlotID = hex.EncodeToString(plotID)[:12]

	fmt.Println("plot id: " + hex.EncodeToString(plotID))

	plotMemo := make([]byte, 0, 128)
	plotMemo = append(plotMemo, poolPk.Bytes()...)   // Len 48
	plotMemo = append(plotMemo, farmerPk.Bytes()...) // Len 48
	plotMemo = append(plotMemo, sk.Bytes()...)       // Len 32

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

	// ChiaCmd = append(ChiaCmd,
	// 	"-p",
	// )

	return ChiaCmd
}

func (p *Plot) MakeFastPos(confYaml Config) []string {
	ChiaCmd := []string{
		"create",
	}

	sk := bls.KeyGen(wallet.TokenBytes(32))
	farmerPk, err := wallet.PublicKeyFromHexString(confYaml.FarmePublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	poolPk, err := wallet.PublicKeyFromHexString(confYaml.PoolPublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	plotPk := farmerPk.Add(sk.LocalSk().GetPublicKey())
	plotID := wallet.CalculatePlotIdPk(poolPk.Bytes(), plotPk.Bytes())

	p.PlotID = hex.EncodeToString(plotID)[:12]

	dt_string := time.Now().Format("2006-01-02-15-04")

	p.LogFile = strings.Join([]string{
		"plot",
		"k" + strconv.Itoa(confYaml.KSize),
		dt_string,
		p.PlotID + ".log",
	}, "-")

	if strings.Compare(confYaml.PoolContractAddress, "") == 0 {
		ChiaCmd = append(ChiaCmd,
			"-p", confYaml.PoolPublicKey,
		)
	} else {
		ChiaCmd = append(ChiaCmd,
			"-c", confYaml.PoolContractAddress,
		)
	}

	ChiaCmd = append(ChiaCmd,
		"-f", confYaml.FarmePublicKey,
		"-r", strconv.Itoa(confYaml.Threads),
		"-K", strconv.Itoa(confYaml.Rmulti2),
		"-u", strconv.Itoa(confYaml.Buckets),
	)

	if strings.Compare(confYaml.TempPath, confYaml.Temp2Path) == 0 || strings.Compare(confYaml.Temp2Path, ".") == 0 {
		ChiaCmd = append(ChiaCmd,
			"-t", confYaml.TempPath+"/",
		)
	} else {
		ChiaCmd = append(ChiaCmd,
			"-t", confYaml.TempPath+"/",
			"-2", confYaml.Temp2Path+"/",
		)
	}

	ChiaCmd = append(ChiaCmd,
		"-d", confYaml.FinalPath+"/",
	)

	return ChiaCmd
}

func (p *Plot) MakeBladeBit(confYaml Config) []string {
	ChiaCmd := []string{
		"-v",
	}

	sk := bls.KeyGen(wallet.TokenBytes(32))
	farmerPk, err := wallet.PublicKeyFromHexString(confYaml.FarmePublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	poolPk, err := wallet.PublicKeyFromHexString(confYaml.PoolPublicKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	plotPk := farmerPk.Add(sk.LocalSk().GetPublicKey())
	plotID := wallet.CalculatePlotIdPk(poolPk.Bytes(), plotPk.Bytes())

	p.PlotID = hex.EncodeToString(plotID)[:12]

	dt_string := time.Now().Format("2006-01-02-15-04")

	p.LogFile = strings.Join([]string{
		"plot",
		"k" + strconv.Itoa(confYaml.KSize),
		dt_string,
		p.PlotID + ".log",
	}, "-")

	if strings.Compare(confYaml.PoolContractAddress, "") == 0 {
		ChiaCmd = append(ChiaCmd,
			"-p", confYaml.PoolPublicKey,
		)
	} else {
		ChiaCmd = append(ChiaCmd,
			"-c", confYaml.PoolContractAddress,
		)
	}

	ChiaCmd = append(ChiaCmd,
		"-f", confYaml.FarmePublicKey,
		"-t", strconv.Itoa(confYaml.Threads),
	)

	ChiaCmd = append(ChiaCmd,
		confYaml.FinalPath+"/",
	)

	return ChiaCmd
}

func (p *Plot) FormatProgressShow(line string) {
	progress := ""
	phaseTime := ""
	phase := ""

	rs := regexp.MustCompile(`Starting phase ([\d]+)/4`).FindStringSubmatch(line)
	if len(rs) > 0 {
		// p.Phase, _ = strconv.Atoi(rs[1])
		p.Phase++
	}

	rs = regexp.MustCompile(`Time for phase ([\d]+) = ([\d.]+) seconds`).FindStringSubmatch(line)
	if len(rs) > 0 {
		// endPhase, _ := strconv.Atoi(rs[1])
		phaseTime, _ := strconv.ParseInt(strings.ReplaceAll(rs[2], ".", ""), 10, 64)
		p.PhaseTime[p.Phase-1] = phaseTime * 1000 * 1000
		progress = fmt.Sprintf("%0.3f", 99/4.0*float64(p.EndPhase))
		p.EndPhase++
	}

	phaseTime = fmt.Sprintf("%s / %s / %s / %s", time.Duration(p.PhaseTime[0]).String(), time.Duration(p.PhaseTime[1]).String(), time.Duration(p.PhaseTime[2]).String(), time.Duration(p.PhaseTime[3]).String())

	switch p.Phase {
	case 0:
		rs = regexp.MustCompile(`Using ([\d]+) buckets`).FindStringSubmatch(line)
		if len(rs) > 0 {
			p.Buckets, _ = strconv.Atoi(rs[1])
		}
		break
	case 1:
		rs := regexp.MustCompile(`Computing table ([\d]+)`).FindStringSubmatch(line)
		if len(rs) > 0 {
			p.Table, _ = strconv.Atoi(rs[1])
			progress = fmt.Sprintf("%0.3f", 99/4.0/8.0*float64(p.Table))
		}
		phase = " [P1]"
		break
	case 2:
		rs := regexp.MustCompile(`Backpropagating on table ([\d]+)`).FindStringSubmatch(line)
		if len(rs) > 0 {
			p.Table, _ = strconv.Atoi(rs[1])
			progress = fmt.Sprintf("%0.3f", (99/4.0/8.0*float64(8-p.Table))+(99/4.0*float64(p.EndPhase-1)))
		}
		phase = " [P2]"
		break
	case 3:
		rs := regexp.MustCompile(`Compressing tables ([\d]+)`).FindStringSubmatch(line)
		if len(rs) > 0 {
			p.Table, _ = strconv.Atoi(rs[1])
			progress = fmt.Sprintf("%0.3f", (99/4.0/7.0*float64(p.Table))+(99/4.0*float64(p.EndPhase-1)))
		}
		phase = " [P3]"
		break
	case 4:
		phase = " [P4]"
		break
	case 5:

		break
	}
	if p.EndPhase >= 4 {
		rs := regexp.MustCompile(`Total time = ([\d.]+) seconds`).FindStringSubmatch(line)
		if len(rs) > 0 {
			totaltime, _ := strconv.ParseInt(strings.ReplaceAll(rs[1], ".", ""), 10, 64)
			fmt.Printf("Plot file Used: %s \n", time.Duration(totaltime*1000).String())
		}
		rs = regexp.MustCompile(`Copy time = ([\d.]+) seconds`).FindStringSubmatch(line)
		if len(rs) > 0 {
			totaltime, _ := strconv.ParseInt(strings.ReplaceAll(rs[1], ".", ""), 10, 64)
			fmt.Printf("Copy file Used: %s \n", time.Duration(totaltime*1000).String())
		}
		if re := regexp.MustCompile("Renamed final file from").MatchString(line); re {
			progress = "100.000"
		}
	}
	if len(progress) > 0 {
		fmt.Printf("Progress: %s ==>%s phase_times %s \n", progress, phase, phaseTime)
	}
	fmt.Println(line)
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
