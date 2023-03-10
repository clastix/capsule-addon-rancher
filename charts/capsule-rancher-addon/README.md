# capsule-rancher-addon

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 1.16.0](https://img.shields.io/badge/AppVersion-1.16.0-informational?style=flat-square)

A Helm chart for Kubernetes

## Requirements

| Repository | Name | Version |
|------------|------|---------|
| https://charts.jetstack.io | cert-manager | 1.11.0 |

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| certManager.generateCertificates | bool | `true` | Specifies whether webhook certificates should be generated using cert-manager |
| fullnameOverride | string | `""` |  |
| image.pullPolicy | string | `"IfNotPresent"` | Set the image pull policy. |
| image.repository | string | `"clastix/capsule-rancher-addon"` | Set the image repository of the capsule. |
| image.tag | string | `""` | Overrides the image tag whose default is the chart appVersion. |
| imagePullSecrets | list | `[]` |  |
| manager.affinity | object | `{}` | The affinity to be applied to the manager pod. |
| manager.nodeSelector | object | `{}` | The nodeSelector to be applied to the manager pod. |
| manager.podAnnotations | object | `{}` | The annotations to be applied to the manager pod. |
| manager.podSecurityContext | object | `{}` | The security context to be applied to the manager pod. |
| manager.replicaCount | int | `1` | The number of replicas of the manager. |
| manager.resources | object | `{}` |  |
| manager.serviceAccount.annotations | object | `{}` |  |
| manager.serviceAccount.create | bool | `true` |  |
| manager.serviceAccount.name | string | `""` |  |
| manager.tolerations | list | `[]` | The tolerations to be applied to the manager pod. |
| manager.webhook.failurePolicy | string | `"Fail"` | The mutating webhook failurePolicy |
| manager.webhook.reinvocationPolicy | string | `"Never"` | ReinvocationPolicy of the mutating webhook |
| manager.webhook.securityContext | object | `{}` |  |
| manager.webhook.service.port | int | `443` | The webhook service port |
| manager.webhook.service.type | string | `"ClusterIP"` | The webhook service type |
| manager.webhook.sideEffects | string | `"None"` | SideEffects of the mutating webhook |
| manager.webhook.timeoutSeconds | int | `30` | Timeout in seconds for the mutating webhook |
| nameOverride | string | `""` |  |
| proxy.CAMountPath | string | `"/tmp/proxy"` | The path in the webhook container to the certificate file of the Certificate Authority that signed the proxy certificate. |
| proxy.CASecretKey | string | `"ca"` | The key in the secret that contains the certificate file of the Certificate Authority that signed the proxy certificate. |
| proxy.CASecretName | string | `"capsule-proxy"` | The name of the secret that contains the certificate file of the Certificate Authority that signed the proxy certificate. |
| proxy.ServicePort | int | `9001` | The port of the proxy Service. |
| proxy.ServiceScheme | string | `"https"` | The scheme of the proxy Service. |
| proxy.ServiceURL | string | `"capsule-proxy.capsule-system.svc"` | The URL of the proxy Service. |
| proxy.configMapKey | string | `"config"` | The key in the ConfigMap that contains the kubeconfig for the proxy. |
| proxy.configMapPrefix | string | `"impersonation-shell-admin-kubeconfig-"` | The name prefix of the ConfigMap that contains the kubeconfig for the proxy. |

