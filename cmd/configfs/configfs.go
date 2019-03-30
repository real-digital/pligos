package configfs

import (
	"fmt"
	"log"

	"k8s.io/helm/pkg/proto/hapi/chart"
	"realcloud.tech/pligos/pkg/configfs"
	"realcloud.tech/pligos/pkg/helm"
	"realcloud.tech/pligos/pkg/pligos"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

type Config struct {
	MountPoint  string
	ConfigPath  string
	ContextName string
}

func ConfigFS(config Config) error {
	c, err := fuse.Mount(
		config.MountPoint,
		fuse.FSName("pligos-configfs"),
		fuse.LocalVolume(),
		fuse.VolumeName("pligos-configfs"),
	)
	if err != nil {
		return err
	}
	defer c.Close()

	helmCreator := func() *chart.Chart {
		createConfig, err := pligos.MakeCreateConfig(config.ConfigPath, config.ContextName)
		if err != nil {
			log.Fatalf("read pligos configuration: %v", err)
		}

		creator := &helm.Creator{}
		c, err := creator.Create(createConfig)
		if err != nil {
			log.Fatalf("create: %v", err)
		}

		for _, e := range c.Files {
			fmt.Println(e.GetTypeUrl())
		}

		return c
	}

	pligosConfig, err := pligos.MakePligosConfig(config.ConfigPath)
	if err != nil {
		return err
	}

	err = fs.Serve(c, configfs.New(pligosConfig.DeploymentConfig.Name, helmCreator))
	if err != nil {
		return err
	}

	<-c.Ready
	return c.MountError
}
