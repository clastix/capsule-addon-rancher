// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/clastix/capsule-addon-rancher/cmd/constants"
	"github.com/clastix/capsule-addon-rancher/internal/controller"
	"github.com/clastix/capsule-addon-rancher/internal/manager"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Starts the MutatingWebhookConfiguration required to integrate Capsule Proxy with Rancher",
		RunE: func(cmd *cobra.Command, args []string) error {
			scheme := runtime.NewScheme()

			if err := clientcmdapi.AddToScheme(scheme); err != nil {
				return errors.Wrap(err, "unable to register clientcmdapi scheme")
			}

			mgr, err := manager.New(scheme)
			if err != nil {
				return errors.Wrap(err, "unable to setup manager")
			}

			ca, err := os.ReadFile(viper.GetString(constants.WebhookProxyCAPath))
			if err != nil {
				return errors.Wrap(err, "unable to read the CA file required by ConfigMapHandler")
			}

			if err = controller.NewConfigMapHandler(
				controller.WithConfigMapPrefix(viper.GetString(constants.WebhookConfigMapPrefix)),
				controller.WithConfigMapKey(viper.GetString(constants.WebhookConfigMapKey)),
				controller.WithProxyScheme(viper.GetString(constants.WebhookProxyServiceScheme)),
				controller.WithProxyURL(viper.GetString(constants.WebhookProxyServiceURL)),
				controller.WithProxyPort(viper.GetInt(constants.WebhookProxyServicePort)),
				controller.WithProxyCA(ca),
			).SetupWithManager(mgr); err != nil {
				return errors.Wrap(err, "unable to setup ConfigMapHandler")
			}

			if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
				return errors.Wrap(err, "unable to start the manager")
			}

			return nil
		},
	}

	cmd.Flags().String(constants.WebhookConfigMapPrefix, "impersonation-shell-admin-kubeconfig-", "Prefix of the Rancher shell kubeconfig ConfigMap that must be mangled")
	cmd.Flags().String(constants.WebhookConfigMapKey, "config", "ConfigMap key name where the admin kubeconfig is stored")

	cmd.Flags().String(constants.WebhookProxyServiceScheme, "https", "HTTP scheme on which Capsule Proxy is listening to")
	cmd.Flags().String(constants.WebhookProxyServiceURL, "capsule-proxy.capsule-system.svc", "Kubernetes Service URL on which Capsule Proxy is waiting for connections")
	cmd.Flags().Int(constants.WebhookProxyServicePort, 9001, "Port on which Capsule Proxy is listening on")
	cmd.Flags().String(constants.WebhookProxyCAPath, "/tmp/ca.crt", "File containing the Certificate Authority used by Capsule Proxy")

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}
