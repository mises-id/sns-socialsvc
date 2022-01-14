module github.com/mises-id/sns-socialsvc

go 1.16

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.12.0
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gogo/protobuf v1.3.3
	github.com/google/go-github/v33 v33.0.0
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/jinzhu/inflection v1.0.0
	github.com/joho/godotenv v1.4.0
	github.com/kr/pretty v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mises-id/mises-tm v0.0.0-20211229053907-7f70cc1f8835 // indirect
	github.com/mises-id/sdk v0.0.0-20220112092616-f50463033a82
	github.com/mises-id/sns-storagesvc/sdk v0.0.0-20211221064425-26bd51fd6a98
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.8.1
	golang.org/x/net v0.0.0-20220105145211-5b0dc2dfae98 // indirect
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/sys v0.0.0-20211216021012-1d35b9e2eb4e // indirect
	google.golang.org/genproto v0.0.0-20211208223120-3a66f561d7aa
	google.golang.org/grpc v1.43.0
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace github.com/tendermint/tendermint => github.com/mises-id/tendermint v0.34.15-0.20211207033151-1f29b59c0edf

replace github.com/cosmos/cosmos-sdk => github.com/mises-id/cosmos-sdk v0.44.6-0.20211209094558-a7c9c77cfc17
