/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"djl.com/DjlD1/moban"
	v1 "k8s.io/api/apps/v1"
	v2 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	webappv1 "djl.com/DjlD1/api/v1"
)

// DjlD1Reconciler reconciles a DjlD1 object
type DjlD1Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// 获取ready pod数量  自定义状态
func (r *DjlD1Reconciler) updateDjlD1Status(ctx context.Context, app *webappv1.DjlD1) error {
	dep := &v1.Deployment{}
	err := r.Client.Get(ctx, client.ObjectKey{
		Namespace: app.Namespace,
		Name:      app.Name,
	}, dep)
	if err != nil {
		return err
	}
	ready := dep.Status.ReadyReplicas
	zong := dep.Status.Replicas
	notready := zong - ready
	app.Status.Ready = ready
	app.Status.Notready = notready
	err = r.Status().Update(ctx, app)
	if err != nil {
		return err
	}
	return nil
}

// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch;update
// +kubebuilder:rbac:groups=webapp.djl.com,resources=djld1s,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=webapp.djl.com,resources=djld1s/status,verbs=get;patch;update
// +kubebuilder:rbac:groups=webapp.djl.com,resources=djld1s/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DjlD1 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.21.0/pkg/reconcile
func (r *DjlD1Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)
	app := webappv1.DjlD1{}
	// TODO(user): your logic here
	err := r.Client.Get(ctx, req.NamespacedName, &app)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Infof("DjlD1 %s not found", req.NamespacedName)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	if app.DeletionTimestamp != nil {
		klog.Infof("Guestbook resource has been deleted:%v-%v", app.Namespace, app.Name)
	}
	olddep := v1.Deployment{}
	//err1 := r.Client.Get(ctx, req.NamespacedName, &olddep)
	err = r.Client.Get(ctx, req.NamespacedName, &olddep)
	obj, _ := moban.Newmoban("dep", app)
	dep, _ := obj.(*v1.Deployment)
	//err = ctrl.SetControllerReference(&app, dep, r.Scheme)
	//if err != nil {
	//	return ctrl.Result{}, err
	//}
	_, err = ctrl.CreateOrUpdate(ctx, r.Client, dep, func() error {
		return ctrl.SetControllerReference(&app, dep, r.Scheme)
	})
	oldsvc := v2.Service{}
	err = r.Client.Get(ctx, req.NamespacedName, &oldsvc)
	obj, _ = moban.Newmoban("svc", app)
	svc, _ := obj.(*v2.Service)
	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, svc, func() error {
		return ctrl.SetControllerReference(&app, svc, r.Scheme)
	})
	olding := networkingv1.Ingress{}
	err = r.Client.Get(ctx, req.NamespacedName, &olding)
	obj, _ = moban.Newmoban("ing", app)
	ing, _ := obj.(*networkingv1.Ingress)
	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, ing, func() error {
		return ctrl.SetControllerReference(&app, ing, r.Scheme)
	})
	//if err1 != nil {
	//	if errors.IsNotFound(err1) {
	//		klog.Infof("add:%v-%v", app.Namespace, app.Name)
	//		err1 = r.Client.Create(ctx, dep)
	//		if err1 != nil {
	//			klog.Infof("创建dep失败")
	//			return ctrl.Result{}, err1
	//		}
	//		return ctrl.Result{}, nil
	//	}
	//	return ctrl.Result{}, err1
	//}
	//if !equality.Semantic.DeepEqual(olddep.Spec, dep.Spec) {
	//	klog.Infof("update djlD1 %s to controller", req.NamespacedName)
	//	olddep.Spec = dep.Spec
	//	if err = r.Client.Update(ctx, &olddep); err != nil {
	//		klog.Infof("更新失败")
	//		return ctrl.Result{}, err
	//	}
	//	return ctrl.Result{}, nil
	//}

	//使用状态方法
	err = r.updateDjlD1Status(ctx, &app)
	if err != nil {
		klog.Errorf("更新 DjlD1 状态失败: %v", err)
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DjlD1Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webappv1.DjlD1{}).
		Named("djld1").Owns(&v1.Deployment{}).
		Complete(r)
}
