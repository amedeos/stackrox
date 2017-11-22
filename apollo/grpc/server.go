package grpc

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	golog "log"
	"net"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/stack-rox/apollo/pkg/logging"
	"bitbucket.org/stack-rox/apollo/pkg/tls/keys"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const grpcPort = ":8080"
const endpoint = "localhost" + grpcPort

var (
	log = logging.New("api")
)

// APIService is the service interface
type APIService interface {
	RegisterServiceServer(server *grpc.Server)
	RegisterServiceHandlerFromEndpoint(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

var apiServices []APIService

// Register adds a new APIService to the list of API services
func Register(service APIService) {
	apiServices = append(apiServices, service)
}

// API listens for new connections on APIPort 8080, and redirects them to the gRPC-Gateway
type API interface {
	Start()
}

type apiImpl struct{}

// NewAPI returns an API object.
func NewAPI() API {
	return &apiImpl{}
}

// Start runs the API in a goroutine.
func (a *apiImpl) Start() {
	go a.run()
}

// httpErrorLogger implements io.Writer interface. It is used to control
// error messages coming from http server which can be logged.
type httpErrorLogger struct {
}

// Write suppresses EOF error messages
func (l httpErrorLogger) Write(p []byte) (n int, err error) {
	if !bytes.Contains(p, []byte("EOF")) {
		return os.Stderr.Write(p)
	}
	return len(p), nil
}

func (a *apiImpl) run() {
	pool, pair, err := getPair()
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewClientTLSFromCert(pool, endpoint)),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_recovery.StreamServerInterceptor(),
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
				grpc_recovery.UnaryServerInterceptor(),
			),
		),
	)

	for _, service := range apiServices {
		service.RegisterServiceServer(grpcServer)
	}

	ctx := context.Background()
	dcreds := credentials.NewTLS(&tls.Config{
		ServerName:         endpoint,
		RootCAs:            pool,
		InsecureSkipVerify: true,
	})
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(dcreds)}

	mux := http.NewServeMux()
	gwmux := runtime.NewServeMux()
	for _, service := range apiServices {
		if err := service.RegisterServiceHandlerFromEndpoint(ctx, gwmux, endpoint, dopts); err != nil {
			panic(err)
		}
	}
	mux.Handle("/", gwmux)
	conn, err := net.Listen("tcp", grpcPort)
	if err != nil {
		panic(err)
	}
	srv := &http.Server{
		Addr:     endpoint,
		Handler:  grpcHandlerFunc(grpcServer, mux),
		ErrorLog: golog.New(httpErrorLogger{}, "", golog.LstdFlags),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{*pair},
			NextProtos:   []string{"h2"},
		},
	}
	fmt.Printf("grpc on port: %v\n", grpcPort)
	err = srv.Serve(tls.NewListener(conn, srv.TLSConfig))
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
	return
}

func getPair() (*x509.CertPool, *tls.Certificate, error) {
	cert, key, err := keys.GenerateStackRoxKeyPair()
	if err != nil {
		return nil, nil, err
	}
	pair, err := tls.X509KeyPair(cert.Key().PEM(), key.Key().PEM())
	if err != nil {
		return nil, nil, err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cert.Key().PEM())
	if !ok {
		return nil, nil, errors.New("Cert is invalid")
	}
	return pool, &pair, nil
}

func grpcHandlerFunc(grpcServer *grpc.Server, httpHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			httpHandler.ServeHTTP(w, r)
		}
	})
}
