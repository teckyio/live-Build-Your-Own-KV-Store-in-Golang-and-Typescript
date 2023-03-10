package main

import "fmt"
import "os"
import "bufio"
import "regexp"
import "net"



var cache map[string]string


var setRegex *regexp.Regexp
var renameRegex *regexp.Regexp

func set(key string, value string) {
	cache[key] = value
}

func get(key string) (string, bool) {
	value, ok := cache[key]
	return value, ok
}

func del(key string) {
	delete(cache, key)
}

func rename(oldKey string, newKey string) {
	value := cache[oldKey]
	delete(cache, oldKey)
	cache[newKey] = value
}

func loopServer(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection: ", err)
			return
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	fmt.Println("Connection established...")
	reader := bufio.NewReader(conn)
	netOutput := NetOutput{conn}
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read from socket, err:", err)
			conn.Close()
			return
		}
		fmt.Print("socket received line: ", text)
		text = text[0 : len(text)-1]
		ok := handleCommand(text, netOutput)
		if !ok {
			conn.Write([]byte("Invalid command\n"))
		}
	}
}

type Output interface {
	WriteLn(line string)
}

func handleCommand(text string, output Output) bool {
	if text[0:3] == "set" {
		match := setRegex.FindStringSubmatch(text)
		key := match[1]
		value := match[2]
		set(key, value)
		return true
	}
	if text[0:6] == "rename" {
		match := renameRegex.FindStringSubmatch(text)
		oldKey := match[1]
		newKey := match[2]
		rename(oldKey, newKey)
		return true
	}
	if text[0:3] == "get" {
		key := text[4:]
		value, ok := get(key)
		if ok {
			output.WriteLn("value:" + value)
		} else {
			output.WriteLn("(none)")
		}
		return true
	}
	if text[0:3] == "del" {
		key := text[4:]
		del(key)
		return true
	}
	return false
}

type CliOutput struct {}

func (cliOutput CliOutput) WriteLn(line string) {
	fmt.Println(line)
}

type NetOutput struct {
	conn net.Conn
}

func (netOutput NetOutput) WriteLn(line string) {
	netOutput.conn.Write([]byte(line+"\n"))
}

func main() {
	fmt.Println("Hello world")
	cache = make(map[string]string)

  setRegex = regexp.MustCompile(`set (\w*) (.*)`)
  renameRegex = regexp.MustCompile(`rename (\w*) (\w*)`)

	reader := bufio.NewReader(os.Stdin)


	protocol := "tcp"
	bindAddress := "localhost:8500"
	listener, err := net.Listen(protocol, bindAddress)
	if err != nil {
		fmt.Println("Failed to listen on", bindAddress, "err:", err)
		return
	}
	fmt.Println("Listening on " + protocol + "://" + bindAddress)

	go loopServer(listener)

	cliOutput := CliOutput{}

	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Failed to readline:", err)
			return
		}
		text = text[0 : len(text)-1]
		if text == "help" {
			fmt.Println("Available commands:")
			fmt.Println("  set <key> <value>")
			fmt.Println("  get <key>")
			fmt.Println("  del <key>")
			fmt.Println("  rename <old-key> <new-key>")
			fmt.Println("  exit")
			continue
		}
		if text == "exit" {
			listener.Close()
			return
		}
		ok := handleCommand(text, cliOutput)
		if !ok {
			fmt.Println("Unknown command. Type \"help\" to see all available commands.")
			fmt.Println("vvvv")
			fmt.Println(text)
			fmt.Println("^^^^")
		}
	}

	// set("alice", "first value")
	// value := get("alice")
	// fmt.Println("value of 'alice':", value)

	// set("alice", "second value")
	// value = get("alice")
	// fmt.Println("value of 'alice':", value)
}
