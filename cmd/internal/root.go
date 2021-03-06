package cmd

import (
	golog "log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"code.cloudfoundry.org/cf-operator/pkg/bosh/converter"
	"code.cloudfoundry.org/cf-operator/pkg/bosh/manifest"
	"code.cloudfoundry.org/cf-operator/pkg/kube/operator"
	"code.cloudfoundry.org/cf-operator/version"
	"code.cloudfoundry.org/quarks-utils/pkg/cmd"
	"code.cloudfoundry.org/quarks-utils/pkg/config"
	"code.cloudfoundry.org/quarks-utils/pkg/ctxlog"
)

const (
	// Port on which the controller-runtime manager listens
	managerPort  = 2999
	namespaceArg = "cf-operator-namespace"
)

var (
	log              *zap.SugaredLogger
	debugGracePeriod = time.Second * 5
)

func wrapError(err error, msg string) error {
	return errors.Wrap(err, "cf-operator command failed. "+msg)
}

var rootCmd = &cobra.Command{
	Use:   "cf-operator",
	Short: "cf-operator manages BOSH deployments on Kubernetes",
	RunE: func(_ *cobra.Command, args []string) error {
		log = cmd.Logger(zap.AddCallerSkip(1))
		defer log.Sync()

		restConfig, err := cmd.KubeConfig(log)
		if err != nil {
			return wrapError(err, "")
		}

		cfg := config.NewDefaultConfig(afero.NewOsFs())

		err = converter.SetupOperatorDockerImage(
			viper.GetString("docker-image-org"),
			viper.GetString("docker-image-repository"),
			viper.GetString("docker-image-tag"),
			corev1.PullPolicy(viper.GetString("docker-image-pull-policy")),
		)
		if err != nil {
			return wrapError(err, "")
		}

		watchNamespace := cmd.Namespaces(cfg, log, namespaceArg)

		manifest.SetBoshDNSDockerImage(viper.GetString("bosh-dns-docker-image"))
		manifest.SetClusterDomain(viper.GetString("cluster-domain"))

		log.Infof("Starting cf-operator %s with namespace %s", version.Version, watchNamespace)
		log.Infof("cf-operator docker image: %s", config.GetOperatorDockerImage())

		serviceHost := viper.GetString("operator-webhook-service-host")
		// Port on which the cf operator webhook kube service listens to.
		servicePort := viper.GetInt32("operator-webhook-service-port")
		useServiceRef := viper.GetBool("operator-webhook-use-service-reference")

		if serviceHost == "" && !useServiceRef {
			return wrapError(errors.New("couldn't determine webhook server"), "operator-webhook-service-host flag is not set (env variable: CF_OPERATOR_WEBHOOK_SERVICE_HOST)")
		}

		cfg.WebhookServerHost = serviceHost
		cfg.WebhookServerPort = servicePort
		cfg.WebhookUseServiceRef = useServiceRef
		cfg.MaxBoshDeploymentWorkers = viper.GetInt("max-boshdeployment-workers")
		cfg.MaxQuarksSecretWorkers = viper.GetInt("max-quarks-secret-workers")
		cfg.MaxQuarksStatefulSetWorkers = viper.GetInt("max-quarks-statefulset-workers")

		cmd.CtxTimeOut(cfg)

		ctx := ctxlog.NewParentContext(log)

		err = cmd.ApplyCRDs(ctx, operator.ApplyCRDs, restConfig)
		if err != nil {
			return wrapError(err, "Couldn't apply CRDs.")
		}

		mgr, err := operator.NewManager(ctx, cfg, restConfig, manager.Options{
			Namespace:          watchNamespace,
			MetricsBindAddress: "0",
			LeaderElection:     false,
			Port:               managerPort,
			Host:               "0.0.0.0",
		})
		if err != nil {
			return wrapError(err, "Failed to create new manager.")
		}

		ctxlog.Info(ctx, "Waiting for configurations to be applied into a BOSHDeployment resource...")

		err = mgr.Start(signals.SetupSignalHandler())
		if err != nil {
			return wrapError(err, "Failed to start cf-operator manager.")
		}
		return nil
	},
	TraverseChildren: true,
}

// NewCFOperatorCommand returns the `cf-operator` command.
func NewCFOperatorCommand() *cobra.Command {
	return rootCmd
}

// Execute the root command, runs the server
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		golog.Fatal(err)
		os.Exit(1)
	}
}

func init() {
	pf := rootCmd.PersistentFlags()

	argToEnv := map[string]string{}

	cmd.CtxTimeOutFlags(pf, argToEnv)
	cmd.KubeConfigFlags(pf, argToEnv)
	cmd.LoggerFlags(pf, argToEnv)
	cmd.NamespacesFlags(pf, argToEnv, namespaceArg)
	cmd.DockerImageFlags(pf, argToEnv, "cf-operator", version.Version)
	cmd.ApplyCRDsFlags(pf, argToEnv)

	pf.StringP("bosh-dns-docker-image", "", "coredns/coredns:1.6.3", "The docker image used for emulating bosh DNS (a CoreDNS image)")
	pf.String("cluster-domain", "cluster.local", "The Kubernetes cluster domain")
	pf.Int("max-boshdeployment-workers", 1, "Maximum number of workers concurrently running BOSHDeployment controller")
	pf.Int("max-quarks-secret-workers", 5, "Maximum number of workers concurrently running QuarksSecret controller")
	pf.Int("max-quarks-statefulset-workers", 1, "Maximum number of workers concurrently running QuarksStatefulSet controller")
	pf.StringP("operator-webhook-service-host", "w", "", "Hostname/IP under which the webhook server can be reached from the cluster")
	pf.StringP("operator-webhook-service-port", "p", "2999", "Port the webhook server listens on")
	pf.BoolP("operator-webhook-use-service-reference", "x", false, "If true the webhook service is targeted using a service reference instead of a URL")

	for _, name := range []string{
		"bosh-dns-docker-image",
		"cluster-domain",
		"max-boshdeployment-workers",
		"max-quarks-secret-workers",
		"max-quarks-statefulset-workers",
		"operator-webhook-service-host",
		"operator-webhook-service-port",
		"operator-webhook-use-service-reference",
	} {
		viper.BindPFlag(name, pf.Lookup(name))
	}
	argToEnv["bosh-dns-docker-image"] = "BOSH_DNS_DOCKER_IMAGE"
	argToEnv["cluster-domain"] = "CLUSTER_DOMAIN"
	argToEnv["max-boshdeployment-workers"] = "MAX_BOSHDEPLOYMENT_WORKERS"
	argToEnv["max-quarks-secret-workers"] = "MAX_QUARKS_SECRET_WORKERS"
	argToEnv["max-quarks-statefulset-workers"] = "MAX_QUARKS_STATEFULSET_WORKERS"
	argToEnv["operator-webhook-service-host"] = "CF_OPERATOR_WEBHOOK_SERVICE_HOST"
	argToEnv["operator-webhook-service-port"] = "CF_OPERATOR_WEBHOOK_SERVICE_PORT"
	argToEnv["operator-webhook-use-service-reference"] = "CF_OPERATOR_WEBHOOK_USE_SERVICE_REFERENCE"

	// Add env variables to help
	cmd.AddEnvToUsage(rootCmd, argToEnv)

	// Do not display cmd usage and errors
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
}
