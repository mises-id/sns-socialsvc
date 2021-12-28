module github.com/mises-id/sns-socialsvc

go 1.16

require (
	github.com/caarlos0/env v3.5.0+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.12.0
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/go-playground/validator v9.31.0+incompatible
	github.com/gogo/protobuf v1.3.2
	github.com/gorilla/mux v1.8.0
	github.com/jinzhu/inflection v1.0.0
	github.com/joho/godotenv v1.4.0
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/mises-id/sdk v0.0.0-20211111082026-f85731f62ba7
	github.com/mises-id/sns-storagesvc/sdk v0.0.0-20211221064425-26bd51fd6a98
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	go.mongodb.org/mongo-driver v1.8.1
	golang.org/x/crypto v0.0.0-20211202192323-5770296d904e // indirect
	golang.org/x/net v0.0.0-20211209124913-491a49abca63 // indirect
	golang.org/x/sys v0.0.0-20211205182925-97ca703d548d // indirect
	google.golang.org/genproto v0.0.0-20211129164237-f09f9a12af12
	google.golang.org/grpc v1.42.0
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

replace github.com/metaverse/truss => github.com/mises-id/truss v0.3.2-0.20211126092701-5f7d5bf015f1
