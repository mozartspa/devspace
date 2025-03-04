package commands

import (
	"context"
	"github.com/loft-sh/devspace/pkg/devspace/config/constants"
	"github.com/loft-sh/devspace/pkg/util/command"
	"github.com/loft-sh/devspace/pkg/util/extract"
	"github.com/mitchellh/go-homedir"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"mvdan.cc/sh/v3/expand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	helmVersion  = "v3.7.1"
	helmDownload = "https://get.helm.sh/helm-" + helmVersion + "-" + runtime.GOOS + "-" + runtime.GOARCH
)

func NewHelmV3Command() Command {
	return &helmv3{}
}

type helmv3 struct{}

func (h *helmv3) Name() string {
	return "helm"
}

func (h *helmv3) InstallPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	installPath := filepath.Join(home, constants.DefaultHomeDevSpaceFolder, "bin", h.Name())
	if runtime.GOOS == "windows" {
		installPath += ".exe"
	}

	return installPath, nil
}

func (h *helmv3) DownloadURL() string {
	url := helmDownload + ".tar.gz"
	if runtime.GOOS == "windows" {
		url = helmDownload + ".zip"
	}

	return url
}

func (h *helmv3) IsValid(ctx context.Context, path string) (bool, error) {
	out, err := command.Output(ctx, "", expand.ListEnviron(os.Environ()...), path, "version")
	if err != nil {
		return false, nil
	}

	return strings.Contains(string(out), `:"v3.`), nil
}

func (h *helmv3) Install(archiveFile string) error {
	installPath, err := h.InstallPath()
	if err != nil {
		return err
	}

	return installHelmBinary(extract.NewExtractor(), archiveFile, installPath, h.DownloadURL())
}

func installHelmBinary(extract extract.Extract, archiveFile, installPath, installFromURL string) error {
	t := filepath.Dir(archiveFile)

	// Extract the binary
	if strings.HasSuffix(installFromURL, ".tar.gz") {
		err := extract.UntarGz(archiveFile, t)
		if err != nil {
			return errors.Wrap(err, "extract tar.gz")
		}
	} else if strings.HasSuffix(installFromURL, ".zip") {
		err := extract.Unzip(archiveFile, t)
		if err != nil {
			return errors.Wrap(err, "extract zip")
		}
	}

	// Copy file to target location
	if runtime.GOOS == "windows" {
		return copy.Copy(filepath.Join(t, runtime.GOOS+"-"+runtime.GOARCH, "helm.exe"), installPath)
	}

	return copy.Copy(filepath.Join(t, runtime.GOOS+"-"+runtime.GOARCH, "helm"), installPath)
}
