package k8s

import (
	"fmt"
	"strings"
)

func makeNameInCluster(rsType, svrType, uid, name string) string {
	if uid == "" {
		uid = DefaultUid
	}

	return strings.ToLower(fmt.Sprintf("%s-%s-%s-%s", rsType, svrType, uid, name))
}

func MakePodName(svrType, uid, name string) string {
	return makeNameInCluster(PodType, svrType, uid, name)
}

func MakeReplicationControllerName(svrType, uid, name string) string {
	return makeNameInCluster(ReplicationControllerType, svrType, uid, name)
}

func MakeReplicaSetName(svrType, uid, name string) string {
	return makeNameInCluster(ReplicaSetType, svrType, uid, name)
}

func MakeDeploymentName(svrType, uid, name string) string {
	return makeNameInCluster(DeploymentType, svrType, uid, name)
}

func MakeStatefulSetName(svrType, uid, name string) string {
	return makeNameInCluster(StatefulSetType, svrType, uid, name)
}

func MakeJobName(svrType, uid, name string) string {
	return makeNameInCluster(JobType, svrType, uid, name)
}

func MakeCronJobName(svrType, uid, name string) string {
	return makeNameInCluster(CronJobType, svrType, uid, name)
}

func MakeDaemonSetName(svrType, uid, name string) string {
	return makeNameInCluster(DaemonSetType, svrType, uid, name)
}

func MakeConfigMapName(svrType, uid, name string) string {
	return makeNameInCluster(ConfigMapType, svrType, uid, name)
}

func MakeSecretName(svrType, uid, name string) string {
	return makeNameInCluster(SecretType, svrType, uid, name)
}

func MakeServiceName(svrType, uid, name string) string {
	return makeNameInCluster(ServiceType, svrType, uid, name)
}

func MakeIngressName(svrType, uid, name string) string {
	return makeNameInCluster(IngressType, svrType, uid, name)
}

func MakeEndpointName(svrType, uid, name string) string {
	return makeNameInCluster(EndpointType, svrType, uid, name)
}

func MakeNamespaceName(svrType, uid, name string) string {
	return makeNameInCluster(NamespaceType, svrType, uid, name)
}

func MakeNetworkPolicyName(svrType, uid, name string) string {
	return makeNameInCluster(NetworkPolicyType, svrType, uid, name)
}

func MakePersistentVolumeName(svrType, uid, name string) string {
	return makeNameInCluster(PersistentVolumeType, svrType, uid, name)
}

func MakePersistentVolumeClaimName(svrType, uid, name string) string {
	return makeNameInCluster(PersistentVolumeClaimType, svrType, uid, name)
}

func MakeStorageClassName(svrType, uid, name string) string {
	return makeNameInCluster(StorageClassType, svrType, uid, name)
}

func MakeServiceAccountName(svrType, uid, name string) string {
	return makeNameInCluster(ServiceAccountType, svrType, uid, name)
}

func MakeClusterName(svrType, uid, name string) string {
	return makeNameInCluster(MySQLClusterType, svrType, uid, name)
}
