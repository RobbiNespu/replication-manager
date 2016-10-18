// replication-manager - Replication Manager Monitoring and CLI for MariaDB
// Authors: Guillaume Lefranc <guillaume.lefranc@mariadb.com>
//          Stephane Varoqui  <stephane.varoqui@mariadb.com>
// This source code is licensed under the GNU General Public License, version 3.
// Redistribution/Reuse of this code is permitted under the GNU v3 license, as
// an additional term, ALL code must carry the original Author(s) credit in comment form.
// See LICENSE in this directory for the integral text.

package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	loglevel   int
	Version    string
	Build      string
	configfile string
)

func init() {

	rootCmd.AddCommand(versionCmd)
	rootCmd.PersistentFlags().StringVar(&user, "user", "", "User for MariaDB login, specified in the [user]:[password] format")
	rootCmd.PersistentFlags().StringVar(&hosts, "hosts", "", "List of MariaDB hosts IP and port (optional), specified in the host:[port] format and separated by commas")
	rootCmd.PersistentFlags().StringVar(&rpluser, "rpluser", "", "Replication user in the [user]:[password] format")
	rootCmd.Flags().StringVar(&keyPath, "keypath", "/etc/replication-manager/.replication-manager.key", "Encryption key file path")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Print detailed execution info")
	rootCmd.PersistentFlags().IntVar(&loglevel, "log-level", 0, "Log verbosity level")
	rootCmd.PersistentFlags().StringVarP(&configfile, "config", "c", "/etc/replication-manager/config.toml", "Config file")

	//rootCmd.PersistentFlags().StringVar(&config, "config", "test", "configuration default /etc/replication-manager/config.toml")
	viper.BindPFlag("config", rootCmd.Flags().Lookup("config"))
	dir, file := filepath.Split(configfile)
	items := strings.Split(file, ".")
	var filename, fileext string
	if len(items) == 1 {
		filename = items[0]
	} else {
		fileext = items[1]
		filename = items[0]
	}

	viper.SetConfigType(fileext)
	viper.SetConfigName(filename)
	viper.AddConfigPath(dir)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if _, ok := err.(viper.ConfigParseError); ok {
		log.Fatalln("ERROR: Could not parse config file:", err)
	}

	viper.BindPFlags(rootCmd.PersistentFlags())
	user = viper.GetString("user")
	hosts = viper.GetString("hosts")
	rpluser = viper.GetString("rpluser")
	keyPath = viper.GetString("keypath")
	loglevel = viper.GetInt("log-level")

	if verbose == true && loglevel == 0 {
		loglevel = 1
	}
	if verbose == false && loglevel > 0 {
		verbose = true
	}
}

func main() {
	//log.Fatalln("main:", configfile)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

}

var rootCmdV *cobra.Command

var rootCmd = &cobra.Command{
	Use:   "replication-manager",
	Short: "MariaDB Replication Manager Utility",
	Long: `replication-manager allows users to monitor interactively MariaDB 10.x GTID replication health
and trigger slave to master promotion (aka switchover), or elect a new master in case of failure (aka failover).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v\n", args)
		cmd.Usage()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the replication manager version number",
	Long:  `All software has versions. This is ours`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MariaDB Replication Manager version", repmgrVersion)
		fmt.Println("Version: ", Version)
		fmt.Println("Build Time: ", Build)
	},
}
