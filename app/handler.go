package app

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gocraft/web"
	"github.com/gofrs/uuid"
	"net/http"
	"strings"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Auth(rw web.ResponseWriter, req *web.Request) (interface{}, error) {
	param := req.URL.Query().Get("user_id")
	userId, err := uuid.FromString(param)
	if err != nil {
		return nil, err
	}
	tokenPair, err := h.service.Auth(req.Context(), &userId)
	if err != nil {
		return nil, err
	}

	tokenPair.Refresh = base64.StdEncoding.EncodeToString([]byte(tokenPair.Refresh))
	return tokenPair, nil
}

func (h *Handler) Refresh(rw web.ResponseWriter, req *web.Request) (interface{}, error) {
	accessToken, err := parseHeader(req)
	if err != nil {
		return nil, err
	}
	refreshToken, err := base64.StdEncoding.DecodeString(req.URL.Query().Get("refreshToken"))
	if err != nil {
		return nil, err
	}
	tokenPair, err := h.service.Refresh(req.Context(), accessToken, refreshToken)
	if err != nil {
		return nil, err
	}

	tokenPair.Refresh = base64.StdEncoding.EncodeToString([]byte(tokenPair.Refresh))
	return tokenPair, nil
}

func parseHeader(req *web.Request) (string, error) {
	header := req.Header.Get("Authorization")
	if header == "" {
		return "", ErrAuthHeader
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" || len(headerParts[1]) == 0 {
		return "", ErrAuthHeader
	}
	return headerParts[1], nil
}

type EndpointHandler func(rw web.ResponseWriter, req *web.Request) (interface{}, error)

func WrapEndpoint(h EndpointHandler) interface{} {
	fn := func(rw web.ResponseWriter, req *web.Request, h EndpointHandler) error {
		result, err := h(rw, req)
		if err != nil {
			return err
		}
		data, err := json.Marshal(result)
		if err != nil {
			return err
		}
		_, err = rw.Write(data)
		return err
	}
	return func(rw web.ResponseWriter, req *web.Request) {
		err := fn(rw, req, h)
		if err != nil {
			writeCode(rw, err)
		}
	}
}

func writeCode(rw web.ResponseWriter, err error) {
	switch err {
	case ErrNotFound:
		rw.WriteHeader(http.StatusNotFound)
	case ErrAuthHeader:
		rw.WriteHeader(http.StatusUnauthorized)
	case ErrWrongToken:
		rw.WriteHeader(http.StatusUnauthorized)
	default:
		rw.WriteHeader(http.StatusInternalServerError)
	}
}
