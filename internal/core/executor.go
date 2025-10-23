package core

import (
	"errors"
	"redis-clone/internal/constant"
	"strconv"
	"syscall"
	"time"
)

// true doi voi string don gian va false doi voi error va string phuc tap
func cmdPING(args []string) []byte {
	var res []byte
	if len(args) > 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		res = Encode("PONG", true)
	} else {
		res = Encode(args[0], false)
	}
	return res
}

func cmdSet(args []string) []byte {
	if len(args) < 2 || len(args) == 3 || len(args) > 4 {
		return Encode(errors.New("ERR wrong number of 'SET' command"), false)
	}

	var key, value string
	key, value = args[0], args[1]
	var ttlMs int64 = -1
	if len(args) > 2 {
		ttlSec, err := strconv.ParseInt(args[3], 10, 64)
		if err != nil {
			return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
		}
		ttlMs = ttlSec * 1000
	}
	dictStore.Set(key, dictStore.NewObject(key, value, ttlMs))
	return constant.RespOk
}

func cmdGet(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'GET' command"), false)
	}
	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.RespNil
	}
	if dictStore.HasExpired(key) {
		return constant.RespNil
	}
	return Encode(obj.Value, false)
}

func cmdTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'TTL' command"), false)
	}
	key := args[0]
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.TtlKeyNotExist
	}
	exp, exist := dictStore.GetExpired(key)
	if !exist {
		return constant.TtlKeyExistNoExpire
	}
	remainms := int64(exp) - time.Now().UnixMilli()
	if remainms < 0 {
		return constant.TtlKeyNotExist
	}
	return Encode(remainms/1000, false)
}

func cmdExpire(args []string) []byte {
	if len(args) != 2 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'EXPIRE' command"), false)
	}
	key := args[0]
	ttlSec, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
	}
	obj := dictStore.Get(key)
	if obj == nil {
		return constant.ExpireKeyNotExist
	}
	if ttlSec <= 0 {
		dictStore.Delete(key)
		return constant.ExpireKeySuccess
	}
	expireAt := time.Now().UnixMilli() + ttlSec*1000
	dictStore.SetExpired(key, expireAt)
	return constant.ExpireKeySuccess
}

func cmdDel(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'Del' command"), false)
	}
	deleteCount := 0
	for _, key := range args {
		obj := dictStore.Get(key)
		if obj != nil {
			dictStore.Delete(key)
			deleteCount++
		}
	}
	return Encode(int64(deleteCount), false)
}

func cmdExists(args []string) []byte {
	if len(args) < 1 {
		return Encode(errors.New("(error) ERR wrong number of arguments for 'Exists' command"), false)
	}
	existCount := 0
	for _, key := range args {
		obj := dictStore.Get(key)
		if obj != nil {
			existCount++
		}
	}
	return Encode(int64(existCount), false)
}

func ExecuteAndResponse(cmd *Command, connFd int) error {
	var res []byte
	switch cmd.Cmd {
	case "PING":
		res = cmdPING(cmd.Args)
	case "SET":
		res = cmdSet(cmd.Args)
	case "GET":
		res = cmdGet(cmd.Args)
	case "TTL":
		res = cmdTTL(cmd.Args)
	case "EXPIRE":
		res = cmdExpire(cmd.Args)
	case "DEL":
		res = cmdDel(cmd.Args)
	case "EXISTS":
		res = cmdExists(cmd.Args)
	case "ZADD":
		res = cmdZADD(cmd.Args)
	case "ZSCORE":
		res = cmdZSCORE(cmd.Args)
	case "ZRANK":
		res = cmdZRANK(cmd.Args)
	case "SADD":
		res = cmdSADD(cmd.Args)
	case "SREM":
		res = cmdSREM(cmd.Args)
	case "SMEMBERS":
		res = cmdSMEMBERS(cmd.Args)
	case "SISMEMBER":
		res = cmdSISMEMBER(cmd.Args)
	case "CMS.INITBYDIM":
		res = cmdCMSINITBYDIM(cmd.Args)
	case "CMS.INITBYPROB":
		res = cmdCMSINITBYPROB(cmd.Args)
	case "CMS.INCRBY":
		res = cmdCMSINCRBY(cmd.Args)
	case "CMS.QUERY":
		res = cmdCMSQUERY(cmd.Args)
	case "BF.RESERVE":
		res = cmdBFRESERVE(cmd.Args)
	case "BF.MADD":
		res = cmdBFMADD(cmd.Args)
	case "BF.EXISTS":
		res = cmdBFEXISTS(cmd.Args)
	default:
		res = []byte("-CMD NOT FOUND\r\n")
	}
	_, err := syscall.Write(connFd, res)
	return err
}
