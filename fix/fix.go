package fix

import "fmt"

type Fix struct {
	dirName string
	LogPath string
	PlotID  string
	LogFile string
}

func New() *Fix {
	return &Fix{}
}

type Config struct {
	FarmerKey string `yaml:"FarmerKey"`
	PoolKey   string `yaml:"PoolKey"`
	LocalSk   string `yaml:"LocalSk"`
	Memo      string `yaml:"Memo"`
}

func (f *Fix) Print(config *Config) error {
	var (
		ppk  string
		fpk  string
		sk   string
		memo string = config.Memo
	)
	if len(config.Memo) > 0 {

		switch len(config.Memo) {
		case 128:
			fpk = memo[:48]
			ppk = memo[48:96]
			sk = memo[96:128]
			break
		case 112:
			break
		default:
			fmt.Println("Invalid memo err")
			return nil
		}

		fmt.Printf("memo: %s \n", memo)
		fmt.Printf("ppk: %s \n", ppk)
		fmt.Printf("fpk: %s \n", fpk)
		fmt.Printf("sk: %s \n", sk)
	}

	return nil
}
