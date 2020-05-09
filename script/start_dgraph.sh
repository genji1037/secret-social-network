# run zero node
nohup dgraph zero --my=localhost:5080 >> zero.log &
# run alpha node
nohup dgraph alpha --zero=localhost:5080 --whitelist 116.235.180.78 -o=2 >> alpha.log &
# run dgraph UI
nohup dgraph-ratel >> ratel.log &