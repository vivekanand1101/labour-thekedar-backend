package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vivekanand/labour-thekedar-backend/internal/models"
	"github.com/vivekanand/labour-thekedar-backend/internal/repository"
	"github.com/vivekanand/labour-thekedar-backend/pkg/otp"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo    *repository.UserRepository
	otpProvider otp.Provider
	jwtSecret   string
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo *repository.UserRepository, otpProvider otp.Provider, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		otpProvider: otpProvider,
		jwtSecret:   jwtSecret,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Phone  string    `json:"phone"`
	jwt.RegisteredClaims
}

// TokenResponse represents the response with JWT tokens
type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         *models.User `json:"user"`
}

// SendOTPRequest represents the request to send OTP
type SendOTPRequest struct {
	Phone string `json:"phone" binding:"required,min=10,max=15"`
}

// VerifyOTPRequest represents the request to verify OTP
type VerifyOTPRequest struct {
	Phone string `json:"phone" binding:"required,min=10,max=15"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

// RefreshTokenRequest represents the request to refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// SendOTP sends an OTP to the given phone number
func (s *AuthService) SendOTP(ctx context.Context, phone string) error {
	_, err := s.otpProvider.SendOTP(ctx, phone)
	return err
}

// VerifyOTP verifies the OTP and returns JWT tokens
func (s *AuthService) VerifyOTP(ctx context.Context, phone, otpCode string) (*TokenResponse, error) {
	// Verify OTP
	valid, err := s.otpProvider.VerifyOTP(ctx, phone, otpCode)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, models.ErrInvalidOTP
	}

	// Get or create user
	user, _, err := s.userRepo.GetOrCreate(ctx, phone)
	if err != nil {
		return nil, err
	}

	// Generate tokens
	return s.generateTokens(user)
}

// RefreshToken refreshes the JWT tokens
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Parse and validate refresh token
	claims, err := s.validateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	return s.generateTokens(user)
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	return s.validateToken(tokenString)
}

func (s *AuthService) validateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (s *AuthService) generateTokens(user *models.User) (*TokenResponse, error) {
	now := time.Now()
	accessExpiry := now.Add(24 * time.Hour)       // 24 hours
	refreshExpiry := now.Add(7 * 24 * time.Hour)  // 7 days

	// Access token
	accessClaims := &Claims{
		UserID: user.ID,
		Phone:  user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	// Refresh token
	refreshClaims := &Claims{
		UserID: user.ID,
		Phone:  user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExpiry),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   user.ID.String(),
		},
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresAt:    accessExpiry,
		User:         user,
	}, nil
}
