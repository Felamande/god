package math

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"math/big"

	"github.com/Felamande/god/lib/jsvm"
	"github.com/robertkrimen/otto"
)

func init() {
	if b := jsvm.Builtin("Math"); b != nil {
		b.Extend("randStr", randStr)
		b.Extend("md5", md5fn)
	}
}

func randStr(call otto.FunctionCall) otto.Value {

	length, err := call.Argument(0).ToInteger()
	if err != nil {
		length = 20
	}
	str := getRandomString(length)
	v, _ := otto.ToValue(str)
	return v
}

func getRandomString(slen int64) string {
	str := []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^*()+=-")
	l := len(str)
	var re []byte
	for i := int64(0); i < slen; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(l)))
		re = append(re, str[n.Int64()])
	}
	return string(re)
}

func md5fn(call otto.FunctionCall) otto.Value {
	arg := call.Argument(0)

	if !arg.IsString() {
		md5str := getMd5("")
		v, _ := otto.ToValue(md5str)
		return v
	}
	str := arg.String()
	md5str := getMd5(str)
	v, _ := otto.ToValue(md5str)
	return v

}

func getMd5(source interface{}) string {
	ctx := md5.New()

	switch ss := source.(type) {
	case io.Reader:
		io.Copy(ctx, ss)
	case string:
		ctx.Write([]byte(ss))
	case []byte:
		ctx.Write(ss)

	}

	return hex.EncodeToString(ctx.Sum(nil))
}
