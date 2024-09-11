package services

import (
	"errors"
	"time"

	db "github.com/ebubekiryigit/golang-mongodb-rest-api-starter/models/db"
	"github.com/golang-jwt/jwt/v4"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/mgm/v3/field"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateToken create a new token record
func CreateToken(user *db.User, expiresAt time.Time) (*string, error) {
	claims := &db.UserClaims{
		UserInfo: *user,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   user.ID.Hex(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(Config.JWTSecretKey))
	if err != nil {
		return nil, errors.New("cannot create access token")
	}

	// tokenModel := db.NewToken(user.ID, tokenString, tokenType, expiresAt)
	// err = mgm.Coll(tokenModel).Create(tokenModel)
	// if err != nil {
	// 	return nil, errors.New("cannot save access token to db")
	// }

	// return tokenModel, nil
	return &tokenString, nil
}

// DeleteTokenById delete token with id
func DeleteTokenById(tokenId primitive.ObjectID) error {
	ctx := mgm.Ctx()
	deleteResult, err := mgm.Coll(&db.Token{}).DeleteOne(ctx, bson.M{field.ID: tokenId})
	if err != nil || deleteResult.DeletedCount <= 0 {
		return errors.New("cannot delete token")
	}

	return nil
}

// GenerateAccessTokens generates "access" and "refresh" token for user
func GenerateAccessTokens(user *db.User) (*string, *string, error) {
	accessExpiresAt := time.Now().Add(time.Duration(Config.JWTAccessExpirationMinutes) * time.Minute)
	refreshExpiresAt := time.Now().Add(time.Duration(Config.JWTRefreshExpirationDays) * time.Hour * 24)

	accessToken, err := CreateToken(user, accessExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := CreateToken(user, refreshExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

// VerifyToken checks jwt validity, expire date, blacklisted
func VerifyToken(token string, tokenType string) (*db.User, error) {
	claims := &db.UserClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Config.JWTSecretKey), nil
	})

	if err != nil {
		return nil, errors.New("not valid token")
	}
	// if time.Now().Sub(claims.ExpiresAt.Time) > 10*time.Second {
	if time.Since(claims.ExpiresAt.Time) > 10*time.Second {
		return nil, errors.New("token is expired")
	}

	return &claims.UserInfo, nil
}
