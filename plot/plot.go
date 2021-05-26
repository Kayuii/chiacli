package plot

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/go-cmd/cmd"
)

type Plot struct {
	jsonIndent func(data interface{}) ([]byte, error)
	toJson     func(data []byte, v interface{}) error
	jsonToYAML func(data []byte) ([]byte, error)
	dirName    string
}

func New() *Plot {
	return &Plot{
		jsonIndent: func(data interface{}) ([]byte, error) {
			return json.MarshalIndent(data, "", "    ")
		},
		toJson: func(data []byte, v interface{}) error {
			return json.Unmarshal(data, &v)
		},
	}
}

type Config struct {
	PlotID     string `yaml:"PlotID"`
	NumPlots   int    `yaml:"NumPlots"`
	KSize      int    `yaml:"KSize"`
	Stripes    int    `yaml:"Stripes"`
	Buffer     int    `yaml:"Buffer"`
	Threads    int    `yaml:"Threads"`
	Buckets    int    `yaml:"Buckets"`
	NoBitfield bool   `yaml:"NoBitfield"`
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

	var (
		ChiaExec string = "/usr/local/bin/chia"
		args     []string
	)

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

	args = MakeChiaPlots(*config)

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
	statusChan := cmd.Start()
	log.Println("Process ID: #%D", cmd.Status().PID)
	go func() {
		for range ticker.C {
			status := cmd.Status()
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
	case finalStatus := <-statusChan:
		n := len(finalStatus.Stdout)
		log.Println(finalStatus.Stdout[n-1])
	default:
	}
	finalStatus := <-statusChan
	if finalStatus.Error != nil {
		log.Println(finalStatus.Error)
	}

	log.Printf("CommandLine Use %s", time.Duration(finalStatus.StopTs-finalStatus.StartTs).String())
	return nil
}

func (p *Plot) Build(config *Config) error {

	var (
		ChiaExec string = "ProofOfSpace"
		args     []string
	)

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

	args = MakeChiaPlots(*config)

	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		sig := <-sigs
		log.Println("signal", sig, "called", ". Terminating...")
		cancel()
	}()

	cmd := cmd.NewCmd(ChiaExec, args...)
	log.Println("running cmd: ", ChiaExec, strings.Join(args, " "))

	ticker := time.NewTicker(2 * time.Second)
	statusChan := cmd.Start()
	go func() {
		for range ticker.C {
			status := cmd.Status()
			n := len(status.Stdout)
			log.Println(status.Stdout[n-1])
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("context done, runner exiting...")
	case <-statusChan:
		// done
	default:
		// no, still running
	}
	// finalStatus := <-statusChan

	log.Println(time.Now().Format("2006-01-02 15:04:05"))

	return nil
}

func RunExec(ChiaExec, LogPath string) {
	LinCmd := strings.Join([]string{`nohup `, ChiaExec, ` > `, LogPath, " 2>&1"}, "")
	cmd := exec.Command("/bin/sh", "-c", LinCmd)
	log.Println(cmd.Args)
	cmd.Start()
	log.Println(cmd.Args)

	pid := cmd.Process.Pid

	// r.activeProcesses[pid] = cmd.Process
	// plotDir.AddPID(pid)
	// farmDir.AddPID(pid)
	// logF("[%d] now plotting. plot dir:%s farm dir:%s\n", pid, plotDir.dirStr, farmDir.dirStr)
	log.Printf("[%d] now plotting. ", pid)
	cmd.Wait()

}
func MakeChiaPlots(confYaml Config) []string {
	ChiaCmd := []string{
		"plots",
		"create",
		"-n", strconv.Itoa(confYaml.NumPlots),
		"-k", strconv.Itoa(confYaml.KSize),
		"-u", strconv.Itoa(confYaml.Buckets),
		"-b", strconv.Itoa(confYaml.Buffer),
		"-r", strconv.Itoa(confYaml.Threads),
		"-f", confYaml.FarmerKey,
		"-p", confYaml.PoolKey,
		"-t", confYaml.TempPath,
		"-2", confYaml.Temp2Path,
		"-d", confYaml.FinalPath,
	}

	if len(confYaml.PlotID) > 0 {
		ChiaCmd = append(ChiaCmd,
			"-i", confYaml.PlotID,
		)
	}
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

func Int2Byte(data int) (ret []byte) {
	var len uintptr = unsafe.Sizeof(data)
	ret = make([]byte, len)
	var tmp int = 0xff
	var index uint = 0
	for index = 0; index < uint(len); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

func Byte2Int(data []byte) int {
	var ret int = 0
	var len int = len(data)
	var i uint = 0
	for i = 0; i < uint(len); i++ {
		ret = ret | (int(data[i]) << (i * 8))
	}
	return ret
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

// func GetCurrentNumber(NumberData string, current int) (n int) {
// 	if IsExist(NumberData) {
// 		number, err := ioutil.ReadFile(NumberData)
// 		if err != nil {
// 			return 0
// 		}
// 		return Byte2Int(number)
// 	} else {
// 		os.Create(NumberData)
// 		number := Int2Byte(current)
// 		ioutil.WriteFile(NumberData, number, 0644)
// 		return current
// 	}
// }

// func WriteCurrentNumber(NumberData string, current int) {
// 	number := Int2Byte(current)
// 	ioutil.WriteFile(NumberData, number, 0644)
// }

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
