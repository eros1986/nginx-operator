package k8s

const (
	PodType                   = "po"
	ReplicationControllerType = "rc"
	ReplicaSetType            = "rs"
	DeploymentType            = "deploy"
	StatefulSetType           = "statefulset"
	JobType                   = "job"
	CronJobType               = "cronjob"
	DaemonSetType             = "ds"
	ConfigMapType             = "cm"
	SecretType                = "secret"
	ServiceType               = "svc"
	IngressType               = "ing"
	EndpointType              = "ep"
	NamespaceType             = "ns"
	NetworkPolicyType         = "netpol"
	PersistentVolumeType      = "pv"
	PersistentVolumeClaimType = "pvc"
	StorageClassType          = "storageclass"
	ServiceAccountType        = "sa"

	MySQLClusterType = "mysqlCluster"
)

const (
	DefaultUid             = "unknown"
	DefaultVendorUid       = "vendor"
	DefaultMySQLRootPasswd = "passWORD"
)

const DefaultNamespace = "default"
const DefaultServiceAccountName = "default"
