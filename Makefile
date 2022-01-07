truss:
	truss proto/socialsvc.proto  --pbpkg github.com/mises-id/sns-socialsvc/proto --svcpkg github.com/mises-id/sns-socialsvc --svcout . -v 

test:
	APP_ENV=test go test -coverprofile coverage.out -v --tags tests -parallel 1  ./...