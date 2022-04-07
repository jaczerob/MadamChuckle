package cmd

import (
	"bytes"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jaczerob/madamchuckle/internal/static"

	log "github.com/sirupsen/logrus"
)

var (
	defaultCfg = []byte(`
discord:
	token: ""
redis:
	addr: 127.0.0.1:6379
`)

	configCmd = &cobra.Command{
		Use:   "config",
		Short: "base config management command",
	}

	setCmd = &cobra.Command{
		Use:   "set [key] [value]",
		Short: "sets a key: value pair in the config",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			key, value := args[0], args[1]
			viper.Set(key, value)
			viper.WriteConfig()

			log.WithFields(log.Fields{
				"key":   key,
				"value": value,
			}).Info("set config")
		},
	}
)

func init() {
	configCmd.AddCommand(setCmd)
	rootCmd.AddCommand(configCmd)
}

func initCfg() {
	viper.SetConfigType("yaml")

	if _, err := os.Stat(static.ConfigPath); os.IsNotExist(err) {
		log.WithField("filename", static.ConfigPath).Info("config file not found, creating")

		err = os.MkdirAll(path.Dir(static.ConfigPath), 0700)
		if err != nil {
			log.WithError(err).Fatal("could not create config directory")
		}

		_, err = os.Create(path.Base(static.ConfigPath))
		if err != nil {
			log.WithError(err).Fatal("could not create config file")
		}

		viper.ReadConfig(bytes.NewBuffer(defaultCfg))
		err = viper.WriteConfigAs(static.ConfigPath)
		if err != nil {
			log.WithError(err).Fatal("could not write to config file")
		}
	} else {
		cfgFile, err := os.Open(static.ConfigPath)
		if err != nil {
			log.WithError(err).Fatal("error opening config file")
		}

		viper.ReadConfig(bytes.NewBuffer(defaultCfg))
		viper.MergeConfig(cfgFile)
	}
}
