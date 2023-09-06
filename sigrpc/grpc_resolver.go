package sigrpc

import (
	"strings"

	"google.golang.org/grpc/resolver"
)

// grpc resolver

type grpcResolverBuilder struct {
	scheme      string
	serviceName string
	addrs       []string
}

func (g *grpcResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &grpcResolver{
		target: target,
		cc:     cc,
		addrsStore: map[string][]string{
			g.serviceName: g.addrs,
		},
	}
	r.start()
	return r, nil
}
func (g *grpcResolverBuilder) Scheme() string { return g.scheme }

type grpcResolver struct {
	target     resolver.Target
	cc         resolver.ClientConn
	addrsStore map[string][]string
}

func removeFirstSlash(input string) string {
	if strings.HasPrefix(input, "/") {
		return input[1:]
	}
	return input
}
func (r *grpcResolver) endpointFromTarget() string {
	if r.target.URL.Path == "" {
		return removeFirstSlash(r.target.URL.Opaque)
	}
	return removeFirstSlash(r.target.URL.Path)
}
func (r *grpcResolver) start() {
	// fmt.Println(r.target.Endpoint, r.target.URL, r.target.URL.Path, r.target.URL.Opaque)
	// addrStrs := r.addrsStore[r.target.Endpoint]
	addrStrs := r.addrsStore[r.endpointFromTarget()]
	addrs := make([]resolver.Address, len(addrStrs))
	for i, s := range addrStrs {
		addrs[i] = resolver.Address{Addr: s}
	}
	r.cc.UpdateState(resolver.State{Addresses: addrs})
}
func (*grpcResolver) ResolveNow(o resolver.ResolveNowOptions) {}
func (*grpcResolver) Close()                                  {}
