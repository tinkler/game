Set-Variable -name SCRIPT_PATH -value $PWD.Path
Set-Location ../../api/proto
protoc.exe --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative game/*.proto
Set-Location $SCRIPT_PATH