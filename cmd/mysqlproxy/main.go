package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ando-masaki/mysqlproxy"
)

var (
	cfg  mysqlproxy.Config
	cfgs = map[bool]mysqlproxy.Config{
		true: mysqlproxy.Config{
			Addr:     "./mysqlproxy.sock",
			Password: "hoge",

			AllowIps:       "@",
			TlsServer:      false,
			TlsClient:      true,
			ClientCertFile: "client.pem",
			ClientKeyFile:  "client.key",
		},
		false: mysqlproxy.Config{
			Addr:     "0.0.0.0:9696",
			Password: "hoge",

			AllowIps:   "",
			TlsServer:  true,
			TlsClient:  false,
			CaCertFile: "ca.pem",
			CaKeyFile:  "ca.key",
		},
	}
	root *bool = flag.Bool("root", false, "Serve as root proxy server.")
)

func init() {
	flag.Parse()
	cfg = cfgs[*root]
	if cfg.TlsServer {
		ca_b, err := ioutil.ReadFile(cfg.CaCertFile)
		if err != nil {
			log.Fatal(err)
		}
		ca, err := x509.ParseCertificate(ca_b)
		if err != nil {
			log.Fatal(err)
		}
		priv_b, err := ioutil.ReadFile(cfg.CaKeyFile)
		if err != nil {
			log.Fatal(err)
		}
		priv, err := x509.ParsePKCS1PrivateKey(priv_b)
		if err != nil {
			log.Fatal(err)
		}
		pool := x509.NewCertPool()
		pool.AddCert(ca)

		cert := tls.Certificate{
			Certificate: [][]byte{ca_b},
			PrivateKey:  priv,
		}
		cfg.TlsServerConf = &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{cert},
			ClientCAs:    pool,
		}
		cfg.TlsServerConf.Rand = rand.Reader
	}
	if cfg.TlsClient {
		cert_b, err := ioutil.ReadFile(cfg.ClientCertFile)
		if err != nil {
			log.Fatal(err)
		}
		priv_b, err := ioutil.ReadFile(cfg.ClientKeyFile)
		if err != nil {
			log.Fatal(err)
		}
		priv, err := x509.ParsePKCS1PrivateKey(priv_b)
		if err != nil {
			log.Fatal(err)
		}
		cfg.TlsClientConf = &tls.Config{
			Certificates: []tls.Certificate{{
				Certificate: [][]byte{cert_b},
				PrivateKey:  priv,
			}},
			InsecureSkipVerify: true,
		}
	}

}

func main() {
	svr, err := mysqlproxy.NewServer(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-sc
		log.Printf("main Got signal: %s", sig)
		svr.Close()
		if cfg.TlsClient {
			if err := os.Remove(cfg.Addr); err != nil {
				log.Fatal(err)
			}
		}
	}()
	svr.Run()
}