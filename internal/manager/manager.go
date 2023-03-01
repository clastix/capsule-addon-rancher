package manager

import (
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func New(scheme *runtime.Scheme) (ctrl.Manager, error) {
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: ":8081",
		Port:               9443,
	})
	if err != nil {
		return nil, errors.Wrap(err, "unable to create manager")
	}

	return mgr, nil
}
