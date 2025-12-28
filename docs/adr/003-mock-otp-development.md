# ADR 003: Mock OTP for Development

## Status
Accepted

## Context
The application requires mobile phone OTP (One-Time Password) authentication. For development and testing:
- Real SMS providers cost money per message
- Testing requires predictable OTP values
- CI/CD pipelines cannot send real SMS
- Local development should not require external services

## Decision
We will implement a **Mock OTP Provider** for development that:
1. Stores OTPs in memory
2. Uses predictable OTP for testing (e.g., "123456")
3. Implements the same interface as real providers
4. Can be swapped for real provider via configuration

## Implementation

### OTP Interface
```go
type OTPProvider interface {
    SendOTP(phone string) error
    VerifyOTP(phone string, otp string) bool
}
```

### Mock Implementation
- Development OTP is always "123456"
- OTPs stored in-memory map
- 5-minute expiration
- Thread-safe with mutex

### Production Ready
When ready for production, implement:
- Twilio adapter
- MSG91 adapter
- Other SMS providers

Configuration determines which provider is used:
```
OTP_PROVIDER=mock    # Development
OTP_PROVIDER=twilio  # Production
```

## Alternatives Considered

### Always use real SMS
- Pros: Tests real integration
- Cons: Expensive, slow, unreliable in CI

### Email OTP
- Pros: Free to send
- Cons: Different user experience, not mobile-first

### Magic Links
- Pros: No code to enter
- Cons: Requires email, complex flow

## Consequences

### Positive
- Fast local development
- Free testing
- Predictable test scenarios
- Easy CI/CD integration

### Negative
- Mock behavior may differ from real providers
- Need to test real provider integration separately
- Risk of deploying with mock provider (mitigate with env checks)

## Security Considerations
- Mock provider MUST NOT be used in production
- Add startup check: fail if `GIN_MODE=release` and `OTP_PROVIDER=mock`
- Log warnings when using mock provider
