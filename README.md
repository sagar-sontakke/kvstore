-----------------
Contents
1. Introduction
2. How to use it?
3. How to use test scripts?
4. How does it works?
5. References
-----------------

1. Introduction
	- This is a simple client-server model for key-value storage and retrieval
	- This is designed in the Go programming language.
	- This is a first assignment for the subject CS-733


2. How to use it?

	1. Run the server side program "tcpServer.go" in one terminal. It will start the service.
	2. Open another terminal(client) and connect to the server. Here server can be a localhost or any other
	   computer where the server program is running.
	3. When the connection gets established, you can run the commands for store, retrieve, delete, display
	   the key-value pairs.
	4. For getting the help about the commands, use the command "help"
	5. It supports following set of commands:

		a) set <key> <value> 			(set a new key)
			- The set command can be used to set/add a new key-value pair
			- eg. "set college iitb" will add a key 'college' with value 'iitb'

		b) get    <key>				(get value of a key)
			- This command can be used to retrieve a value for the corresponding key.
			- It will print the value if the corresponding key is present else will return a message.
			- eg. "get college" command will fetch the value 'iitb'

		c) update <key> <value>			(update value of a key)
			- This will update the value for an already existing key.	
			- If the key is not already present, it will add that key-value pair.
			- eg. "update college iitd" command will update the value to 'iitd'

		d) rename <key> <newkey>		(renames a key)
			- It will rename a keyy to a new key name
			- eg. "rename college school" will rename key 'college' to 'school'
				
		e) delete/remove <key>			(delete a key value pair)
			- This will delete the key-value pair.
			- eg. "delete college" will remove the 'college-iitb' pair.

		f) list   <keys/all>			(list all keys/key-value pairs)
			- It will list all the keys or all key-value pairs
			- eg. "list keys" will list all the keys present on the server side
			      "list all" will list all the key-value pairs on the server side

		g) exit/quit/abort			(close the server)
			- Closes the server (only to be used on server side)

		g) help					(display this help)

	6. The key and the value both must be a single word value without any spaces.

3. How to use test scripts?

	- There are different bash test cases provided to test this assignment. These test files are
	  present in the directory "test-cases"
	- To run these tests follow following steps:
		a) Start the server
		b) Start one or more clients
		c) From the test-cases directory run the bash scripts as "bash <scrip-name>"
		d) The tests will check the basic functionality of this model.
	- Run all the test script with a single go
		- follow the steps above till (b)
		- Run the bash script "run-all.sh"
		- It will run all the test scripts in the test-cases directory and notify when the test completed.

4. How does it works?

	1. Mechanism
		- It works on a simple mechanism of taking the key-value pairs from client and storing, 
		  rerieving, displaying it.
		- A text file 'key-value.txt' is used to store the key-value pairs. This is a single main file
		  containing the data.
		- A temporary file 'recovery.txt' is used to recover the data in case of any failure occurs.

	2. In-memory data
		- All the key-value pairs from the key-value file is copied in memory in the map data structure.
		  This data structure will be updated when any changes (insertion/updation/deletion) occurs to
		  the key-value pairs.
		- Two different maps are used for this purpose. One map for previous data and one for the newly
		  added data.
		- If any changes occur (delete,update) to the previous data, a flag "IS_DIRTY" will set to indicate
		  that the data file needs modification
		- If there are no changes in the previous data, the flag IS_DIRTY is not set and only the newly added
		  data is appended to the data file.

	3. Recovery mechanism
		- A separate file 'recovery.txt' is used to help in case of any failure. We log all the operations
		  in the recovery file specifying whether it is an SET or DELETE operation.
		- This dumping of data to recovery file is continuous, meaning after every set, update and delete
		  commands, these operations are logged.
		- We sync this file with the data file at the time when the server starts and at the time the server
		  closes.
		- If the server fails (eg. abnormal termination, power off), the operations logged in the recovery
		  file will be useful to recover the data that was not previously saved.
		- The recovery file contains a flag at the first line to indicate whether the sync with data file
		  is needed. It it is set to DIRTY, it neans that the data is modified and needs to be sync'ed. If
		  is CLEAN then no sync needed.
		- This sync operation, if needed will occur at the start and close of server.

5. References

	a) Go online tutorial for Go language introduction (http://tour.golang.org/#1)
	b) E-book for connection concepts (Network programming with Go)
	c) Online Go portal for package help (http://golang.org/doc/effective_go.html)
	d) Instagram website for key-value idea (http://instagram-engineering.tumblr.com/post/12202313862/storing-hundreds-of-millions-of-simple-key-value-pairs)
