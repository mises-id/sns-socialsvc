

# ssh config mises_alpha
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/cli/main.go
upload:
	scp ./main mises_alpha:/apps/sns-socialsvc/
replace:
	ssh mises_alpha "mv /apps/sns-socialsvc/main /apps/sns-socialsvc/sns-socialsvc"
restart:
	ssh mises_alpha "sudo supervisorctl restart socialsvc"
deploy: build \
	upload \
	replace \
	restart

truss:
	truss proto/socialsvc.proto  --pbpkg github.com/mises-id/sns-socialsvc/proto --svcpkg github.com/mises-id/sns-socialsvc --svcout . -v 

test:
	APP_ENV=test go test -coverprofile coverage.out -v --tags tests -parallel 1  ./...

