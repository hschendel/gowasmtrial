all: client server

run: server
	cd server;\
	./server;\
	cd ..

server: server/server

server/server: server/server.go shared/some_entity.go
	cd server;\
	go build;\
	cd ..

client: server/_web/client.wasm server/_web/wasm_exec.js

server/_web/client.wasm: client/client.go shared/some_entity.go
	cd client;\
	GOOS=js GOARCH=wasm go build -o ../server/_web/client.wasm;\
	cd ..

server/_web/wasm_exec.js:
	 cp "$(GOROOT)/misc/wasm/wasm_exec.js" server/_web

clean:
	rm server/server;\
	rm server/_web/wasm_exec.js;\
	rm server/_web/client.wasm
