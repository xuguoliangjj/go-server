protoc.exe --plugin=protoc-gen-go=%GOPATH%\bin\protoc-gen-go.exe  --go_out=../src/message  --proto_path=../src/message ../src/message/message.proto ../src/message/client.proto ../src/message/game.proto
//protoc.exe --js_out=../src/message  --proto_path=../src/message ../src/message/message.proto ../src/message/client.proto ../src/message/game.proto
protoc.exe --cpp_out=../src/message/c++  --proto_path=../src/message ../src/message/message.proto ../src/message/client.proto ../src/message/game.proto