package pgr

import (
	"github.com/apache/thrift/lib/go/thrift"
)

func init() {
	transportSocket := thrift.NewTSocketConf("localhost:9090", nil)
	transport := thrift.NewTBufferedTransport(transportSocket, 8192)
	err := transport.Open()
	if err != nil {
		panic(err)
	}
	defer func(transport *thrift.TBufferedTransport) {
		err := transport.Close()
		if err != nil {
			panic(err)
		}
	}(transport)

	protocol := thrift.THeaderProtocolBinary

}
