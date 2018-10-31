.PHONY: web-dep agent frontend server all

web-dep:
	cd web && npm install

frontend:
	cd web && npm run build

server:
	cd server && packr && go build -o ../artifacts/power-start-server *.go

agent:
	go build -o artifacts/power-start-agent agent/agent.go
