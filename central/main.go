package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/NYTimes/gziphandler"
	alertDatastore "github.com/stackrox/rox/central/alert/datastore"
	alertService "github.com/stackrox/rox/central/alert/service"
	apiTokenExpiration "github.com/stackrox/rox/central/apitoken/expiration"
	apiTokenService "github.com/stackrox/rox/central/apitoken/service"
	"github.com/stackrox/rox/central/audit"
	authService "github.com/stackrox/rox/central/auth/service"
	"github.com/stackrox/rox/central/auth/userpass"
	authProviderDS "github.com/stackrox/rox/central/authprovider/datastore"
	authProviderSvc "github.com/stackrox/rox/central/authprovider/service"
	authProviderTelemetry "github.com/stackrox/rox/central/authprovider/telemetry"
	centralHealthService "github.com/stackrox/rox/central/centralhealth/service"
	"github.com/stackrox/rox/central/certgen"
	"github.com/stackrox/rox/central/cli"
	clusterDataStore "github.com/stackrox/rox/central/cluster/datastore"
	clusterService "github.com/stackrox/rox/central/cluster/service"
	"github.com/stackrox/rox/central/clusterinit/backend"
	clusterInitService "github.com/stackrox/rox/central/clusterinit/service"
	clustersHelmConfig "github.com/stackrox/rox/central/clusters/helmconfig"
	clustersZip "github.com/stackrox/rox/central/clusters/zip"
	complianceDatastore "github.com/stackrox/rox/central/compliance/datastore"
	complianceHandlers "github.com/stackrox/rox/central/compliance/handlers"
	complianceManagerService "github.com/stackrox/rox/central/compliance/manager/service"
	complianceService "github.com/stackrox/rox/central/compliance/service"
	configDS "github.com/stackrox/rox/central/config/datastore"
	configService "github.com/stackrox/rox/central/config/service"
	credentialExpiryService "github.com/stackrox/rox/central/credentialexpiry/service"
	clusterCveCsv "github.com/stackrox/rox/central/cve/cluster/csv"
	clusterCVEService "github.com/stackrox/rox/central/cve/cluster/service"
	"github.com/stackrox/rox/central/cve/csv"
	"github.com/stackrox/rox/central/cve/fetcher"
	imageCveCsv "github.com/stackrox/rox/central/cve/image/csv"
	imageCVEService "github.com/stackrox/rox/central/cve/image/service"
	nodeCveCsv "github.com/stackrox/rox/central/cve/node/csv"
	nodeCVEService "github.com/stackrox/rox/central/cve/node/service"
	cveService "github.com/stackrox/rox/central/cve/service"
	"github.com/stackrox/rox/central/cve/suppress"
	debugService "github.com/stackrox/rox/central/debug/service"
	"github.com/stackrox/rox/central/declarativeconfig"
	deploymentDatastore "github.com/stackrox/rox/central/deployment/datastore"
	deploymentService "github.com/stackrox/rox/central/deployment/service"
	detectionService "github.com/stackrox/rox/central/detection/service"
	developmentService "github.com/stackrox/rox/central/development/service"
	"github.com/stackrox/rox/central/docs"
	"github.com/stackrox/rox/central/endpoints"
	"github.com/stackrox/rox/central/enrichment"
	_ "github.com/stackrox/rox/central/externalbackups/plugins/all" // Import all of the external backup plugins
	backupService "github.com/stackrox/rox/central/externalbackups/service"
	featureFlagService "github.com/stackrox/rox/central/featureflags/service"
	"github.com/stackrox/rox/central/globaldb"
	dbAuthz "github.com/stackrox/rox/central/globaldb/authz"
	globaldbHandlers "github.com/stackrox/rox/central/globaldb/handlers"
	backupRestoreService "github.com/stackrox/rox/central/globaldb/v2backuprestore/service"
	graphqlHandler "github.com/stackrox/rox/central/graphql/handler"
	groupDataStore "github.com/stackrox/rox/central/group/datastore"
	groupService "github.com/stackrox/rox/central/group/service"
	"github.com/stackrox/rox/central/grpc/metrics"
	"github.com/stackrox/rox/central/helmcharts"
	imageDatastore "github.com/stackrox/rox/central/image/datastore"
	imageService "github.com/stackrox/rox/central/image/service"
	iiDatastore "github.com/stackrox/rox/central/imageintegration/datastore"
	iiService "github.com/stackrox/rox/central/imageintegration/service"
	iiStore "github.com/stackrox/rox/central/imageintegration/store"
	integrationHealthService "github.com/stackrox/rox/central/integrationhealth/service"
	"github.com/stackrox/rox/central/jwt"
	licenseService "github.com/stackrox/rox/central/license/service"
	logimbueHandler "github.com/stackrox/rox/central/logimbue/handler"
	logimbueStore "github.com/stackrox/rox/central/logimbue/store"
	metadataService "github.com/stackrox/rox/central/metadata/service"
	mitreService "github.com/stackrox/rox/central/mitre/service"
	namespaceService "github.com/stackrox/rox/central/namespace/service"
	networkBaselineDataStore "github.com/stackrox/rox/central/networkbaseline/datastore"
	networkBaselineService "github.com/stackrox/rox/central/networkbaseline/service"
	networkEntityDataStore "github.com/stackrox/rox/central/networkgraph/entity/datastore"
	"github.com/stackrox/rox/central/networkgraph/entity/gatherer"
	networkFlowService "github.com/stackrox/rox/central/networkgraph/service"
	networkPolicyService "github.com/stackrox/rox/central/networkpolicies/service"
	nodeService "github.com/stackrox/rox/central/node/service"
	notifierDS "github.com/stackrox/rox/central/notifier/datastore"
	"github.com/stackrox/rox/central/notifier/processor"
	notifierService "github.com/stackrox/rox/central/notifier/service"
	_ "github.com/stackrox/rox/central/notifiers/all" // These imports are required to register things from the respective packages.
	"github.com/stackrox/rox/central/option"
	pingService "github.com/stackrox/rox/central/ping/service"
	podService "github.com/stackrox/rox/central/pod/service"
	policyDataStore "github.com/stackrox/rox/central/policy/datastore"
	policyService "github.com/stackrox/rox/central/policy/service"
	policyCategoryService "github.com/stackrox/rox/central/policycategory/service"
	probeUploadService "github.com/stackrox/rox/central/probeupload/service"
	processBaselineDataStore "github.com/stackrox/rox/central/processbaseline/datastore"
	processBaselineService "github.com/stackrox/rox/central/processbaseline/service"
	processIndicatorService "github.com/stackrox/rox/central/processindicator/service"
	processListeningOnPorts "github.com/stackrox/rox/central/processlisteningonport/service"
	"github.com/stackrox/rox/central/pruning"
	rbacService "github.com/stackrox/rox/central/rbac/service"
	reportConfigurationService "github.com/stackrox/rox/central/reportconfigurations/service"
	vulnReportScheduleManager "github.com/stackrox/rox/central/reports/manager"
	reportService "github.com/stackrox/rox/central/reports/service"
	"github.com/stackrox/rox/central/reprocessor"
	collectionService "github.com/stackrox/rox/central/resourcecollection/service"
	"github.com/stackrox/rox/central/risk/handlers/timeline"
	"github.com/stackrox/rox/central/role"
	roleDataStore "github.com/stackrox/rox/central/role/datastore"
	"github.com/stackrox/rox/central/role/mapper"
	"github.com/stackrox/rox/central/role/resources"
	roleService "github.com/stackrox/rox/central/role/service"
	centralSAC "github.com/stackrox/rox/central/sac"
	"github.com/stackrox/rox/central/scanner"
	scannerDefinitionsHandler "github.com/stackrox/rox/central/scannerdefinitions/handler"
	searchService "github.com/stackrox/rox/central/search/service"
	secretService "github.com/stackrox/rox/central/secret/service"
	sensorService "github.com/stackrox/rox/central/sensor/service"
	"github.com/stackrox/rox/central/sensor/service/connection"
	"github.com/stackrox/rox/central/sensor/service/pipeline/all"
	sensorUpgradeControlService "github.com/stackrox/rox/central/sensorupgrade/controlservice"
	sensorUpgradeService "github.com/stackrox/rox/central/sensorupgrade/service"
	serviceAccountService "github.com/stackrox/rox/central/serviceaccount/service"
	siStore "github.com/stackrox/rox/central/serviceidentities/datastore"
	siService "github.com/stackrox/rox/central/serviceidentities/service"
	signatureIntegrationDS "github.com/stackrox/rox/central/signatureintegration/datastore"
	signatureIntegrationService "github.com/stackrox/rox/central/signatureintegration/service"
	"github.com/stackrox/rox/central/splunk"
	summaryService "github.com/stackrox/rox/central/summary/service"
	"github.com/stackrox/rox/central/systeminfo/listener"
	"github.com/stackrox/rox/central/telemetry/centralclient"
	"github.com/stackrox/rox/central/telemetry/gatherers"
	telemetryService "github.com/stackrox/rox/central/telemetry/service"
	"github.com/stackrox/rox/central/tlsconfig"
	"github.com/stackrox/rox/central/ui"
	userService "github.com/stackrox/rox/central/user/service"
	"github.com/stackrox/rox/central/version"
	vStore "github.com/stackrox/rox/central/version/store"
	versionUtils "github.com/stackrox/rox/central/version/utils"
	vulnRequestManager "github.com/stackrox/rox/central/vulnerabilityrequest/manager/requestmgr"
	vulnRequestService "github.com/stackrox/rox/central/vulnerabilityrequest/service"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/auth/authproviders"
	"github.com/stackrox/rox/pkg/auth/authproviders/iap"
	"github.com/stackrox/rox/pkg/auth/authproviders/oidc"
	"github.com/stackrox/rox/pkg/auth/authproviders/openshift"
	"github.com/stackrox/rox/pkg/auth/authproviders/saml"
	authProviderUserpki "github.com/stackrox/rox/pkg/auth/authproviders/userpki"
	"github.com/stackrox/rox/pkg/auth/permissions"
	"github.com/stackrox/rox/pkg/clientconn"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/config"
	"github.com/stackrox/rox/pkg/devbuild"
	"github.com/stackrox/rox/pkg/devmode"
	"github.com/stackrox/rox/pkg/env"
	pkgGRPC "github.com/stackrox/rox/pkg/grpc"
	"github.com/stackrox/rox/pkg/grpc/authn"
	"github.com/stackrox/rox/pkg/grpc/authn/service"
	"github.com/stackrox/rox/pkg/grpc/authn/servicecerttoken"
	"github.com/stackrox/rox/pkg/grpc/authn/tokenbased"
	authnUserpki "github.com/stackrox/rox/pkg/grpc/authn/userpki"
	"github.com/stackrox/rox/pkg/grpc/authz"
	"github.com/stackrox/rox/pkg/grpc/authz/allow"
	"github.com/stackrox/rox/pkg/grpc/authz/or"
	"github.com/stackrox/rox/pkg/grpc/authz/perrpc"
	"github.com/stackrox/rox/pkg/grpc/authz/user"
	"github.com/stackrox/rox/pkg/grpc/errors"
	"github.com/stackrox/rox/pkg/grpc/routes"
	"github.com/stackrox/rox/pkg/httputil"
	"github.com/stackrox/rox/pkg/httputil/proxy"
	"github.com/stackrox/rox/pkg/logging"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	"github.com/stackrox/rox/pkg/migrations"
	"github.com/stackrox/rox/pkg/osutils"
	"github.com/stackrox/rox/pkg/postgres/pgadmin"
	"github.com/stackrox/rox/pkg/postgres/pgconfig"
	"github.com/stackrox/rox/pkg/premain"
	"github.com/stackrox/rox/pkg/sac"
	"github.com/stackrox/rox/pkg/sac/observe"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
	pkgVersion "github.com/stackrox/rox/pkg/version"
)

var (
	log = logging.CreatePersistentLogger(logging.CurrentModule(), 0)

	authProviderBackendFactories = map[string]authproviders.BackendFactoryCreator{
		oidc.TypeName:                oidc.NewFactory,
		"auth0":                      oidc.NewFactory, // legacy
		saml.TypeName:                saml.NewFactory,
		authProviderUserpki.TypeName: authProviderUserpki.NewFactoryFactory(tlsconfig.ManagerInstance()),
		iap.TypeName:                 iap.NewFactory,
	}

	imageIntegrationContext = sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Integration),
		))
)

const (
	ssoURLPathPrefix = "/sso/"
	//#nosec G101 -- This is a false positive
	tokenRedirectURLPath = "/auth/response/generic"

	grpcServerWatchdogTimeout = 20 * time.Second

	maxServiceCertTokenLeeway = 1 * time.Minute

	proxyConfigPath = "/run/secrets/stackrox.io/proxy-config"
	proxyConfigFile = "config.yaml"
)

func init() {
	if !proxy.UseWithDefaultTransport() {
		log.Warn("Failed to use proxy transport with default HTTP transport. Some proxy features may not work.")
	}
}

func runSafeMode() {
	log.Info("Started Central up in safe mode. Sleeping forever...")

	signalsC := make(chan os.Signal, 1)
	signal.Notify(signalsC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-signalsC
	log.Infof("Caught %s signal", sig)
	log.Info("Central terminated")
}

func main() {
	premain.StartMain()

	conf := config.GetConfig()
	if conf == nil || conf.Maintenance.SafeMode {
		if conf == nil {
			log.Error("cannot get central configuration. Starting up in safe mode")
		}
		runSafeMode()
		return
	}

	clientconn.SetUserAgent(clientconn.Central)

	ctx := context.Background()
	proxy.WatchProxyConfig(ctx, proxyConfigPath, proxyConfigFile, true)

	devmode.StartOnDevBuilds("central")

	log.Infof("Running StackRox Version: %s", pkgVersion.GetMainVersion())
	log.Warn("The following permission resources have been replaced:\n" +
		"	Access replaces AuthProvider, Group, Licenses, and User\n" +
		"	Administration replaces AllComments, Config, DebugLogs, NetworkGraphConfig, ProbeUpload, ScannerBundle, ScannerDefinitions, SensorUpgradeConfig, and ServiceIdentity\n" +
		"	Cluster also covers ClusterCVE\n" +
		"	Compliance replaces ComplianceRuns\n" +
		"	DeploymentExtension replaces Indicator, NetworkBaseline, ProcessWhitelist, and Risk\n" +
		"	Integration replaces APIToken, BackupPlugins, ImageIntegration, Notifier, and SignatureIntegration\n" +
		"	Image now also covers ImageComponent\n" +
		"The following permission resources will be replaced in the upcoming versions:\n" +
		"	Access will replace Role\n" +
		"	WorkflowAdministration will replace Policy and VulnerabilityReports.")
	ensureDB(ctx)

	// Need to remove the backup clone and set the current version
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		sourceMap, config, err := pgconfig.GetPostgresConfig()
		if err != nil {
			log.Errorf("Unable to get Postgres DB config: %v", err)
		}

		err = pgadmin.DropDB(sourceMap, config, migrations.GetBackupClone())
		if err != nil {
			log.Errorf("Failed to remove backup DB: %v", err)
		}
		versionUtils.SetCurrentVersionPostgres(globaldb.GetPostgres())
	} else {
		// Update last associated software version on DBs.
		migrations.SetCurrent(option.CentralOptions.DBPathBase)
	}

	// Now that we verified that the DB can be loaded, remove the .backup directory
	if err := migrations.SafeRemoveDBWithSymbolicLink(filepath.Join(migrations.DBMountPath(), migrations.GetBackupClone())); err != nil {
		log.Fatalf("Failed to remove backup DB: %v", err)
	}

	// Start the prometheus metrics server
	pkgMetrics.NewDefaultHTTPServer(pkgMetrics.CentralSubsystem).RunForever()
	pkgMetrics.GatherThrottleMetricsForever(pkgMetrics.CentralSubsystem.String())

	go startGRPCServer()

	waitForTerminationSignal()
}

func ensureDB(ctx context.Context) {
	var versionStore vStore.Store
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		versionStore = vStore.NewPostgres(globaldb.InitializePostgres(ctx))
	} else {
		versionStore = vStore.New(globaldb.GetGlobalDB(), globaldb.GetRocksDB())
	}

	err := version.Ensure(versionStore)
	if err != nil {
		log.Panicf("DB version check failed. You may need to run migrations: %v", err)
	}
}

func startServices() {
	reprocessor.Singleton().Start()
	suppress.Singleton().Start()
	pruning.Singleton().Start()
	gatherer.Singleton().Start()
	vulnRequestManager.Singleton().Start()

	if env.PostgresDatastoreEnabled.BooleanSetting() {
		apiTokenExpiration.Singleton().Start()
	}

	go registerDelayedIntegrations(iiStore.DelayedIntegrations)
}

func servicesToRegister(registry authproviders.Registry, authzTraceSink observe.AuthzTraceSink) []pkgGRPC.APIService {
	// PLEASE KEEP THE FOLLOWING LIST SORTED.
	servicesToRegister := []pkgGRPC.APIService{
		alertService.Singleton(),
		apiTokenService.Singleton(),
		authService.New(),
		authProviderSvc.New(registry, groupDataStore.Singleton()),
		backupRestoreService.Singleton(),
		backupService.Singleton(),
		centralHealthService.Singleton(),
		certgen.ServiceSingleton(),
		clusterInitService.Singleton(),
		clusterService.Singleton(),
		complianceManagerService.Singleton(),
		complianceService.Singleton(),
		configService.Singleton(),
		credentialExpiryService.Singleton(),
		debugService.New(
			clusterDataStore.Singleton(),
			connection.ManagerSingleton(),
			gatherers.Singleton(),
			logimbueStore.Singleton(),
			authzTraceSink,
			registry,
			groupDataStore.Singleton(),
			roleDataStore.Singleton(),
			configDS.Singleton(),
			notifierDS.Singleton(),
		),
		deploymentService.Singleton(),
		detectionService.Singleton(),
		featureFlagService.Singleton(),
		groupService.Singleton(),
		helmcharts.NewService(),
		imageService.Singleton(),
		iiService.Singleton(),
		licenseService.New(),
		integrationHealthService.Singleton(),
		metadataService.New(),
		mitreService.Singleton(),
		namespaceService.Singleton(),
		networkBaselineService.Singleton(),
		networkFlowService.Singleton(),
		networkPolicyService.Singleton(),
		nodeService.Singleton(),
		notifierService.Singleton(),
		pingService.Singleton(),
		podService.Singleton(),
		policyService.Singleton(),
		probeUploadService.Singleton(),
		processIndicatorService.Singleton(),
		processBaselineService.Singleton(),
		rbacService.Singleton(),
		reportConfigurationService.Singleton(),
		reportService.Singleton(),
		roleService.Singleton(),
		searchService.Singleton(),
		secretService.Singleton(),
		sensorService.New(connection.ManagerSingleton(), all.Singleton(), clusterDataStore.Singleton()),
		sensorUpgradeControlService.Singleton(),
		sensorUpgradeService.Singleton(),
		serviceAccountService.Singleton(),
		signatureIntegrationService.Singleton(),
		siService.Singleton(),
		summaryService.Singleton(),
		telemetryService.Singleton(),
		userService.Singleton(),
		vulnRequestService.Singleton(),
	}
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		servicesToRegister = append(servicesToRegister, clusterCVEService.Singleton())
		servicesToRegister = append(servicesToRegister, imageCVEService.Singleton())
		servicesToRegister = append(servicesToRegister, nodeCVEService.Singleton())
		servicesToRegister = append(servicesToRegister, collectionService.Singleton())
		servicesToRegister = append(servicesToRegister, policyCategoryService.Singleton())
		servicesToRegister = append(servicesToRegister, processListeningOnPorts.Singleton())
	} else {
		servicesToRegister = append(servicesToRegister, cveService.Singleton())
	}

	autoTriggerUpgrades := sensorUpgradeService.Singleton().AutoUpgradeSetting()
	if err := connection.ManagerSingleton().Start(
		clusterDataStore.Singleton(),
		networkEntityDataStore.Singleton(),
		policyDataStore.Singleton(),
		processBaselineDataStore.Singleton(),
		networkBaselineDataStore.Singleton(),
		autoTriggerUpgrades,
	); err != nil {
		log.Panicf("Couldn't start sensor connection manager: %v", err)
	}

	if !env.OfflineModeEnv.BooleanSetting() {
		go fetcher.SingletonManager().Start()
	}

	if devbuild.IsEnabled() {
		servicesToRegister = append(servicesToRegister, developmentService.Singleton())
	}

	return servicesToRegister
}

func watchdog(signal *concurrency.Signal, timeout time.Duration) {
	if !concurrency.WaitWithTimeout(signal, timeout) {
		log.Errorf("API server failed to start within %v!", timeout)
		log.Error("This usually means something is *very* wrong. Terminating ...")
		if err := syscall.Kill(syscall.Getpid(), syscall.SIGABRT); err != nil {
			panic(err)
		}
	}
}

func startGRPCServer() {
	// Temporarily elevate permissions to modify auth providers.
	authProviderRegisteringCtx := sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS, storage.Access_READ_WRITE_ACCESS),
			sac.ResourceScopeKeys(resources.Access)))

	// Create the registry of applied auth providers.
	registry, err := authproviders.NewStoreBackedRegistry(
		ssoURLPathPrefix, tokenRedirectURLPath,
		authProviderDS.Singleton(), jwt.IssuerFactorySingleton(),
		mapper.FactorySingleton())
	if err != nil {
		log.Panicf("Could not create auth provider registry: %v", err)
	}

	// env.EnableOpenShiftAuth signals the desire but does not guarantee Central
	// is configured correctly to talk to the OpenShift's OAuth server. If this
	// is the case, we can be setting up an auth providers which won't work.
	if env.EnableOpenShiftAuth.BooleanSetting() {
		authProviderBackendFactories[openshift.TypeName] = openshift.NewFactory
	}

	for typeName, factoryCreator := range authProviderBackendFactories {
		if err := registry.RegisterBackendFactory(authProviderRegisteringCtx, typeName, factoryCreator); err != nil {
			log.Panicf("Could not register %s auth provider factory: %v", typeName, err)
		}
	}
	if err := registry.Init(); err != nil {
		log.Panicf("Could not initialize auth provider registry: %v", err)
	}

	basicAuthMgr, err := userpass.CreateManager(roleDataStore.Singleton())
	if err != nil {
		log.Panicf("Could not create basic auth manager: %v", err)
	}

	basicAuthProvider := userpass.RegisterAuthProviderOrPanic(authProviderRegisteringCtx, basicAuthMgr, registry)

	if env.DeclarativeConfiguration.BooleanSetting() {
		declarativeconfig.ManagerSingleton(registry).ReconcileDeclarativeConfigurations()
	}

	clusterInitBackend := backend.Singleton()
	serviceMTLSExtractor, err := service.NewExtractorWithCertValidation(clusterInitBackend)
	if err != nil {
		log.Panicf("Could not create mTLS-based service identity extractor: %v", err)
	}

	serviceTokenExtractor, err := servicecerttoken.NewExtractorWithCertValidation(maxServiceCertTokenLeeway, clusterInitBackend)
	if err != nil {
		log.Panicf("Could not create ServiceCert token-based identity extractor: %v", err)
	}

	idExtractors := []authn.IdentityExtractor{
		serviceMTLSExtractor, // internal services
		tokenbased.NewExtractor(roleDataStore.Singleton(), jwt.ValidatorSingleton()), // JWT tokens
		userpass.IdentityExtractorOrPanic(roleDataStore.Singleton(), basicAuthMgr, basicAuthProvider),
		serviceTokenExtractor,
		authnUserpki.NewExtractor(tlsconfig.ManagerInstance()),
	}

	endpointCfgs, err := endpoints.InstantiateAll(tlsconfig.ManagerInstance())
	if err != nil {
		log.Panicf("Could not instantiate endpoint configs: %v", err)
	}

	config := pkgGRPC.Config{
		CustomRoutes:       customRoutes(),
		IdentityExtractors: idExtractors,
		AuthProviders:      registry,
		Auditor:            audit.New(processor.Singleton()),
		GRPCMetrics:        metrics.GRPCSingleton(),
		HTTPMetrics:        metrics.HTTPSingleton(),
		Endpoints:          endpointCfgs,
	}

	if devbuild.IsEnabled() {
		config.UnaryInterceptors = append(config.UnaryInterceptors,
			errors.LogInternalErrorInterceptor,
			errors.PanicOnInvariantViolationUnaryInterceptor,
		)
		config.StreamInterceptors = append(config.StreamInterceptors,
			errors.LogInternalErrorStreamInterceptor,
			errors.PanicOnInvariantViolationStreamInterceptor,
		)
	}

	// This adds an on-demand global tracing for the built-in authorization.
	authzTraceSink := observe.NewAuthzTraceSink()
	config.UnaryInterceptors = append(config.UnaryInterceptors,
		observe.AuthzTraceInterceptor(authzTraceSink),
	)
	config.HTTPInterceptors = append(config.HTTPInterceptors, observe.AuthzTraceHTTPInterceptor(authzTraceSink))

	// Before authorization is checked, we want to inject the sac client into the context.
	config.PreAuthContextEnrichers = append(config.PreAuthContextEnrichers,
		centralSAC.GetEnricher().GetPreAuthContextEnricher(authzTraceSink),
	)

	telemetryCtx := sac.WithGlobalAccessScopeChecker(context.Background(),
		sac.AllowFixedScopes(
			sac.AccessModeScopeKeys(storage.Access_READ_ACCESS),
			sac.ResourceScopeKeys(resources.Administration)))

	if cds, err := configDS.Singleton().GetConfig(telemetryCtx); err == nil || cds == nil {
		if t := cds.GetPublicConfig().GetTelemetry(); t == nil || t.GetEnabled() {
			if cfg := centralclient.Enable(); cfg.Enabled() {
				centralclient.RegisterCentralClient(&config, basicAuthProvider.ID())
				gs := cfg.Gatherer()
				gs.AddGatherer(authProviderTelemetry.Gather)
				gs.AddGatherer(signatureIntegrationDS.Gather)
				gs.AddGatherer(roleDataStore.Gather)
				gs.AddGatherer(clusterDataStore.Gather)
			}
		}
	}

	server := pkgGRPC.NewAPI(config)
	server.Register(servicesToRegister(registry, authzTraceSink)...)

	startServices()
	startedSig := server.Start()

	go watchdog(startedSig, grpcServerWatchdogTimeout)
}

func registerDelayedIntegrations(integrationsInput []iiStore.DelayedIntegration) {
	integrationManager := enrichment.ManagerSingleton()

	integrations := make(map[int]iiStore.DelayedIntegration, len(integrationsInput))
	for k, v := range integrationsInput {
		integrations[k] = v
	}
	ds := iiDatastore.Singleton()
	for len(integrations) > 0 {
		for idx, integration := range integrations {
			_, exists, _ := ds.GetImageIntegration(imageIntegrationContext, integration.Integration.GetId())
			if exists {
				delete(integrations, idx)
				continue
			}
			ready := integration.Trigger()
			if !ready {
				continue
			}
			// add the integration first, which is more likely to fail. If it does, no big deal -- you can still try to
			// manually add it and get the error message.
			err := integrationManager.Upsert(integration.Integration)
			if err == nil {
				err = ds.UpdateImageIntegration(imageIntegrationContext, integration.Integration)
				if err != nil {
					// so, we added the integration to the set but we weren't able to save it.
					// This is ok -- the image scanner will "work" and after a restart we'll try to save it again.
					log.Errorf("We added the %q integration, but saving it failed with: %v. We'll try again next restart", integration.Integration.GetName(), err)
				} else {
					log.Infof("Registered integration %q", integration.Integration.GetName())
				}
				reprocessor.Singleton().ShortCircuit()
			} else {
				log.Errorf("Unable to register integration %q: %v", integration.Integration.GetName(), err)
			}
			// either way, time to stop watching this entry
			delete(integrations, idx)
		}
		time.Sleep(5 * time.Second)
	}
	log.Debug("All dynamic integrations registered, exiting")
}

func uiRoute() routes.CustomRoute {
	return routes.CustomRoute{
		Route:         "/",
		Authorizer:    allow.Anonymous(),
		ServerHandler: ui.Mux(),
		Compression:   true,
	}
}

func customRoutes() (customRoutes []routes.CustomRoute) {
	customRoutes = []routes.CustomRoute{
		uiRoute(),
		{
			Route:         "/api/extensions/clusters/zip",
			Authorizer:    or.SensorOrAuthorizer(user.With(permissions.View(resources.Cluster), permissions.View(resources.Administration))),
			ServerHandler: clustersZip.Handler(clusterDataStore.Singleton(), siStore.Singleton()),
			Compression:   false,
		},
		{
			Route:         "/api/extensions/scanner/zip",
			Authorizer:    user.With(permissions.View(resources.Administration)),
			ServerHandler: scanner.Handler(),
			Compression:   false,
		},
		{
			Route:         "/api/cli/download/",
			Authorizer:    user.With(),
			ServerHandler: cli.Handler(),
			Compression:   true,
		},
		{
			Route:         "/api/docs/swagger",
			Authorizer:    user.With(permissions.View(resources.Integration)),
			ServerHandler: docs.Swagger(),
			Compression:   true,
		},
		{
			Route:         "/api/graphql",
			Authorizer:    user.With(), // graphql enforces permissions internally
			ServerHandler: graphqlHandler.Handler(),
			Compression:   true,
		},
		{
			Route:         "/api/compliance/export/csv",
			Authorizer:    user.With(permissions.View(resources.Compliance)),
			ServerHandler: complianceHandlers.CSVHandler(),
			Compression:   true,
		},
		{
			Route:         "/api/risk/timeline/export/csv",
			Authorizer:    user.With(permissions.View(resources.Deployment), permissions.View(resources.DeploymentExtension)),
			ServerHandler: timeline.CSVHandler(),
			Compression:   true,
		},
		{
			Route:         "/api/vm/export/csv",
			Authorizer:    user.With(permissions.View(resources.Image), permissions.View(resources.Deployment), permissions.View(resources.Node)),
			ServerHandler: csv.CVECSVHandler(),
			Compression:   true,
		},
		{
			Route:         "/api/splunk/ta/vulnmgmt",
			Authorizer:    user.With(permissions.View(resources.Image), permissions.View(resources.Deployment)),
			ServerHandler: splunk.NewVulnMgmtHandler(deploymentDatastore.Singleton(), imageDatastore.Singleton()),
			Compression:   true,
		},
		{
			Route:         "/api/splunk/ta/compliance",
			Authorizer:    user.With(permissions.View(resources.Compliance)),
			ServerHandler: splunk.NewComplianceHandler(complianceDatastore.Singleton()),
			Compression:   true,
		},
		{
			Route:         "/api/splunk/ta/violations",
			Authorizer:    user.With(permissions.View(resources.Alert)),
			ServerHandler: splunk.NewViolationsHandler(alertDatastore.Singleton()),
			Compression:   true,
		},
		{
			Route:         "/db/v2/restore",
			Authorizer:    dbAuthz.DBWriteAccessAuthorizer(),
			ServerHandler: backupRestoreService.Singleton().RestoreHandler(),
		},
		{
			Route:         "/db/v2/resumerestore",
			Authorizer:    dbAuthz.DBWriteAccessAuthorizer(),
			ServerHandler: backupRestoreService.Singleton().ResumeRestoreHandler(),
		},
		{
			Route:         "/api/logimbue",
			Authorizer:    user.With(),
			ServerHandler: logimbueHandler.Singleton(),
			Compression:   false,
		},
	}

	if env.PostgresDatastoreEnabled.BooleanSetting() {
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:      "/db/backup",
			Authorizer: dbAuthz.DBReadAccessAuthorizer(),
			ServerHandler: notImplementedOnManagedServices(
				globaldbHandlers.BackupDB(nil, nil, globaldb.GetPostgres(), listener.Singleton(), false),
			),
			Compression: true,
		})
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:      "/api/extensions/backup",
			Authorizer: user.WithRole(role.Admin),
			ServerHandler: notImplementedOnManagedServices(
				globaldbHandlers.BackupDB(nil, nil, globaldb.GetPostgres(), listener.Singleton(), true),
			),
			Compression: true,
		})
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         "/api/export/csv/node/cve",
			Authorizer:    user.With(permissions.View(resources.Node)),
			ServerHandler: nodeCveCsv.NodeCVECSVHandler(),
			Compression:   true,
		})
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         "/api/export/csv/image/cve",
			Authorizer:    user.With(permissions.View(resources.Image), permissions.View(resources.Deployment)),
			ServerHandler: imageCveCsv.ImageCVECSVHandler(),
			Compression:   true,
		})
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         "/api/export/csv/cluster/cve",
			Authorizer:    user.With(permissions.View(resources.Cluster)),
			ServerHandler: clusterCveCsv.ClusterCVECSVHandler(),
			Compression:   true,
		})
	} else {
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         "/db/backup",
			Authorizer:    dbAuthz.DBReadAccessAuthorizer(),
			ServerHandler: notImplementedOnManagedServices(globaldbHandlers.BackupDB(globaldb.GetGlobalDB(), globaldb.GetRocksDB(), nil, listener.Singleton(), false)),
			Compression:   true,
		})
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         "/api/extensions/backup",
			Authorizer:    user.WithRole(role.Admin),
			ServerHandler: notImplementedOnManagedServices(globaldbHandlers.BackupDB(globaldb.GetGlobalDB(), globaldb.GetRocksDB(), nil, listener.Singleton(), true)),
			Compression:   true,
		})

		// v1 style restore endpoint, not supported for Postgres
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         "/db/restore",
			Authorizer:    dbAuthz.DBWriteAccessAuthorizer(),
			ServerHandler: globaldbHandlers.RestoreDB(globaldb.GetGlobalDB(), globaldb.GetRocksDB()),
		})
	}

	customRoutes = append(customRoutes, routes.CustomRoute{
		Route:         "/api/extensions/clusters/helm-config.yaml",
		Authorizer:    or.SensorOrAuthorizer(user.With(permissions.View(resources.Cluster))),
		ServerHandler: clustersHelmConfig.Handler(clusterDataStore.Singleton()),
		Compression:   true,
	})

	scannerDefinitionsRoute := "/api/extensions/scannerdefinitions"
	// Only grant compression to well-known content types. It should capture files
	// worthy of compression in definition's bundle. Ignore all other types (e.g.,
	// `.zip` for the bundle itself).
	definitionsFileGzipHandler, err := gziphandler.GzipHandlerWithOpts(gziphandler.ContentTypes([]string{
		"application/json",
		"application/yaml",
		"text/plain",
	}))
	utils.CrashOnError(err)
	customRoutes = append(customRoutes,
		routes.CustomRoute{
			Route: scannerDefinitionsRoute,
			Authorizer: perrpc.FromMap(map[authz.Authorizer][]string{
				or.SensorOrAuthorizer(
					or.ScannerOr(
						user.With(permissions.View(resources.Administration)))): {
					routes.RPCNameForHTTP(scannerDefinitionsRoute, http.MethodGet),
				},
				user.With(permissions.Modify(resources.Administration)): {
					routes.RPCNameForHTTP(scannerDefinitionsRoute, http.MethodPost),
				},
			}),
			ServerHandler: definitionsFileGzipHandler(scannerDefinitionsHandler.Singleton()),
		},
	)

	customRoutes = append(customRoutes, debugRoutes()...)
	return
}

func notImplementedOnManagedServices(fn http.Handler) http.Handler {
	return utils.IfThenElse[http.Handler](
		env.ManagedCentral.BooleanSetting(), httputil.NotImplementedHandler("api is not supported in a managed central environment."),
		fn)
}

func debugRoutes() []routes.CustomRoute {
	customRoutes := make([]routes.CustomRoute, 0, len(routes.DebugRoutes))

	for r, h := range routes.DebugRoutes {
		customRoutes = append(customRoutes, routes.CustomRoute{
			Route:         r,
			Authorizer:    user.WithRole(role.Admin),
			ServerHandler: h,
			Compression:   true,
		})
	}
	return customRoutes
}

type stoppable interface {
	Stop()
}

type stoppableWithName struct {
	obj  stoppable
	name string
}

func waitForTerminationSignal() {
	signalsC := make(chan os.Signal, 1)
	signal.Notify(signalsC, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	sig := <-signalsC
	log.Infof("Caught %s signal", sig)

	stoppables := []stoppableWithName{
		{reprocessor.Singleton(), "reprocessor loop"},
		{suppress.Singleton(), "cve unsuppress loop"},
		{pruning.Singleton(), "gargage collector"},
		{gatherer.Singleton(), "network graph default external sources gatherer"},
		{vulnReportScheduleManager.Singleton(), "vuln reports schedule manager"},
		{vulnRequestManager.Singleton(), "vuln deferral requests expiry loop"},
		{centralclient.InstanceConfig().Gatherer(), "telemetry gatherer"},
		{centralclient.InstanceConfig().Telemeter(), "telemetry client"},
	}
	if env.PostgresDatastoreEnabled.BooleanSetting() {
		stoppables = append(stoppables, stoppableWithName{obj: apiTokenExpiration.Singleton(), name: "api token expiration notifier"})
	}

	var wg sync.WaitGroup
	for _, stoppable := range stoppables {
		wg.Add(1)
		go func(s stoppableWithName) {
			defer wg.Done()
			s.obj.Stop()
			log.Infof("Stopped %s", s.name)
		}(stoppable)
	}
	wg.Wait()

	globaldb.Close()

	if sig == syscall.SIGHUP {
		log.Info("Restarting central")
		osutils.Restart()
	}
	log.Info("Central terminated")
}
