package moban

import (
	webappv1 "djl.com/DjlD1/api/v1"
	"errors"
	appv1 "k8s.io/api/apps/v1"
	v2 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Newmoban(lei string, app webappv1.DjlD1) (client.Object, error) {
	if lei == "dep" {
		depsletr := &metav1.LabelSelector{
			MatchLabels: app.Labels,
		}
		pods := []v2.Container{
			{
				Name:            app.Name,
				Image:           app.Spec.Image,
				ImagePullPolicy: app.Spec.ImagePullPolicy,
			},
		}
		return &appv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,
				Labels:    app.Labels,
			},
			Spec: appv1.DeploymentSpec{
				Replicas: app.Spec.Size,
				Selector: depsletr,
				Template: v2.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: app.Labels,
					},
					Spec: v2.PodSpec{
						Containers: pods,
					},
				},
			},
		}, nil
	} else if lei == "svc" {
		return &v2.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,
				Labels:    app.Labels,
			},
			Spec: v2.ServiceSpec{
				Selector: app.Labels,
				Ports: []v2.ServicePort{
					{
						Port:       app.Spec.Ports[0].Port,
						Protocol:   app.Spec.Ports[0].Protocol,
						TargetPort: app.Spec.Ports[0].TargetPort,
					},
				},
				Type: v2.ServiceTypeClusterIP,
			},
		}, nil
	} else if lei == "ing" {
		return &networkingv1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      app.Name,
				Namespace: app.Namespace,
				Labels:    app.Labels,
			},
			Spec: networkingv1.IngressSpec{
				Rules: []networkingv1.IngressRule{
					{
						Host: app.Spec.ServerName,
						IngressRuleValue: networkingv1.IngressRuleValue{
							HTTP: &networkingv1.HTTPIngressRuleValue{
								Paths: []networkingv1.HTTPIngressPath{
									{
										Path: "/",
										PathType: func() *networkingv1.PathType {
											pt := networkingv1.PathTypePrefix
											return &pt
										}(),
										Backend: networkingv1.IngressBackend{
											Service: &networkingv1.IngressServiceBackend{
												Name: app.Name,
												Port: networkingv1.ServiceBackendPort{
													Number: app.Spec.Ports[0].Port,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}, nil
	} else {
		return nil, errors.New("参数不对")
	}
}
