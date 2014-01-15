(sleep 1
# set a new key with some value
echo "set City Mumbai"
sleep 1

# retrive the value
echo "get City"
sleep 1

# modify the value
echo "update City Bangalore"
sleep 1

# retrieve updated value
echo "get City"
sleep 1

# delete the key-value
echo "delete City"
sleep 1

# retrieve updated value
echo "get City"
sleep 1
) | telnet localhost 1201
