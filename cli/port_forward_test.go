package main

import (
	"bytes"
	"flag"
	"fmt"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/client-go/util/homedir"
	"net/http"
	URL "net/url"
	"path/filepath"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig",
			filepath.Join(".kube", "config"),
			"(optional absolute path to the kubeconfig file)")
	} else {
		kubeconfig = flag.String("kubeconfig", "",
			"absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		panic(err)
	}

	roundTripper, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", "todo-prod", "todo-backend-f498bfd-mx59j")
	hostIP := strings.TrimLeft(config.Host, "htps:/")
	serverURL := URL.URL{Scheme: "https", Path: path, Host: hostIP}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)

	forwarder, err := portforward.New(dialer, []string{"3000"}, stopChan, readyChan, out, errOut)
	if err != nil {
		panic(err)
	}

	go func() {
		for range readyChan { // Kubernetes will close this channel when it has something to tell us.
		}
		if len(errOut.String()) != 0 {
			panic(errOut.String())
		} else if len(out.String()) != 0 {
			fmt.Println(out.String())
		}
	}()

	if err = forwarder.ForwardPorts(); err != nil { // Locks until stopChan is closed.
		panic(err)
	}
}
