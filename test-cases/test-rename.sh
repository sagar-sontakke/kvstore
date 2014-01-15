(sleep 1
# set a new key with some value
echo "set news IndiaTimes"
sleep 1

# retrive the value
echo "get news"
sleep 1

# rename the key
echo "rename news newspaper"
sleep 1

# retrieve with old value
echo "get news"
sleep 1

# retrieve with updated value
echo "get newspaper"
sleep 1
) | telnet localhost 1201
