package dgraph_ldap

import (
	"google.golang.org/grpc"
	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"
	"context"
	"encoding/json"
)

type DGraphLDAPService struct {
	conn *grpc.ClientConn
	dg *client.Dgraph
	Status *DGraphLDAPStatus
}

func NewService(address string, options ...grpc.DialOption) (*DGraphLDAPService, error) {
	var conn *grpc.ClientConn
	var err error
	if options == nil {
		conn, err = grpc.Dial("127.0.0.1:9080", grpc.WithInsecure())
	} else {
		conn, err = grpc.Dial("127.0.0.1:9080", options...)
	}
	if err != nil {
		return nil, err
	}

	dc := api.NewDgraphClient(conn)
	dg := client.NewDgraphClient(dc)
	ldap := &DGraphLDAPService{conn: conn, dg: dg}
	if err = ldap.init(); err != nil {
		return nil, err
	}
	return ldap, nil
}

func (ldap *DGraphLDAPService) init() error {
	// defer ldap.conn.Close()
	if err, ok := ldap.RefreshStatus().(*Error); ok && err.Code() == ErrorGraphNotInitialized {
		if err := ldap.initDGraphDatabase(); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (ldap *DGraphLDAPService) Close() {
	defer ldap.conn.Close()
}

type statusQueryResult struct {
	Status []DGraphLDAPStatus `json:"status"`
}

func (ldap *DGraphLDAPService) RefreshStatus() error {
	resp, err := ldap.dg.NewTxn().Query(context.Background(), `
		{
		  status(func: uid(0x1)) {
		    expand(_all_)
		  }
		}`)
	if err != nil {
		return err // TODO !!!!!!!
	}
	var ret statusQueryResult // TODO !!!!!!! *statusQueryResult
	if err = json.Unmarshal(resp.Json, &ret); err != nil {
		return err  // TODO: More detail: 必须是空数据库
	}
	if len(ret.Status) > 0 {
		ldap.Status = &ret.Status[0]
		return nil
	} else {
		return NewError(ErrorGraphNotInitialized, "DGraph Database Not Initialized")
	}
}



/*
{
  "data": {
    "aaa": [
      {
        "uid": "0x1",
        "name": "Alice",
        "car": "MA0123",
        "mobile": "040123456"
      }
    ]
  }
}
 */

