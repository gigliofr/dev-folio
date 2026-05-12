package server

import (
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

func createToken(subject string, ttl time.Duration) (string, error) {
    secret := os.Getenv("DEVFOLIO_JWT_SECRET")
    if secret == "" {
        return "", nil
    }
    claims := jwt.RegisteredClaims{
        Subject:   subject,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func parseAndValidateToken(tokenStr string) (*jwt.RegisteredClaims, error) {
    secret := os.Getenv("DEVFOLIO_JWT_SECRET")
    if secret == "" {
        return nil, nil
    }
    token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
        return claims, nil
    }
    return nil, nil
}
