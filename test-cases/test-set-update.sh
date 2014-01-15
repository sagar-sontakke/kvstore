(sleep 1
# set a new key with some value
echo "set country India"
sleep 1

# retrive the value
echo "get country"
sleep 1

# modify the value
echo "update country America"
sleep 1

# retrieve updated value
echo "get country"
sleep 1

# modify value again
echo "update country England"
sleep 1

# retrieve updated value
echo "get country"
sleep 1
) | telnet localhost 1201
