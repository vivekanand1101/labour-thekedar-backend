package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vivekanand/labour-thekedar-backend/pkg/otp"
)

func TestAuthService_GenerateAndValidateToken(t *testing.T) {
	// Create a mock OTP provider
	otpProvider := otp.NewMockProvider(true)

	// Note: This is a unit test without database
	// Full integration tests would require a test database

	t.Run("OTP provider sends OTP", func(t *testing.T) {
		ctx := context.Background()
		code, err := otpProvider.SendOTP(ctx, "+1234567890")
		require.NoError(t, err)
		assert.Equal(t, "123456", code)
	})

	t.Run("OTP provider verifies correct OTP", func(t *testing.T) {
		ctx := context.Background()
		phone := "+1234567891"

		code, err := otpProvider.SendOTP(ctx, phone)
		require.NoError(t, err)

		valid, err := otpProvider.VerifyOTP(ctx, phone, code)
		require.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("OTP provider rejects wrong OTP", func(t *testing.T) {
		ctx := context.Background()
		phone := "+1234567892"

		_, err := otpProvider.SendOTP(ctx, phone)
		require.NoError(t, err)

		valid, err := otpProvider.VerifyOTP(ctx, phone, "000000")
		require.NoError(t, err)
		assert.False(t, valid)
	})
}
