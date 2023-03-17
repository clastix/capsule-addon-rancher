// Copyright 2020-2021 Clastix Labs
// SPDX-License-Identifier: Apache-2.0

package webhook

import (
	"flag"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/clastix/capsule-addon-rancher/cmd/constants"
	"github.com/clastix/capsule-addon-rancher/internal/controller"
	"github.com/clastix/capsule-addon-rancher/internal/manager"
)

type options struct {
	zo *zap.Options
}

func New() *cobra.Command {
	zo := &zap.Options{
		EncoderConfigOptions: append([]zap.EncoderConfigOption{}, func(config *zapcore.EncoderConfig) {
			config.EncodeTime = zapcore.ISO8601TimeEncoder
		}),
	}

	o := &options{zo: zo}

	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "Starts the MutatingWebhookConfiguration required to integrate Capsule Proxy with Rancher",
		RunE:  o.Run,
	}

	// Add ConfigMap options.
	cmd.Flags().String(constants.WebhookConfigMapPrefix, "impersonation-shell-admin-kubeconfig-", "Prefix of the Rancher shell kubeconfig ConfigMap that must be mangled")
	cmd.Flags().String(constants.WebhookConfigMapKey, "config", "ConfigMap key name where the admin kubeconfig is stored")

	// Add Capsule Proxy options.
	cmd.Flags().String(constants.WebhookProxyServiceScheme, "https", "HTTP scheme on which Capsule Proxy is listening to")
	cmd.Flags().String(constants.WebhookProxyServiceURL, "capsule-proxy.capsule-system.svc", "Kubernetes Service URL on which Capsule Proxy is waiting for connections")
	cmd.Flags().Int(constants.WebhookProxyServicePort, 9001, "Port on which Capsule Proxy is listening on")
	cmd.Flags().String(constants.WebhookProxyCAPath, "/tmp/ca.crt", "File containing the Certificate Authority used by Capsule Proxy")

	// Add and bind Zap options flags.
	var fs flag.FlagSet
	o.zo.BindFlags(&fs)
	cmd.Flags().AddGoFlagSet(&fs)

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}

func (o *options) Run(cmd *cobra.Command, agrs []string) error {
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
		controller.WithZap(o.zo)).SetupWithManager(mgr); err != nil {
		return errors.Wrap(err, "unable to setup ConfigMapHandler")
	}

	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		return errors.Wrap(err, "unable to start the manager")
	}

	return nil
}
