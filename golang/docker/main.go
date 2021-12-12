package main

func main() {
	id := startContainer()
	if id != "" {
		getContainerLogs(id)
	}
}
