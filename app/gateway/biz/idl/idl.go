package idl

import (
	"fmt"
	"log"
	"strings"
	"sync/atomic"

	"clicky.website/clicky/gateway/biz/utils"
	"clicky.website/clicky/gateway/conf"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client/genericclient"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

var (
	SvcMapManagerInstance  SvcMapManager
	IDLTreeManagerInstance IDLTreeManager
)

type UpdateMsgSvc struct {
	Op    string // "add" / "del" / update
	Key   string
	Value genericclient.Client
}

type UpdateMsgIDL struct {
	Op    string // "add"/"del"
	Key   string
	Value uint64
}

type (
	// svc:genericClient
	SvcClient     map[string]genericclient.Client
	SvcMapManager struct {
		value atomic.Value
		ch    chan UpdateMsgSvc
	}
)

type (
	// path:modifyIndex
	IDLKeyModifyIndex map[string]uint64
	IDLTreeManager    struct {
		value atomic.Value
		ch    chan UpdateMsgIDL
	}
)

func InitSvcMap() {
	SvcMapManagerInstance = SvcMapManager{
		ch: make(chan UpdateMsgSvc, 100),
	}
	SvcMapManagerInstance.value.Store(make(SvcClient))

	IDLTreeManagerInstance = IDLTreeManager{
		ch: make(chan UpdateMsgIDL, 100),
	}
	IDLTreeManagerInstance.value.Store(make(IDLKeyModifyIndex))
	go SvcMapManagerInstance.svcMapUpdater()
	go IDLTreeManagerInstance.idlTreeUpdater()
	go listenIdlTree()
}

func listenIdlTree() {
	params := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": "idl/clicky/service/",
	}
	plan, err := watch.Parse(params)
	if err != nil {
		log.Fatalln(err)
	}

	plan.Handler = func(idx uint64, data interface{}) {
		if kvs, ok := data.(consulapi.KVPairs); ok {

			newSet := make(map[string]uint64)
			fmt.Printf("Watch (Index=%d), total %d KV :\n", idx, len(kvs))
			for _, kv := range kvs {
				fmt.Printf("  * %s (size=%d)\n", kv.Key, len(kv.Value))
				// split path
				splited := strings.Split(kv.Key, "/")
				// get idl file name without suffix
				idlFile := strings.Split(splited[len(splited)-1], ".")

				// idl/clicky/service/svcName/svcMain.thrift
				// pass svcMain.thrift does not exist
				if len(splited) < 5 {
					continue
				}
				// pass non idl file
				if len(idlFile) < 2 || idlFile[1] != "thrift" {
					continue
				}
				// total path is key
				// idl/clicky/service/svcName/svcMain.thrift
				newSet[kv.Key] = kv.ModifyIndex

			}

			// has newSet
			if len(newSet) > 0 {
				oldSet := IDLTreeManagerInstance.value.Load().(IDLKeyModifyIndex)
				// local cache key hash not equal with consul
				if utils.HashKeys(newSet) != utils.HashKeys(oldSet) {
					fmt.Println("local cache key hash not equal")
					hlog.Debug("local cache key hash not equal")

					// local cache elements equal with consul
					if len(oldSet) == len(newSet) {
						// Iterate over oldSet
						// Delete the ones that do not exist in the local cache
						for k := range oldSet {
							if _, ok := newSet[k]; !ok {
								// remove the old one
								fmt.Printf("remove key: %s\n", k)
								IDLTreeManagerInstance.DelIdlTree(k)
							}
						}

						// Iterate over newSet
						// Add the ones that do not exist in the local cache
						// todo delete can be deleted after the client is completely updated
						// todo later provide an interface to delete
						// todo here can consider using a map[string]bool to do the delete mark
						for k, newVal := range newSet {
							if _, ok := oldSet[k]; !ok {
								// add the new one
								fmt.Printf("add key: %s\n", k)
								IDLTreeManagerInstance.AddIdlTree(k, newVal)

							}
						}
					} else {

						fmt.Println("svc path structure changed")
						// Iterate over newSet
						// Add the ones that do not exist in the local cache
						for k, newVal := range newSet {
							if _, ok := oldSet[k]; !ok {
								// add new one
								fmt.Printf("add key: %s\n", k)
								IDLTreeManagerInstance.AddIdlTree(k, newVal)
							}
						}

						// Iterate over oldSet
						// Delete the ones that do not exist in the local cache
						for k := range oldSet {
							if _, ok := newSet[k]; !ok {
								fmt.Printf("remove key: %s\n", k)
								IDLTreeManagerInstance.DelIdlTree(k)
							}
						}

					}

				} else {
					// local cache key hash equal with consul
					// compare modifyIndex
					fmt.Println("compare modifyIndex")
					for k, newVal := range newSet {
						if oldVal, ok := oldSet[k]; ok {
							if newVal != oldVal {
								fmt.Printf("key: %s, modifyIndex changed! old: %d, new: %d\n", k, oldVal, newVal)
								IDLTreeManagerInstance.AddIdlTree(k, newVal)
							}
						}
					}
				}

			}

		}
	}

	// Run will block and keep Watch in the background
	if err := plan.Run(conf.GetConf().Registry.RegistryAddress[0]); err != nil {
		log.Fatalln(err)
	}
}

func (svc *SvcMapManager) svcMapUpdater() {
	for msg := range svc.ch {
		oldMap := svc.value.Load().(map[string]genericclient.Client)

		// copy on write
		newMap := make(map[string]genericclient.Client, len(oldMap))
		for k, v := range oldMap {
			newMap[k] = v
		}

		switch msg.Op {
		case "add":
			newMap[msg.Key] = msg.Value
		case "del":
			delete(newMap, msg.Key)

		}

		// atomically update
		svc.value.Store(newMap)
	}
}

func (idl *IDLTreeManager) idlTreeUpdater() {
	for msg := range idl.ch {
		old := idl.value.Load().(IDLKeyModifyIndex)

		// copy on write
		newTree := make(IDLKeyModifyIndex, len(old))

		for k, v := range old {
			newTree[k] = v
		}

		switch msg.Op {
		case "add":
			newTree[msg.Key] = msg.Value
			idlc := NewIDLContent()
			// filter service
			if "service" == strings.Split(msg.Key, "/")[2] {
				// parse the IDL file dependency
				idlc.pharse(msg.Key)
				// insert the IDL generic client
				idlc.getGenericClient()

			}

		case "del":
			delete(newTree, msg.Key)

			// find the IDL generic client
			_, exist := SvcMapManagerInstance.GetSvc(msg.Key)
			if exist {
				// remove the IDL generic client
				fmt.Printf("remove SvcMapManagerInstance key: %s\n", msg.Key)
				SvcMapManagerInstance.DelSvc(msg.Key)
			}

			// double check
			_, exist = SvcMapManagerInstance.GetSvc(msg.Key)

			if exist {
				// remove the IDL generic client
				fmt.Printf("remove SvcMapManagerInstance key: %s fail\n", msg.Key)
			}

		}

		// atomically update
		idl.value.Store(newTree)

		// print IDLTreeManagerInstance
		// print len
		fmt.Printf("IDLTreeManagerInstance length: %d\n", len(newTree))
		for k, v := range newTree {
			fmt.Printf("key: %s, modifyIndex: %d\n", k, v)
		}
	}
}

func (svc *SvcMapManager) AddSvc(key string, client genericclient.Client) {
	svc.ch <- UpdateMsgSvc{Op: "add", Key: key, Value: client}
}

func (svc *SvcMapManager) DelSvc(key string) {
	svc.ch <- UpdateMsgSvc{Op: "del", Key: key}
}

func (svc *SvcMapManager) GetSvc(key string) (genericclient.Client, bool) {
	m := svc.value.Load().(map[string]genericclient.Client)
	v, ok := m[key]
	return v, ok
}

func (idl *IDLTreeManager) AddIdlTree(key string, val uint64) {
	idl.ch <- UpdateMsgIDL{Op: "add", Key: key, Value: val}
}

func (idl *IDLTreeManager) DelIdlTree(key string) {
	idl.ch <- UpdateMsgIDL{Op: "del", Key: key}
}

func (idl *IDLTreeManager) GetIdlTree(key string) (genericclient.Client, bool) {
	m := idl.value.Load().(map[string]genericclient.Client)
	v, ok := m[key]
	return v, ok
}
