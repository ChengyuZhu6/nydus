// Copyright 2023 Nydus Developers. All rights reserved.
//
// SPDX-License-Identifier: Apache-2.0

package tool

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

type BinaryContext struct {
	Builder                      string
	Nydusd                       string
	Nydusify                     string
	NydusifyChecker              string
	NydusifyOnlySupportV5        bool
	NydusifyNotSupportCompressor bool
}

type BuildContext struct {
	FSVersion  string
	Compressor string
	ChunkSize  string
	OCIRef     bool
	OCIRefGzip bool
}

type RuntimeContext struct {
	CacheType       string
	CacheDir        string
	MountDir        string
	CacheCompressed bool
	RafsMode        string
	EnablePrefetch  bool
}

type EnvContext struct {
	WorkDir       string
	BlobDir       string
	CacheDir      string
	MountDir      string
	BootstrapPath string
}

type Context struct {
	Binary  BinaryContext
	Build   BuildContext
	Runtime RuntimeContext
	Env     EnvContext
}

func DefaultContext() *Context {
	return &Context{
		Binary: BinaryContext{
			Builder:               "nydus-image",
			Nydusd:                "nydusd",
			Nydusify:              "nydusify",
			NydusifyChecker:       "nydusify",
			NydusifyOnlySupportV5: false,
		},
		Build: BuildContext{
			FSVersion:  "6",
			Compressor: "zstd",
			ChunkSize:  "0x100000",
		},
		Runtime: RuntimeContext{
			CacheType:       "blobcache",
			CacheCompressed: false,
			RafsMode:        "direct",
			EnablePrefetch:  true,
		},
	}
}

func (ctx *Context) PrepareWorkDir(t *testing.T) {
	tempDir := os.Getenv("WORK_DIR")
	if tempDir == "" {
		tempDir = os.TempDir()
	}
	workDir, err := ioutil.TempDir(tempDir, "nydus-smoke-")
	require.NoError(t, err)

	blobDir := filepath.Join(workDir, "blobs")
	err = os.MkdirAll(blobDir, 0755)
	require.NoError(t, err)
	cacheDir := filepath.Join(workDir, "cache")
	err = os.MkdirAll(cacheDir, 0755)
	require.NoError(t, err)
	mountDir := filepath.Join(workDir, "mnt")
	err = os.MkdirAll(mountDir, 0755)
	require.NoError(t, err)

	ctx.Env = EnvContext{
		WorkDir:  workDir,
		BlobDir:  blobDir,
		CacheDir: cacheDir,
		MountDir: mountDir,
	}
}

func (ctx *Context) Destroy(t *testing.T) {
	os.RemoveAll(ctx.Env.WorkDir)
}