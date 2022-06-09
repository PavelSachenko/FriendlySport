proto:
	protoc -I=. --go_out=paths=source_relative:.  --go-grpc_out=paths=source_relative:. pkg/**/pb/**/*.proto \
	&& sed -i "s/,omitempty//" pkg/**/pb/**/*.pb.go