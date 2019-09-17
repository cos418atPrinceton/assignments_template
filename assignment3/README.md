# COS418 Assignment 3: Raft Leader Election

<h2>Introduction</h2>

<p>
  This is the first in a series of assignments in which you'll build a
  fault-tolerant key/value storage system. You'll start in this
  assignment by implementing the leader election features of Raft,
  a replicated state machine protocol. In Assignment 3 you will complete
  Raft's log consensus agreement features. You will implement Raft as a
  Go object with associated methods, available to be used as a module in
  a larger service. Once you have completed Raft, the course assignments
  will conclude with such a service: a key/value service built on top of Raft.
</p>

<h2>Raft Overview</h2>
<p>
  The Raft protocol is used to manage replica servers for services
  that must continue operation in the face of failure (e.g.
  server crashes, broken or flaky networks). The challenge is that,
  in the face of these failures, the replicas won't always hold identical data.
  The Raft protocol helps sort out what the correct data is.
</p>

<p>
  Raft's basic approach for this is to implement a replicated state
  machine. Raft organizes client requests into a sequence, called
  the log, and ensures that all the replicas agree on the the
  contents of the log. Each replica executes the client requests
  in the log in the order they appear in the log, applying those
  requests to the service's state. Since all the live replicas
  see the same log contents, they all execute the same requests
  in the same order, and thus continue to have identical service
  state. If a server fails but later recovers, Raft takes care of
  bringing its log up to date. Raft will continue to operate as
  long as at least a majority of the servers are alive and can
  talk to each other. If there is no such majority, Raft will
  make no progress, but will pick up where it left off as soon as
  a majority is alive again.
</p>

<p>
  You should consult the
  <a href="https://raft.github.io/raft.pdf">extended Raft paper</a>
  and the Raft lecture notes. You may also find this
  <a href="http://thesecretlivesofdata.com/raft/">illustrated Raft guide</a>
  useful to get a sense of the high-level workings of Raft. For a
  wider perspective, have a look at Paxos, Chubby, Paxos Made
  Live, Spanner, Zookeeper, Harp, Viewstamped Replication, and
  <a href="http://static.usenix.org/event/nsdi11/tech/full_papers/Bolosky.pdf">Bolosky et al.</a>
</p>

<h2>Software</h2>
<p>
  For this assignment, we will focus primarily on the code and tests for the Raft implementation in
  <tt>src/raft</tt> and the simple RPC-like system in <tt>src/labrpc</tt>. It is worth your while to
  read and digest the code in these packages.
</p>

<p>
  Before you have implemented anything, your raft tests will fail, but this behavior is a sign that you
  have everything properly configured and are ready to begin:
<pre>
# Go needs $GOPATH to be set to the directory containing "src"
$ cd 418/assignment3
$ export GOPATH="$PWD"
$ cd "$GOPATH/src/raft"
$ go test -run Election
Test: initial election ...
--- FAIL: TestInitialElection (5.00s)
config.go:286: expected one leader, got none
Test: election after network failure ...
--- FAIL: TestReElection (5.00s)
config.go:286: expected one leader, got none
FAIL
exit status 1</pre>
</p>

<p>
  You should implement Raft by adding code to
  <tt>raft/raft.go</tt> (only). In that file you'll find a bit of
  skeleton code, plus some examples of how to send and receive
  RPCs, and examples of how to save and restore persistent state.
</p>


<h2>Your Task: Leader Election</h2>

<p>
  You should start by reading the code to determine which
  functions are responsible for conducting Raft leader election, if
  you haven't already.
</p>

<p>
  The natural first task is to fill in the <tt>RequestVoteArgs</tt> and
  <tt>RequestVoteReply</tt> structs, and modify
  <tt>Make()</tt> to create a background goroutine that
  starts an election (by sending out <tt>RequestVote</tt>
  RPCs) when it hasn't heard from another peer for a
  while. For election to work, you will also need to
  implement the <tt>RequestVote()</tt> RPC handler so
  that servers will vote for one another.
</p>

<p>
  To implement heartbeats, you will need to define an
  <tt>AppendEntries</tt> RPC struct (though you will not need
  any real payload yet), and have the leader send
  them out periodically. You will also have to write an
  <tt>AppendEntries</tt> RPC handler method that resets
  the election timeout so that other servers don't step
  forward as leaders when one has already been elected.
</p>

<p>
  Make sure the timers in different Raft peers are not
  synchronized.  In particular, make sure the election
  timeouts don't always fire at the same time, or else
  all peers will vote for themselves and no one will
  become leader.
</p>

<p>
  Your Raft implementation must support the following interface, which
  the tester and (eventually) your key/value server will use.
  You'll find more details in comments in <tt>raft.go</tt>.

<pre>
// create a new Raft server instance:
rf := Make(peers, me, persister, applyCh)

// start agreement on a new log entry:
rf.Start(command interface{}) (index, term, isleader)

// ask a Raft peer for its current term, and whether it thinks it is leader
rf.GetState() (term, isLeader)

// each time a new entry is committed to the log, each Raft peer
// should send an ApplyMsg to the service (or tester).
type ApplyMsg</pre>

<p>
  A service calls <tt>Make(peers,me,&hellip;)</tt> to create a
  Raft peer. The peers argument is an array of established RPC
  connections, one to each Raft peer (including this one). The
  <tt>me</tt> argument is the index of this peer in the peers
  array. <tt>Start(command)</tt> asks Raft to start the processing
  to append the command to the replicated log. <tt>Start()</tt>
  should return immediately, without waiting for for this process
  to complete. The service expects your implementation to send an
  <tt>ApplyMsg</tt> for each new committed log entry to the
  <tt>applyCh</tt> argument to <tt>Make()</tt>.

<p>
  Your Raft peers should exchange RPCs using the labrpc Go
  package that we provide to you. It is modeled after Go's
  <a href="https://golang.org/pkg/net/rpc/">rpc library</a>, but
  internally uses Go channels rather than sockets.
  <tt>raft.go</tt> contains some example code that sends an RPC
  (<tt>sendRequestVote()</tt>) and that handles an incoming RPC
  (<tt>RequestVote()</tt>).
</p>

<p>
  Implementing leader election and heartbeats (empty
  <tt>AppendEntries</tt> calls) should be sufficient for a
  single leader to be elected and -- in the absence of failures -- stay the leader,
  as well as redetermine leadership after failures.
  Once you have this working, you should be
  able to pass the two Election "go test" tests:
<pre>
$ go test -run Election
Test: initial election ...
  ... Passed
Test: election after network failure ...
  ... Passed
PASS
  ok  raft7.008s</pre>
</p>

<h2>Resources and Advice</h2>

<ul class="hints">
  <li>
    Start early. Although the amount of code to implement
    isn't large, getting it to work correctly will be very
    challenging. Both the algorithm and the code is tricky
    and there are many corner cases to consider. When one
    of the tests fails, it may take a bit of puzzling to
    understand in what scenario your solution isn't
    correct, and how to fix your solution.
  </li>

  <li>
    Read and understand the
    <a href="https://raft.github.io/raft.pdf">extended Raft paper</a>
    and the Raft lecture notes before you start. Your
    implementation should follow the paper's description
    closely, since that's what the tests expect. Figure 2 may
    be useful as a pseudocode reference.
  </li>

  <li>
    Add any state you need to keep to the <tt>Raft</tt>
    struct in <tt>raft.go</tt>. Figure 2 in the paper may
    provide a good guideline.
  </li>
</ul>

## Point Distribution

<table>
<tr><th>Test</th><th>Points</th></tr>
<tr><td>InitialElection</td><td>50</td></tr>
<tr><td>ReElection</td><td>50</td></tr>
</table>

## Submitting Assignment

You hand in your assignment as before.

```bash
$ git commit -am "[you fill me in]"
$ git tag -a -m "i finished assignment 3" a3-handin
$ git push origin master a3-handin
```

<p>Recall, in order to overwrite a tag use the force flag as follows.</p>

```bash
$ git tag -fam "i finished assignment 3" a3-handin
$ git push -f --tags
```

You should verify that you are able to see your final commit and tags
on the Github page of your repository for this assignment.

<p>
  Before submitting, please run the full tests given above for both parts one final time. 
  You will receive full credit for the leader election component if your software passes
  the Election tests (as run by the <tt>go test</tt> commands above) on the CS servers.
</p>

<h2>Acknowledgements</h2>
<p>This assignment is adapted from MIT's 6.824 course. Thanks to Frans Kaashoek, Robert Morris, and Nickolai Zeldovich for their support.</p>
