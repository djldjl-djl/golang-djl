package main

import (
	webappv1 "djl.com/DjlD1/api/v1"
	"djl.com/DjlD1/httpapi"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"
	"os"
	ctrl "sigs.k8s.io/controller-runtime"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(webappv1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}
func main() {
	metrics := metricsserver.Options{
		BindAddress: "0", //metrics不监听
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:  scheme,
		Metrics: metrics,
	})
	if err != nil {
		klog.Error(err, "unable to start manager")
		os.Exit(1)
	}
	app := httpapi.Apphttp{
		Client: mgr.GetClient(),
		Router: gin.Default(),
	}
	err = mgr.Add(&app)
	if err != nil {
		klog.Error(err, "unable to add controller")
		os.Exit(1)
	}
	klog.Info("starting manager")
	if err = mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		klog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
