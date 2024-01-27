build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/new_sms new_sms/main.go

.PHONY: test
test:
	go test ./new_sms

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose --aws-profile chatbot -r us-east-1 --force

.PHONY:
deploy_dev: clean build
	sls deploy -s dev --verbose --aws-profile chatbot -r us-east-1 --force

.PHONY:
deploy_prod: clean build
	sls deploy -s prod --verbose --aws-profile chatbot -r us-east-1 --force


# TODO:

# Get the complete golang-serveress project at: 
# https://github.com/Amo-Addai/api-feature-development/tree/main/serverless-apis/golang-serverless
