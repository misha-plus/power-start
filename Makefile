.PHONY: web-dep agent frontend server all

web-dep:
	cd web && npm install

frontend: web-dep
	cd web && npm run build

server: frontend
	packr build -o artifacts/power-start-server server/*.go

agent:
	go build -o artifacts/power-start-agent agent/agent.go

all: agent server
