// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package k8sclient

import (
	"log"
	"testing"
	"time"

	"gotest.tools/v3/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	v1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
)

var endpointsArray = []interface{}{
	&discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "guestbook",
			GenerateName:    "",
			Namespace:       "default",
			SelfLink:        "/api/v1/namespaces/default/endpoints/guestbook",
			UID:             "a885b78c-5573-11e9-b47e-066a7a20bac8",
			ResourceVersion: "1550348",
			Generation:      0,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			Labels: map[string]string{
				"app": "guestbook",
			},
			ClusterName: "",
		},
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"192.168.122.125"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-76-61.eu-west-1.compute.internal"),
				TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "default",
					Name:            "guestbook-qjqnz",
					UID:             "9ca74e86-5573-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "1550311",
					FieldPath:       "",
				},
			},
			{
				Addresses: []string{"192.168.176.235"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-153-1.eu-west-1.compute.internal"),
				TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "default",
					Name:            "guestbook-92wmq",
					UID:             "9ca662bb-5573-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "1550311",
					FieldPath:       "",
				},
			},
			{
				Addresses: []string{"192.168.251.65"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-76-61.eu-west-1.compute.internal"),
				TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "default",
					Name:            "guestbook-qbdv8",
					UID:             "9ca74e86-5573-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "1550311",
					FieldPath:       "",
				},
			},
		},
	},
	&discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "kubernetes",
			GenerateName:    "",
			Namespace:       "default",
			SelfLink:        "/api/v1/namespaces/default/endpoints/kubernetes",
			UID:             "4daf1688-4c0a-11e9-b47e-066a7a20bac8",
			ResourceVersion: "5807557",
			Generation:      0,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			ClusterName: "",
		},
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"192.168.174.242", "192.168.82.3"},
			},
		},
	},
	&discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "redis-master",
			GenerateName:    "",
			Namespace:       "default",
			SelfLink:        "/api/v1/namespaces/default/endpoints/redis-master",
			UID:             "74ac431b-5573-11e9-b47e-066a7a20bac8",
			ResourceVersion: "1550146",
			Generation:      0,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			Labels: map[string]string{
				"app":  "redis",
				"role": "master",
			},
			ClusterName: "",
		},
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"192.168.108.68"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-76-61.eu-west-1.compute.internal"),
				TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "default",
					Name:            "redis-master-rh2bd",
					UID:             "5d7825f3-5573-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "1550097",
					FieldPath:       "",
				},
			},
		},
	},
	&discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "redis-slave",
			GenerateName:    "",
			Namespace:       "default",
			SelfLink:        "/api/v1/namespaces/default/endpoints/redis-slave",
			UID:             "8dee375e-5573-11e9-b47e-066a7a20bac8",
			ResourceVersion: "1550242",
			Generation:      0,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			Labels: map[string]string{
				"app":  "redis",
				"role": "slave",
			},
			ClusterName: "",
		},
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"192.168.186.217"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-153-1.eu-west-1.compute.internal"), TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "default",
					Name:            "redis-slave-mdjsj",
					UID:             "5d7825f3-5573-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "1550097",
					FieldPath:       "",
				},
			},
			{
				Addresses: []string{"192.168.68.108"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-76-61.eu-west-1.compute.internal"), TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "default",
					Name:            "redis-slave-gtd5x",
					UID:             "5d7825f3-5573-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "1550097",
					FieldPath:       "",
				},
			},
		},
	},
	&discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "kube-controller-manager",
			GenerateName:    "",
			Namespace:       "kube-system",
			SelfLink:        "/api/v1/namespaces/kube-system/endpoints/kube-controller-manager",
			UID:             "4f77dc4b-4c0a-11e9-b47e-066a7a20bac8",
			ResourceVersion: "6461574",
			Generation:      0,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			Annotations: map[string]string{
				"control-plane.alpha.kubernetes.io/leader": "{\"holderIdentity\":\"ip-10-0-189-120.eu-west-1.compute.internal_89407f85-57e1-11e9-b6ea-02eb484bead6\",\"leaseDurationSeconds\":15,\"acquireTime\":\"2019-04-05T20:34:54Z\",\"renewTime\":\"2019-05-06T20:04:02Z\",\"leaderTransitions\":1}",
			},
			ClusterName: "",
		},
	},
	&discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "kube-dns",
			GenerateName:    "",
			Namespace:       "kube-system",
			SelfLink:        "/api/v1/namespaces/kube-system/endpoints/kube-dns",
			UID:             "5049bf97-4c0a-11e9-b47e-066a7a20bac8",
			ResourceVersion: "5847",
			Generation:      0,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			Labels: map[string]string{
				"eks.amazonaws.com/component":   "kube-dns",
				"k8s-app":                       "kube-dns",
				"kubernetes.io/cluster-service": "true",
				"kubernetes.io/name":            "CoreDNS",
			},
			ClusterName: "",
		},
		Endpoints: []discoveryv1.Endpoint{
			{
				Addresses: []string{"192.168.212.227"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-200-63.eu-west-1.compute.internal"),
				TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "kube-system",
					Name:            "coredns-7554568866-26jdf",
					UID:             "503e1eae-4c0a-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "5842",
					FieldPath:       "",
				},
			},
			{
				Addresses: []string{"192.168.222.250"},
				Hostname:  aws.String(""),
				NodeName:  aws.String("ip-192-168-200-63.eu-west-1.compute.internal"),
				TargetRef: &v1.ObjectReference{
					Kind:            "Pod",
					Namespace:       "kube-system",
					Name:            "coredns-7554568866-shwn6",
					UID:             "503e1eae-4c0a-11e9-b47e-066a7a20bac8",
					APIVersion:      "",
					ResourceVersion: "5842",
					FieldPath:       "",
				},
			},
		},
	},
	&discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "kube-scheduler",
			GenerateName:    "",
			Namespace:       "kube-system",
			SelfLink:        "/api/v1/namespaces/kube-system/endpoints/kube-scheduler",
			UID:             "4e8782bc-4c0a-11e9-b47e-066a7a20bac8",
			ResourceVersion: "6461575",
			Generation:      0,
			CreationTimestamp: metav1.Time{
				Time: time.Now(),
			},
			Annotations: map[string]string{
				"control-plane.alpha.kubernetes.io/leader": "{\"holderIdentity\":\"ip-10-0-189-120.eu-west-1.compute.internal_949a4400-57e1-11e9-a7bb-02eb484bead6\",\"leaseDurationSeconds\":15,\"acquireTime\":\"2019-04-05T20:34:57Z\",\"renewTime\":\"2019-05-06T20:04:02Z\",\"leaderTransitions\":1}",
			},
			ClusterName: "",
		},
	},
}

func setUpEndpointClient() (*epClient, chan struct{}) {
	stopChan := make(chan struct{})

	client := &epClient{
		stopChan: stopChan,
		store:    NewObjStore(transformFuncEndpoint),
		inited:   true, //make it true to avoid further initialization invocation.
	}
	return client, stopChan
}

func TestEpClient_PodKeyToServiceNames(t *testing.T) {
	client, stopChan := setUpEndpointClient()
	defer close(stopChan)

	client.store.Replace(endpointsArray, "")

	expectedMap := map[string][]string{
		"namespace:default,podName:redis-master-rh2bd":           {"redis-master"},
		"namespace:default,podName:redis-slave-mdjsj":            {"redis-slave"},
		"namespace:default,podName:redis-slave-gtd5x":            {"redis-slave"},
		"namespace:kube-system,podName:coredns-7554568866-26jdf": {"kube-dns"},
		"namespace:kube-system,podName:coredns-7554568866-shwn6": {"kube-dns"},
		"namespace:default,podName:guestbook-qjqnz":              {"guestbook"},
		"namespace:default,podName:guestbook-92wmq":              {"guestbook"},
		"namespace:default,podName:guestbook-qbdv8":              {"guestbook"},
	}
	resultMap := client.PodKeyToServiceNames()
	log.Printf("PodKeyToServiceNames (len=%v): %v", len(resultMap), awsutil.Prettify(resultMap))
	assert.DeepEqual(t, expectedMap, resultMap)
}

func TestEpClient_ServiceNameToPodNum(t *testing.T) {
	client, stopChan := setUpEndpointClient()
	defer close(stopChan)

	client.store.Replace(endpointsArray, "")

	expectedMap := map[Service]int{
		NewService("redis-slave", "default"):  2,
		NewService("kube-dns", "kube-system"): 2,
		NewService("redis-master", "default"): 1,
		NewService("guestbook", "default"):    3,
	}
	resultMap := client.ServiceToPodNum()
	log.Printf("ServiceNameToPodNum (len=%v): %v", len(resultMap), awsutil.Prettify(resultMap))
	assert.DeepEqual(t, expectedMap, resultMap)
}
