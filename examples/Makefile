.PHONY: all clean
all: echo msg task
	@echo "All demos built successfully."

clean:
	@echo "Cleaning demos..."
	rm -f echo_server.exe echo_client.exe
	rm -f msg_server.exe msg_client.exe
	rm -f task_server.exe task_client.exe

.PHONY: echo echo_clean
echo:
	@echo "Building echo demo..."
	go build -o echo_server.exe echo/server/main.go
	go build -o echo_client.exe echo/client/main.go
echo_clean:
	@echo "Cleaning echo demo..."
	rm -f echo_server.exe echo_client.exe

.PHONY: msg msg_clean
msg:
	@echo "Building msg demo..."
	go build -o msg_server.exe msg/server/main.go
	go build -o msg_client.exe msg/client/main.go
msg_clean:
	@echo "Cleaning msg demo..."
	rm -f msg_server.exe msg_client.exe

.PHONY: task task_clean
task:
	@echo "Building task demo..."
	go build -o task_server.exe task/server/main.go task/server/calculate.go
	go build -o task_client.exe task/client/main.go
task_clean:
	@echo "Cleaning task demo..."
	rm -f task_server.exe task_client.exe
