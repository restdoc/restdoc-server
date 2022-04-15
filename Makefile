all:  pack


pack:
	sh cmd.sh
	GOOS=linux GOARCH=amd64 go build -o restdoc-server main.go

opack:
	sh cmd.sh
	GOOS=linux GOARCH=amd64 go build -o restdoc-server main.go

saas: pack qcloud

deploy:
	scp restdoc-server server1:/data/www/restdoc-server/_tmp
	ssh server1 "cd /data/www/restdoc-server/ && mv _tmp restdoc-server && supervisorctl -c /data/www/restdoc_supervisord.conf restart restdoc-server"


fmt:
	go fmt ./...
