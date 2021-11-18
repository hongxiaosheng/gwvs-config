package controllers

import (
	"cmit.com/crd/gwvs-config/api/v1alpha1"
	"context"
	"github.com/go-logr/logr"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log"
)

var (
//services                 []*networkingv1alpha3.Server
//httpRouteList            []*networkingv1alpha3.HTTPRoute
//HTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
//HTTPMatchRequestList     []*networkingv1alpha3.HTTPMatchRequest
)

func CreateGateway(namespace string, reqLogger logr.Logger, istioClient *versioned.Clientset) {

	reqLogger.Info("istio  ---  Gateway-------操作开始-----------")
	gwname := namespace + "-gateway"
	reqLogger.Info("istio ", " gatewayname:", gwname)
	service := &networkingv1alpha3.Server{
		Port: &networkingv1alpha3.Port{
			Number: 80,
			// MUST BE one of HTTP|HTTPS|GRPC|HTTP2|MONGO|TCP|TLS.
			Protocol: "HTTP",
			Name:     "http",
		},
		Hosts: []string{"*"},
	}
	var services []*networkingv1alpha3.Server
	services = append(services, service)
	gw := &v1alpha3.Gateway{
		ObjectMeta: v1.ObjectMeta{
			Name:      gwname,
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.Gateway{
			Servers: services,
			//Selector:             nil,
		},
	}
	reqLogger.Info("istio  create gateway:", gwname, gw.Spec.GetServers())
	sgw, _ := istioClient.NetworkingV1alpha3().Gateways(namespace).Get(context.TODO(), gwname, v1.GetOptions{})
	if sgw.Name == gwname {
		reqLogger.Info("istio   gateway exists:", gwname, sgw)

	} else {
		cgw, err := istioClient.NetworkingV1alpha3().Gateways(namespace).Create(context.TODO(), gw, v1.CreateOptions{})
		if err != nil {
			// Request object not found.
			log.Println("create vs  failed", err.Error())
		}
		// 打印VS
		reqLogger.Info("istio  created cgw:", gwname, cgw)
	}
	reqLogger.Info("istio  ---  Gateway-------操作结束-----------")
}

func CreateVirtualService(namespace string, hosts []string, route []v1alpha1.Route, reqLogger logr.Logger, istioClient *versioned.Clientset) {
	reqLogger.Info("istio  --- vs-------操作开始-----------")
	var httpRouteList []*networkingv1alpha3.HTTPRoute
	var HTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
	var HTTPMatchRequestList []*networkingv1alpha3.HTTPMatchRequest
	//var httpRouteSign 			 *networkingv1alpha3.HTTPRoute
	gwname := namespace + "-gateway"
	vsname := route[0].Service + "-" + namespace + "-vs"
	for i := 0; i < len(route); i++ {
		// 定义http路由
		HTTPRouteDestination := &networkingv1alpha3.HTTPRouteDestination{
			Destination: &networkingv1alpha3.Destination{
				Host: route[i].Service,
				Port: &networkingv1alpha3.PortSelector{
					Number: route[i].Port,
				},
				//Subset:               "v2",
			},
			// 定义权重
			Weight: 100,
		}
		reqLogger.Info("istio  ", "HTTPRouteDestination :", HTTPRouteDestination)

		HTTPRouteDestinationList = append(HTTPRouteDestinationList, HTTPRouteDestination)
		for j := 0; j < len(route[i].Uri); j++ {
			HTTPMatchRequest := &networkingv1alpha3.HTTPMatchRequest{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Prefix{
						Prefix: route[i].Uri[j],
					},
				},
			}
			reqLogger.Info("istio  ", "HTTPMatchRequest :", HTTPMatchRequest)
			HTTPMatchRequestList = append(HTTPMatchRequestList, HTTPMatchRequest)
		}

		httpRouteSign := &networkingv1alpha3.HTTPRoute{
			Route: HTTPRouteDestinationList,
			Match: HTTPMatchRequestList,
		}
		reqLogger.Info("istio ", "httpRouteSign :", httpRouteSign)
		httpRouteList = append(httpRouteList, httpRouteSign)
	}

	virtualService := &v1alpha3.VirtualService{

		ObjectMeta: v1.ObjectMeta{
			Name:      vsname, // 定义vs的名称
			Namespace: namespace,
		},
		Spec: networkingv1alpha3.VirtualService{
			Hosts:    hosts, // 定义可访问的hosts
			Gateways: []string{gwname},
			Http:     httpRouteList, // 为hosts 绑定路由
		},
	}

	//select vs
	reqLogger.Info("istio  ", "visturalservice :", virtualService)
	vs, _ := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Get(context.TODO(), vsname, v1.GetOptions{})
	if vs.Name == vsname {
		istioClient.NetworkingV1alpha3().VirtualServices(namespace).Delete(context.TODO(), vsname, v1.DeleteOptions{})
	}
	// 创建VS
	cvs, err := istioClient.NetworkingV1alpha3().VirtualServices(namespace).Create(context.TODO(), virtualService, v1.CreateOptions{})
	if err != nil {
		// Request object not found.
		log.Println("create vs  failed", err.Error())
	}
	// 打印VS
	reqLogger.Info("istio  created cvs:", vsname, cvs)
	reqLogger.Info("istio  --- vs-------操作结束-----------")
}
