package idl

import (
	"fmt"
	"strings"

	"clicky.website/clicky/gateway/conf"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

// Store IDL file path(consul key) ane content
type IDLContent struct {
	// Main IDL file path
	MainIdlPath string
	// Main IDL and include IDL path:content map
	PathContent map[string]string
}

func NewIDLContent() *IDLContent {
	return &IDLContent{
		PathContent: make(map[string]string),
	}
}

// parse IDL file content from consul
//
// key: main idl file path
func (idlc *IDLContent) pharse(key string) {
	// from consul get idl file content
	pair, _, err := conf.ConsulOfficeClient.KV().Get(key, nil)
	if err != nil {
		panic(err)
	}

	// if pair is nil, it means the key is not found
	if pair == nil {
		fmt.Printf("Key:%s not found\n", key)
		hlog.Warn("Key not found: %s", key)
	} else {
		fmt.Printf("Key: %s\nValue: %s\n", pair.Key, pair.Value)

		// main IDL file
		if "service" == strings.Split(pair.Key, "/")[2] {
			hlog.Debug("Main IDL file: %s", pair.Key)

			lines := strings.Split(string(pair.Value), "\n")

			hasInclude := false
			// find include path
			for _, line := range lines {
				if strings.HasPrefix(line, "include") {
					// get include path
					includePath := strings.TrimSpace(strings.Split(line, " ")[1])
					// remove double quotes
					includePath = strings.Trim(includePath, "\"")
					// recursive parse include path
					idlc.pharse(includePath)
					// cache include path and content
					idlc.PathContent[pair.Key] = string(pair.Value)
					idlc.MainIdlPath = pair.Key
					hlog.Debug("MainPathContent: %s", pair.Key)
					hasInclude = true

				}
			}

			// main IDL file without include
			if !hasInclude {

				idlc.PathContent[pair.Key] = string(pair.Value)
				idlc.MainIdlPath = pair.Key

			}

			// not main IDL file
		} else {
			idlc.PathContent[pair.Key] = string(pair.Value)
		}
	}
}

func (idle *IDLContent) getGenericClient() {
	provider, err := generic.NewThriftContentWithAbsIncludePathProviderWithDynamicGo(
		idle.MainIdlPath,
		idle.PathContent)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		hlog.Errorf("error: %v\n", err)
	}

	// get generic client
	g, err := generic.HTTPThriftGeneric(provider)
	if err != nil {
		panic(err)
	}

	// get service name from main idl path
	svcName := strings.Split(idle.MainIdlPath, "/")[3]

	client, err := genericclient.NewClient(
		svcName,
		g,
		client.WithResolver(*conf.ConsulResolver),
		client.WithTransportProtocol(transport.TTHeader),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
	)
	if err != nil {
		panic(err)
	}

	SvcMapManagerInstance.AddSvc(svcName, client)
}
