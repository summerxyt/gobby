// TODO: need to deal with corner cases: master failure in several timings

package chubbyclient

import (
	"errors"
	"github.com/gobby/src/config"
	"github.com/gobby/src/rpc/chubbyrpc"
	"math/rand"
	"net/rpc"
	"strconv"
)

type chubbyclient struct {
	masterHostPort string
	masterConn     *rpc.Client
}

func NewClient(numNodes int, idx int) (Chubbyclient, error) {
	client := new(chubbyclient)
	//TODO:How to get master?
	/*for _, hostport := range config.Hostports {
	    if conn, err := rpc.DialHTTP("tcp", hostport); err == nil {
	        args := new(chubbyrpc.GetMasterArgs)
	        reply := new(chubbyrpc.GetMasterReply)
	        if err = conn.Call("ChubbyServer.GetMasterHostport", args, reply); err == nil {
	            client.masterHostPort = reply.Hostport
	            break
	        }
	    }
	}*/

	if idx >= 0 {
		client.masterHostPort = config.Nodes[idx].Address + ":" + strconv.Itoa(config.Nodes[idx].Port)

		if conn, err := rpc.DialHTTP("tcp", client.masterHostPort); err == nil {
			client.masterConn = conn
		}
		return client, nil
	}

	//Now random choose one
	i := rand.Int31n(int32(numNodes))
	client.masterHostPort = config.Nodes[i].Address + ":" + strconv.Itoa(config.Nodes[i].Port)

	if conn, err := rpc.DialHTTP("tcp", client.masterHostPort); err == nil {
		client.masterConn = conn
	}
	return client, nil
}

func (client *chubbyclient) Put(key, value string) error {
	args := new(chubbyrpc.PutArgs)
	args.Key = key
	args.Value = value
	//reply := new(chubbyrpc.PutReply)
	reply := new(chubbyrpc.ChubbyReply)
	if err := client.masterConn.Call("ChubbyServer.Put", args, reply); err == nil {
		if reply.Status == chubbyrpc.OK {
			return nil
		} else {
			return errors.New("Put error")
		}
	} else {
		return err
	}
}

func (client *chubbyclient) Get(key string) (string, error) {
	args := new(chubbyrpc.GetArgs)
	args.Key = key
	//reply := new(chubbyrpc.GetReply)
	reply := new(chubbyrpc.ChubbyReply)
	if err := client.masterConn.Call("ChubbyServer.Get", args, reply); err == nil {
		if reply.Status == chubbyrpc.OK {
			return reply.Value, nil
		} else {
			return "", errors.New("Get error")
		}
	} else {
		return "", err
	}
}

func (client *chubbyclient) Acquire(key string) (string, error) {
	args := new(chubbyrpc.AcquireArgs)
	args.Key = key
	//reply := new(chubbyrpc.AquireReply)
	reply := new(chubbyrpc.ChubbyReply)
	if err := client.masterConn.Call("ChubbyServer.Acquire", args, reply); err == nil {
		if reply.Status == chubbyrpc.OK {
			return reply.Value, nil
		} else {
			return reply.Value, errors.New("Acquire error")
		}
	} else {
		return reply.Value, err
	}
}

func (client *chubbyclient) Release(key, lockstamp string) error {
	args := new(chubbyrpc.ReleaseArgs)
	args.Key = key
	args.Lockstamp = lockstamp
	//reply := new(chubbyrpc.ReleaseReply)
	reply := new(chubbyrpc.ChubbyReply)
	if err := client.masterConn.Call("ChubbyServer.Release", args, reply); err == nil {
		if reply.Status == chubbyrpc.OK {
			return nil
		} else {
			return errors.New("Release error")
		}
	} else {
		return err
	}
}
