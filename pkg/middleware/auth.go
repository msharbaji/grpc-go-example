package middleware

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"reflect"
)

var (
	ErrMissingMetadata  = status.Errorf(codes.InvalidArgument, "missing metadata")
	ErrMissingHmac      = status.Errorf(codes.InvalidArgument, "missing x-hmac-signature metadata")
	ErrMissingHmacKeyID = status.Errorf(codes.InvalidArgument, "missing x-hmac-key-id metadata")
	ErrUnauthorized     = status.Errorf(codes.Unauthenticated, "unauthorized")
)

type clientAuthInterceptor struct {
	hmacKeyID  string
	hmacSecret string
}

type serverAuthInterceptor struct {
	secrets map[string]string
}

func NewClientAuthInterceptor(hmacKeyID, hmacSecret string) grpc.UnaryClientInterceptor {
	c := &clientAuthInterceptor{
		hmacKeyID:  hmacKeyID,
		hmacSecret: hmacSecret,
	}

	return c.clientInterceptor
}

func NewServerAuthInterceptor(secrets map[string]string) grpc.UnaryServerInterceptor {
	s := &serverAuthInterceptor{
		secrets: secrets,
	}

	return s.serverInterceptor
}

func (c *clientAuthInterceptor) clientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	//tflog.Info(ctx, "authenticating request")
	ctx = metadata.AppendToOutgoingContext(ctx, "x-hmac-key-id", c.hmacKeyID)
	plaintext, err := plainText(req, method)

	if err != nil {
		return err
	}
	ctx = metadata.AppendToOutgoingContext(ctx, "x-hmac-signature", signature(c.hmacSecret, plaintext))
	//tflog.Info(ctx, "sending request")
	return invoker(ctx, method, req, reply, cc, opts...)
}

func (s *serverAuthInterceptor) serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger := log.With().Str("method", info.FullMethod).Logger()

	logger.Debug().Msg("authenticating request")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Debug().Msg("missing metadata")
		return nil, ErrMissingMetadata
	}

	hmacSign, ok := md["x-hmac-signature"]
	if !ok || len(hmacSign) != 1 {
		logger.Debug().Msg("missing or invalid x-hmac-signature metadata")
		return nil, ErrMissingHmac
	}

	hmacKeyID, ok := md["x-hmac-key-id"]
	if !ok || len(hmacKeyID) != 1 {
		logger.Debug().Msg("missing or invalid x-hmac-key-id metadata")
		return nil, ErrMissingHmacKeyID
	}

	secretKey, err := s.getHMACSecretKey(hmacKeyID[0])
	if err != nil {
		logger.Debug().Err(err).Msg("failed to get HMAC secret key")
		return nil, status.Errorf(codes.Internal, "failed to get HMAC secret key")
	}

	plaintext, err := plainText(req, info.FullMethod)
	if err != nil {
		logger.Debug().Err(err).Msg("failed to get plaintext")
		return nil, status.Errorf(codes.Internal, "failed to get plaintext")
	}

	// Compare HMAC signatures.
	if !hmac.Equal([]byte(hmacSign[0]), signatureBytes(secretKey, plaintext)) {
		logger.Debug().Msg("invalid HMAC signature")
		return nil, status.Errorf(codes.Unauthenticated, "invalid HMAC signature")
	}

	// Call the handler to process the request
	return handler(ctx, req)
}

func (s *serverAuthInterceptor) getHMACSecretKey(key string) (string, error) {
	log.Debug().Str("key", key).Msg("getting HMAC secret key")
	if secretKey, ok := s.secrets[key]; ok {
		return secretKey, nil
	}
	return "", ErrUnauthorized
}

func plainText(req interface{}, method string) (string, error) {
	// Use a bytes.Buffer to efficiently build the plain text representation
	var buf bytes.Buffer

	// Check if the request struct has any exported fields
	t := reflect.TypeOf(req).Elem()
	hasExportedFields := false
	for i := 0; i < t.NumField(); i++ {
		if t.Field(i).PkgPath == "" {
			hasExportedFields = true
			break
		}
	}

	// Encode the request only if it has exported fields
	if hasExportedFields {
		enc := gob.NewEncoder(&buf)
		if err := enc.Encode(req); err != nil {
			return "", fmt.Errorf("failed to encode request: %w", err)
		}
	}

	// Append the method to the buffer
	buf.WriteString("method=" + method)

	// Convert the buffer content to a string and return it
	return buf.String(), nil
}

func signature(secretKey string, message string) string {
	mac := hmac.New(sha512.New512_256, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func signatureBytes(secretKey string, message string) []byte {
	return []byte(signature(secretKey, message))
}
