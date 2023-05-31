package grpcserver

import (
	"bytes"
	"context"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"github.com/chai2010/webp"
	"github.com/whitekid/goxp"
	"github.com/whitekid/goxp/fx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"qrcodeapi/pkg/qrcode"
	"qrcodeapi/proto"
)

type v1alpha1ServiceImpl struct {
	proto.UnimplementedQRCodeServer
}

func New() proto.QRCodeServer {
	return &v1alpha1ServiceImpl{}
}

func (s *v1alpha1ServiceImpl) Version(context.Context, *emptypb.Empty) (*wrapperspb.StringValue, error) {
	return wrapperspb.String("v1alpha1"), nil
}

func (s *v1alpha1ServiceImpl) Generate(ctx context.Context, in *proto.Request) (*proto.Response, error) {
	q, err := qrcode.Text(in.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	width := goxp.Ternary(in.Width < 20, 200, int(in.Width))
	height := goxp.Ternary(in.Width < 20, 200, int(in.Height))

	width = fx.Min(fx.Of(fx.Max(fx.Of(20, width)), 200))
	height = fx.Min(fx.Of(fx.Max(fx.Of(20, height)), 200))

	img, err := q.Render(width, height)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var buf bytes.Buffer
	contentType := "image/png"
	accepts := strings.Split(strings.ToLower(in.Accept), ",")
	for _, accept := range accepts {
		switch strings.ToLower(accept) {
		case "image/jpeg", "image/jpg":
			contentType = "image/jpeg"
			err = jpeg.Encode(&buf, img, nil)
		case "image/gif":
			contentType = "image/gif"
			err = gif.Encode(&buf, img, nil)
		case "image/webp":
			contentType = "image/webp"
			err = webp.Encode(&buf, img, nil)
		default:
			png.Encode(&buf, img)
		}
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &proto.Response{
		ContentType: contentType,
		Image:       buf.Bytes(),
		Width:       int32(width),
		Height:      int32(height),
	}, nil
}
