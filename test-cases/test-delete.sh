(sleep 1
# set a new key with some value
echo "set compression tar"
sleep 1

# retrive the value
echo "get compression"
sleep 1

# delete the key-value
echo "delete compression"
sleep 1

# try to retrieve deleted value
echo "get compression"
sleep 1
) | telnet localhost 1201
