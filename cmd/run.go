// Copyright 2025 snowy-jaguar
// Contact: @snowyjaguar (Discord)
// Contact: contact@snowyjaguar.xyz (Email)
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
	"github.com/spf13/cobra"

	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/config"
	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/log"
	"github.com/snowy-jaguar/adguardhomesync-swarm/pkg/sync"
)

// runCmd represents the run command.
var doCmd = &cobra.Command{
	Use:   "run",
	Short: "Start a synchronization from origin to replica",
	Long:  `Synchronizes the configuration form an origin instance to a replica`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logger = log.GetLogger("run")
		cfg, err := config.Get(cfgFile, cmd.Flags())
		if err != nil {
			logger.Error(err)
			return err
		}

		if err := cfg.Init(); err != nil {
			logger.Error(err)
			return err
		}

		if cfg.PrintConfigOnly() {
			if err := cfg.Print(); err != nil {
				logger.Error(err)
				return err
			}

			return nil
		}

		return sync.Sync(cfg.Get())
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
	doCmd.PersistentFlags().String(config.FlagCron, "", "The cron expression to run in daemon mode")
	doCmd.PersistentFlags().Bool(config.FlagRunOnStart, true, "Run the sync job on start.")
	doCmd.PersistentFlags().Bool(config.FlagPrintConfigOnly, false, "Prints the configuration only and exists. "+
		"Can be used to debug the config E.g: when having authentication issues.")
	doCmd.PersistentFlags().Bool(config.FlagContinueOnError, false, "If enabled, the synchronization task "+
		"will not fail on single errors, but will log the errors and continue.")

	doCmd.PersistentFlags().
		Int(config.FlagAPIPort, 8080, "Sync API Port, the API endpoint will be started to enable remote triggering; if 0 port API is disabled.")
	doCmd.PersistentFlags().String(config.FlagAPIUsername, "", "Sync API username")
	doCmd.PersistentFlags().String(config.FlagAPIPassword, "", "Sync API password")
	doCmd.PersistentFlags().String(config.FlagAPIDarkMode, "", "API UI in dark mode")

	doCmd.PersistentFlags().Bool(config.FlagFeatureDhcpServerConfig, true, "Enable DHCP server config feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureDhcpStaticLeases, true, "Enable DHCP server static leases feature")

	doCmd.PersistentFlags().Bool(config.FlagFeatureDNSServerConfig, true, "Enable DNS server config feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureDNSAccessLists, true, "Enable DNS server access lists feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureDNSRewrites, true, "Enable DNS rewrites feature")

	doCmd.PersistentFlags().Bool(config.FlagFeatureGeneral, true, "Enable general settings feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureQueryLog, true, "Enable query log config feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureStats, true, "Enable stats config feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureClient, true, "Enable client settings feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureServices, true, "Enable services sync feature")
	doCmd.PersistentFlags().Bool(config.FlagFeatureFilters, true, "Enable filters sync feature")

	doCmd.PersistentFlags().String(config.FlagOriginURL, "", "Origin instance url")
	doCmd.PersistentFlags().
		String(config.FlagOriginWebURL, "", "Origin instance web url used in the web interface (default: <origin-url>)")
	doCmd.PersistentFlags().String(config.FlagOriginAPIPath, "/control", "Origin instance API path")
	doCmd.PersistentFlags().String(config.FlagOriginUsername, "", "Origin instance username")
	doCmd.PersistentFlags().String(config.FlagOriginPassword, "", "Origin instance password")
	doCmd.PersistentFlags().String(config.FlagOriginCookie, "", "If Set, uses a cookie for authentication")
	doCmd.PersistentFlags().Bool(config.FlagOriginISV, false, "Enable Origin instance InsecureSkipVerify")

	doCmd.PersistentFlags().String(config.FlagReplicaURL, "", "Replica instance url")
	doCmd.PersistentFlags().
		String(config.FlagReplicaWebURL, "", "Replica instance web url used in the web interface (default: <replica-url>)")
	doCmd.PersistentFlags().String(config.FlagReplicaAPIPath, "/control", "Replica instance API path")
	doCmd.PersistentFlags().String(config.FlagReplicaUsername, "", "Replica instance username")
	doCmd.PersistentFlags().String(config.FlagReplicaPassword, "", "Replica instance password")
	doCmd.PersistentFlags().String(config.FlagReplicaCookie, "", "If Set, uses a cookie for authentication")
	doCmd.PersistentFlags().Bool(config.FlagReplicaISV, false, "Enable Replica instance InsecureSkipVerify")
	doCmd.PersistentFlags().
		Bool(config.FlagReplicaAutoSetup, false, "Enable automatic setup of new AdguardHome instances. This replaces the setup wizard.")
	doCmd.PersistentFlags().
		String(config.FlagReplicaInterfaceName, "", "Optional change the interface name of the replica if it differs from the master")
}
