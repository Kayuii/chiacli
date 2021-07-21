package massbindingtarget

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/massnetorg/mass-core/logging"
	"github.com/massnetorg/mass-core/massutil"
	"github.com/massnetorg/mass-core/poc"
	"github.com/massnetorg/mass-core/poc/chiawallet"
	"github.com/urfave/cli/v2"
)

var (
	getBindingListArgFilename     string
	getBindingListFlagOverwrite   bool
	getBindingListFlagListAll     bool
	getBindingListFlagKeystore    string
	getBindingListFlagDirectories []string
)

func getOfflineBindingList() (list *massutil.BindingList, err error) {
	var absDirectories []string
	for _, dir := range getBindingListFlagDirectories {
		absDir, err := filepath.Abs(dir)
		if err != nil {
			return nil, err
		}
		absDirectories = append(absDirectories, absDir)
	}

	interruptCh := make(chan os.Signal, 2)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)

	logging.CPrint(logging.INFO, "searching for plot files from disk, this may take a while (enter CTRL+C to cancel running)",
		logging.LogFormat{"dir_count": len(absDirectories)})

	var plots []massutil.BindingPlot
	var defaultCount, chiaCount uint64

	plots, err = getOfflineBindingListV2(interruptCh, absDirectories, getBindingListFlagListAll, getBindingListFlagKeystore)
	chiaCount = uint64(len(plots))

	if err != nil {
		logging.CPrint(logging.ERROR, "fail to get offline binding list", logging.LogFormat{"err": err})
		return
	}

	list = &massutil.BindingList{
		Plots:        plots,
		TotalCount:   defaultCount + chiaCount,
		DefaultCount: defaultCount,
		ChiaCount:    chiaCount,
	}
	return list, nil
}

func getOfflineBindingListV2(interruptCh chan os.Signal, dirs []string, all bool, keystoreFile string) ([]massutil.BindingPlot, error) {
	regStrB, suffixB := `^PLOT-K\d{2}-\d{4}(-\d{2}){4}-[A-F0-9]{64}\.PLOT$`, ".PLOT"
	regExpB, err := regexp.Compile(regStrB)
	if err != nil {
		return nil, err
	}

	var keystore *chiawallet.Keystore
	if keystoreFile != "" {
		if keystore, err = chiawallet.NewKeystoreFromFile(keystoreFile); err != nil {
			return nil, err
		}
	}

	var ownablePlot = func(info *massutil.MassDBInfoV2) bool {
		if keystore == nil {
			return true
		}
		if _, err := keystore.GetPoolPrivateKey(info.PoolPublicKey); err != nil {
			return false
		}
		if _, err := keystore.GetFarmerPrivateKey(info.FarmerPublicKey); err != nil {
			return false
		}
		return true
	}

	var plots []massutil.BindingPlot
	var totalSearched int

	for _, dbDir := range dirs {
		dirFileInfos, err := ioutil.ReadDir(dbDir)
		if err != nil {
			return nil, err
		}

		logging.CPrint(logging.INFO, "searching for chia plot files", logging.LogFormat{"dir": dbDir})

		dirSearched := 0
		for _, fi := range dirFileInfos {
			select {
			case <-interruptCh:
				logging.CPrint(logging.WARN, "cancel searching plot files")
				return nil, nil
			default:
			}

			fileName := fi.Name()
			if !strings.HasSuffix(strings.ToUpper(fileName), suffixB) || !regExpB.MatchString(strings.ToUpper(fileName)) {
				continue
			}

			info, err := massutil.NewMassDBInfoV2FromFile(filepath.Join(dbDir, fileName))
			if err != nil {
				logging.CPrint(logging.WARN, "fail to read chia plot info", logging.LogFormat{"err": err})
				continue
			}

			if !ownablePlot(info) {
				continue
			} else {
				target, err := massutil.GetChiaPlotBindingTarget(info.PlotID, info.K)
				if err != nil {
					return nil, err
				}
				plots = append(plots, massutil.BindingPlot{
					Target: target,
					Type:   uint8(poc.ProofTypeChia),
					Size:   uint8(info.K),
				})
				dirSearched += 1
			}
		}

		logging.CPrint(logging.INFO, "loaded chia plot files from directory", logging.LogFormat{
			"dir":      dbDir,
			"db_count": dirSearched,
		})
		totalSearched += dirSearched
	}

	logging.CPrint(logging.INFO, "loaded chia plot files from all directories", logging.LogFormat{
		"dir_count":      len(dirs),
		"total_db_count": totalSearched,
	})

	return plots, err
}

func Target(c *cli.Context) error {
	if c.NArg() < 1 {
		return cli.ShowAppHelp(c)
	}
	abs, err := filepath.Abs(c.Args().First())
	if err != nil {
		logging.CPrint(logging.ERROR, "wrong filename format", logging.LogFormat{"err": err, "filename": c.Args().First()})
		return err
	}
	fi, err := os.Stat(abs)
	if err == nil && fi.IsDir() {
		logging.CPrint(logging.ERROR, "filename is a directory", logging.LogFormat{"filename": c.Args().First()})
		return err
	}
	getBindingListArgFilename = abs
	getBindingListFlagOverwrite = c.Bool("overwrite")
	getBindingListFlagListAll = c.Bool("all")
	getBindingListFlagKeystore = c.String("keystore")
	getBindingListFlagDirectories = c.StringSlice("dirs")

	// main logics
	_, err = os.Stat(getBindingListArgFilename)
	if !os.IsNotExist(err) && !getBindingListFlagOverwrite {
		logging.CPrint(logging.ERROR, "cannot overwrite existed file, try again with --overwrite", logging.LogFormat{
			"filename": getBindingListArgFilename,
		})
		return err
	}

	list, err := getOfflineBindingList()
	if err != nil {
		logging.CPrint(logging.ERROR, "fail to get binding list", logging.LogFormat{"err": err})
		return err
	}
	list = list.RemoveDuplicate()

	if len(list.Plots) == 0 {
		fmt.Println("saved nothing in the binding list")
		return nil
	}

	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		logging.CPrint(logging.ERROR, "fail to marshal json", logging.LogFormat{
			"err":         err,
			"total_count": list.TotalCount,
		})
		return err
	}

	if err = ioutil.WriteFile(getBindingListArgFilename, data, 0666); err != nil {
		logging.CPrint(logging.ERROR, "fail to write into binding list file", logging.LogFormat{
			"err":         err,
			"total_count": list.TotalCount,
			"byte_size":   len(data),
		})
		return err
	}

	fmt.Printf("collected %d plot files.\n", list.TotalCount)
	return nil
}
