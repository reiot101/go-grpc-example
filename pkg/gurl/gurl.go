package gurl

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/fullstorydev/grpcurl"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

type FormatOptions struct {
	EmitJSONDefaultFields bool
	IncludeTextSeparator  bool
	AllowUnknownFields    bool
}

type Options struct {
	// grpc server info
	network     string
	address     string
	dialOptions []grpc.DialOption

	// grpcurl format options
	formatOptions  FormatOptions
	verbosityLevel int
	verbose        bool
	grpcHeaders    []string
}

func defaultOptions() Options {
	return Options{
		network: "tcp",
		address: ":0",
		dialOptions: []grpc.DialOption{
			grpc.WithInsecure(),
		},
		formatOptions: FormatOptions{
			EmitJSONDefaultFields: false,
			IncludeTextSeparator:  true,
			AllowUnknownFields:    false,
		},
		verbosityLevel: 0,
		grpcHeaders:    []string{},
	}
}

type Option func(*Options)

func Network(s string) Option {
	return func(o *Options) {
		o.network = s
	}
}

func Addr(s string) Option {
	return func(o *Options) {
		o.address = s
	}
}

func DialOption(args ...grpc.DialOption) Option {
	return func(o *Options) {
		o.dialOptions = append(o.dialOptions, args...)
	}
}

func Formatter(arg FormatOptions) Option {
	return func(o *Options) {
		o.formatOptions = arg
	}
}

func Headers(args ...string) Option {
	return func(o *Options) {
		o.grpcHeaders = append(o.grpcHeaders, args...)
	}
}

// VerbosityLevel
// 0 = default
// 1 = verbose
// 2 = very verbose
func VerbosityLevel(n int) Option {
	return func(o *Options) {
		o.verbosityLevel = n
	}
}

func Verbose() Option {
	return func(o *Options) {
		o.verbose = true
	}
}

func Call(ctx context.Context, serviceMethod string, payload string, opts ...Option) ([]byte, error) {
	// default options
	options := defaultOptions()

	// apply options
	for i := range opts {
		opts[i](&options)
	}

	// input
	var in []byte
	if payload != "" {
		in = []byte(payload)
	}

	// dial
	cc, err := grpcurl.BlockingDial(ctx, options.network, options.address, nil, options.dialOptions...)
	if err != nil {
		return nil, err
	}

	// get descriptor source from target server
	desc := grpcurl.DescriptorSourceFromServer(ctx, grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(cc)))

	// request messages
	rf, formatter, err := grpcurl.RequestParserAndFormatter(
		grpcurl.FormatJSON,
		desc,
		bytes.NewReader(in),
		grpcurl.FormatOptions{
			EmitJSONDefaultFields: options.formatOptions.EmitJSONDefaultFields,
			IncludeTextSeparator:  options.formatOptions.IncludeTextSeparator,
			AllowUnknownFields:    options.formatOptions.AllowUnknownFields,
		},
	)
	if err != nil {
		return nil, err
	}

	// output
	out := new(bytes.Buffer)
	// handler
	h := grpcurl.NewDefaultEventHandler(out, desc, formatter, options.verbose)

	// invoke grpc
	if err := grpcurl.InvokeRPC(ctx, desc, cc, serviceMethod, options.grpcHeaders, h, rf.Next); err != nil {
		if errStatus, ok := status.FromError(err); ok {
			h.Status = errStatus
		} else {
			return nil, fmt.Errorf("invalid service method %q", serviceMethod)
		}
		return nil, fmt.Errorf("failed to call %v", err)
	}

	if h.Status.Code() != codes.OK {
		grpcurl.PrintStatus(os.Stderr, h.Status, formatter)
	}

	return out.Bytes(), nil
}
