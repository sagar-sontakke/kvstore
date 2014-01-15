/*
 * Assignment Number: 1
 * Subject: CS-733
 * Submitted by: Sagar Sontakke, Roll - 133050074, email: sagarsb@cse,iitb.ac.in
 */

package main
import (
	"net"
	"os"
	"fmt"
	"strings"
	"bufio"
	"io"
)

var oldKeyValue = make(map[string]string)
var newKeyValue = make(map[string]string)
var fmain, ftemp *os.File
var IS_DIRTY bool = true
var REMOVE_RECOVERY bool = true
var KEEP_RECOVERY bool = false
var DIRTY string = "DIRTY"
var CLEAN string = "CLEAN"
var find bool = false

func main() {

	/*
	 * Reference: Took the client-server connection mechanism idea from the book "Network programming with Go"
	 * which is available online
	 */

	service := ":1201"
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	/*
	 * At the start of server, recover data from recovery if in case there was any failure previously
	 */

	syncFiles(REMOVE_RECOVERY)

	/*
	 * Copy the contents of the previous file in memory so that it can used for store/retrieval operations
	 */

	IS_DIRTY = false
	fmain, _ = os.OpenFile("key-value.txt",os.O_CREATE|os.O_RDWR,0600)
	ftemp, _ = os.OpenFile("recovery.txt",os.O_CREATE|os.O_RDWR,0600)
	ftemp.Write([]byte(DIRTY))

	rd := bufio.NewReader(fmain)
	for  {
		line,err := rd.ReadString('\n')
		if err != nil {
				break
		}
		pair1 := strings.Split(line,"\n")
		pair2 := strings.Split(pair1[0]," ")
		oldKeyValue[pair2[0]] = pair2[1]
	}

	/*
	 * Listen for the requests from the client or server.
	 */

	for {
		go handleServer()
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		/*
		 * Run as a goroutine. Allows multi-threading so as to allow multiple clients simultaneously
		 */

		go handleClient(conn)
	}
}


/*
 * This for recoverying the data from the recovery file. All the key-value change transactions are logged
 * in the recovery.txt file. This logged data helps for the persistence of the data
 */

func syncFiles(removeRecovery bool) {
	fmain, _ := os.OpenFile("key-value.txt",os.O_CREATE|os.O_RDWR,0600)
	ftemp, _ := os.OpenFile("recovery.txt",os.O_CREATE|os.O_RDWR,0600)
	rdmain := bufio.NewReader(fmain)
	rdtemp := bufio.NewReader(ftemp)
	oldKeyValue := make(map[string]string)

	nline, err := rdtemp.ReadString('\n')
	if err == nil {
		line1 := strings.Split(nline,"\n")
		flag := string(line1[0])

		/*
		 * The "DIRTY" keyword at the start of recovery file shows that the key-value file was not
		 * sync'ed before. The IS_DIRTY flag shows that there happened some update/delete
		 * operations for the already existing key-values and these changes are still in memory.
		 * and need to be written back to key-value file
		 */

		if flag == "DIRTY" && IS_DIRTY == true {
			for  {
				line,err := rdmain.ReadString('\n')
				if err != nil {
					break
				}
				pair1 := strings.Split(line,"\n")
				pair2 := strings.Split(pair1[0]," ")
				key1 := string(pair2[0])
				val1 := string(pair2[1])
				oldKeyValue[key1] = val1
			}

			for  {
				line,err := rdtemp.ReadString('\n')

				if line == "" {
					break
				}
				if err == io.EOF {
					pair2 := strings.Split(line," ")
					operation := string(pair2[0])
					key1 := string(pair2[1])
					val1 := string(pair2[2])
					if operation == "SET" {
						oldKeyValue[key1] = val1
					} else if operation == "DEL" {
						delete(oldKeyValue,key1)
					}
					break
				}
				if err != nil{
					break
				}
				pair1 := strings.Split(line,"\n")
				pair2 := strings.Split(pair1[0]," ")
				operation := string(pair2[0])
				key1 := string(pair2[1])
				val1 := string(pair2[2])
				if operation == "SET" {
					oldKeyValue[key1] = val1
				} else if operation == "DEL" {
					delete(oldKeyValue,key1)
				}
			}

			os.Remove("key-value.txt")
			fmain, _ := os.OpenFile("key-value.txt",os.O_CREATE|os.O_RDWR,0600)
			for key, value := range oldKeyValue {
				fmain.Write([]byte(key+" "+value+"\n"))
			}
			ftemp.Seek(0, os.SEEK_SET)
			ftemp.Write([]byte(CLEAN))
		}
	}
	fmain.Close()
	ftemp.Close()
	if removeRecovery {
		os.Remove("recovery.txt")
	}
}

/*
 * This function is to handle any requests from the server
 */

func handleServer() {
	servRead := bufio.NewReader(os.Stdin)
	text, _ := servRead.ReadString('\n')
	actual := strings.Split(text,"\n")
	command := strings.ToLower(strings.TrimSpace(string(actual[0])))

	if command == "exit" || command == "quit" || command == "abort" {
		syncFiles(KEEP_RECOVERY)
		os.Exit(1)
	}
}

/*
 * Handle requests from client
 */

func handleClient(conn net.Conn) {

	 /* This function will execute at the end. It closes connection on exit */

	defer conn.Close()
	var command, key, value string = "", "", ""
	var err2 error
	var buf [512]byte

	for {
		/* Read upto 512 bytes from client */

		buf = [512]byte{}
		_, err := conn.Read(buf[0:])
		if err != nil {
			return
		}

		/*
		 * We will receive the requests from client in the form of a query. Here we parse
		 * the query and split into the command (set,get,update,delete), key and value.
		 * Reference: Took idea for set/get commands interface from Instagram website
		 * URL: http://instagram-engineering.tumblr.com/post/12202313862/storing-hundreds-of-millions-of-simple-key-value-pairs
		 */

		raw1 := strings.Split(string(buf[0:]), "\n")
		query := strings.Split(raw1[0], " ")
		
		if len(query) > 0 {
			command = strings.ToLower(strings.TrimSpace(query[0]))

			if len(query) > 1 {
				key = strings.TrimSpace(query[1])
				if command == "set" || command == "update" || command == "rename" {
					value = strings.TrimSpace(query[2])
				}
			}
		}

		if command == "update" {
			command = "set"
		}

		if command == "remove" {
			command = "delete"
		}

		switch command {

			case "set":
					_, found := oldKeyValue[key]
					oldKeyValue[key] = value
					if found {
						IS_DIRTY = true
						oldKeyValue[key] = value
						ftemp.Write([]byte("\nSET"+" "+key+" "+value))
						_, err2 = conn.Write([]byte("KEY UPDATED\n"))
					} else {
						newKeyValue[key] = value
						fmain.Write([]byte(key+" "+value+"\n"))
						ftemp.Write([]byte("\nSET"+" "+key+" "+value))
						_, err2 = conn.Write([]byte("KEY ADDED\n"))
					}
					
			case "get":
					val, found := oldKeyValue[key]
					if found {
						_, err2 = conn.Write([]byte(val+"\n"))
					} else {
						val, found := newKeyValue[key]
						if found {
							_, err2 = conn.Write([]byte(val+"\n"))
						} else {
							conn.Write([]byte("KEY NOT PRESENT\n"))
						}
					}

			case "rename":
					find = false
					val, found := oldKeyValue[key]
					if found {
						find = true
						delete(oldKeyValue, key)
						oldKeyValue[value] = val
						ftemp.Write([]byte("\nDEL"+" "+key+" "+value))
						ftemp.Write([]byte("\nSET"+" "+value+" "+val))
					}
					val, found = newKeyValue[key]
					if found {
						if find == false {
							ftemp.Write([]byte("\nDEL"+" "+key+" "+value))
							ftemp.Write([]byte("\nSET"+" "+value+" "+val))
						}
						find = true
						delete(newKeyValue, key)
						newKeyValue[value] = val
					}
					if find == false {
						conn.Write([]byte("KEY NOT PRESENT\n"))
					} else {
						_, err2 = conn.Write([]byte("RENAMED\n"))
					}

			case "delete":
					find = false
					findold := false
					value, found := oldKeyValue[key]
					if found {
						IS_DIRTY = true
						find = true
						findold = true
						delete(oldKeyValue, key)
						ftemp.Write([]byte("\nDEL"+" "+key+" "+value))
						_, err2 = conn.Write([]byte("KEY DELETED\n"))
					}
					value, found = newKeyValue[key]
					if found {
						find = true
						delete(newKeyValue, key)
						ftemp.Write([]byte("\nDEL"+" "+key+" "+value))
						if findold != true {
							_, err2 = conn.Write([]byte("KEY DELETED\n"))
						}
					}
							

					if find == false {
						_, err2 = conn.Write([]byte("KEY NOT PRESENT\n"))
					}
				
			case "list":
					if strings.ToLower(key) == "all" {   
						_, err2 = conn.Write([]byte("KEYS -----> VALUES\n--------------\n"))
						for key1, value1 := range oldKeyValue {
							_, err2 = conn.Write([]byte(key1+" -----> "+value1+"\n"))
						}
					} else if strings.ToLower(key) == "keys" {	
						_, err2 = conn.Write([]byte("ALL KEYS\n--------\n"))
						for key1, _ := range oldKeyValue {
							_, err2 = conn.Write([]byte(key1+"\n"))
						}
					}
					_, err2 = conn.Write([]byte("\n"))

			case "help":
					_, err2 = conn.Write([]byte("Following set of commands\n"))
					_, err2 = conn.Write([]byte("-------------------------\n"))
					_, err2 = conn.Write([]byte("set    <key> <value>\t-> set a new key\n"))
					_, err2 = conn.Write([]byte("get    <key>     \t-> get value of a key\n"))
					_, err2 = conn.Write([]byte("update <key> <value> \t-> update value of a key\n"))
					_, err2 = conn.Write([]byte("rename <key> <newkey> \t-> renames a key\n"))
					_, err2 = conn.Write([]byte("delete <key> \t\t-> delete a key value pair\n"))
					_, err2 = conn.Write([]byte("list   <keys/all>\t-> list all keys/key-value pairs\n"))
					_, err2 = conn.Write([]byte("exit/quit \t\t-> close connection\n"))
					_, err2 = conn.Write([]byte("help \t\t\t-> display this help\n\n"))
					
		}

		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
