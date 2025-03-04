package generic

import (
	"context"
	"fmt"
	devspacecontext "github.com/loft-sh/devspace/pkg/devspace/context"
	"github.com/loft-sh/devspace/pkg/util/command"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/loft-sh/devspace/pkg/devspace/config/versions/latest"
	"github.com/loft-sh/devspace/pkg/util/downloader"
	"github.com/loft-sh/devspace/pkg/util/downloader/commands"
	"github.com/loft-sh/devspace/pkg/util/extract"
	"github.com/loft-sh/devspace/pkg/util/log"

	"github.com/pkg/errors"
)

const stableChartRepo = "https://charts.helm.sh/stable"

type Client interface {
	Exec(ctx devspacecontext.Context, args []string) ([]byte, error)
	WriteValues(values map[string]interface{}) (string, error)
}

func NewGenericClient(command commands.Command, log log.Logger) Client {
	c := &client{
		log:     log,
		extract: extract.NewExtractor(),
	}

	c.downloader = downloader.NewDownloader(command, log)
	return c
}

type client struct {
	log        log.Logger
	extract    extract.Extract
	downloader downloader.Downloader

	helmPath string
}

func (c *client) WriteValues(values map[string]interface{}) (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer f.Close()
	out, err := yaml.Marshal(values)
	if err != nil {
		return "", errors.Wrap(err, "marshal values")
	}

	_, err = f.Write(out)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

func (c *client) Exec(ctx devspacecontext.Context, args []string) ([]byte, error) {
	err := c.ensureHelmBinary(ctx.Context())
	if err != nil {
		return nil, err
	}

	if !ctx.KubeClient().IsInCluster() {
		args = append(args, "--kube-context", ctx.KubeClient().CurrentContext())
	}

	// disable log for list, because it prints same command multiple times if we've multiple deployments.
	if args[0] != "list" {
		c.log.Debugf("Execute '%s %s'", c.helmPath, strings.Join(args, " "))
	}

	result, err := command.Output(ctx.Context(), ctx.WorkingDir(), ctx.Environ(), c.helmPath, args...)
	if err != nil {
		return nil, fmt.Errorf("%s %v", string(result), err)
	}

	return result, nil
}

func (c *client) ensureHelmBinary(ctx context.Context) error {
	if c.helmPath != "" {
		return nil
	}

	var err error
	c.helmPath, err = c.downloader.EnsureCommand(ctx)
	return err
}

func ChartNameAndRepo(helmConfig *latest.HelmConfig) (string, string) {
	chartName := strings.TrimSpace(helmConfig.Chart.Name)
	chartRepo := helmConfig.Chart.RepoURL
	if strings.HasPrefix(chartName, "stable/") && chartRepo == "" {
		chartName = chartName[7:]
		chartRepo = stableChartRepo
	}

	return chartName, chartRepo
}
