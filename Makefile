

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
	APP_ENV=test go test -coverprofile coverage.out  -count=1 --tags tests  -coverpkg=./app/... ./tests/...


 #backup
upload-backup:
	scp ./main mises_backup:/apps/sns-socialsvc/
replace-backup:
	ssh mises_backup "mv /apps/sns-socialsvc/main /apps/sns-socialsvc/sns-socialsvc"
restart-backup:
	ssh mises_backup "sudo supervisorctl restart socialsvc"
deploy-backup: build \
	upload-backup \
	replace-backup \
	restart-backup
coverage:
	go tool cover -html=coverage.out
