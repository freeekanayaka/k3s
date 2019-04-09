package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	dqlite "github.com/CanonicalLtd/go-dqlite"
	"github.com/rancher/k3s/pkg/cli/cmds"
	"github.com/urfave/cli"
	yaml "gopkg.in/yaml.v2"
)

// Bootstrap a new node.
func Bootstrap(app *cli.Context) error {
	id := app.Args().Get(0)
	if id == "" {
		return fmt.Errorf("No ID given")
	}
	address := app.Args().Get(1)
	if address == "" {
		return fmt.Errorf("No address given")
	}
	serverID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	info := dqlite.ServerInfo{
		ID:      uint64(serverID),
		Address: address,
	}
	dir := filepath.Join(cmds.BootstrapConfig.DataDir, "server", "db")
	err = os.MkdirAll(dir, 0700)
	if err != nil {
		return err
	}
	server, err := dqlite.NewServer(info, dir)
	if err != nil {
		return err
	}
	infos := make([]dqlite.ServerInfo, len(cmds.BootstrapConfig.Members))
	for i, member := range cmds.BootstrapConfig.Members {
		parts := strings.Split(member, ",")
		if len(parts) != 2 {
			return fmt.Errorf("Bad member format, should be <ID>,<IP>")
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
	data, err := yaml.Marshal(&info)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dir, "info"), data, 0644)
	if err != nil {
		return err
	}
	store, err := dqlite.DefaultServerStore(filepath.Join(dir, "store"))
	if err != nil {
		return err
	}
	err = store.Set(context.Background(), infos)
	if err != nil {
		return err
	}
	return nil
}
