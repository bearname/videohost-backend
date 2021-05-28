package domain

import (
	commonDto "github.com/bearname/videohost/pkg/common/dto"
	"github.com/bearname/videohost/pkg/video-scaler/domain"
	"io"
	"net/http"
)

type AuthService interface {
	CreateUser(requestBody io.ReadCloser) (domain.Token, error, int)
	Login(requestBody io.ReadCloser) (domain.Token, error, int)
	ValidateToken(authorizationHeader string) (commonDto.UserDto, int)
	RefreshToken(request *http.Request) (domain.Token, error, int)
}
