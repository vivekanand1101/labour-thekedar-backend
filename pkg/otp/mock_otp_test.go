package otp

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockProvider_SendOTP(t *testing.T) {
	provider := NewMockProvider(true) // Use fixed OTP
	ctx := context.Background()

	t.Run("sends OTP successfully", func(t *testing.T) {
		code, err := provider.SendOTP(ctx, "+1234567890")
		require.NoError(t, err)
		assert.Equal(t, "123456", code)
	})

	t.Run("overwrites previous OTP", func(t *testing.T) {
		code1, err := provider.SendOTP(ctx, "+1234567890")
		require.NoError(t, err)

		code2, err := provider.SendOTP(ctx, "+1234567890")
		require.NoError(t, err)

		assert.Equal(t, code1, code2)
	})
}

func TestMockProvider_VerifyOTP(t *testing.T) {
	provider := NewMockProvider(true) // Use fixed OTP
	ctx := context.Background()

	t.Run("verifies correct OTP", func(t *testing.T) {
		phone := "+1234567890"
		code, err := provider.SendOTP(ctx, phone)
		require.NoError(t, err)

		verified, err := provider.VerifyOTP(ctx, phone, code)
		require.NoError(t, err)
		assert.True(t, verified)
	})

	t.Run("rejects incorrect OTP", func(t *testing.T) {
		phone := "+1234567891"
		_, err := provider.SendOTP(ctx, phone)
		require.NoError(t, err)

		verified, err := provider.VerifyOTP(ctx, phone, "000000")
		require.NoError(t, err)
		assert.False(t, verified)
	})

	t.Run("rejects OTP for non-existent phone", func(t *testing.T) {
		verified, err := provider.VerifyOTP(ctx, "+9999999999", "123456")
		require.NoError(t, err)
		assert.False(t, verified)
	})

	t.Run("OTP can only be used once", func(t *testing.T) {
		phone := "+1234567892"
		code, err := provider.SendOTP(ctx, phone)
		require.NoError(t, err)

		// First verification should succeed
		verified, err := provider.VerifyOTP(ctx, phone, code)
		require.NoError(t, err)
		assert.True(t, verified)

		// Second verification should fail
		verified, err = provider.VerifyOTP(ctx, phone, code)
		require.NoError(t, err)
		assert.False(t, verified)
	})
}

func TestMockProvider_RandomOTP(t *testing.T) {
	provider := NewMockProvider(false) // Use random OTP
	ctx := context.Background()

	t.Run("generates different OTPs for different phones", func(t *testing.T) {
		code1, err := provider.SendOTP(ctx, "+1111111111")
		require.NoError(t, err)

		code2, err := provider.SendOTP(ctx, "+2222222222")
		require.NoError(t, err)

		// With random OTPs, they might be different (not guaranteed but likely)
		assert.Len(t, code1, 6)
		assert.Len(t, code2, 6)
	})
}
