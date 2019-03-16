package helmfs

import (
	"realcloud.tech/pligos/cmd/create"
	"realcloud.tech/pligos/pkg/helmfs"
	"realcloud.tech/pligos/pkg/pligos"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type Config struct {
	MountPoint  string
	ConfigPath  string
	ContextName string
}

func HelmFS(config Config) error {
	c, err := fuse.Mount(
		config.MountPoint,
		fuse.FSName("pligos-helmfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("pligos-helmfs"),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	helmCreator := func(path string) error {
		createConfig := create.Config{
			ConfigPath:  config.ConfigPath,
			ContextName: config.ContextName,
			ChartPath:   path,
		}

		return create.Create(createConfig)
	}

	pligosConfig, err := pligos.OpenPligosConfig(config.ConfigPath)
	if err != nil {
		return err
	}

	err = fs.Serve(c, helmfs.New(pligosConfig.DeploymentConfig.Name, helmCreator))
	if err != nil {
		return err
	}

	<-c.Ready
	return c.MountError
}
