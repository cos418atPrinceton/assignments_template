#!/bin/bash
# Run from assignment5/src/kvraft

cat > ../raft/raft.go <<- EOM
//go:binary-only-package

package raft
EOM

# unset GOPATH
cd ../
go mod init src
go mod edit -go=1.17
go mod tidy

replace_full () {
	python3 -c "#!/usr/bin/env python
import sys
import shutil
import tempfile

tmp=tempfile.mkstemp()

with open(sys.argv[3]) as fd1, open(tmp[1],'w') as fd2:
    mod_to_replace = sys.argv[1]
    replace_to = sys.argv[2]
    for line in fd1:
        line = line.replace(mod_to_replace,replace_to)
        fd2.write(line)

shutil.move(tmp[1],sys.argv[3])" "$1" "$2" "$3"
}

replace () {
	replace_full "\"$1\"" "\"src/$1\"" $2
}

for y in "raft/config.go" "kvraft/config.go" "kvraft/client.go" "kvraft/server.go"
do
	replace labrpc $y
done

for y in "kvraft/server.go" "kvraft/config.go"
do
	replace raft $y
done

cd kvraft