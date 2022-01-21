package encryption

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
	"time"
)

//可以创建唯一id
//用于token等

var nTools_md5_sequeue uint64 = 0

func GetNewSeq() uint64 {
	nSeq := atomic.AddUint64(&nTools_md5_sequeue, 1)
	if nSeq == 0 {
		nSeq = atomic.AddUint64(&nTools_md5_sequeue, 1)
	}
	return nSeq
}

func IntToBytes(n uint64) []byte {
	x := int64(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func GetUniqueId(prifix string) string {
	toBytes := IntToBytes(uint64(time.Now().UnixNano()))
	toBytes = append(toBytes, IntToBytes(GetNewSeq())...)
	hash := md5.New()
	hash.Write(toBytes)
	return prifix + hex.EncodeToString(hash.Sum(nil)[0:16])
}
