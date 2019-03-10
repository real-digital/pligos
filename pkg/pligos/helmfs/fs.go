package helmfs

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	yaml "gopkg.in/yaml.v2"
	"realcloud.tech/cloud-tools/pkg/pligos/helm"
)

func New(chartName string, helmCreator func() *helm.Helm) *FileSystem {
	return &FileSystem{
		helmCreator: helmCreator,
		chartName:   chartName,
	}
}

type FileSystem struct {
	helmCreator func() *helm.Helm
	chartName   string
}

func (f *FileSystem) Root() (fs.Node, error) {
	return &Dir{isConfigRoot: true, FileSystem: f}, nil
}

func (f *FileSystem) buildConfigTree() (map[string]interface{}, error) {
	dir, err := ioutil.TempDir("", "pligs-helmfs")
	if err != nil {
		return nil, err
	}

	if err := f.helmCreator().Create(filepath.Join(dir, f.chartName)); err != nil {
		return nil, err
	}

	resources, err := template(filepath.Join(dir, f.chartName), release)
	if err != nil {
		return nil, err
	}

	if err := os.RemoveAll(dir); err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	for k, e := range resources {
		if !isNthLevel(k, 3) {
			continue
		}

		manifests := strings.Split(e, "---")

		for _, m := range manifests {
			kind, err := findConfigKind(m)
			if err != nil {
				return nil, err
			}

			if kind == unknown {
				continue
			}

			name, err := findResourceName(kind, m)
			if err != nil {
				return nil, err
			}

			var dataProbe struct {
				Data map[string]interface{} `yaml:"data"`
			}

			if err := yaml.Unmarshal([]byte(m), &dataProbe); err != nil {
				return nil, err
			}

			if kind == secret {
				dec, err := decode(dataProbe.Data)
				if err != nil {
					return nil, err
				}

				res[name] = dec
				continue
			}

			res[name] = dataProbe.Data
		}
	}

	return res, nil
}

func isNthLevel(path string, n int) bool {
	for i := 0; i < n; i++ {
		path = filepath.Dir(path)
	}

	return path == "."
}

func decode(secrets map[string]interface{}) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	for k, v := range secrets {
		n, err := base64.StdEncoding.DecodeString(v.(string))
		if err != nil {
			return nil, err
		}

		res[k] = string(n)
	}

	return res, nil
}

type Dir struct {
	*FileSystem

	isConfigRoot bool
	name         string
}

func (*Dir) Attr(ctx context.Context, a *fuse.Attr) error {
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

const release = "configuration"

func (d *Dir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var res []fuse.Dirent

	tree, err := d.buildConfigTree()
	if err != nil {
		return nil, err
	}

	if d.isConfigRoot {
		for e := range tree {
			res = append(
				res,
				fuse.Dirent{
					Type: fuse.DT_Dir,
					Name: e,
				},
			)
		}

		return res, nil
	}

	for e := range tree[d.name].(map[string]interface{}) {
		res = append(res, fuse.Dirent{
			Type: fuse.DT_File,
			Name: e,
		})
	}

	return res, nil
}

func (d *Dir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	if d.isConfigRoot {
		return &Dir{
			FileSystem: &FileSystem{chartName: d.chartName, helmCreator: d.helmCreator},

			name: name,
		}, nil
	}

	tree, err := d.buildConfigTree()
	if err != nil {
		return nil, err
	}

	if _, ok := tree[d.name]; !ok {
		return nil, fuse.ENOENT
	}

	if _, ok := tree[d.name].(map[string]interface{})[name]; ok {
		return &File{
			FileSystem: &FileSystem{chartName: d.chartName, helmCreator: d.helmCreator},

			manifest: d.name,
			name:     name,
		}, nil
	}

	return nil, fuse.ENOENT
}

type File struct {
	*FileSystem

	manifest, name string
}

func (f *File) Attr(ctx context.Context, a *fuse.Attr) error {
	tree, err := f.buildConfigTree()
	if err != nil {
		return err
	}

	a.Inode = 2
	a.Mode = 0444
	a.Size = uint64(len(tree[f.manifest].(map[string]interface{})[f.name].(string)))
	return nil
}

func (f *File) ReadAll(ctx context.Context) ([]byte, error) {
	tree, err := f.buildConfigTree()
	if err != nil {
		return nil, err
	}

	return []byte(tree[f.manifest].(map[string]interface{})[f.name].(string)), nil
}
