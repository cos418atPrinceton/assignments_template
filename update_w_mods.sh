#!/bin/bash
rm go.mod
cd assignment1-1
go mod init src
go mod tidy
cd ../assignment1-2/src
for x in assignment1-{2,3} assignment{2,3,4,5}
do
	cd ../../"$x"/src
	go mod init src
	go mod edit -go=1.17
	go mod tidy
done
cd ../../

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

for y in "assignment1-2/src/main/ii.go" "assignment1-2/src/main/wc.go"
do 
	replace mapreduce $y
done

for y in "assignment3/src/raft/config.go" "assignment3/src/raft/raft.go" "assignment5/src/raft/config.go" "assignment5/src/kvraft/config.go" "assignment5/src/kvraft/client.go" "assignment5/src/kvraft/server.go"
do
	replace labrpc $y
done

for y in "assignment5/src/kvraft/server.go" "assignment5/src/kvraft/config.go"
do
	replace raft $y
done

for y in assignment1-{2,3} assignment{2,3,4,5}
do
	replace_full " Remember to set your <tt>GOPATH</tt> first." "" "$y/README.md"
	replace_full "\$GOPATH/src/main" "$y/src/main" "$y/README.md"
done

replace_full "cd 418/assignment1-2" "cd 418labsf22-*/assignment1-2/src" assignment1-2/README.md
replace_full "\$ ls"$'\n' "" assignment1-2/README.md
replace_full "README.md src"$'\n' "" assignment1-2/README.md
replace_full "# Go needs \$GOPATH to be set to the directory containing \"src\""$'\n' "" assignment1-2/README.md
replace_full "\$ export GOPATH=\"\$PWD\"" "\$ go test -run Sequential src/mapreduce/..." assignment1-2/README.md
replace_full "\$ cd src" "" assignment1-2/README.md
replace_full "\$ go test -run Sequential mapreduce/..." "" assignment1-2/README.md

replace_full "Remember to set <tt>GOPATH</tt> before continuing. " "" assignment1-2/README.md

replace_full "cd 418" "cd 418labsf22-*" assignment1-3/README.md
replace_full "cp -r assignment1-2/src/* assignment1-3/src/" "cp -r assignment1-2/src/main assignment1-3/src/; cp -r assignment1-2/src/mapreduce assignment1-3/src/" assignment1-3/README.md
replace_full "main      mapreduce" "go.mod      main      mapreduce" assignment1-3/README.md

replace_full "go test -v -run TestBasic mapreduce/..." "go test -v -run TestBasic src/mapreduce/..."  assignment1-3/README.md

replace_full "go test -run Failure mapreduce/..." "go test -run Failure src/mapreduce/..."  assignment1-3/README.md

replace_full "cd 418/assignment3" "cd 418labsf22-*/assignment3/src/raft" assignment3/README.md
replace_full "# Go needs \$GOPATH to be set to the directory containing \"src\""$'\n' "" assignment3/README.md
replace_full "\$ export GOPATH=\"\$PWD\"" "" assignment3/README.md
replace_full "\$ cd \"\$GOPATH/src/raft\"" "" assignment3/README.md


replace_full "  # Go needs \$GOPATH to be set to the directory containing \"src\"" "" assignment5/README.md

replace_full "  \$ cd 418/assignment5" ""  assignment5/README.md

replace_full "  \$ export GOPATH=\"\$PWD\"" "" assignment5/README.md
replace_full "  \$ cd \"\$GOPATH/src/kvraft\"" "" assignment5/README.md
