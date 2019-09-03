package main // import "github.com/wickedev/cetus"

import (
	"fmt"

	"k8s.io/helm/pkg/proto/hapi/version"
)

func main() {
	var v version.Version = version.Version{
		SemVer:       "SemVer",
		GitCommit:    "GitCommit",
		GitTreeState: "GitTreeState",
	}

	fmt.Println(v.String())
}
