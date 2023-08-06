package config

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
)

func TestInitConfig(t *testing.T) {
	var configDir = "../"
	absDir, _ := filepath.Abs(configDir)

	fmt.Print(absDir)
}

func TestDbDnsInit(t *testing.T) {
	InitConfig()
	fmt.Printf("host is : %v \n", viper.GetString("mysql.host"))
	dns := DbDnsInit()
	fmt.Printf("dns is : %v", dns)
}
