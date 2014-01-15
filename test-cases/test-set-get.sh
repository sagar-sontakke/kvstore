(sleep 1
# set a new key with some value
echo "set department Computer"
sleep 1

# retrive the value
echo "get department"
sleep 1

# set the value
echo "set location India"
sleep 1

# retrieve updated value
echo "get location"
sleep 1
) | telnet localhost 1201
