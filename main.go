package main

import (
	"net"

	"github.com/siddontang/go-log/log"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/siddontang/go-mysql/test_util/test_keys"

	"crypto/tls"
	"time"
)

type RemoteThrottleProvider struct {
	*server.InMemoryProvider
	delay int // in milliseconds
}

func (m *RemoteThrottleProvider) GetCredential(username string) (password string, found bool, err error) {
	time.Sleep(time.Millisecond * time.Duration(m.delay))
	return m.InMemoryProvider.GetCredential(username)
}

func main() {
	l, _ := net.Listen("tcp", "127.0.0.1:3309")
	// user either the in-memory credential provider or the remote credential provider (you can implement your own)
	//inMemProvider := server.NewInMemoryProvider()
	//inMemProvider.AddUser("root", "123")
	remoteProvider := &RemoteThrottleProvider{server.NewInMemoryProvider(), 10 + 50}
	remoteProvider.AddUser("root", "123")

	var tlsConf = server.NewServerTLSConfig(test_keys.CaPem,
		test_keys.CertPem, test_keys.KeyPem, tls.VerifyClientCertIfGiven)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			// Create a connection with user root and an empty password.
			// You can use your own handler to handle command here.
			// const version = "8.0.12"
			const version = "5.7.0"

			svr := server.NewServer(version,
				mysql.DEFAULT_COLLATION_ID, mysql.AUTH_CACHING_SHA2_PASSWORD, test_keys.PubPem, tlsConf)
			conn, err := server.NewCustomizedConn(c, svr, remoteProvider, &MysqlHandler{})

			if err != nil {
				log.Errorf("Connection error: %v", err)
				return
			}

			for {
				if herr := conn.HandleCommand(); herr != nil {
					log.Fatal(herr)
				}
			}
		}()
	}
}