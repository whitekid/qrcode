package grpcserver

import (
	"bytes"
	"context"
	"image/png"
	"net"
	"testing"

	"qrcodeapi/proto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func newTestClient(ctx context.Context, t *testing.T) proto.QRCodeClient {
	service := New()

	g := grpc.NewServer()
	proto.RegisterQRCodeServer(g, service)

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	go g.Serve(ln)
	go func() {
		<-ctx.Done()
		ln.Close()
	}()

	conn, err := grpc.DialContext(ctx, ln.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)

	return proto.NewQRCodeClient(conn)
}

func TestVersion(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := newTestClient(ctx, t)

	got, err := client.Version(ctx, &emptypb.Empty{})
	require.NoError(t, err)

	require.Equal(t, "v1alpha1", got.Value)
}

func TestGenerate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := newTestClient(ctx, t)

	type args struct {
		req *proto.Request
	}
	tests := [...]struct {
		name     string
		args     args
		wantErr  bool
		wantResp *proto.Response
	}{
		{`valid`, args{&proto.Request{Content: "hello world"}}, false, &proto.Response{ContentType: "image/png"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := client.Generate(ctx, tt.args.req)
			require.Truef(t, (err != nil) == tt.wantErr, `Generate() failed: error = %+v, wantErr = %v`, err, tt.wantErr)
			if tt.wantErr {
				return
			}

			require.Equal(t, tt.wantResp.ContentType, got.ContentType)
			require.NotEmpty(t, got.Image)

			img, err := png.Decode(bytes.NewReader(got.Image))
			require.NoError(t, err)
			require.Equal(t, 200, img.Bounds().Dx())
		})
	}
}
