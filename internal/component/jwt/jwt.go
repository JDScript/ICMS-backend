package jwt

import (
	"crypto/ecdsa"
	"errors"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt"
)

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌格式有误")
	ErrTokenInvalid           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("需要认证才能访问！")
	ErrHeaderMalformed        error = errors.New("请求头中 Authorization 格式有误")
)

type Jwt struct {
	Issuer     string
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
	Expire     time.Duration
	MaxRefresh time.Duration
}

type Config struct {
	Issuer           string
	PublicKey        string
	PrivateKey       string
	ExpireMinute     uint64
	MaxRefreshMinute uint64
}

type JwtCustomClaims struct {
	UserId string `json:"userId"`
}

type JwtFullClaims struct {
	JwtCustomClaims
	jwtpkg.StandardClaims
}

func New(config *Config) (*Jwt, error) {
	if config == nil {
		return nil, nil
	}

	priv, err := jwtpkg.ParseECPrivateKeyFromPEM([]byte(config.PrivateKey))
	if err != nil {
		return nil, err
	}
	pub, err := jwtpkg.ParseECPublicKeyFromPEM([]byte(config.PublicKey))
	if err != nil {
		return nil, err
	}

	if config.ExpireMinute == 0 {
		config.ExpireMinute = 120
	}

	if config.MaxRefreshMinute == 0 {
		config.MaxRefreshMinute = 43200
	}

	return &Jwt{
		Issuer:     config.Issuer,
		PublicKey:  pub,
		PrivateKey: priv,
		Expire:     time.Duration(config.ExpireMinute) * time.Minute,
		MaxRefresh: time.Duration(config.MaxRefreshMinute) * time.Minute,
	}, nil
}

func (jwt *Jwt) IssueToken(claims *JwtCustomClaims) (string, error) {
	expireAt := time.Now().Add(jwt.Expire)
	fullClaims := JwtFullClaims{
		JwtCustomClaims: *claims,
		StandardClaims: jwtpkg.StandardClaims{
			NotBefore: time.Now().Unix(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireAt.Unix(),
			Issuer:    jwt.Issuer,
		},
	}

	token, err := jwtpkg.NewWithClaims(jwtpkg.SigningMethodES256, fullClaims).SignedString(jwt.PrivateKey)

	if err != nil {
		return "", err
	}

	return token, nil
}

// ParserToken 解析 Token，中间件中调用
func (jwt *Jwt) ParserToken(c *gin.Context) (*JwtFullClaims, error) {
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}

	// 1. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 2. 解析出错
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	// 3. 将 token 中的 claims 信息解析出来和 JWTCustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*JwtFullClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 更新 Token，用以提供 refresh token 接口
func (jwt *Jwt) RefreshToken(c *gin.Context) (string, error) {
	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		// 满足 refresh 的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}

	// 4. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*JwtFullClaims)

	// 5. 检查是否过了『最大允许刷新的时间』
	x := time.Now().Add(-jwt.MaxRefresh).Unix()
	if claims.IssuedAt > x {
		// 修改过期时间
		claims.StandardClaims.ExpiresAt = time.Now().Add(jwt.Expire).Unix()

		token, err := jwtpkg.NewWithClaims(jwtpkg.SigningMethodES256, claims).SignedString(jwt.PrivateKey)

		if err != nil {
			return "", err
		}

		return token, nil
	}

	return "", ErrTokenExpiredMaxRefresh
}

func (jwt *Jwt) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &JwtFullClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.PublicKey, nil
	})
}

// getTokenFromHeader 使用 jwtpkg.ParseWithClaims 解析 Token
// Authorization:Bearer xxxxx
func (jwt *Jwt) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}
	return parts[1], nil
}
