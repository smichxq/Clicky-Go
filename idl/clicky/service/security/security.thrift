namespace go security

include "idl/clicky/include/common/empty.thrift"

struct Code2SessionReq {
    1: optional string code (api.query = 'code')
}

struct Code2SessionResp {
    1: optional string key (api.body = 'key')
    2: optional string token (api.body = 'token')
}

service security {
    Code2SessionResp Code2Session(1: Code2SessionReq req) (
        api.get = '/security/code2session',
        api.param = 'true',
    )
    empty.EmptyResp Demo(1: empty.EmptyResp req)

}
