package idl

import (
	"strings"

	"clicky.website/clicky/gateway/conf"
	"clicky.website/clicky/gateway/suite"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/generic"
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
		hlog.Warnf("Key not found: %s", key)
	} else {
		hlog.Debugf("Key: %s\nValue: %s\n", pair.Key, pair.Value)

		// main IDL file
		if "service" == strings.Split(pair.Key, "/")[2] {
			hlog.Debugf("Main IDL file: %s", pair.Key)

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
					hlog.Debugf("MainPathContent: %s", pair.Key)
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
		hlog.Errorf("error: %v\n", err)
	}

	cbMap := buildCBConfigMap(idle)

	for k := range cbMap {
		hlog.Debugf("cbm %s \n", k)
	}

	// get generic client
	g, err := generic.HTTPThriftGeneric(provider)
	if err != nil {
		panic(err)
	}

	// get service name from idl
	svcName := g.IDLServiceName()

	hlog.Debugf("Service name: %s", svcName)

	client, err := genericclient.NewClient(
		svcName,
		g,
		suite.GenericSuite(cbMap)...,
	)
	if err != nil {
		panic(err)
	}

	SvcMapManagerInstance.AddSvc(svcName, client)
}

func buildCBConfigMap(idle *IDLContent) map[string]circuitbreak.CBConfig {
	cbConfig := make(map[string]circuitbreak.CBConfig)
	thrift, err := generic.ParseContent(idle.MainIdlPath, idle.PathContent[idle.MainIdlPath], idle.PathContent, true)
	if err != nil {
		hlog.Error(err)
	}

	svc := thrift.GetServices()

	for _, v := range svc {
		for _, f := range v.Functions {
			key := "gateway" + "/" + v.Name + "/" + f.Name
			cbConfig[key] = circuitbreak.CBConfig{
				Enable:    true,
				ErrRate:   0.2,
				MinSample: 10,
			}

		}
	}

	return cbConfig
}
