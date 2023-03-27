.PHONY: plan apply destroy

plan:
	cd sandbox && go run ../main.go plan

apply:
	cd sandbox && go run ../main.go apply

destroy:
	cd sandbox && go run ../main.go destroy

install:
	go install github.com/orangekame3/tftarget@latest

localstack:
	cd sandbox && docker compose -f compose.yml up -d localstack

test:
	cd cmd && go test

coverage:
	cd cmd && go test -covermode=count -coverprofile=c.out && go tool cover -html=c.out -o coverage.html

