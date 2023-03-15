// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package controller

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/config-map,mutating=true,sideEffects=None,admissionReviewVersions=v1,failurePolicy=ignore,groups="",resources=configmaps,verbs=create;update,versions=v1,name=configmap.rancher.addon.capsule.clastix.io

type ConfigMapHandler struct {
	ConfigMapPrefix string
	ConfigMapKey    string

	ProxyScheme string
	ProxyURL    string
	ProxyPort   int
	ProxyCA     []byte

	decoder *admission.Decoder
	encoder *json.Serializer
}

func (c *ConfigMapHandler) SetupWithManager(mgr manager.Manager) error {
	c.encoder = json.NewSerializerWithOptions(json.SimpleMetaFactory{}, mgr.GetScheme(), mgr.GetScheme(), json.SerializerOptions{})

	mgr.GetWebhookServer().Register("/configmap", &webhook.Admission{
		Handler:      c,
		RecoverPanic: true,
	})

	return nil
}

func (c *ConfigMapHandler) InjectDecoder(d *admission.Decoder) error {
	c.decoder = d

	return nil
}

func (c *ConfigMapHandler) Handle(ctx context.Context, request admission.Request) admission.Response {
	cm := &corev1.ConfigMap{}
	if err := c.decoder.Decode(request, cm); err != nil {
		return admission.Errored(http.StatusInternalServerError, errors.Wrap(err, "unable to decode to *corev1.ConfigMap"))
	}

	if !strings.HasPrefix(cm.GetGenerateName(), c.ConfigMapPrefix) {
		return admission.Allowed("")
	}

	if len(cm.Data) == 0 {
		return admission.Errored(http.StatusBadRequest, fmt.Errorf("missing data to mangle"))
	}

	rawKubeConfig, ok := cm.Data[c.ConfigMapKey]
	if !ok {
		return admission.Errored(http.StatusBadRequest, fmt.Errorf("missing ConfigMap key, expected %s", c.ConfigMapKey))
	}

	kubeconfig, err := clientcmd.Load([]byte(rawKubeConfig))
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, errors.Wrap(err, "unable to load Config"))
	}

	for key := range kubeconfig.Clusters {
		kubeconfig.Clusters[key].Server = fmt.Sprintf("%s://%s:%d", c.ProxyScheme, c.ProxyURL, c.ProxyPort)
		kubeconfig.Clusters[key].CertificateAuthority = ""
		kubeconfig.Clusters[key].CertificateAuthorityData = c.ProxyCA
	}

	mangled, err := clientcmd.Write(*kubeconfig)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, errors.Wrap(err, "unable to marshal *clientcmdapi.Config"))
	}

	cm.Data[c.ConfigMapKey] = string(mangled)

	buf := bytes.NewBuffer([]byte{})
	if err = c.encoder.Encode(cm, buf); err != nil {
		return admission.Errored(http.StatusInternalServerError, errors.Wrap(err, "unable to encode mangled *corev1.ConfigMap"))
	}

	return admission.PatchResponseFromRaw(request.Object.Raw, buf.Bytes())
}
