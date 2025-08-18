protoc:
	docker run --rm -v $$(pwd):/workspace -w /workspace 94peter/grpc-gateway-builder \
		protoc -I. -I /proto -I/proto/validate \
		--go_out=. \
		--go-grpc_out=. \
		--grpc-gateway_out=. \
		--validate_out="lang=go:." \
		--openapiv2_out=openapi \
		proto/*

test:
	docker run -it --rm 94peter/grpc-gateway-builder sh
		echo $PATH         # 確保包含 /usr/local/bin/plugins
		which protoc-gen-validate

update-submodules:
	git submodule update --init --recursive