(sleep 1
# set a new key with some value
echo "set Language Hindi"
sleep 1

# set the value
echo "set Country India"
sleep 1

# set the value
echo "update Color Blue"
sleep 1

# retrieve updated value
echo "set Music Classical"
sleep 1

# List the keys
echo "list keys"
sleep 1

# List the key-values
echo "list all"
sleep 1

# show help
echo "help"
sleep 1
) | telnet localhost 1201
