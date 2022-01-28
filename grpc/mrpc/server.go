package mrpc

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"strings"
)

// ServiceInfo contains unary RPC method info, streaming RPC method info and metadata for a service.
type ServiceInfo struct {
	// Contains the implementation for the methods in this service.
	serviceImpl interface{}
	methods     map[string]*grpc.MethodDesc
	//streams     map[string]*StreamDesc
	mdata interface{}
}

// Server is a gRPC server to serve RPC requests.
type Server struct {
	services map[string]*ServiceInfo // service name -> service info
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]*ServiceInfo),
	}
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	//if ss != nil {
	//ht := reflect.TypeOf(sd.HandlerType).Elem()
	//st := reflect.TypeOf(ss)
	//if !st.Implements(ht) {
	//logger.Fatalf("grpc: Server.RegisterService found the handler of type %v that does not satisfy %v", st, ht)
	//}
	//}
	s.register(sd, ss)
}

func (s *Server) register(sd *grpc.ServiceDesc, ss interface{}) {
	//s.printf("RegisterService(%q)", sd.ServiceName)
	//if s.serve {
	//	logger.Fatalf("grpc: Server.RegisterService after Server.Serve for %q", sd.ServiceName)
	//}
	//if _, ok := s.services[sd.ServiceName]; ok {
	//logger.Fatalf("grpc: Server.RegisterService found duplicate service registration for %q", sd.ServiceName)
	//}

	info := &ServiceInfo{
		serviceImpl: ss,
		methods:     make(map[string]*grpc.MethodDesc),
		mdata:       sd.Metadata,
	}
	for i := range sd.Methods {
		d := &sd.Methods[i]
		info.methods[d.MethodName] = d
	}

	s.services[sd.ServiceName] = info
}

func (s *Server) HandleMessage(head string, data []byte) ([]byte, error) {
	if head != "" && head[0] == '/' {
		head = head[1:]
	}

	pos := strings.LastIndex(head, "/")
	if pos == -1 {
		return nil, fmt.Errorf("pos = -1")
	}

	service := head[:pos]
	method := head[pos+1:]

	srv, knownService := s.services[service]
	if knownService {
		if md, ok := srv.methods[method]; ok {
			return s.processRPC(srv, md, data)
		}
	}

	var errDesc string
	if !knownService {
		errDesc = fmt.Sprintf("unknown service %v", service)
	} else {
		errDesc = fmt.Sprintf("unknown method %v for service %v", method, service)
	}

	return nil, fmt.Errorf(errDesc)
}
func (s *Server) processRPC(info *ServiceInfo, md *grpc.MethodDesc, data []byte) ([]byte, error) {
	callback := func(v interface{}) error {
		vv, ok := v.(proto.Message)
		if !ok {
			return fmt.Errorf("faild to unmarshal, message is %T,want proto.Message", v)
		}
		return proto.Unmarshal(data, vv)
	}

	reply, err := md.Handler(info.serviceImpl, nil, callback, nil)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(reply.(proto.Message))
}

//type methodHandler
//func(srv interface{}, ctx context.Context, dec func(interface{}) error) (interface{}, error)
