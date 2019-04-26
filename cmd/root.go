// Copyright © 2019 Sai Kothapalle <ephemeral972@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "shr",
	Short: "shr cmd line tool to export system's user information",
	Long:  `shr command tool exports server's usernames, IDs, home directories as JSON`,

	Run: func(cmd *cobra.Command, args []string) {
		path, _ := cmd.Flags().GetString("path")
		if path == "" {
			path = "/users/"
		}
		fmt.Println("Exporting server's User info to the file" + path)
		getUserInfo(path)
		fmt.Println("User Info export is completed, please check the file")
	},
}

type hr struct {
	User  string `json:"user"`
	ID    string `json:"id"`
	Home  string `json:"home"`
	Shell string `json:"shell"`
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.shr.json)")
	rootCmd.PersistentFlags().StringP("format", "f", "json", "User data export format")
	rootCmd.PersistentFlags().StringP("path", "p", "", "file path to export user data")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".shr" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".shr")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

// getUserInfo gets the server's user stats
func getUserInfo(path string) {
	// get the user stats
	f, err := os.Create(path)
	check(err)
	defer f.Close()
	cmd := exec.Command("cat", "/etc/passwd")
	out, err := cmd.CombinedOutput()
	check(err)
	bhr := collectUsers(string(out))
	_, err = f.WriteString(string(bhr))
	check(err)
}

// mangle sifts through the command output to get in req format
func collectUsers(s string) []byte {
	var hrs []hr
	userSlice := strings.Split(s, "\n")
	// iterate through the array
	for i := 0; i < len(userSlice)-1; i++ {
		ui := strings.Split(userSlice[i], ":")
		sysID, err := strconv.Atoi(ui[2])
		if err != nil {
			check(err)
		}
		if sysID > 1000 {
			shr := hr{ui[0], ui[2], ui[5], ui[6]}
			hrs = append(hrs, shr)
		}
	}
	rd, _ := json.MarshalIndent(hrs, "", " ")
	return rd
}

func check(e error) {
	if e != nil {
		fmt.Println(e)
	}
	return
}
