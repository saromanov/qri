package cmd

import (
	"os"
	"path/filepath"

	"github.com/qri-io/analytics"
	ipfs "github.com/qri-io/cafs/ipfs"
	"github.com/qri-io/cafs/memfs"
	"github.com/qri-io/qri/api"
	"github.com/qri-io/qri/repo"
	"github.com/qri-io/qri/repo/profile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverCmdPort  string
	serverMemOnly  bool
	serverOffline  bool
	serverInitIpfs bool
)

// serverCmd represents the run command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start a qri server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if serverInitIpfs {
			err := initRepoIfEmpty(viper.GetString(IpfsFsPath), "")
			ExitIfErr(err)
		}

		var (
			r   repo.Repo
			err error
		)

		if serverMemOnly {
			// TODO - refine, adding better identity generation
			// or options for BYO user profile
			r, err = repo.NewMemRepo(
				&profile.Profile{
					Username: "mem user",
				},
				memfs.NewMapstore(),
				repo.MemPeers{},
				&analytics.Memstore{})
			ExitIfErr(err)
		} else {
			r = GetRepo(true)
		}

		s, err := api.New(r, func(cfg *api.Config) {
			cfg.Logger = log
			cfg.Port = serverCmdPort
			cfg.MemOnly = serverMemOnly
			cfg.Online = !serverOffline
			cfg.BoostrapAddrs = viper.GetStringSlice("bootstrap")
		})
		ExitIfErr(err)

		err = s.Serve()
		ExitIfErr(err)
	},
}

func init() {
	serverCmd.Flags().StringVarP(&serverCmdPort, "port", "p", "3000", "port to start server on")
	serverCmd.Flags().BoolVarP(&serverInitIpfs, "init-ipfs", "", false, "initialize a new default ipfs repo if empty")
	serverCmd.Flags().BoolVarP(&serverMemOnly, "mem-only", "", false, "run qri entirely in-memory")
	serverCmd.Flags().BoolVarP(&serverOffline, "offline", "", false, "disable networking")
	RootCmd.AddCommand(serverCmd)
}

func initRepoIfEmpty(repoPath, configPath string) error {
	if repoPath != "" {
		if _, err := os.Stat(filepath.Join(repoPath, "config")); os.IsNotExist(err) {
			if err := os.MkdirAll(repoPath, os.ModePerm); err != nil {
				return err
			}
			return ipfs.InitRepo(repoPath, configPath)
		}
	}
	return nil
}
