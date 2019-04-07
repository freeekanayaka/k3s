package server

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	dqlite "github.com/CanonicalLtd/go-dqlite"
	"github.com/rancher/k3s/pkg/cli/cmds"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

const port = 6444

// Bootstrap a new node.
func Bootstrap(app *cli.Context) error {
	id := app.Args().Get(0)
	if id == "" {
		return fmt.Errorf("No ID given")
	}
	ip := app.Args().Get(1)
	if ip == "" {
		return fmt.Errorf("No IP given")
	}
	serverID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	info := dqlite.ServerInfo{
		ID:      uint64(serverID),
		Address: fmt.Sprintf("%s:%d", ip, port),
	}
	server, err := dqlite.NewServer(info, cmds.BootstrapConfig.DataDir)
	if err != nil {
		return err
	}
	infos := make([]dqlite.ServerInfo, len(cmds.BootstrapConfig.Members))
	for i, member := range cmds.BootstrapConfig.Members {
		parts := strings.Split(member, ":")
		if len(parts) != 2 {
			return fmt.Errorf("Bad member format, should be <ID>:<IP>")
		}
		memberID, err := strconv.Atoi(parts[0])
		if err != nil {
			return err
		}
		infos[i].ID = uint64(memberID)
		infos[i].Address = parts[1]
	}
	err = server.Bootstrap(infos)
	if err != nil {
		return err
	}
	bytes, err := yaml.Marshal(&info)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(cmds.BootstrapConfig.DataDir, "info"), bytes, 0644)
	if err != nil {
		return err
	}
	return nil
}
