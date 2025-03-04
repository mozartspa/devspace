package latest

import (
	"strings"

	"github.com/loft-sh/devspace/pkg/util/yamlutil"

	"encoding/json"

	"github.com/loft-sh/devspace/pkg/devspace/config/versions/config"
	"gopkg.in/yaml.v3"
	k8sv1 "k8s.io/api/core/v1"
)

// Version is the current api version
const Version string = "v2beta1"

// GetVersion returns the version
func (c *Config) GetVersion() string {
	return Version
}

// New creates a new config object
func New() config.Config {
	return NewRaw()
}

// NewRaw creates a new config object
func NewRaw() *Config {
	return &Config{
		Version: Version,
		Images:  map[string]*Image{},
	}
}

func (c *Config) Clone() *Config {
	out, _ := yaml.Marshal(c)
	n := &Config{}
	_ = yamlutil.Unmarshal(out, n)
	return n
}

// Config defines the configuration
type Config struct {
	// Version holds the config version. DevSpace will always convert older configs to the current latest
	// config version, which makes it possible to use the newest DevSpace version also with older config
	// versions.
	Version string `yaml:"version" json:"version" jsonschema:"required" jsonschema_description:"Version holds the config version. DevSpace will always convert older configs to the current latest config version, which makes it possible to use the newest DevSpace version also with older config versions."`

	// Name specifies the name of the DevSpace project and uniquely identifies a project.
	// DevSpace will not allow multiple active projects with the same name in the same Kubernetes namespace.
	Name string `yaml:"name" json:"name" jsonschema:"required" jsonschema_description:"Name specifies the name of the DevSpace project and uniquely identifies a project."`

	// Imports merges specified config files into this one. This is very useful to split up your DevSpace configuration
	// into multiple files and reuse those through git, a remote url or common local path.
	Imports []Import `yaml:"imports,omitempty" json:"imports,omitempty" jsonschema_description:"Imports merges specified config files into this one. This is very useful to split up your DevSpace configuration into multiple files and reuse those through git, a remote url or common local path."`

	// Functions are POSIX functions that can be used within pipelines. Those functions can also be imported by
	// imports.
	Functions map[string]string `yaml:"functions,omitempty" json:"functions,omitempty" jsonschema_description:"Functions are POSIX functions that can be used within pipelines. Those functions can also be imported by imports."`

	// Pipelines are the work blocks that DevSpace should execute when devspace dev, devspace build, devspace deploy or devspace purge
	// is called. Pipelines are defined through a special POSIX script that allows you to use special commands
	// such as create_deployments, start_dev, build_images etc. to signal DevSpace you want to execute
	// a specific functionality. The pipelines dev, build, deploy and purge are special and will override
	// the default functionality of the respective command if defined. All other pipelines can be either run
	// via the devspace run-pipeline command or used within another pipeline through run_pipelines.
	Pipelines map[string]*Pipeline `yaml:"pipelines,omitempty" json:"pipelines,omitempty" jsonschema_description:" Pipelines are the work blocks that DevSpace should execute when devspace dev, devspace build, devspace deploy or devspace purge is called."`

	// Images holds configuration of how DevSpace should build images. By default, DevSpace will build all defined images.
	// If you are using a custom pipeline, you can dynamically define which image is built at which time during the
	// execution.
	Images map[string]*Image `yaml:"images,omitempty" json:"images,omitempty" jsonschema_description:"Images holds configuration of how DevSpace should build images. By default, DevSpace will build all defined images."`

	// Deployments holds configuration of how DevSpace should deploy resources to Kubernetes. By default, DevSpace will deploy all defined deployments.
	// If you are using a custom pipeline, you can dynamically define which deployment is deployed at which time during the
	// execution.
	Deployments map[string]*DeploymentConfig `yaml:"deployments,omitempty" json:"deployments,omitempty" jsonschema_description:"Deployments holds configuration of how DevSpace should deploy resources to Kubernetes. By default, DevSpace will deploy all defined deployments."`

	// Dev holds development configuration. Each dev configuration targets a single pod and enables certain dev services on that pod
	// or even rewrites it if certain changes are requested, such as adding an environment variable or changing the entrypoint.
	// Dev allows you to:
	// - sync local folders to the Kubernetes pod
	// - port forward remote ports to your local computer
	// - forward local ports into the Kubernetes pod
	// - configure an ssh tunnel to the Kubernetes pod
	// - proxy local commands to the container
	// - restart the container on file changes
	Dev map[string]*DevPod `yaml:"dev,omitempty" json:"dev,omitempty" jsonschema_description:"Dev holds development configuration. Each dev configuration targets a single pod and enables certain dev services on that pod or even rewrites it if certain changes are requested, such as adding an environment variable or changing the entrypoint."`

	// Vars are config variables that can be used inside other config sections to replace certain values dynamically
	Vars map[string]*Variable `yaml:"vars,omitempty" json:"vars,omitempty" jsonschema_description:"Vars are config variables that can be used inside other config sections to replace certain values dynamically."`

	// Commands are custom commands that can be executed via 'devspace run COMMAND'. These commands are run within a pseudo bash
	// that also allows executing special commands such as run_watch or is_equal.
	Commands map[string]*CommandConfig `yaml:"commands,omitempty" json:"commands,omitempty" jsonschema_description:"Commands are custom commands that can be executed via 'devspace run COMMAND'. These commands are run within a pseudo bash that also allows executing special commands such as run_watch or is_equal."`

	// Dependencies are sub devspace projects that lie in a local folder or remote git repository that can be executed
	// from within the pipeline. In contrast to imports, these projects pose as separate fully functional DevSpace projects
	// that typically lie including source code in a different folder and can be used to compose a full microservice
	// application that will be deployed by DevSpace. Each dependency name can only be used once and if you want to use
	// the same project multiple times, make sure to use a different name for each of those instances.
	Dependencies map[string]*DependencyConfig `yaml:"dependencies,omitempty" json:"dependencies,omitempty" jsonschema_description:"Dependencies are sub devspace projects that lie in a local folder or remote git repository that can be executed from within the pipeline."`

	// PullSecrets are image pull secrets that will be created by devspace in the target namespace
	// during devspace dev or devspace deploy. DevSpace will merge all defined pull secrets into a single
	// one or the one specified.
	PullSecrets map[string]*PullSecretConfig `yaml:"pullSecrets,omitempty" json:"pullSecrets,omitempty" jsonschema_description:"PullSecrets are image pull secrets that will be created by devspace in the target namespace during devspace dev or devspace deploy. DevSpace will merge all defined pull secrets into a single one or the one specified."`

	// Require defines what DevSpace, plugins and command versions are required to use this config and if a condition is not
	// fulfilled, DevSpace will fail.
	Require RequireConfig `yaml:"require,omitempty" json:"require,omitempty" jsonschema_description:"Require defines what DevSpace, plugins and command versions are required to use this config and if a condition is not fulfilled, DevSpace will fail."`

	// Profiles can be used to change the current configuration and change the behavior of devspace. They are deprecated and
	// imports should be used instead.
	Profiles []*ProfileConfig `yaml:"profiles,omitempty" json:"profiles,omitempty" jsonschema:"-" jsonschema_description:"Profiles can be used to change the current configuration and change the behavior of devspace. They are deprecated and imports should be used instead."`

	// Hooks are actions that are executed at certain points within the pipeline. Hooks are ordered and are executed
	// in the order they are specified. They are deprecated and pipelines should be used instead.
	Hooks []*HookConfig `yaml:"hooks,omitempty" json:"hooks,omitempty" jsonschema:"-" jsonschema_description:"Hooks are actions that are executed at certain points within the pipeline. Hooks are ordered and are executed in the order they are specified. They are deprecated and pipelines should be used instead."`
}

// Import specifies the source of the devspace config to merge
type Import struct {
	// Enabled specifies if the given import should be enabled
	Enabled *bool `yaml:"enabled,omitempty" json:"enabled,omitempty" jsonschema:"required" jsonschema_description:"Enabled specifies if the given import should be enabled."`

	// SourceConfig defines the source for this import
	SourceConfig `yaml:",inline" json:",inline" jsonschema_description:"SourceConfig defines the source for this import."`
}

// Pipeline defines what DevSpace should do. A pipeline consists of one or more
// jobs that are run in parallel and can depend on each other. Each job consists
// of one or more conditional steps that are executed in order.
type Pipeline struct {
	// Name of the pipeline, will be filled automatically
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"enum=dev,enum=deploy,enum=build,enum=purge,enum=.*" jsonschema_description:"Name of the pipeline, will be filled automatically."`

	// Run is the actual shell command that should be executed during this pipeline
	Run string `yaml:"run,omitempty" json:"run,omitempty" jsonschema:"required" jsonschema_description:"Run is the actual shell command that should be executed during this pipeline."`

	// Flags are extra flags that can be used for running the pipeline via
	// devspace run-pipeline.
	Flags []PipelineFlag `yaml:"flags,omitempty" json:"flags,omitempty" jsonschema_description:"Flags are extra flags that can be used for running the pipeline via devspace run-pipeline."`

	// ContinueOnError will not fail the whole job and pipeline if
	// a call within the step fails.
	ContinueOnError bool `yaml:"continueOnError,omitempty" json:"continueOnError,omitempty" jsonschema_description:"ContinueOnError will not fail the whole job and pipeline if a call within the step fails."`
}

// PipelineFlag defines an extra pipeline flag
type PipelineFlag struct {
	// Name is the name of the flag
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name is the name of the flag."`

	// Short is optional and is the shorthand name for this flag. E.g. 'g' converts to '-g'
	Short string `yaml:"short,omitempty" json:"short,omitempty" jsonschema_description:"Short is optional and is the shorthand name for this flag. E.g. 'g' converts to '-g'."`

	// Type is the type of the flag. Defaults to `bool`
	Type PipelineFlagType `yaml:"type,omitempty" json:"type,omitempty" jsonschema:"enum=bool,enum=int,enum=string,enum=stringArray" jsonschema_description:"Type is the type of the flag. Defaults to 'bool'."`

	// Default is the default value for this flag
	Default interface{} `yaml:"default,omitempty" json:"default,omitempty" jsonschema_description:"Default is the default value for this flag."`

	// Description is the description as shown in `devspace run-pipeline my-pipe -h`
	Description string `yaml:"description,omitempty" json:"description,omitempty" jsonschema_description:"Description is the description as shown in 'devspace run-pipeline my-pipe -h'."`
}

type PipelineFlagType string

const (
	PipelineFlagTypeString      = "string"
	PipelineFlagTypeBoolean     = "bool"
	PipelineFlagTypeInteger     = "int"
	PipelineFlagTypeStringArray = "stringArray"
)

func (p *Pipeline) UnmarshalYAML(unmarshal func(interface{}) error) error {
	pipelineString := ""
	err := unmarshal(&pipelineString)
	if err != nil {
		m := map[string]interface{}{}
		err := unmarshal(m)
		if err != nil {
			return err
		}

		out, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return yamlutil.UnmarshalStrictJSON(out, p)
	}

	p.Run = pipelineString
	return nil
}

type RequireConfig struct {
	// DevSpace specifies the DevSpace version constraint that is needed to use this config
	DevSpace string `yaml:"devspace,omitempty" json:"devspace,omitempty" jsonschema_description:"DevSpace specifies the DevSpace version constraint that is needed to use this config."`

	// Commands specifies an array of commands that need to be installed locally to use this config
	Commands []RequireCommand `yaml:"commands,omitempty" json:"commands,omitempty" jsonschema_description:"Commands specifies an array of commands that need to be installed locally to use this config."`

	// Plugins specifies an array of plugins that need to be installed locally
	Plugins []RequirePlugin `yaml:"plugins,omitempty" json:"plugins,omitempty" jsonschema_description:"Plugins specifies an array of plugins that need to be installed locally."`
}

type RequirePlugin struct {
	// Name of the plugin that should be installed
	Name string `yaml:"name" json:"name" jsonschema:"required" jsonschema_description:"Name of the plugin that should be installed."`

	// Version constraint of the plugin that should be installed
	Version string `yaml:"version" json:"version" jsonschema:"required" jsonschema_description:"Version constraint of the plugin that should be installed."`
}

type RequireCommand struct {
	// Name is the name of the command that should be installed
	Name string `yaml:"name" json:"name" jsonschema:"required" jsonschema_description:"Name is the name of the command that should be installed."`

	// VersionArgs are the arguments to retrieve the version of the command
	VersionArgs []string `yaml:"versionArgs,omitempty" json:"versionArgs,omitempty" jsonschema_description:"VersionArgs are the arguments to retrieve the version of the command."`

	// VersionRegEx is the regex that is used to parse the version
	VersionRegEx string `yaml:"versionRegEx,omitempty" json:"versionRegEx,omitempty" jsonschema_description:"VersionRegEx is the regex that is used to parse the version."`

	// Version constraint of the command that should be installed
	Version string `yaml:"version,omitempty" json:"version,omitempty" jsonschema_description:"Version constraint of the command that should be installed."`
}

// Image defines the image specification
type Image struct {
	// Name of the image, will be filled automatically
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name of the image, will be filled automatically."`

	// Image is the complete image name including registry and repository
	// for example myregistry.com/mynamespace/myimage
	Image string `yaml:"image" json:"image" jsonschema:"required" jsonschema_description:"Image is the complete image name including registry and repository for example myregistry.com/mynamespace/myimage."`

	// Tags is an array that specifies all tags that should be build during
	// the build process. If this is empty, devspace will generate a random tag
	Tags []string `yaml:"tags,omitempty" json:"tags,omitempty" jsonschema_description:"Tags is an array that specifies all tags that should be build during the build process. If this is empty, devspace will generate a random tag."`

	// Dockerfile specifies a path (relative or absolute) to the dockerfile. Defaults
	// to ./Dockerfile
	Dockerfile string `yaml:"dockerfile" json:"dockerfile" jsonschema:"default=./Dockerfile" jsonschema_extras:"group=buildConfig" jsonschema_description:"Dockerfile specifies a path (relative or absolute) to the dockerfile. Defaults to ./Dockerfile."`

	// Context is the context path to build with. Defaults to the current working directory
	Context string `yaml:"context,omitempty" json:"context,omitempty" jsonschema:"default=./" jsonschema_extras:"group=buildConfig" jsonschema_description:"Context is the context path to build with. Defaults to the current working directory."`

	// Entrypoint specifies an entrypoint that will be appended to the dockerfile during
	// image build in memory. Example: ["sleep", "99999"]
	Entrypoint []string `yaml:"entrypoint,omitempty" json:"entrypoint,omitempty" jsonschema_extras:"group=overwrites,group_name=In-Memory Overwrites" jsonschema_description:"Entrypoint specifies an entrypoint that will be appended to the dockerfile during image build in memory. Example: [\"sleep\", \"99999\"]."`

	// Cmd specifies the arguments for the entrypoint that will be appended
	// during build in memory to the dockerfile
	Cmd []string `yaml:"cmd,omitempty" json:"cmd,omitempty" jsonschema_extras:"group=overwrites" jsonschema_description:"Cmd specifies the arguments for the entrypoint that will be appended during build in memory to the dockerfile."`

	// AppendDockerfileInstructions are instructions that will be appended to the Dockerfile that is build
	// at the current build target and are appended before the entrypoint and cmd instructions
	AppendDockerfileInstructions []string `yaml:"appendDockerfileInstructions,omitempty" json:"appendDockerfileInstructions,omitempty" jsonschema_extras:"group=overwrites" jsonschema_description:"AppendDockerfileInstructions are instructions that will be appended to the Dockerfile that is build at the current build target and are appended before the entrypoint and cmd instructions."`

	// BuildArgs are the build args that are to the build
	BuildArgs map[string]*string `yaml:"buildArgs,omitempty" json:"buildArgs,omitempty" jsonschema_extras:"group=buildConfig,group_name=Build Configuration" jsonschema_description:"BuildArgs are the build args that are to the build."`

	// Target is the target that should get used during the build. Only works if the dockerfile supports this
	Target string `yaml:"target,omitempty" json:"target,omitempty" jsonschema_extras:"group=buildConfig" jsonschema_description:"Target is the target that should get used during the build. Only works if the dockerfile supports this."`

	// Network is the network that should get used to build the image
	Network string `yaml:"network,omitempty" json:"network,omitempty" jsonschema_extras:"group=buildConfig" jsonschema_description:"Network is the network that should get used to build the image."`

	// RebuildStrategy is used to determine when DevSpace should rebuild an image. By default, devspace will
	// rebuild an image if one of the following conditions is true:
	// - The dockerfile has changed
	// - The configuration within the devspace.yaml for the image has changed
	// - A file within the docker context (excluding .dockerignore rules) has changed
	// This option is ignored for custom builds.
	RebuildStrategy RebuildStrategy `yaml:"rebuildStrategy,omitempty" json:"rebuildStrategy,omitempty" jsonschema:"enum=default,enum=always,enum=ignoreContextChanges" jsonschema_extras:"group=buildConfig" jsonschema_description:"RebuildStrategy is used to determine when DevSpace should rebuild an image."`

	// SkipPush will not push the image to a registry if enabled. Only works if docker or buildkit is chosen
	// as build method
	SkipPush bool `yaml:"skipPush,omitempty" json:"skipPush,omitempty" jsonschema_extras:"group=pushPull,group_name=Push & Pull" jsonschema_description:"SkipPush will not push the image to a registry if enabled. Only works if docker or buildkit is chosen as build method."`

	// CreatePullSecret specifies if a pull secret should be created for this image in the
	// target namespace. Defaults to true
	CreatePullSecret *bool `yaml:"createPullSecret,omitempty" json:"createPullSecret,omitempty" jsonschema:"required" jsonschema_extras:"group=pushPull" jsonschema_description:"CreatePullSecret specifies if a pull secret should be created for this image in the target namespace. Defaults to true."`

	// BuildKit if buildKit is specified, DevSpace will build the image either in-cluster or locally with BuildKit
	BuildKit *BuildKitConfig `yaml:"buildKit,omitempty" json:"buildKit,omitempty" jsonschema_extras:"group=engines,group_name=Build Engines" jsonschema_description:"BuildKit if buildKit is specified, DevSpace will build the image either in-cluster or locally with BuildKit."`

	// Docker if docker is specified, DevSpace will build the image using the local docker daemon
	Docker *DockerConfig `yaml:"docker,omitempty" json:"docker,omitempty" jsonschema_extras:"group=engines" jsonschema_description:"Docker if docker is specified, DevSpace will build the image using the local docker daemon."`

	// Kaniko if kaniko is specified, DevSpace will build the image in-cluster with kaniko
	Kaniko *KanikoConfig `yaml:"kaniko,omitempty" json:"kaniko,omitempty" jsonschema_extras:"group=engines" jsonschema_description:"Kaniko if kaniko is specified, DevSpace will build the image in-cluster with kaniko."`

	// Custom if custom is specified, DevSpace will build the image with the help of
	// a custom script.
	Custom *CustomConfig `yaml:"custom,omitempty" json:"custom,omitempty" jsonschema_extras:"group=engines" jsonschema_description:"Custom if custom is specified, DevSpace will build the image with the help of a custom script."`

	// InjectRestartHelper will inject a small restart script into the container and wraps the entrypoint of that
	// container, so that devspace is able to restart the complete container during sync.
	// Please make sure you either have an Entrypoint defined in the devspace config or in the
	// dockerfile for this image, otherwise devspace will fail.
	InjectRestartHelper bool `yaml:"injectRestartHelper,omitempty" json:"injectRestartHelper,omitempty" jsonschema:"-" jsonschema_description:"InjectRestartHelper will inject a small restart script into the container and wraps the entrypoint of that container, so that devspace is able to restart the complete container during sync."`

	// RestartHelperPath will load the restart helper from this location instead of using the bundled
	// one within DevSpace. Can be either a local path or an URL where to find the restart helper.
	RestartHelperPath string `yaml:"restartHelperPath,omitempty" json:"restartHelperPath,omitempty" jsonschema:"-" jsonschema_description:"RestartHelperPath will load the restart helper from this location instead of using the bundled one within DevSpace. Can be either a local path or an URL where to find the restart helper."`
}

// RebuildStrategy is the type of a image rebuild strategy
type RebuildStrategy string

// List of values that source can take
const (
	RebuildStrategyDefault              RebuildStrategy = "default"
	RebuildStrategyAlways               RebuildStrategy = "always"
	RebuildStrategyIgnoreContextChanges RebuildStrategy = "ignoreContextChanges"
)

// DockerConfig tells the DevSpace CLI to build with Docker on Minikube or on localhost
type DockerConfig struct {
	// DisableFallback allows you to turn off kaniko building if docker isn't installed
	DisableFallback *bool `yaml:"disableFallback,omitempty" json:"disableFallback,omitempty" jsonschema_description:"DisableFallback allows you to turn off kaniko building if docker isn't installed."`
	// PreferMinikube allows you to turn off using the minikube docker daemon if the minikube
	// context is used.
	PreferMinikube *bool `yaml:"preferMinikube,omitempty" json:"preferMinikube,omitempty" jsonschema_description:"PreferMinikube allows you to turn off using the minikube docker daemon if the minikube context is used."`
	// UseCLI specifies if DevSpace should use the docker cli for building
	UseCLI bool `yaml:"useCli,omitempty" json:"useCli,omitempty" jsonschema_description:"// UseCLI specifies if DevSpace should use the docker cli for building."`
	// Args are additional arguments to pass to the docker cli
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are additional arguments to pass to the docker cli."`

	// DEPRECATED: UseBuildKit
	UseBuildKit bool `yaml:"useBuildKit,omitempty" json:"useBuildKit,omitempty" jsonschema:"-" jsonschema_description:"UseBuildKit"`
}

// BuildKitConfig tells the DevSpace CLI to
type BuildKitConfig struct {
	// InCluster if specified, DevSpace will use BuildKit to build the image within the cluster
	InCluster *BuildKitInClusterConfig `yaml:"inCluster,omitempty" json:"inCluster,omitempty" jsonschema_description:"InCluster if specified, DevSpace will use BuildKit to build the image within the cluster."`

	// PreferMinikube if false, will not try to use the minikube docker daemon to build the image
	PreferMinikube *bool `yaml:"preferMinikube,omitempty" json:"preferMinikube,omitempty" jsonschema_description:"PreferMinikube if false, will not try to use the minikube docker daemon to build the image."`

	// Args are additional arguments to call docker buildx build with
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are additional arguments to call docker buildx build with."`

	// Command to override the base command to create a builder and build images. Defaults to ["docker", "buildx"]
	Command []string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command to override the base command to create a builder and build images. Defaults to [\"docker\", \"buildx\"]."`
}

// BuildKitInClusterConfig holds the buildkit builder config
type BuildKitInClusterConfig struct {
	// Name is the name of the builder to use. If omitted, DevSpace will try to create
	// or reuse a builder in the form devspace-$NAMESPACE
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name is the name of the builder to use. If omitted, DevSpace will try to create or reuse a builder in the form devspace-$NAMESPACE."`

	// Namespace where to create the builder deployment in. Defaults to the current
	// active namespace.
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema_description:"Namespace where to create the builder deployment in. Defaults to the current active namespace."`

	// Rootless if enabled will create a rootless builder deployment.
	Rootless bool `yaml:"rootless,omitempty" json:"rootless,omitempty" jsonschema_description:"Rootless if enabled will create a rootless builder deployment."`

	// Image is the docker image to use for the BuildKit deployment
	Image string `yaml:"image,omitempty" json:"image,omitempty" jsonschema_description:"Image is the docker image to use for the BuildKit deployment."`

	// NodeSelector is the node selector to use for the BuildKit deployment
	NodeSelector string `yaml:"nodeSelector,omitempty" json:"nodeSelector,omitempty" jsonschema_description:"NodeSelector is the node selector to use for the BuildKit deployment."`

	// NoCreate. By default, DevSpace will try to create a new builder if it cannot be found.
	// If this is true, DevSpace will fail if the specified builder cannot be found.
	NoCreate bool `yaml:"noCreate,omitempty" json:"noCreate,omitempty" jsonschema_description:"NoCreate. By default, DevSpace will try to create a new builder if it cannot be found. If this is true, DevSpace will fail if the specified builder cannot be found."`

	// NoRecreate. By default, DevSpace will try to recreate the builder if the builder configuration
	// in the devspace.yaml differs from the actual builder configuration. If this is
	// true, DevSpace will not try to do that.
	NoRecreate bool `yaml:"noRecreate,omitempty" json:"noRecreate,omitempty" jsonschema_description:"NoRecreate. By default, DevSpace will try to recreate the builder if the builder configuration in the devspace.yaml differs from the actual builder configuration. If this is true, DevSpace will not try to do that."`

	// NoLoad if enabled, DevSpace will not try to load the built image into the local docker
	// daemon if skip push is defined
	NoLoad bool `yaml:"noLoad,omitempty" json:"noLoad,omitempty" jsonschema_description:"NoLoad if enabled, DevSpace will not try to load the built image into the local docker daemon if skip push is defined."`

	// CreateArgs are additional args to create the builder with.
	CreateArgs []string `yaml:"createArgs,omitempty" json:"createArgs,omitempty" jsonschema_description:"CreateArgs are additional args to create the builder with."`
}

// KanikoConfig tells the DevSpace CLI to build with Docker on Minikube or on localhost
type KanikoConfig struct {
	// Cache tells DevSpace if a cache repository should be used. defaults to false
	Cache bool `yaml:"cache,omitempty" json:"cache,omitempty" jsonschema_description:"Cache tells DevSpace if a cache repository should be used. defaults to false."`

	// SnapshotMode tells DevSpace which snapshot mode kaniko should use. defaults to time
	SnapshotMode string `yaml:"snapshotMode,omitempty" json:"snapshotMode,omitempty" jsonschema_description:"SnapshotMode tells DevSpace which snapshot mode kaniko should use. defaults to time."`

	// Image is the image name of the kaniko pod to use
	Image string `yaml:"image,omitempty" json:"image,omitempty" jsonschema_description:"Image is the image name of the kaniko pod to use."`

	// InitImage to override the init image of the kaniko pod
	InitImage string `yaml:"initImage,omitempty" json:"initImage,omitempty" jsonschema_description:"InitImage to override the init image of the kaniko pod."`

	// Args for additional arguments that should be passed to kaniko
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args for additional arguments that should be passed to kaniko."`

	// Command to replace the starting command for the kaniko container
	Command []string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command to replace the starting command for the kaniko container."`

	// Namespace is the namespace where the kaniko pod should be run
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema_description:"Namespace is the namespace where the kaniko pod should be run."`

	// Insecure allows pushing to insecure registries
	Insecure *bool `yaml:"insecure,omitempty" json:"insecure,omitempty" jsonschema_description:"Insecure allows pushing to insecure registries."`

	// PullSecret is the pull secret to mount by default
	PullSecret string `yaml:"pullSecret,omitempty" json:"pullSecret,omitempty" jsonschema_description:"PullSecret is the pull secret to mount by default."`

	// SkipPullSecretMount will skip mounting the pull secret
	SkipPullSecretMount bool `yaml:"skipPullSecretMount,omitempty" json:"skipPullSecretMount,omitempty" jsonschema_description:"SkipPullSecretMount will skip mounting the pull secret."`

	// NodeSelector is the node selector to use for the kaniko pod
	NodeSelector map[string]string `yaml:"nodeSelector,omitempty" json:"nodeSelector,omitempty" jsonschema_description:"NodeSelector is the node selector to use for the kaniko pod."`

	// Tolerations is a tolerations list to use for the kaniko pod
	Tolerations []k8sv1.Toleration `yaml:"tolerations,omitempty" json:"tolerations,omitempty" jsonschema_description:"Tolerations is a tolerations list to use for the kaniko pod."`

	// ServiceAccount the service account to use for the kaniko pod
	ServiceAccount string `yaml:"serviceAccount,omitempty" json:"serviceAccount,omitempty" jsonschema_description:"ServiceAccount the service account to use for the kaniko pod."`

	// Annotations are extra annotations that will be added to the build pod
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty" jsonschema_description:"Annotations are extra annotations that will be added to the build pod."`

	// Labels are extra labels that will be added to the build pod
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty" jsonschema_description:"Labels are extra labels that will be added to the build pod."`

	// InitEnv are extra environment variables that will be added to the build init container
	InitEnv map[string]string `yaml:"initEnv,omitempty" json:"initEnv,omitempty" jsonschema_description:"InitEnv are extra environment variables that will be added to the build init container."`

	// Env are extra environment variables that will be added to the build kaniko container
	// Will populate the env.value field.
	Env map[string]string `yaml:"env,omitempty" json:"env,omitempty" jsonschema_description:"Env are extra environment variables that will be added to the build kaniko container will populate the env.value field."`

	// EnvFrom are extra environment variables from configmap or secret that will be added to the build kaniko container
	// Will populate the env.valueFrom field.
	EnvFrom map[string]map[string]interface{} `yaml:"envFrom,omitempty" json:"envFrom,omitempty" jsonschema_description:"EnvFrom are extra environment variables from configmap or secret that will be added to the build kaniko container will populate the env.valueFrom field."`

	// AdditionalMounts are additional mounts that will be added to the build pod
	AdditionalMounts []KanikoAdditionalMount `yaml:"additionalMounts,omitempty" json:"additionalMounts,omitempty" jsonschema_description:"AdditionalMounts are additional mounts that will be added to the build pod."`

	// Resources are the resources that should be set on the kaniko pod
	Resources *PodResources `yaml:"resources,omitempty" json:"resources,omitempty" jsonschema_description:"Resources are the resources that should be set on the kaniko pod."`
}

// PodResources describes the resources section of the started kaniko pod
type PodResources struct {
	// Requests are the requests part of the resources
	Requests map[string]string `yaml:"requests,omitempty" json:"requests,omitempty" jsonschema_description:"Requests are the requests part of the resources."`

	// Limits are the limits part of the resources
	Limits map[string]string `yaml:"limits,omitempty" json:"limits,omitempty" jsonschema_description:"Limits are the limits part of the resources."`
}

// KanikoAdditionalMount tells devspace how the additional mount of the kaniko pod should look like
type KanikoAdditionalMount struct {
	// The secret that should be mounted
	Secret *KanikoAdditionalMountSecret `yaml:"secret,omitempty" json:"secret,omitempty" jsonschema_description:"The secret that should be mounted."`

	// The configMap that should be mounted
	ConfigMap *KanikoAdditionalMountConfigMap `yaml:"configMap,omitempty" json:"configMap,omitempty" jsonschema_description:"The configMap that should be mounted."`

	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `yaml:"readOnly,omitempty" json:"readOnly,omitempty" jsonschema_description:"Mounted read-only if true, read-write otherwise (false or unspecified). Defaults to false."`

	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	MountPath string `yaml:"mountPath,omitempty" json:"mountPath,omitempty" jsonschema_description:"Path within the container at which the volume should be mounted. Must not contain ':'."`

	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	SubPath string `yaml:"subPath,omitempty" json:"subPath,omitempty" jsonschema_description:"Path within the volume from which the container's volume should be mounted. Defaults to \"\" (volume's root)."`
}

type KanikoAdditionalMountConfigMap struct {
	// Name of the configmap
	// +optional
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name of the configmap."`

	// If unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	Items []KanikoAdditionalMountKeyToPath `yaml:"items,omitempty" json:"items,omitempty" jsonschema_description:"If unspecified, each key-value pair in the Data field of the referenced ConfigMap will be projected into the volume as a file whose name is the key and content is the value."`

	// Optional: mode bits to use on created files by default. Must be a
	// value between 0 and 0777. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `yaml:"defaultMode,omitempty" json:"defaultMode,omitempty" jsonschema_description:"mode bits to use on created files by default. Must be a value between 0 and 0777. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set."`
}

type KanikoAdditionalMountSecret struct {
	// Name of the secret in the pod's namespace to use.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	// +optional
	Name string `yaml:"name" json:"name" jsonschema_description:"Name of the secret in the pod's namespace to use."`

	// If unspecified, each key-value pair in the Data field of the referenced
	// Secret will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the Secret,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	Items []KanikoAdditionalMountKeyToPath `yaml:"items,omitempty" json:"items,omitempty" jsonschema_description:"If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value."`

	// Optional: mode bits to use on created files by default. Must be a
	// value between 0 and 0777. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `yaml:"defaultMode,omitempty" json:"defaultMode,omitempty" jsonschema_description:"mode bits to use on created files by default. Must be a value between 0 and 0777. Defaults to 0644."`
}

type KanikoAdditionalMountKeyToPath struct {
	// The key to project.
	Key string `yaml:"key" json:"key" jsonschema_description:"The key to project."`

	// The relative path of the file to map the key to.
	// May not be an absolute path.
	// May not contain the path element '..'.
	// May not start with the string '..'.
	Path string `yaml:"path" json:"path" jsonschema_description:"The relative path of the file to map the key to."`

	// Optional: mode bits to use on this file, must be a value between 0
	// and 0777. If not specified, the volume defaultMode will be used.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	Mode *int32 `yaml:"mode,omitempty" json:"mode,omitempty" jsonschema_description:"mode bits to use on this file, must be a value between 0 and 0777. If not specified, the volume defaultMode will be used."`
}

// CustomConfig tells the DevSpace CLI to build with a custom build script
type CustomConfig struct {
	// Command to execute to build the image. You can use ${runtime.images.my-image.image} and ${runtime.image.my-image.tag}
	// to reference the image and tag that should get built.
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command to execute to build the image. You can use ${runtime.images.my-image.image} and ${runtime.image.my-image.tag} to reference the image and tag that should get built."`
	// OnChange will determine when the command should be rerun
	OnChange []string `yaml:"onChange,omitempty" json:"onChange,omitempty" jsonschema_description:"OnChange will determine when the command should be rerun."`

	// DEPRECATED: Commands
	Commands []CustomConfigCommand `yaml:"commands,omitempty" json:"commands,omitempty" jsonschema:"-" jsonschema_description:"Commands"`
	// DEPRECATED: Args
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema:"-" jsonschema_description:"Args"`
	// DEPRECATED: AppendArgs
	AppendArgs []string `yaml:"appendArgs,omitempty" json:"appendArgs,omitempty" jsonschema:"-" jsonschema_description:"AppendArgs"`
	// DEPRECATED: ImageFlag
	ImageFlag string `yaml:"imageFlag,omitempty" json:"imageFlag,omitempty" jsonschema:"-" jsonschema_description:"ImageFlag"`
	// DEPRECATED: ImageTagOnly
	ImageTagOnly bool `yaml:"imageTagOnly,omitempty" json:"imageTagOnly,omitempty" jsonschema:"-" jsonschema_description:"ImageTagOnly"`
	// DEPRECATED: SkipImageArg
	SkipImageArg *bool `yaml:"skipImageArg,omitempty" json:"skipImageArg,omitempty" jsonschema:"-" jsonschema_description:"SkipImageArg"`
}

// CustomConfigCommand holds the information about a command on a specific operating system
type CustomConfigCommand struct {
	// Command to run
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command to run."`
	// OperatingSystem to run this command on
	OperatingSystem string `yaml:"os,omitempty" json:"os,omitempty" jsonschema_description:"OperatingSystem to run this command on."`
}

// DeploymentConfig defines the configuration how the devspace should be deployed
type DeploymentConfig struct {
	// Name of the deployment
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name of the deployment."`

	// Helm tells DevSpace to deploy this deployment via helm
	Helm *HelmConfig `yaml:"helm,omitempty" json:"helm,omitempty" jsonschema_description:"Helm tells DevSpace to deploy this deployment via helm."`
	// Kubectl tells DevSpace to deploy this deployment via kubectl or kustomize
	Kubectl *KubectlConfig `yaml:"kubectl,omitempty" json:"kubectl,omitempty" jsonschema_description:"Kubectl tells DevSpace to deploy this deployment via kubectl or kustomize."`

	// UpdateImageTags lets you define if DevSpace should update the tags of the images defined in the
	// images section with their most recent built tag.
	UpdateImageTags *bool `yaml:"updateImageTags,omitempty" json:"updateImageTags,omitempty" jsonschema_description:"UpdateImageTags lets you define if DevSpace should update the tags of the images defined in the images section with their most recent built tag."`

	// Namespace where to deploy this deployment
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema_description:"Namespace where to deploy this deployment."`
}

// ComponentConfig holds the component information
type ComponentConfig struct {
	InitContainers      []*ContainerConfig       `yaml:"initContainers,omitempty" json:"initContainers,omitempty"`
	Containers          []*ContainerConfig       `yaml:"containers,omitempty" json:"containers,omitempty"`
	Labels              map[string]string        `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations         map[string]string        `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	Volumes             []*VolumeConfig          `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	Service             *ServiceConfig           `yaml:"service,omitempty" json:"service,omitempty"`
	ServiceName         string                   `yaml:"serviceName,omitempty" json:"serviceName,omitempty"`
	Ingress             *IngressConfig           `yaml:"ingress,omitempty" json:"ingress,omitempty"`
	Replicas            *int                     `yaml:"replicas,omitempty" json:"replicas,omitempty"`
	Autoscaling         *AutoScalingConfig       `yaml:"autoScaling,omitempty" json:"autoScaling,omitempty"`
	RollingUpdate       *RollingUpdateConfig     `yaml:"rollingUpdate,omitempty" json:"rollingUpdate,omitempty"`
	PullSecrets         []*string                `yaml:"pullSecrets,omitempty" json:"pullSecrets,omitempty"`
	Tolerations         []map[string]interface{} `yaml:"tolerations,omitempty" json:"tolerations,omitempty"`
	Affinity            map[string]interface{}   `yaml:"affinity,omitempty" json:"affinity,omitempty"`
	NodeSelector        map[string]interface{}   `yaml:"nodeSelector,omitempty" json:"nodeSelector,omitempty"`
	NodeName            string                   `yaml:"nodeName,omitempty" json:"nodeName,omitempty"`
	PodManagementPolicy string                   `yaml:"podManagementPolicy,omitempty" json:"podManagementPolicy,omitempty"`

	DNSConfig                     map[string]interface{}   `yaml:"dnsConfig,omitempty" json:"dnsConfig,omitempty"`
	HostAliases                   []map[string]interface{} `yaml:"hostAliases,omitempty" json:"hostAliases,omitempty"`
	Overhead                      map[string]interface{}   `yaml:"overhead,omitempty" json:"overhead,omitempty"`
	ReadinessGates                []map[string]interface{} `yaml:"readinessGates,omitempty" json:"readinessGates,omitempty"`
	SecurityContext               map[string]interface{}   `yaml:"securityContext,omitempty" json:"securityContext,omitempty"`
	TopologySpreadConstraints     []map[string]interface{} `yaml:"topologySpreadConstraints,omitempty" json:"topologySpreadConstraints,omitempty"`
	ActiveDeadlineSeconds         *int                     `yaml:"activeDeadlineSeconds,omitempty" json:"activeDeadlineSeconds,omitempty"`
	AutomountServiceAccountToken  *bool                    `yaml:"automountServiceAccountToken,omitempty" json:"automountServiceAccountToken,omitempty"`
	DNSPolicy                     *string                  `yaml:"dnsPolicy,omitempty" json:"dnsPolicy,omitempty"`
	EnableServiceLinks            *bool                    `yaml:"enableServiceLinks,omitempty" json:"enableServiceLinks,omitempty"`
	HostIPC                       *bool                    `yaml:"hostIPC,omitempty" json:"hostIPC,omitempty"`
	HostNetwork                   *bool                    `yaml:"hostNetwork,omitempty" json:"hostNetwork,omitempty"`
	HostPID                       *bool                    `yaml:"hostPID,omitempty" json:"hostPID,omitempty"`
	Hostname                      *string                  `yaml:"hostname,omitempty" json:"hostname,omitempty"`
	PreemptionPolicy              *string                  `yaml:"preemptionPolicy,omitempty" json:"preemptionPolicy,omitempty"`
	Priority                      *int                     `yaml:"priority,omitempty" json:"priority,omitempty"`
	PriorityClassName             *string                  `yaml:"priorityClassName,omitempty" json:"priorityClassName,omitempty"`
	RestartPolicy                 *string                  `yaml:"restartPolicy,omitempty" json:"restartPolicy,omitempty"`
	RuntimeClassName              *string                  `yaml:"runtimeClassName,omitempty" json:"runtimeClassName,omitempty"`
	SchedulerName                 *string                  `yaml:"schedulerName,omitempty" json:"schedulerName,omitempty"`
	ServiceAccount                *string                  `yaml:"serviceAccount,omitempty" json:"serviceAccount,omitempty"`
	ServiceAccountName            *string                  `yaml:"serviceAccountName,omitempty" json:"serviceAccountName,omitempty"`
	SetHostnameAsFQDN             *bool                    `yaml:"setHostnameAsFQDN,omitempty" json:"setHostnameAsFQDN,omitempty"`
	ShareProcessNamespace         *bool                    `yaml:"shareProcessNamespace,omitempty" json:"shareProcessNamespace,omitempty"`
	Subdomain                     *string                  `yaml:"subdomain,omitempty" json:"subdomain,omitempty"`
	TerminationGracePeriodSeconds *int                     `yaml:"terminationGracePeriodSeconds,omitempty" json:"terminationGracePeriodSeconds,omitempty"`
	EphemeralContainers           []map[string]interface{} `yaml:"ephemeralContainers,omitempty" json:"ephemeralContainers,omitempty"`
}

// ContainerConfig holds the configurations of a container
type ContainerConfig struct {
	Name                     string                   `yaml:"name,omitempty" json:"name,omitempty"`
	Image                    string                   `yaml:"image,omitempty" json:"image,omitempty"`
	Command                  []string                 `yaml:"command,omitempty" json:"command,omitempty"`
	Args                     []string                 `yaml:"args,omitempty" json:"args,omitempty"`
	Stdin                    bool                     `yaml:"stdin,omitempty" json:"stdin,omitempty"`
	TTY                      bool                     `yaml:"tty,omitempty" json:"tty,omitempty"`
	Env                      []map[string]interface{} `yaml:"env,omitempty" json:"env,omitempty"`
	EnvFrom                  []map[string]interface{} `yaml:"envFrom,omitempty" json:"envFrom,omitempty"`
	VolumeMounts             []*VolumeMountConfig     `yaml:"volumeMounts,omitempty" json:"volumeMounts,omitempty"`
	Resources                map[string]interface{}   `yaml:"resources,omitempty" json:"resources,omitempty"`
	LivenessProbe            map[string]interface{}   `yaml:"livenessProbe,omitempty" json:"livenessProbe,omitempty"`
	ReadinessProbe           map[string]interface{}   `yaml:"readinessProbe,omitempty" json:"readinessProbe,omitempty"`
	StartupProbe             map[string]interface{}   `yaml:"startupProbe,omitempty" json:"startupProbe,omitempty"`
	SecurityContext          map[string]interface{}   `yaml:"securityContext,omitempty" json:"securityContext,omitempty"`
	Lifecycle                map[string]interface{}   `yaml:"lifecycle,omitempty" json:"lifecycle,omitempty"`
	VolumeDevices            []map[string]interface{} `yaml:"volumeDevices,omitempty" json:"volumeDevices,omitempty"`
	ImagePullPolicy          string                   `yaml:"imagePullPolicy,omitempty" json:"imagePullPolicy,omitempty"`
	WorkingDir               string                   `yaml:"workingDir,omitempty" json:"workingDir,omitempty"`
	StdinOnce                bool                     `yaml:"stdinOnce,omitempty" json:"stdinOnce,omitempty"`
	TerminationMessagePath   string                   `yaml:"terminationMessagePath,omitempty" json:"terminationMessagePath,omitempty"`
	TerminationMessagePolicy string                   `yaml:"terminationMessagePolicy,omitempty" json:"terminationMessagePolicy,omitempty"`
}

// VolumeMountConfig holds the configuration for a specific mount path
type VolumeMountConfig struct {
	ContainerPath string                   `yaml:"containerPath,omitempty" json:"containerPath,omitempty"`
	Volume        *VolumeMountVolumeConfig `yaml:"volume,omitempty" json:"volume,omitempty"`
}

// VolumeMountVolumeConfig holds the configuration for a specific mount path volume
type VolumeMountVolumeConfig struct {
	Name     string `yaml:"name,omitempty" json:"name,omitempty"`
	SubPath  string `yaml:"subPath,omitempty" json:"subPath,omitempty"`
	ReadOnly *bool  `yaml:"readOnly,omitempty" json:"readOnly,omitempty"`
	Shared   *bool  `yaml:"shared,omitempty" json:"shared,omitempty"`
}

// VolumeConfig holds the configuration for a specific volume
type VolumeConfig struct {
	Name             string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Labels           map[string]string      `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations      map[string]string      `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	Size             string                 `yaml:"size,omitempty" json:"size,omitempty"`
	ConfigMap        map[string]interface{} `yaml:"configMap,omitempty" json:"configMap,omitempty"`
	Secret           map[string]interface{} `yaml:"secret,omitempty" json:"secret,omitempty"`
	EmptyDir         map[string]interface{} `yaml:"emptyDir,omitempty" json:"emptyDir,omitempty"`
	StorageClassName string                 `yaml:"storageClassName,omitempty" json:"storageClassName,omitempty"`
	VolumeMode       string                 `yaml:"volumeMode,omitempty" json:"volumeMode,omitempty"`
	VolumeName       string                 `yaml:"volumeName,omitempty" json:"volumeName,omitempty"`
	DataSource       map[string]interface{} `yaml:"dataSource,omitempty" json:"dataSource,omitempty"`
	AccessModes      []string               `yaml:"accessModes,omitempty" json:"accessModes,omitempty"`
}

// ServiceConfig holds the configuration of a component service
type ServiceConfig struct {
	Name                     string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Labels                   map[string]string      `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations              map[string]string      `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	Type                     string                 `yaml:"type,omitempty" json:"type,omitempty"`
	Ports                    []*ServicePortConfig   `yaml:"ports,omitempty" json:"ports,omitempty"`
	ExternalIPs              []string               `yaml:"externalIPs,omitempty" json:"externalIPs,omitempty"`
	ClusterIP                string                 `yaml:"clusterIP,omitempty" json:"clusterIP,omitempty"`
	ExternalName             string                 `yaml:"externalName,omitempty" json:"externalName,omitempty"`
	ExternalTrafficPolicy    string                 `yaml:"externalTrafficPolicy,omitempty" json:"externalTrafficPolicy,omitempty"`
	HealthCheckNodePort      int                    `yaml:"healthCheckNodePort,omitempty" json:"healthCheckNodePort,omitempty"`
	IPFamily                 *string                `yaml:"ipFamily,omitempty" json:"ipFamily,omitempty"`
	LoadBalancerIP           *string                `yaml:"loadBalancerIP,omitempty" json:"loadBalancerIP,omitempty"`
	LoadBalancerSourceRanges []string               `yaml:"loadBalancerSourceRanges,omitempty" json:"loadBalancerSourceRanges,omitempty"`
	PublishNotReadyAddresses bool                   `yaml:"publishNotReadyAddresses,omitempty" json:"publishNotReadyAddresses,omitempty"`
	SessionAffinity          map[string]interface{} `yaml:"sessionAffinity,omitempty" json:"sessionAffinity,omitempty"`
	SessionAffinityConfig    map[string]interface{} `yaml:"sessionAffinityConfig,omitempty" json:"sessionAffinityConfig,omitempty"`
	TopologyKeys             []string               `yaml:"topologyKeys,omitempty" json:"topologyKeys,omitempty"`
}

// ServicePortConfig holds the port configuration of a component service
type ServicePortConfig struct {
	Name          string `yaml:"name,omitempty" json:"name,omitempty"`
	Port          *int   `yaml:"port,omitempty" json:"port,omitempty"`
	ContainerPort *int   `yaml:"containerPort,omitempty" json:"containerPort,omitempty"`
	Protocol      string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

// IngressConfig holds the configuration of a component ingress
type IngressConfig struct {
	Name             string                 `yaml:"name,omitempty" json:"name,omitempty"`
	Labels           map[string]string      `yaml:"labels,omitempty" json:"labels,omitempty"`
	Annotations      map[string]string      `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	TLS              string                 `yaml:"tls,omitempty" json:"tls,omitempty"`
	TLSClusterIssuer string                 `yaml:"tlsClusterIssuer,omitempty" json:"tlsClusterIssuer,omitempty"`
	IngressClass     string                 `yaml:"ingressClass,omitempty" json:"ingressClass,omitempty"`
	Rules            []*IngressRuleConfig   `yaml:"rules,omitempty" json:"rules,omitempty"`
	Backend          map[string]interface{} `yaml:"backend,omitempty" json:"backend,omitempty"`
	IngressClassName *string                `yaml:"ingressClassName,omitempty" json:"ingressClassName,omitempty"`
}

// IngressRuleConfig holds the port configuration of a component service
type IngressRuleConfig struct {
	Host        string `yaml:"host,omitempty" json:"host,omitempty"`
	TLS         string `yaml:"tls,omitempty" json:"tls,omitempty"` // DEPRECATED
	Path        string `yaml:"path,omitempty" json:"path,omitempty"`
	ServicePort *int   `yaml:"servicePort,omitempty" json:"servicePort,omitempty"`
	ServiceName string `yaml:"serviceName,omitempty" json:"serviceName,omitempty"`
	PathType    string `yaml:"pathType,omitempty" json:"pathType,omitempty"`
}

// AutoScalingConfig holds the autoscaling config of a component
type AutoScalingConfig struct {
	Horizontal *AutoScalingHorizontalConfig `yaml:"horizontal,omitempty" json:"horizontal,omitempty"`
}

// AutoScalingHorizontalConfig holds the horizontal autoscaling config of a component
type AutoScalingHorizontalConfig struct {
	MaxReplicas           *int   `yaml:"maxReplicas,omitempty" json:"maxReplicas,omitempty"`
	AverageCPU            string `yaml:"averageCPU,omitempty" json:"averageCPU,omitempty"`
	AverageRelativeCPU    string `yaml:"averageRelativeCPU,omitempty" json:"averageRelativeCPU,omitempty"`
	AverageMemory         string `yaml:"averageMemory,omitempty" json:"averageMemory,omitempty"`
	AverageRelativeMemory string `yaml:"averageRelativeMemory,omitempty" json:"averageRelativeMemory,omitempty"`
}

// RollingUpdateConfig holds the configuration for rolling updates
type RollingUpdateConfig struct {
	Enabled        *bool  `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	MaxSurge       string `yaml:"maxSurge,omitempty" json:"maxSurge,omitempty"`
	MaxUnavailable string `yaml:"maxUnavailable,omitempty" json:"maxUnavailable,omitempty"`
	Partition      *int   `yaml:"partition,omitempty" json:"partition,omitempty"`
}

// HelmConfig defines the specific helm options used during deployment
type HelmConfig struct {
	// ReleaseName of the helm configuration
	ReleaseName string `yaml:"releaseName,omitempty" json:"releaseName,omitempty" jsonschema_description:"ReleaseName of the helm configuration."`
	// Chart holds the chart configuration and where DevSpace can find the chart
	Chart *ChartConfig `yaml:"chart,omitempty" json:"chart,omitempty" jsonschema:"required" jsonschema_description:"Chart holds the chart configuration and where DevSpace can find the chart."`
	// Values are additional values that should get passed to deploying this chart
	Values map[string]interface{} `yaml:"values,omitempty" json:"values,omitempty" jsonschema_description:"Values are additional values that should get passed to deploying this chart."`
	// ValuesFiles are additional files that hold values for deploying this chart
	ValuesFiles []string `yaml:"valuesFiles,omitempty" json:"valuesFiles,omitempty" jsonschema_description:"ValuesFiles are additional files that hold values for deploying this chart."`
	// DisplayOutput allows you to display the helm output to the console
	DisplayOutput bool `yaml:"displayOutput,omitempty" json:"output,omitempty" jsonschema_description:"DisplayOutput allows you to display the helm output to the console."`

	// UpgradeArgs are additional arguments to pass to `helm upgrade`
	UpgradeArgs []string `yaml:"upgradeArgs,omitempty" json:"upgradeArgs,omitempty" jsonschema_description:"UpgradeArgs are additional arguments to pass to 'helm upgrade'."`
	// TemplateArgs are additional arguments to pass to `helm template`
	TemplateArgs []string `yaml:"templateArgs,omitempty" json:"templateArgs,omitempty" jsonschema_description:"TemplateArgs are additional arguments to pass to 'helm template'."`
}

// ChartConfig defines the helm chart options
type ChartConfig struct {
	// Name is the name of the helm chart to deploy. Can also be a local path
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"required" jsonschema_extras:"group=repo,group_name=Source: Helm Repository" jsonschema_description:"Name is the name of the helm chart to deploy. Can also be a local path."`
	// Version is the version of the helm chart to deploy
	Version string `yaml:"version,omitempty" json:"version,omitempty" jsonschema_extras:"group=repo" jsonschema_description:"Version is the version of the helm chart to deploy."`
	// RepoURL is the url of the repo to deploy the chart from
	RepoURL string `yaml:"repo,omitempty" json:"repo,omitempty" jsonschema_extras:"group=repo" jsonschema_description:"RepoURL is the url of the repo to deploy the chart from."`
	// Username is the username to authenticate to the chart repo
	Username string `yaml:"username,omitempty" json:"username,omitempty" jsonschema_extras:"group=repo" jsonschema_description:"Username is the username to authenticate to the chart repo."`
	// Password is the password to authenticate to the chart repo
	Password string `yaml:"password,omitempty" json:"password,omitempty" jsonschema_extras:"group=repo" jsonschema_description:"Password is the password to authenticate to the chart repo."`
	// Source can be used to reference an helm chart from a distant location
	// such as a git repository
	Source *SourceConfig `yaml:",inline" json:",inline" jsonschema_description:"Source can be used to reference an helm chart from a distant location such as a git repository."`
}

// KubectlConfig defines the specific kubectl options used during deployment
type KubectlConfig struct {
	// Manifests is a list of files or folders that will be deployed by DevSpace using kubectl
	// or kustomize
	Manifests []string `yaml:"manifests,omitempty" json:"manifests,omitempty" jsonschema:"required" jsonschema_description:"Manifests is a list of files or folders that will be deployed by DevSpace using kubectl or kustomize."`
	// ApplyArgs are extra arguments for `kubectl apply`
	ApplyArgs []string `yaml:"applyArgs,omitempty" json:"applyArgs,omitempty" jsonschema_description:"ApplyArgs are extra arguments for 'kubectl apply'."`
	// CreateArgs are extra arguments for `kubectl create` which will be run before `kubectl apply`
	CreateArgs []string `yaml:"createArgs,omitempty" json:"createArgs,omitempty" jsonschema_description:"CreateArgs are extra arguments for 'kubectl create' which will be run before 'kubectl apply'."`
	// KubectlBinaryPath is the optional path where to find the kubectl binary
	KubectlBinaryPath string `yaml:"kubectlBinaryPath,omitempty" json:"kubectlBinaryPath,omitempty" jsonschema_description:"KubectlBinaryPath is the optional path where to find the kubectl binary."`

	// Kustomize can be used to enable kustomize instead of kubectl
	Kustomize *bool `yaml:"kustomize,omitempty" json:"kustomize,omitempty" jsonschema_extras:"group=kustomize,group_name=Kustomize" jsonschema_description:"Kustomize can be used to enable kustomize instead of kubectl."`
	// KustomizeArgs are extra arguments for `kustomize build` which will be run before `kubectl apply`
	KustomizeArgs []string `yaml:"kustomizeArgs,omitempty" json:"kustomizeArgs,omitempty" jsonschema_extras:"group=kustomize" jsonschema_description:"KustomizeArgs are extra arguments for 'kustomize build' which will be run before 'kubectl apply'."`
	// KustomizeBinaryPath is the optional path where to find the kustomize binary
	KustomizeBinaryPath string `yaml:"kustomizeBinaryPath,omitempty" json:"kustomizeBinaryPath,omitempty" jsonschema_extras:"group=kustomize" jsonschema_description:"KustomizeBinaryPath is the optional path where to find the kustomize binary."`
}

// DevPod holds configurations for selecting a pod and starting dev services for that pod
type DevPod struct {
	// Name of the dev configuration
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name of the dev configuration."`
	// ImageSelector to select a pod
	ImageSelector string `yaml:"imageSelector,omitempty" json:"imageSelector,omitempty" jsonschema_extras:"group=selector" jsonschema_description:"ImageSelector to select a pod."`
	// LabelSelector to select a pod
	LabelSelector map[string]string `yaml:"labelSelector,omitempty" json:"labelSelector,omitempty" jsonschema_extras:"group=selector" jsonschema_description:"LabelSelector to select a pod."`
	// Namespace where to select the pod
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema_extras:"group=selector" jsonschema_description:"Namespace where to select the pod."`

	// DevContainer can either be defined inline if the pod only has a single container or
	// containers can be used to define configurations for multiple containers in the same
	// pod.
	DevContainer `yaml:",inline" json:",inline" jsonschema_description:"DevContainer can either be defined inline if the pod only has a single container or containers can be used to define configurations for multiple containers in the same pod."`

	// Ports defines port mappings from the remote pod that should be forwarded to your local
	// computer
	Ports []*PortMapping `yaml:"ports,omitempty" json:"ports,omitempty" jsonschema_extras:"group=ports" jsonschema_description:"Ports defines port mappings from the remote pod that should be forwarded to your local computer."`

	// PersistenceOptions are additional options for persisting paths within this pod
	PersistenceOptions *PersistenceOptions `yaml:"persistenceOptions,omitempty" json:"persistenceOptions,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"PersistenceOptions are additional options for persisting paths within this pod."`

	// Patches are additional changes to the pod spec that should be applied
	Patches []*PatchConfig `yaml:"patches,omitempty" json:"patches,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"Patches are additional changes to the pod spec that should be applied."`

	// Open defines urls that should be opened as soon as they are reachable
	Open []*OpenConfig `yaml:"open,omitempty" json:"open,omitempty" jsonschema_extras:"group=workflows_background,group_name=Background Dev Workflows" jsonschema_description:"Open defines urls that should be opened as soon as they are reachable."`

	Containers map[string]*DevContainer `yaml:"containers,omitempty" json:"containers,omitempty" jsonschema_extras:"group=selector" jsonschema_description:"DevContainer holds options for dev services that should get started within a certain container of the selected pod."`
}

// DevContainer holds options for dev services that should
// get started within a certain container of the selected pod
type DevContainer struct {
	// Container is the container name these services should get started.
	Container string `yaml:"container,omitempty" json:"container,omitempty" jsonschema_extras:"group=selector,group_name=Selector" jsonschema_description:"Container is the container name these services should get started."`

	// Target Container architecture to use for the devspacehelper (currently amd64 or arm64). Defaults to amd64, but
	// devspace tries to find out the architecture by itself by looking at the node this container runs on.
	Arch ContainerArchitecture `yaml:"arch,omitempty" json:"arch,omitempty" jsonschema_extras:"group=selector" jsonschema_description:" Target Container architecture to use for the devspacehelper (currently amd64 or arm64). Defaults to amd64, but devspace tries to find out the architecture by itself by looking at the node this container runs on."`

	// DevImage is the image to use for this container and will replace the existing image
	// if necessary.
	DevImage string `yaml:"devImage,omitempty" json:"devImage,omitempty" jsonschema_extras:"group=modifications,group_name=Modifications" jsonschema_description:"DevImage is the image to use for this container and will replace the existing image if necessary."`
	// Command can be used to override the entrypoint of the container
	Command []string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"Command can be used to override the entrypoint of the container."`
	// Args can be used to override the args of the container
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"Args can be used to override the args of the container."`
	// WorkingDir can be used to override the working dir of the container
	WorkingDir string `yaml:"workingDir,omitempty" json:"workingDir,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"WorkingDir can be used to override the working dir of the container."`
	// Env can be used to add environment variables to the container. DevSpace will
	// not replace existing environment variables if an environment variable is defined here.
	Env []EnvVar `yaml:"env,omitempty" json:"env,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"Env can be used to add environment variables to the container. DevSpace will not replace existing environment variables if an environment variable is defined here."`
	// Resources can be used to override the resource definitions of the container
	Resources *PodResources `yaml:"resources,omitempty" json:"resources,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"Resources can be used to override the resource definitions of the container."`

	// ReversePorts are port mappings to make local ports available inside the container
	ReversePorts []*PortMapping `yaml:"reversePorts,omitempty" json:"reversePorts,omitempty" jsonschema_extras:"group=ports,group_name=Port Forwarding" jsonschema_description:"ReversePorts are port mappings to make local ports available inside the container."`

	// Sync allows you to sync certain local paths with paths inside the container
	Sync []*SyncConfig `yaml:"sync,omitempty" json:"sync,omitempty" jsonschema_extras:"group=sync,group_name=File Sync" jsonschema_description:"Sync allows you to sync certain local paths with paths inside the container."`
	// SSH allows you to create an SSH tunnel to this container
	PersistPaths []PersistentPath `yaml:"persistPaths,omitempty" json:"persistPaths,omitempty" jsonschema_extras:"group=modifications" jsonschema_description:"SSH allows you to create an SSH tunnel to this container."`

	// Terminal allows you to tell DevSpace to open a terminal with screen support to this container
	Terminal *Terminal `yaml:"terminal,omitempty" json:"terminal,omitempty" jsonschema_extras:"group=workflows,group_name=Foreground Dev Workflows" jsonschema_description:"Terminal allows you to tell DevSpace to open a terminal with screen support to this container."`
	// Logs allows you to tell DevSpace to stream logs from this container to the console
	Logs *Logs `yaml:"logs,omitempty" json:"logs,omitempty" jsonschema_extras:"group=workflows" jsonschema_description:"Logs allows you to tell DevSpace to stream logs from this container to the console."`
	// Attach allows you to tell DevSpace to attach to this container
	Attach *Attach `yaml:"attach,omitempty" json:"attach,omitempty" jsonschema_extras:"group=workflows" jsonschema_description:"Attach allows you to tell DevSpace to attach to this container."`
	// PersistPaths allows you to persist certain paths within this container with a persistent volume claim
	SSH *SSH `yaml:"ssh,omitempty" json:"ssh,omitempty" jsonschema_extras:"group=workflows_background" jsonschema_description:"PersistPaths allows you to persist certain paths within this container with a persistent volume claim."`
	// ProxyCommands allow you to proxy certain local commands to the container
	ProxyCommands []*ProxyCommand `yaml:"proxyCommands,omitempty" json:"proxyCommands,omitempty" jsonschema_extras:"group=workflows_background" jsonschema_description:"ProxyCommands allow you to proxy certain local commands to the container."`
	// RestartHelper holds restart helper specific configuration. The restart helper is used to delay starting of
	// the container and restarting it and is injected via an annotation in the replaced pod.
	RestartHelper *RestartHelper `yaml:"restartHelper,omitempty" json:"restartHelper,omitempty" jsonschema_extras:"group=workflows_background" jsonschema_description:"RestartHelper holds restart helper specific configuration. The restart helper is used to delay starting of the container and restarting it and is injected via an annotation in the replaced pod."`
}

type RestartHelper struct {
	// Path defines the path to the restart helper that might be used if certain config
	// options are enabled
	Path string `yaml:"path,omitempty" json:"path,omitempty" jsonschema_description:"Path defines the path to the restart helper that might be used if certain config options are enabled."`
	// Inject signals DevSpace to inject the restart helper
	Inject *bool `yaml:"inject,omitempty" json:"inject,omitempty" jsonschema_description:"Inject signals DevSpace to inject the restart helper."`
}

type ProxyCommand struct {
	// GitCredentials configures a git credentials helper inside the container that proxies local git credentials
	GitCredentials bool `yaml:"gitCredentials,omitempty" json:"gitCredentials,omitempty" jsonschema_description:"GitCredentials configures a git credentials helper inside the container that proxies local git credentials."`

	// Command is the name of the command that should be available in the remote container. DevSpace
	// will create a small script for that inside the container that redirect command execution to
	// the local computer.
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command is the name of the command that should be available in the remote container. DevSpace will create a small script for that inside the container that redirect command execution to the local computer."`

	// LocalCommand can be used to run a different command than specified via the command option. By
	// default, this will be assumed to be the same as command.
	LocalCommand string `yaml:"localCommand,omitempty" json:"localCommand,omitempty" jsonschema_description:"LocalCommand can be used to run a different command than specified via the command option. By default, this will be assumed to be the same as command."`

	// SkipContainerEnv will not forward the container environment variables to the local command
	SkipContainerEnv bool `yaml:"skipContainerEnv,omitempty" json:"skipContainerEnv,omitempty" jsonschema_description:"SkipContainerEnv will not forward the container environment variables to the local command."`

	// Env are extra environment variables to set for the command
	Env map[string]string `yaml:"env,omitempty" json:"env,omitempty" jsonschema_description:"Env are extra environment variables to set for the command."`
}

type SSH struct {
	// Enabled can be used to enable the ssh server within the container. By default,
	// DevSpace will generate the required keys and create an entry in your ~/.ssh/config
	// for this container that can be used via `ssh dev-config-name.dev-project-name.devspace`
	Enabled bool `yaml:"enabled,omitempty" json:"enabled,omitempty" jsonschema_description:"Enabled can be used to enable the ssh server within the container. By default,	DevSpace will generate the required keys and create an entry in your ~/.ssh/config for this container that can be used via 'ssh dev-config-name.dev-project-name.devspace'."`

	// LocalHostname is the local ssh host to write to the ~/.ssh/config
	LocalHostname string `yaml:"localHostname,omitempty" json:"localHostname,omitempty" jsonschema_description:"LocalHostname is the local ssh host to write to the ~/.ssh/config."`

	// LocalPort is the local port to forward from, if empty will be random
	LocalPort int `yaml:"localPort,omitempty" json:"localPort,omitempty" jsonschema_description:"LocalPort is the local port to forward from, if empty will be random."`

	// RemoteAddress is the address to listen to inside the container
	RemoteAddress string `yaml:"remoteAddress,omitempty" json:"remoteAddress,omitempty" jsonschema_description:"RemoteAddress is the address to listen to inside the container."`
}

type EnvVar struct {
	// Name of the environment variable
	Name string `yaml:"name" json:"name" jsonschema_description:"Name of the environment variable."`
	// Value of the environment variable
	Value string `yaml:"value" json:"value" jsonschema_description:"Value of the environment variable."`
}

type Attach struct {
	// Enabled can be used to enable attaching to a container
	Enabled *bool `yaml:"enabled,omitempty" json:"enabled,omitempty" jsonschema_description:"Enabled can be used to enable attaching to a container."`

	// DisableReplace prevents DevSpace from actually replacing the pod with modifications so that
	// the pod starts up correctly.
	DisableReplace bool `yaml:"disableReplace,omitempty" json:"disableReplace,omitempty" jsonschema_description:"DisableReplace prevents DevSpace from actually replacing the pod with modifications so that the pod starts up correctly."`

	// DisableTTY is used to tell DevSpace to not use a TTY connection for attaching
	DisableTTY bool `yaml:"disableTTY,omitempty" json:"disableTTY,omitempty" jsonschema_description:"DisableTTY is used to tell DevSpace to not use a TTY connection for attaching."`
}

type Logs struct {
	// Enabled can be used to enable printing container logs
	Enabled *bool `yaml:"enabled,omitempty" json:"enabled,omitempty" jsonschema_description:"Enabled can be used to enable printing container logs."`

	// LastLines is the amount of lines to print of the running container initially
	LastLines int64 `yaml:"lastLines,omitempty" json:"lastLines,omitempty" jsonschema_description:"LastLines is the amount of lines to print of the running container initially."`
}

// PersistenceOptions are general persistence options DevSpace should use for all persistent paths
// within a single dev configuration
type PersistenceOptions struct {
	// Size is the size of the created persistent volume in Kubernetes size notation like 5Gi
	Size string `yaml:"size,omitempty" json:"size,omitempty" jsonschema_description:"Size is the size of the created persistent volume in Kubernetes size notation like 5Gi."`
	// StorageClassName is the storage type DevSpace should use for this persistent volume
	StorageClassName string `yaml:"storageClassName,omitempty" json:"storageClassName,omitempty" jsonschema_description:"StorageClassName is the storage type DevSpace should use for this persistent volume."`
	// AccessModes are the access modes DevSpace should use for the persistent volume
	AccessModes []string `yaml:"accessModes,omitempty" json:"accessModes,omitempty" jsonschema_description:"AccessModes are the access modes DevSpace should use for the persistent volume."`
	// ReadOnly specifies if the volume should be read only
	ReadOnly bool `yaml:"readOnly,omitempty" json:"readOnly,omitempty" jsonschema_description:"ReadOnly specifies if the volume should be read only."`
	// Name is the name of the PVC that should be created. If a PVC with that name
	// already exists, DevSpace will use that PVC instead of creating one.
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name is the name of the PVC that should be created. If a PVC with that name already exists, DevSpace will use that PVC instead of creating one."`
}

// PersistentPath holds options to configure persistence for DevSpace
type PersistentPath struct {
	// Path is the container path that should get persisted. By default, DevSpace will create an init container
	// that will copy over the contents of this folder from the existing image.
	Path string `yaml:"path,omitempty" json:"path,omitempty" jsonschema_description:"Path is the container path that should get persisted. By default, DevSpace will create an init container that will copy over the contents of this folder from the existing image."`
	// VolumePath is the sub path on the volume that is mounted as persistent volume for this path
	VolumePath string `yaml:"volumePath,omitempty" json:"volumePath,omitempty" jsonschema_description:"VolumePath is the sub path on the volume that is mounted as persistent volume for this path."`
	// ReadOnly will make the persistent path read only to the user
	ReadOnly bool `yaml:"readOnly,omitempty" json:"readOnly,omitempty" jsonschema_description:"ReadOnly will make the persistent path read only to the user."`
	// SkipPopulate will not create an init container to copy over the existing contents if true
	SkipPopulate bool `yaml:"skipPopulate,omitempty" json:"skipPopulate,omitempty" jsonschema_description:"SkipPopulate will not create an init container to copy over the existing contents if true."`

	// InitContainer holds additional options for the persistent path init container
	InitContainer *PersistentPathInitContainer `yaml:"initContainer,omitempty" json:"initContainer,omitempty" jsonschema_description:"InitContainer holds additional options for the persistent path init container."`
}

// PersistentPathInitContainer defines additional options for the persistent path init container
type PersistentPathInitContainer struct {
	// Resources are the resources used by the persistent path init container
	Resources *PodResources `yaml:"resources,omitempty" json:"resources,omitempty" jsonschema_description:"Resources are the resources used by the persistent path init container."`
}

// PortMapping defines the ports for a PortMapping
type PortMapping struct {
	// Port is a port mapping that maps the localPort:remotePort. So if
	// you port forward the remote port will be available at the local port.
	// If you do reverse port forwarding, the local port will be available
	// at the remote port in the container. If only port is specified, local and
	// remote port are the same.
	Port string `yaml:"port" json:"port" jsonschema_description:"Port is a port mapping that maps the localPort:remotePort."`

	// BindAddress is the address DevSpace should listen on. Optional and defaults
	// to localhost.
	BindAddress string `yaml:"bindAddress,omitempty" json:"bindAddress,omitempty" jsonschema_description:"BindAddress is the address DevSpace should listen on. Optional and defaults to localhost."`
}

// OpenConfig defines what to open after services have been started
type OpenConfig struct {
	// URL is the url to open in the browser after it is available
	URL string `yaml:"url,omitempty" json:"url,omitempty" jsonschema_description:"URL is the url to open in the browser after it is available."`
}

// SyncConfig defines the paths for a SyncFolder
type SyncConfig struct {
	// Path is the path to sync. This can be defined in the form localPath:remotePath. You can also use '.'
	// to specify either the local or remote working directory. This is valid for example: .:.
	Path string `yaml:"path,omitempty" json:"path,omitempty" jsonschema_description:"Path is the path to sync. This can be defined in the form localPath:remotePath. You can also use '.' to specify either the local or remote working directory. This is valid for example: .:."`

	// ExcludePaths is an array of file patterns in gitignore format to exclude.
	ExcludePaths []string `yaml:"excludePaths,omitempty" json:"excludePaths,omitempty" jsonschema_extras:"group=exclude,group_name=Exclude Paths From File Sync" jsonschema_description:"ExcludePaths is an array of file patterns in gitignore format to exclude."`
	// ExcludeFile loads the file patterns to exclude from a file.
	ExcludeFile string `yaml:"excludeFile,omitempty" json:"excludeFile,omitempty" jsonschema_extras:"group=exclude" jsonschema_description:"ExcludeFile loads the file patterns to exclude from a file."`
	// DownloadExcludePaths is an array of file patterns in gitignore format to exclude from downloading
	DownloadExcludePaths []string `yaml:"downloadExcludePaths,omitempty" json:"downloadExcludePaths,omitempty" jsonschema_extras:"group=exclude" jsonschema_description:"DownloadExcludePaths is an array of file patterns in gitignore format to exclude from downloading."`
	// DownloadExcludeFile loads the file patterns to exclude from downloading from a file.
	DownloadExcludeFile string `yaml:"downloadExcludeFile,omitempty" json:"downloadExcludeFile,omitempty" jsonschema_extras:"group=exclude" jsonschema_description:"DownloadExcludeFile loads the file patterns to exclude from downloading from a file."`
	// UploadExcludePaths is an array of file patterns in gitignore format to exclude from uploading
	UploadExcludePaths []string `yaml:"uploadExcludePaths,omitempty" json:"uploadExcludePaths,omitempty" jsonschema_extras:"group=exclude" jsonschema_description:"UploadExcludePaths is an array of file patterns in gitignore format to exclude from uploading."`
	// UploadExcludeFile loads the file patterns to exclude from uploading from a file.
	UploadExcludeFile string `yaml:"uploadExcludeFile,omitempty" json:"uploadExcludeFile,omitempty" jsonschema_extras:"group=exclude" jsonschema_description:"UploadExcludeFile loads the file patterns to exclude from uploading from a file."`

	// StartContainer will start the container after initial sync is done. This will
	// inject a devspacehelper into the pod and you need to define dev.*.command for
	// this to work.
	StartContainer bool `yaml:"startContainer,omitempty" json:"startContainer,omitempty" jsonschema_extras:"group=actions,group_name=Sync-Triggered Actions" jsonschema_description:"StartContainer will start the container after initial sync is done. This will inject a devspacehelper into the pod and you need to define dev.*.command for this to work."`

	// OnUpload can be used to execute certain commands on uploading either in the container or locally as
	// well as restart the container after a file changed has happened.
	OnUpload *SyncOnUpload `yaml:"onUpload,omitempty" json:"onUpload,omitempty" jsonschema_extras:"group=actions" jsonschema_description:"OnUpload can be used to execute certain commands on uploading either in the container or locally as well as restart the container after a file changed has happened."`

	// InitialSync defines the initial sync strategy to use when this sync starts. Defaults to mirrorLocal
	InitialSync InitialSyncStrategy `yaml:"initialSync,omitempty" json:"initialSync,omitempty" jsonschema_extras:"group=initial_sync,group_name=Initial Sync" jsonschema_description:"InitialSync defines the initial sync strategy to use when this sync starts. Defaults to mirrorLocal."`

	// WaitInitialSync can be used to tell DevSpace to not wait until the initial sync is done
	WaitInitialSync *bool `yaml:"waitInitialSync,omitempty" json:"waitInitialSync,omitempty" jsonschema_extras:"group=initial_sync" jsonschema_description:"WaitInitialSync can be used to tell DevSpace to not wait until the initial sync is done."`

	// InitialSyncCompareBy defines if the sync should only compare by the given type. Either mtime or size are possible
	InitialSyncCompareBy InitialSyncCompareBy `yaml:"initialSyncCompareBy,omitempty" json:"initialSyncCompareBy,omitempty" jsonschema_extras:"group=initial_sync" jsonschema_description:"InitialSyncCompareBy defines if the sync should only compare by the given type. Either mtime or size are possible."`

	// DisableDownload will disable downloading completely
	DisableDownload bool `yaml:"disableDownload,omitempty" json:"disableDownload,omitempty" jsonschema_extras:"group=one_direction,group_name=One-Directional Sync" jsonschema_description:"DisableDownload will disable downloading completely."`
	// DisableUpload will disable uploading completely
	DisableUpload bool `yaml:"disableUpload,omitempty" json:"disableUpload,omitempty" jsonschema_extras:"group=one_direction" jsonschema_description:"DisableUpload will disable uploading completely."`

	// BandwidthLimits can be used to limit the amount of bytes that are transferred by DevSpace with this
	// sync configuration
	BandwidthLimits *BandwidthLimits `yaml:"bandwidthLimits,omitempty" json:"bandwidthLimits,omitempty" jsonschema_description:"BandwidthLimits can be used to limit the amount of bytes that are transferred by DevSpace with this sync configuration."`

	// Polling will tell the remote container to use polling instead of inotify
	Polling bool `yaml:"polling,omitempty" json:"polling,omitempty" jsonschema_description:"Polling will tell the remote container to use polling instead of inotify."`

	// NoWatch will terminate the sync after the initial sync is done
	NoWatch bool `yaml:"noWatch,omitempty" json:"noWatch,omitempty" jsonschema_description:"NoWatch will terminate the sync after the initial sync is done."`

	// PrintLogs defines if sync logs should be displayed on the terminal
	PrintLogs bool `yaml:"printLogs,omitempty" json:"printLogs,omitempty" jsonschema_description:"PrintLogs defines if sync logs should be displayed on the terminal."`
}

type ContainerArchitecture string

const (
	ContainerArchitectureAmd64 ContainerArchitecture = "amd64"
	ContainerArchitectureArm64 ContainerArchitecture = "arm64"
)

// SyncOnUpload defines the struct for the command that should be executed when files / folders are uploaded
type SyncOnUpload struct {
	// If true restart container will try to restart the container after a change has been made. Make sure that
	// images.*.injectRestartHelper is enabled for the container that should be restarted or the devspace-restart-helper
	// script is present in the container root folder.
	RestartContainer bool `yaml:"restartContainer,omitempty" json:"restartContainer,omitempty" jsonschema_description:"If true restart container will try to restart the container after a change has been made. Make sure that images.*.injectRestartHelper is enabled for the container that should be restarted or the devspace-restart-helper script is present in the container root folder."`

	// Exec will execute the given commands in order after a sync operation
	Exec []SyncExec `yaml:"exec,omitempty" json:"exec,omitempty" jsonschema_description:"Exec will execute the given commands in order after a sync operation."`

	// Defines what commands should be executed on the container side if a change is uploaded and applied in the target
	// container
	ExecRemote *SyncExecCommand `yaml:"execRemote,omitempty" json:"execRemote,omitempty" jsonschema_description:"Defines what commands should be executed on the container side if a change is uploaded and applied in the target container."`
}

type SyncExec struct {
	// Name is the name to show for this exec in the logs
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name is the name to show for this exec in the logs."`

	// Command is the command to execute. If no args are specified this is executed
	// within a shell.
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command is the command to execute. If no args are specified this is executed within a shell."`

	// Args are arguments to pass to the command
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are arguments to pass to the command."`

	// FailOnError specifies if the sync should fail if the command fails
	FailOnError bool `yaml:"failOnError,omitempty" json:"failOnError,omitempty" jsonschema_description:"FailOnError specifies if the sync should fail if the command fails."`

	// Local specifies if the command should be executed locally instead of within the
	// container
	Local bool `yaml:"local,omitempty" json:"local,omitempty" jsonschema_description:"Local specifies if the command should be executed locally instead of within the container."`

	// OnChange is an array of file patterns that trigger this command execution
	OnChange []string `yaml:"onChange,omitempty" json:"onChange,omitempty" jsonschema_description:"OnChange is an array of file patterns that trigger this command execution."`
}

// SyncExecCommand holds the configuration of commands that should be executed when files / folders are change
type SyncExecCommand struct {
	// Command is the command that should get executed
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command is the command that should get executed."`
	// Args are arguments that should get appended to the command
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are arguments that should get appended to the command."`

	// OnFileChange is invoked after every file change. DevSpace will wait for the command to successfully finish
	// and then will continue to upload files & create folders
	OnFileChange *SyncCommand `yaml:"onFileChange,omitempty" json:"onFileChange,omitempty" jsonschema_description:"OnFileChange is invoked after every file change. DevSpace will wait for the command to successfully finish and then will continue to upload files & create folders."`

	// OnDirCreate is invoked after every directory that is created. DevSpace will wait for the command to successfully finish
	// and then will continue to upload files & create folders
	OnDirCreate *SyncCommand `yaml:"onDirCreate,omitempty" json:"onDirCreate,omitempty" jsonschema_description:"OnDirCreate is invoked after every directory that is created. DevSpace will wait for the command to successfully finish and then will continue to upload files & create folders."`

	// OnBatch executes the given command after a batch of changes has been processed. DevSpace will wait for the command to finish
	// and then will continue execution. This is useful for commands
	// that shouldn't be executed after every single change that may take a little bit longer like recompiling etc.
	OnBatch *SyncCommand `yaml:"onBatch,omitempty" json:"onBatch,omitempty" jsonschema_description:"OnBatch executes the given command after a batch of changes has been processed. DevSpace will wait for the command to finish and then will continue execution. This is useful for commands that shouldn't be executed after every single change that may take a little bit longer like recompiling etc."`
}

// SyncCommand holds a command definition
type SyncCommand struct {
	// Command is the command that should get executed
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command is the command that should get executed."`

	// Args are arguments that should get appended to the command
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are arguments that should get appended to the command."`
}

// InitialSyncStrategy is the type of a initial sync strategy
type InitialSyncStrategy string

// List of values that source can take
const (
	InitialSyncStrategyMirrorLocal  InitialSyncStrategy = "mirrorLocal"
	InitialSyncStrategyMirrorRemote InitialSyncStrategy = "mirrorRemote"
	InitialSyncStrategyPreferLocal  InitialSyncStrategy = "preferLocal"
	InitialSyncStrategyPreferRemote InitialSyncStrategy = "preferRemote"
	InitialSyncStrategyPreferNewest InitialSyncStrategy = "preferNewest"
	InitialSyncStrategyKeepAll      InitialSyncStrategy = "keepAll"
)

// InitialSyncCompareBy is the type of how a change should be determined during the initial sync
type InitialSyncCompareBy string

// List of values that compare by can take
const (
	InitialSyncCompareByMTime InitialSyncCompareBy = "mtime"
	InitialSyncCompareBySize  InitialSyncCompareBy = "size"
)

// BandwidthLimits defines the struct for specifying the sync bandwidth limits
type BandwidthLimits struct {
	// Download is the download limit in kilo bytes per second
	Download *int64 `yaml:"download,omitempty" json:"download,omitempty" jsonschema_description:"Download is the download limit in kilo bytes per second."`

	// Upload is the upload limit in kilo bytes per second
	Upload *int64 `yaml:"upload,omitempty" json:"upload,omitempty" jsonschema_description:"Upload is the upload limit in kilo bytes per second."`
}

// Terminal describes the terminal options
type Terminal struct {
	// Command is the command that should be executed on terminal start.
	// This command is executed within a shell.
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command is the command that should be executed on terminal start. This command is executed within a shell."`

	// WorkDir is the working directory that is used to execute the command in.
	WorkDir string `yaml:"workDir,omitempty" json:"workDir,omitempty" jsonschema_description:"WorkDir is the working directory that is used to execute the command in."`

	// If enabled is true, DevSpace will use the terminal. Can be also
	// used to disable the terminal if set to false. DevSpace makes sure
	// that within a pipeline only one dev configuration can open a terminal
	// at a time and subsequent dev terminals will fail.
	Enabled *bool `yaml:"enabled,omitempty" json:"enabled,omitempty" jsonschema_description:" If enabled is true, DevSpace will use the terminal. It can be also used to disable the terminal if set to false. DevSpace makes sure that within a pipeline only one dev configuration can open a terminal at a time and subsequent dev terminals will fail."`

	// DisableReplace tells DevSpace to not replace the pod or adjust its settings
	// to make sure the pod is sleeping when opening a terminal
	DisableReplace bool `yaml:"disableReplace,omitempty" json:"disableReplace,omitempty" jsonschema_description:"DisableReplace tells DevSpace to not replace the pod or adjust its settings to make sure the pod is sleeping when opening a terminal."`

	// DisableScreen will disable screen which is used by DevSpace by default to preserve
	// sessions if connections interrupt or the session is lost.
	DisableScreen bool `yaml:"disableScreen,omitempty" json:"disableScreen,omitempty" jsonschema_description:"DisableScreen will disable screen which is used by DevSpace by default to preserve sessions if connections interrupt or the session is lost."`
}

// DependencyConfig defines the devspace dependency
type DependencyConfig struct {
	// Name is used internally
	Name string `yaml:"name" json:"name" jsonschema_description:"Name is used internally."`

	// Source holds the dependency project
	Source *SourceConfig `yaml:",inline" json:",inline" jsonschema_description:"Source holds the dependency project."`

	// Pipeline is the pipeline to deploy by default. Defaults to 'deploy'
	Pipeline string `yaml:"pipeline,omitempty" json:"pipeline,omitempty" jsonschema:"default=deploy" jsonschema_extras:"group=execution,group_name=Execution" jsonschema_description:"Pipeline is the pipeline to deploy by default. Defaults to 'deploy'."`

	// Vars are variables that should be passed to the dependency
	Vars map[string]string `yaml:"vars,omitempty" json:"vars,omitempty" jsonschema_extras:"group=execution" jsonschema_description:"Vars are variables that should be passed to the dependency."`

	// OverwriteVars specifies if DevSpace should pass the parent variables to the dependency
	OverwriteVars bool `yaml:"overwriteVars,omitempty" json:"overwriteVars,omitempty" jsonschema_extras:"group=execution" jsonschema_description:"OverwriteVars specifies if DevSpace should pass the parent variables to the dependency."`

	// IgnoreDependencies defines if dependencies of the dependency should be excluded
	IgnoreDependencies bool `yaml:"ignoreDependencies,omitempty" json:"ignoreDependencies,omitempty" jsonschema_extras:"group=execution" jsonschema_description:"IgnoreDependencies defines if dependencies of the dependency should be excluded."`

	// Namespace specifies the namespace this dependency should be deployed to
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema_extras:"group=execution" jsonschema_description:"Namespace specifies the namespace this dependency should be deployed to."`

	// Profiles specifies which profiles should be applied while loading the dependency
	Profiles []string `yaml:"profiles,omitempty" json:"profiles,omitempty" jsonschema:"-" jsonschema_description:"Profiles specifies which profiles should be applied while loading the dependency."`

	// DisableProfileActivation disabled automatic profile activation of dependency profiles
	DisableProfileActivation bool `yaml:"disableProfileActivation,omitempty" json:"disableProfileActivation,omitempty" jsonschema:"-" jsonschema_description:"DisableProfileActivation disabled automatic profile activation of dependency profiles."`
}

// SourceConfig defines an artifact source
type SourceConfig struct {
	// Path is the local path where DevSpace can find the artifact.
	// This option is mutually exclusive with the git option.
	Path string `yaml:"path,omitempty" json:"path,omitempty" jsonschema_extras:"group=path,group_name=Source: Local Filesystem" jsonschema_description:"Path is the local path where DevSpace can find the artifact. This option is mutually exclusive with the git option."`

	// Git is the remote repository to download the artifact from. You can either use
	// https projects or ssh projects here, but need to make sure git can pull the project.
	// This option is mutually exclusive with the path option.
	Git string `yaml:"git,omitempty" json:"git,omitempty" jsonschema_extras:"group=git,group_name=Source: Git Repository" jsonschema_description:"Git is the remote repository to download the artifact from. You can either use https projects or ssh projects here, but need to make sure git can pull the project. This option is mutually exclusive with the path option."`

	// SubPath is a path within the git repository where the artifact lies in
	SubPath string `yaml:"subPath,omitempty" json:"subPath,omitempty" jsonschema_extras:"group=git" jsonschema_description:"SubPath is a path within the git repository where the artifact lies in."`

	// Branch is the git branch to pull
	Branch string `yaml:"branch,omitempty" json:"branch,omitempty" jsonschema_extras:"group=git" jsonschema_description:"Branch is the git branch to pull."`

	// Tag is the tag to pull
	Tag string `yaml:"tag,omitempty" json:"tag,omitempty" jsonschema_extras:"group=git" jsonschema_description:"Tag is the tag to pull."`

	// Revision is the git revision to pull
	Revision string `yaml:"revision,omitempty" json:"revision,omitempty" jsonschema_extras:"group=git" jsonschema_description:"Revision is the git revision to pull."`

	// CloneArgs are additional arguments that should be supplied to the git CLI
	CloneArgs []string `yaml:"cloneArgs,omitempty" json:"cloneArgs,omitempty" jsonschema_extras:"group=git" jsonschema_description:"CloneArgs are additional arguments that should be supplied to the git CLI."`

	// DisableShallow can be used to turn off shallow clones as these are the default used
	// by devspace
	DisableShallow bool `yaml:"disableShallow,omitempty" json:"disableShallow,omitempty" jsonschema_extras:"group=git" jsonschema_description:"DisableShallow can be used to turn off shallow clones as these are the default used by devspace."`

	// DisablePull will disable pulling every time DevSpace is reevaluating this source
	DisablePull bool `yaml:"disablePull,omitempty" json:"disablePull,omitempty" jsonschema_extras:"group=git" jsonschema_description:"DisablePull will disable pulling every time DevSpace is reevaluating this source."`
}

// HookConfig defines a hook
type HookConfig struct {
	// Name is the name of the hook
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name is the name of the hook."`

	// Disabled can be used to disable the hook
	Disabled bool `yaml:"disabled,omitempty" json:"disabled,omitempty" jsonschema_description:"Disabled can be used to disable the hook."`

	// Events are the events when the hook should be executed
	Events []string `yaml:"events" json:"events" jsonschema_description:"Events are the events when the hook should be executed."`

	// Command is the base command that is either executed locally or in a remote container.
	// Command is mutually exclusive with other hook actions. In the case this is defined
	// together with where.container, DevSpace will until the target container is running and
	// only then execute the command. If the container does not start in time, DevSpace will fail.
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command is the base command that is either executed locally or in a remote container. Command is mutually exclusive with other hook actions. In the case this is defined together with where.container, DevSpace will until the target container is running and only then execute the command. If the container does not start in time, DevSpace will fail."`

	// Args are additional arguments passed together with the command to execute.
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are additional arguments passed together with the command to execute."`

	// If an operating system is defined, the hook will only be executed for the given os.
	// All supported golang OS types are supported and multiple can be combined with ','.
	OperatingSystem string `yaml:"os,omitempty" json:"os,omitempty" jsonschema_description:"If an operating system is defined, the hook will only be executed for the given os. All supported golang OS types are supported and multiple can be combined with ','."`

	// If Upload is specified, DevSpace will upload certain local files or folders into a
	// remote container.
	Upload *HookSyncConfig `yaml:"upload,omitempty" json:"upload,omitempty" jsonschema_description:" If Upload is specified, DevSpace will upload certain local files or folders into a remote container."`
	// Same as Upload, but with this option DevSpace will download files or folders from
	// a remote container.
	Download *HookSyncConfig `yaml:"download,omitempty" json:"download,omitempty" jsonschema_description:"Same as Upload, but with this option DevSpace will download files or folders from a remote container."`
	// If logs is defined will print the logs of the target container. This is useful for containers
	// that should finish like init containers or job pods. Otherwise this hook will never terminate.
	Logs *HookLogsConfig `yaml:"logs,omitempty" json:"logs,omitempty" jsonschema_description:"If logs is defined will print the logs of the target container. This is useful for containers that should finish like init containers or job pods. Otherwise this hook will never terminate."`
	// If wait is defined the hook will wait until the matched pod or container is running or is terminated
	// with a certain exit code.
	Wait *HookWaitConfig `yaml:"wait,omitempty" json:"wait,omitempty" jsonschema_description:"If wait is defined the hook will wait until the matched pod or container is running or is terminated with a certain exit code."`

	// If true, the hook will be executed in the background.
	Background bool `yaml:"background,omitempty" json:"background,omitempty" jsonschema_description:"If true, the hook will be executed in the background."`
	// If true, the hook will not output anything to the standard out of DevSpace except
	// for the case when the hook fails, where DevSpace will show the error including
	// the captured output streams of the hook.
	Silent bool `yaml:"silent,omitempty" json:"silent,omitempty" jsonschema_description:"If true, the hook will not output anything to the standard out of DevSpace except for the case when the hook fails, where DevSpace will show the error including the captured output streams of the hook."`

	// Container specifies where the hook should be run. If this is omitted DevSpace expects a
	// local command hook.
	Container *HookContainer `yaml:"container,omitempty" json:"container,omitempty" jsonschema_description:"Container specifies where the hook should be run. If this is omitted DevSpace expects a local command hook."`
}

// HookWaitConfig defines a hook wait config
type HookWaitConfig struct {
	// If running is true, will wait until the matched containers are running. Can be used together with terminatedWithCode.
	Running bool `yaml:"running,omitempty" json:"running,omitempty" jsonschema_description:" If running is true, will wait until the matched containers are running. Can be used together with terminatedWithCode."`

	// If terminatedWithCode is not nil, will wait until the matched containers are terminated with the given exit code.
	// If the container has exited with a different exit code, the hook will fail. Can be used together with running.
	TerminatedWithCode *int32 `yaml:"terminatedWithCode,omitempty" json:"terminatedWithCode,omitempty" jsonschema_description:"If terminatedWithCode is not nil, will wait until the matched containers are terminated with the given exit code. If the container has exited with a different exit code, the hook will fail. Can be used together with running."`

	// The amount of seconds to wait until the hook will fail. Defaults to 150 seconds.
	Timeout int64 `yaml:"timeout,omitempty" json:"timeout,omitempty" jsonschema_description:"The amount of seconds to wait until the hook will fail. Defaults to 150 seconds."`
}

// HookLogsConfig defines a hook logs config
type HookLogsConfig struct {
	// If set, the number of lines from the end of the logs to show. If not specified,
	// logs are shown from the creation of the container
	TailLines *int64 `yaml:"tailLines,omitempty" json:"tailLines,omitempty" jsonschema_description:"If set, the number of lines from the end of the logs to show. If not specified, logs are shown from the creation of the container."`
}

// HookSyncConfig defines a hook upload config
type HookSyncConfig struct {
	// LocalPath to sync files from
	LocalPath string `yaml:"localPath,omitempty" json:"localPath,omitempty" jsonschema_description:"LocalPath to sync files from."`

	// ContainerPath to sync files to
	ContainerPath string `yaml:"containerPath,omitempty" json:"containerPath,omitempty" jsonschema_description:"ontainerPath to sync files to."`
}

// HookContainer defines how to select one or more containers to execute a hook in
type HookContainer struct {
	// LabelSelector to select a container
	LabelSelector map[string]string `yaml:"labelSelector,omitempty" json:"labelSelector,omitempty" jsonschema_description:"HookContainer defines how to select one or more containers to execute a hook in."`

	// Pod name to use
	Pod string `yaml:"pod,omitempty" json:"pod,omitempty" jsonschema_description:"Pod name to use."`

	// Namespace to use
	Namespace string `yaml:"namespace,omitempty" json:"namespace,omitempty" jsonschema_description:"Namespace to use."`

	// ImageSelector to select a container
	ImageSelector string `yaml:"imageSelector,omitempty" json:"imageSelector,omitempty" jsonschema_description:"ImageSelector to select a container."`

	// ContainerName to use
	ContainerName string `yaml:"containerName,omitempty" json:"containerName,omitempty" jsonschema_description:"ContainerName to use."`

	// Wait can be used to disable waiting
	Wait *bool `yaml:"wait,omitempty" json:"wait,omitempty" jsonschema_description:"Wait can be used to disable waiting."`

	// Timeout how long to wait
	Timeout int64 `yaml:"timeout,omitempty" json:"timeout,omitempty" jsonschema_description:"Timeout how long to wait."`

	// Once only executes an hook once in the container until it is restarted
	Once *bool `yaml:"once,omitempty" json:"once,omitempty" jsonschema_description:"Once only executes an hook once in the container until it is restarted."`
}

// CommandConfig defines the command specification
type CommandConfig struct {
	// Name is the name of a command that is used via `devspace run NAME`
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name is the name of a command that is used via 'devspace run NAME'."`

	// Section can be used to group similar commands together in `devspace list commands`
	Section string `yaml:"section,omitempty" json:"section,omitempty" jsonschema_description:"Section can be used to group similar commands together in 'devspace list commands'."`

	// Command is the command that should be executed. For example: 'echo 123'
	Command string `yaml:"command" json:"command" jsonschema:"required" jsonschema_description:"Command is the command that should be executed. For example: 'echo 123'."`

	// Args are optional and if defined, command is not executed within a shell
	// and rather directly.
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are optional and if defined, command is not executed within a shell and rather directly."`

	// AppendArgs will append arguments passed to the DevSpace command automatically to
	// the specified command.
	AppendArgs bool `yaml:"appendArgs,omitempty" json:"appendArgs,omitempty" jsonschema_description:"AppendArgs will append arguments passed to the DevSpace command automatically to the specified command."`

	// Description describes what the command is doing and can be seen in `devspace list commands`
	Description string `yaml:"description,omitempty" json:"description,omitempty" jsonschema_description:"Description describes what the command is doing and can be seen in 'devspace list commands'."`

	// Internal commands are not show in list and are usable through run_command
	Internal bool `yaml:"internal,omitempty" json:"internal,omitempty" jsonschema_description:"Internal commands are not show in list and are usable through run_command."`

	// After is executed after the command was run. It is executed also when
	// the command was interrupted which will set the env variable COMMAND_INTERRUPT
	// to true as well as when the command errored which will set the error string to
	// COMMAND_ERROR.
	After string `yaml:"after,omitempty" json:"after,omitempty" jsonschema_description:"After is executed after the command was run. It is executed also when the command was interrupted which will set the env variable COMMAND_INTERRUPT to true as well as when the command errored which will set the error string to COMMAND_ERROR."`
}

type CommandFlag struct {
}

func (c *CommandConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	commandString := ""
	err := unmarshal(&commandString)
	if err != nil {
		m := map[string]interface{}{}
		err := unmarshal(m)
		if err != nil {
			return err
		}

		out, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return yamlutil.UnmarshalStrictJSON(out, c)
	}

	c.Command = commandString
	return nil
}

// Variable describes the var definition
type Variable struct {
	// Name is the name of the variable
	Name string `yaml:"name" json:"name" jsonschema_description:"Name is the name of the variable."`

	// Value is a shortcut for using source: none and default: my-value
	Value interface{} `yaml:"value,omitempty" json:"value,omitempty" jsonschema:"oneof_type=string;integer;boolean" jsonschema_extras:"group=static,group_name=Static Value" jsonschema_description:"Value is a shortcut for using source: none and default: my-value."`

	// Question can be used to define a custom question if the variable was not yet used
	Question string `yaml:"question,omitempty" json:"question,omitempty" jsonschema_extras:"group=question,group_name=Value From Input (Question)" jsonschema_description:"Question can be used to define a custom question if the variable was not yet used."`

	// Default is the default value the variable should have if not set by the user
	Default interface{} `yaml:"default,omitempty" json:"default,omitempty" jsonschema:"oneof_type=string;integer;boolean" jsonschema_extras:"group=question" jsonschema_description:"Default is the default value the variable should have if not set by the user."`

	// Options are options that can be selected when the variable question is asked
	Options []string `yaml:"options,omitempty" json:"options,omitempty" jsonschema_extras:"group=question" jsonschema_description:"Options are options that can be selected when the variable question is asked."`

	// Password signals that this variable should not be visible if entered
	Password bool `yaml:"password,omitempty" json:"password,omitempty" jsonschema_extras:"group=question" jsonschema_description:"Password signals that this variable should not be visible if entered."`

	// ValidationPattern can be used to verify the user input
	ValidationPattern string `yaml:"validationPattern,omitempty" json:"validationPattern,omitempty" jsonschema_extras:"group=question" jsonschema_description:"ValidationPattern can be used to verify the user input."`

	// ValidationMessage can be used to tell the user the format of the variable value
	ValidationMessage string `yaml:"validationMessage,omitempty" json:"validationMessage,omitempty" jsonschema_extras:"group=question" jsonschema_description:"ValidationMessage can be used to tell the user the format of the variable value."`

	// NoCache can be used to prompt the user on every run for this variable
	NoCache bool `yaml:"noCache,omitempty" json:"noCache,omitempty" jsonschema_extras:"group=question" jsonschema_description:"NoCache can be used to prompt the user on every run for this variable."`

	// Command is the command how to retrieve the variable. If args is omitted, command is parsed as a shell
	// command.
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_extras:"group=execution,group_name=Value From Command" jsonschema_description:"Command is the command how to retrieve the variable. If args is omitted, command is parsed as a shell command."`

	// Args are optional args that will be used for the command
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_extras:"group=execution" jsonschema_description:"Args are optional args that will be used for the command."`

	// Commands are additional commands that can be used to run a different command on a different operating
	// system.
	Commands []VariableCommand `yaml:"commands,omitempty" json:"commands,omitempty" jsonschema_extras:"group=execution" jsonschema_description:"Commands are additional commands that can be used to run a different command on a different operating system."`

	// AlwaysResolve makes sure this variable will always be resolved and not only if it is used somewhere. Defaults to true.
	AlwaysResolve *bool `yaml:"alwaysResolve,omitempty" json:"alwaysResolve,omitempty" jsonschema_description:"AlwaysResolve makes sure this variable will always be resolved and not only if it is used somewhere."`

	// Source defines where the variable should be taken from
	Source VariableSource `yaml:"source,omitempty" json:"source,omitempty" jsonschema_description:"Source defines where the variable should be taken from." jsonschema:"enum=all,enum=env,enum=input,enum=command,enum=none"`
}

func (v *Variable) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// try string next
	varString := ""
	err := unmarshal(&varString)
	if err != nil {
		m := map[string]interface{}{}
		err := unmarshal(m)
		if err != nil {
			return err
		}

		out, err := json.Marshal(m)
		if err != nil {
			return err
		}

		return yamlutil.UnmarshalStrictJSON(out, v)
	}
	if strings.HasPrefix(varString, "$(") && strings.HasSuffix(varString, ")") {
		varString = strings.TrimPrefix(strings.TrimSuffix(varString, ")"), "$(")
		v.Command = varString
		return nil
	}

	v.Value = varString
	return nil
}

type VariableCommand struct {
	// OperatingSystem is optional and defines the operating system this
	// command should be executed on
	OperatingSystem string `yaml:"os,omitempty" json:"os,omitempty" jsonschema_description:"OperatingSystem is optional and defines the operating system this command should be executed on."`

	// Command is the command to use to retrieve the value for this variable. If no
	// args are specified the command is run within a pseudo shell.
	Command string `yaml:"command,omitempty" json:"command,omitempty" jsonschema_description:"Command is the command to use to retrieve the value for this variable. If no args are specified the command is run within a pseudo shell."`

	// Args are optional arguments for the command
	Args []string `yaml:"args,omitempty" json:"args,omitempty" jsonschema_description:"Args are optional arguments for the command."`
}

// VariableSource is type of a variable source
type VariableSource string

// List of values that source can take
const (
	VariableSourceDefault VariableSource = ""
	VariableSourceAll     VariableSource = "all"
	VariableSourceEnv     VariableSource = "env"
	VariableSourceInput   VariableSource = "input"
	VariableSourceCommand VariableSource = "command"
	VariableSourceNone    VariableSource = "none"
)

// ProfileConfig defines a profile config
type ProfileConfig struct {
	// Name is the name of the profile
	Name string `yaml:"name" json:"name" jsonschema_description:"Args are optional arguments for the command."`

	// Description is the profile description
	Description string `yaml:"description,omitempty" json:"description,omitempty" jsonschema_description:"Description is the profile description."`

	// Parent is the profile parent within this config
	Parent string `yaml:"parent,omitempty" json:"parent,omitempty" jsonschema_description:"Parent is the profile parent within this config."`

	// Parents are additional profile parents that should get loaded
	Parents []*ProfileParent `yaml:"parents,omitempty" json:"parents,omitempty" jsonschema_description:"Parents are additional profile parents that should get loaded."`

	// Patches are patches to apply as part of this profile
	Patches []*PatchConfig `yaml:"patches,omitempty" json:"patches,omitempty" jsonschema_description:"Patches are patches to apply as part of this profile."`

	// Replace are config replacements as part of this profile
	Replace *ProfileConfigStructure `yaml:"replace,omitempty" json:"replace,omitempty" jsonschema_description:"Replace are config replacements as part of this profile."`

	// Merge will merge these parts with the current config as part of this profile
	Merge *ProfileConfigStructure `yaml:"merge,omitempty" json:"merge,omitempty" jsonschema_description:"Merge will merge these parts with the current config as part of this profile."`

	// Activation defines how and when the profile should get activated automatically
	Activation []*ProfileActivation `yaml:"activation,omitempty" json:"activation,omitempty" jsonschema_description:"Activation defines how and when the profile should get activated automatically."`
}

// ProfileConfigStructure is the base structure used to validate profiles
type ProfileConfigStructure struct {
	// Vars references variables
	Vars *map[string]interface{} `yaml:"vars,omitempty" json:"vars,omitempty" jsonschema_description:"Vars references variables."`

	// PullSecrets references pull secrets
	PullSecrets *map[string]interface{} `yaml:"pullSecrets,omitempty" json:"pullSecrets,omitempty" jsonschema_description:"PullSecrets references pull secrets."`

	// Images references images
	Images *map[string]interface{} `yaml:"images,omitempty" json:"images,omitempty" jsonschema_description:"Images references images."`

	// Deployments references the deployments section
	Deployments *map[string]interface{} `yaml:"deployments,omitempty" json:"deployments,omitempty" jsonschema_description:"Deployments references the deployments section."`

	// Dev references the dev section
	Dev *map[string]interface{} `yaml:"dev,omitempty" json:"dev,omitempty" jsonschema_description:"Dev references the dev section."`

	// Commands references the commands section
	Commands *map[string]interface{} `yaml:"commands,omitempty" json:"commands,omitempty" jsonschema_description:"Commands references the commands section."`

	// Dependencies references the dependencies section
	Dependencies *map[string]interface{} `yaml:"dependencies,omitempty" json:"dependencies,omitempty" jsonschema_description:"Dependencies references the dependencies section."`

	// Hooks references the hooks section
	Hooks *[]interface{} `yaml:"hooks,omitempty" json:"hooks,omitempty" jsonschema_description:"Hooks references the hooks section."`

	// DEPRECATED: OldDeployments references the old deployments
	OldDeployments *[]interface{} `yaml:"oldDeployments,omitempty" json:"oldDeployments,omitempty" jsonschema:"-" jsonschema_description:"DEPRECATED: OldDeployments references the old deployments."`

	// DEPRECATED: OldDependencies references the old dependencies
	OldDependencies *[]interface{} `yaml:"oldDependencies,omitempty" json:"oldDependencies,omitempty" jsonschema:"-" jsonschema_description:"DEPRECATED: OldDependencies references the old dependencies."`

	// DEPRECATED: OldCommands references the old commands
	OldCommands *[]interface{} `yaml:"oldCommands,omitempty" json:"oldCommands,omitempty" jsonschema:"-" jsonschema_description:"DEPRECATED: OldCommands references the old commands."`

	// DEPRECATED: OldPullSecrets references the old pullsecrets
	OldPullSecrets *[]interface{} `yaml:"oldPullSecrets,omitempty" json:"oldPullSecrets,omitempty" jsonschema:"-" jsonschema_description:"DEPRECATED: OldPullSecrets references the old pullsecrets."`

	// DEPRECATED: OldVars references the old vars
	OldVars *[]interface{} `yaml:"oldVars,omitempty" json:"oldVars,omitempty" jsonschema:"-" jsonschema_description:"DEPRECATED: OldVars references the old vars."`
}

// ProfileParent defines where to load the profile from
type ProfileParent struct {
	// Source holds the source configuration for this profile parent
	Source *SourceConfig `yaml:"source,omitempty" json:"source,omitempty" jsonschema_description:"Source holds the source configuration for this profile parent."`

	// Profile holds the profile to load from this parent
	Profile string `yaml:"profile" json:"profile" jsonschema_description:"Profile holds the profile to load from this parent."`
}

// ProfileActivation defines rules that automatically activate a profile when evaluated to true
type ProfileActivation struct {
	// Environment defines key/value pairs where the key is the name of the environment variable and the value is a regular expression used to match the variable's value.
	// When multiple keys are specified, they must all evaluate to true to activate the profile.
	Environment map[string]string `yaml:"env,omitempty" json:"env,omitempty" jsonschema_description:"Environment defines key/value pairs where the key is the name of the environment variable and the value is a regular expression used to match the variable's value. When multiple keys are specified, they must all evaluate to true to activate the profile."`

	// Vars defines key/value pairs where the key is the name of the variable and the value is a regular expression used to match the variable's value.
	// When multiple keys are specified, they must all evaluate to true to activate the profile.
	Vars map[string]string `yaml:"vars,omitempty" json:"vars,omitempty" jsonschema_description:"Vars defines key/value pairs where the key is the name of the variable and the value is a regular expression used to match the variable's value. When multiple keys are specified, they must all evaluate to true to activate the profile."`
}

// PatchConfig describes a config patch and how it should be applied
type PatchConfig struct {
	// Operation is the path operation to do. Can be either replace, add or remove
	Operation string `yaml:"op" json:"op" jsonschema_description:"Operation is the path operation to do. Can be either replace, add or remove."`

	// Path is the config path to apply the patch to
	Path string `yaml:"path" json:"path" jsonschema_description:"Path is the config path to apply the patch to."`

	// Value is the value to use for this patch.
	Value interface{} `yaml:"value,omitempty" json:"value,omitempty" jsonschema_description:"Value is the value to use for this patch."`
}

// PullSecretConfig defines a pull secret that should be created by DevSpace
type PullSecretConfig struct {
	// Name is the pull secret name to deploy
	Name string `yaml:"name,omitempty" json:"name,omitempty" jsonschema_description:"Name is the pull secret name to deploy."`

	// The registry to create the image pull secret for.
	// Empty string == docker hub
	// e.g. gcr.io
	Registry string `yaml:"registry,omitempty" json:"registry" jsonschema:"required" jsonschema_description:"The registry to create the image pull secret for."`

	// The username of the registry. If this is empty, devspace will try
	// to receive the auth data from the local docker
	Username string `yaml:"username,omitempty" json:"username,omitempty" jsonschema_description:"The username of the registry. If this is empty, devspace will try to receive the auth data from the local docker."`

	// The password to use for the registry. If this is empty, devspace will
	// try to receive the auth data from the local docker
	Password string `yaml:"password,omitempty" json:"password,omitempty" jsonschema_description:"The password to use for the registry. If this is empty, devspace will try to receive the auth data from the local docker."`

	// The optional email to use
	Email string `yaml:"email,omitempty" json:"email,omitempty" jsonschema_description:"The optional email to use."`

	// The secret to create
	Secret string `yaml:"secret,omitempty" json:"secret,omitempty" jsonschema_description:"The secret to create."`

	// The service account to add the secret to
	ServiceAccounts []string `yaml:"serviceAccounts,omitempty" json:"serviceAccounts,omitempty" jsonschema_description:"The service account to add the secret to."`
}
