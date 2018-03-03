package dgraph_ldap

import (
	"github.com/dgraph-io/dgraph/protos/api"
	"context"
)

type DGraphLDAPStatus struct {
	Version uint64 `json:"version"`
	RootUID uint64 `json:"root_uid"`
	Initialized bool `json:"initialized"`
}

func (ldap *DGraphLDAPService) initDGraphDatabase() error {
	// schema
	if err := ldap.dg.Alter(context.Background(), &api.Operation{
		Schema: `
# global name:
gn: string @index(hash) .

# type: enum for: document, file, tag, user, ...
t: int @index(int) .

# name:
n: string @index(hash, exact, fulltext) @count .

# children:
ch: uid @reverse .

# create time:
ct: dateTime @index(day) .

# modify time:
mt: dateTime @index(hour) .

# user create:
uc: uid .

# owner:
uo: uid .

# group owner:
go: uid .

# permissions:
# roles = owner, group, user, tmp_user, guest
# permissions: 
# r: read data, w: write data,     // data: data, tags, d-xxx
# l: list children, c: append/remove children, 
# m: modify name/parrent, d: delete, 
# p: modify permission, owner, group
#
# excute !!!!!!!!!!!!!
#	EX|210clwr|43210clwr|765pdm43210clwr|765pdm43210clwr|765pdm43210clwr    // 63 bits -> int64
# E: has extra perm
# +/- no/ has extra perm ?  补码 XXXX deprecated !!!!!!!!
# guest can't be owner !!!!!!!!!!!!!!
# np.binary_repr(np.bitwise_and(np.int8(-0b1010110), np.int(0b1111111)))  wrong !!!
# np.binary_repr(np.bitwise_and(np.abs(np.int8(-0b1010110)), np.int(0b1111111)))  right
pm: int .
		`,
	}); err != nil {
		return err
	}

	initTxn := ldap.dg.NewTxn()
	defer initTxn.Discard(context.Background())




	return nil
}


