# COS418 Assignment 5: Key-Value Storage Service

<h2>Introduction</h2>

<p>
  In this assignment you will build a fault-tolerant key-value storage
  service using your Raft library from the previous assignments.
  Your key-value service will be structured as a replicated state machine 
  with several key-value servers that coordinate their activities
  through the Raft log. Your key/value service should continue to
  process client requests as long as a majority of the servers
  are alive and can communicate, in spite of other failures or
  network partitions.
</p>

<p>
  Your system will consist of clients and key/value servers,
  where each key/value server also acts as a Raft peer. Clients
  send <tt>Put()</tt>, <tt>Append()</tt>, and <tt>Get()</tt> RPCs
  to key/value servers (called kvraft servers), who then place
  those calls into the Raft log and execute them in order. A
  client can send an RPC to any of the kvraft servers, but if that 
  server is not currently a Raft leader, or if there's a failure, the 
  client should retry by sending to a different server. If the 
  operation is committed to the Raft log (and hence applied to
  the key/value state machine), its result is reported to the
  client. If the operation failed to commit (for example, if the
  leader was replaced), the server reports an error, and the
  client retries with a different server.
</p>

<h2>Software</h2>

<p>
  We have supplied you with skeleton code and tests under this directory. You will need to modify
  <tt>kvraft/client.go</tt>, <tt>kvraft/server.go</tt>, and perhaps <tt>kvraft/common.go</tt>. 
  (Even if you don't modify <tt>common.go</tt>, you should submit it as-provided.) For this assignment
  we give you the option to either use your own implementation of Raft from HW4 or use our solution,
  which we recommend. Our solution will be provided to you as a binary that
  you will import in your project.
</p>

<p>
  To get up and running, execute the following commands, as in the previous assignments, and change into the <tt>src/kvraft</tt> directory:
  <pre>
  # Go needs $GOPATH to be set to the directory containing "src"
  $ cd 418/assignment5
  $ export GOPATH="$PWD"
  $ cd "$GOPATH/src/kvraft"
  </pre>
  To apply the binary, follow these instructions:
  <pre>
  $ tar -vzxf raft-binary*.tgz
  $ rm raft-binary*.tgz
  $ mv raft.go src/raft/raft.go
  $ cd src/kvraft
  $ go test # this should now print "Creating RAFT instance from binary"
  </pre>
</p>

<h2>Part I</h2>
<p>
  The service supports three RPCs: <tt>Put(key, value)</tt>,
  <tt>Append(key, arg)</tt>, and <tt>Get(key)</tt>. It maintains
  a simple database of key/value pairs. <tt>Put()</tt> replaces
  the value for a particular key in the database, <tt>Append(key,
    arg)</tt> appends arg to key's value, and <tt>Get()</tt>
  fetches the current value for a key. An <tt>Append</tt> to
  a non-existant key should act like <tt>Put</tt>.
</p>

<p>
  You will implement the service as a replicated state machine
  consisting of several kvservers. Your kvraft client code
  (<tt>Clerk</tt> in <tt>src/kvraft/client.go</tt>) should try
  different kvservers it knows about until one responds
  positively. As long as a client can contact a kvraft server
  that is a Raft leader in a majority partition, its operations
  should eventually succeed.
</p>

<p>
  Your kvraft servers should not directly communicate; they
  should only interact with each other through the Raft log.
</p>

<p>
  Your first task is to implement a solution that works
  when there are no dropped messages, and no failed
  servers. Note that your service must provide
  <em>sequential consistency</em> to applications that
  use its client interface. That is, completed
  application calls to the <tt>Clerk.Get()</tt>,
  <tt>Clerk.Put()</tt>, and <tt>Clerk.Append()</tt>
  methods in <tt>kvraft/client.go</tt> must appear to
  have affected all kvservers in the same order, and have
  at-most-once semantics. A <tt>Clerk.Get(key)</tt>
  should see the value written by the most recent
  <tt>Clerk.Put(key, &hellip;)</tt> or
  <tt>Clerk.Append(key, &hellip;)</tt> (in the total
  order).
</p>

<p>
  A reasonable plan of attack may be to first fill in the
  <tt>Op</tt> struct in <tt>server.go</tt> with the
  "value" information that kvraft will use Raft to agree
  on (remember that <tt>Op</tt>
  field names must start with capital letters, since they
  will be sent through RPC), and then implement the
  <tt>PutAppend()</tt> and <tt>Get()</tt> handlers in
  <tt>server.go</tt>. The handlers should enter an
  <tt>Op</tt> in the Raft log using <tt>Start()</tt>, and
  should reply to the client when that log entry is committed.
  Note that you <strong>cannot</strong> execute
  an operation until the point at which it is committed in
  the log (i.e., when it arrives on the Raft
  <tt>applyCh</tt>).
</p>

<p>
  After calling <tt>Start()</tt>, your kvraft
  servers will need to wait for Raft to complete
  agreement. Commands that have been agreed upon arrive
  on the <tt>applyCh</tt>. You should think carefully
  about how to arrange your code so that your code will
  keep reading <tt>applyCh</tt>, while
  <tt>PutAppend()</tt> and <tt>Get()</tt> handlers submit
  commands to the Raft log using <tt>Start()</tt>. It is
  easy to achieve deadlock between the kvserver and its
  Raft library.
</p>

<p>
  Your solution needs to handle the case in
  which a leader has called Start() for a client
  RPC, but loses its leadership before the
  request is committed to the log. In this case
  you should arrange for the client to re-send
  the request to other servers until it finds
  the new leader. One way to do this is for the
  server to detect that it has lost leadership,
  by noticing that a different request has
  appeared at the index returned by Start(), or
  that the term reported by Raft.GetState() has
  changed.
  If the ex-leader is partitioned by
  itself, it won't know about new leaders; but
  any client in the same partition won't be able
  to talk to a new leader either, so it's OK in
  this case for the server and client to wait
  indefinitely until the partition heals.  More 
  generally, a kvraft server should not complete 
  a <tt>Get()</tt> RPC if it is not part of a majority.
</p>


<p>
  You have completed Part I when you
  <strong>reliably</strong> pass the first test in the
  test suite: "One client". You may also find that you
  can pass the "concurrent clients" test, depending on
  how sophisticated your implementation is.
  From the <tt>src/kvraft</tt> directory:
  <pre>
$ go test -v -run TestBasic
=== RUN   TestBasic
Test: One client ...
  ... Passed
--- PASS: TestBasic (15.18s)
PASS
ok  kvraft 15.190s</pre>
</p>



<h2>Part II</h2>
<p>
  In the face of unreliable connections and node failures, your
  clients may send RPCs multiple times until it finds a kvraft
  server that replies positively. One consequence of this is that
  you must ensure that each application call to
  <tt>Clerk.Put()</tt> or <tt>Clerk.Append()</tt> must appear in
  that order just once (i.e., write the key/value database just
  once).
</p>

<p>
  Thus, your task in Part II is to cope with duplicate client requests, including
  situations where the client sends a request to a kvraft leader
  in one term, times out waiting for a reply, and re-sends the
  request to a new leader in another term. The client request
  should always execute just once. 
</p>

<p>
  You will need to uniquely identify client operations to
  ensure that they execute just once. You can assume that
  each clerk has only one outstanding <tt>Put</tt>,
  <tt>Get</tt>, or <tt>Append</tt>.
</p>

<p>
  For stability, you must make sure that your scheme for duplicate
  detection frees server memory quickly, for example by
  having the client tell the servers which RPCs it has
  heard a reply for. It's OK to piggyback this
  information on the next client request.
</p>

<p>
  You have completed Part II when you
  <strong>reliably</strong> pass all tests through
  <tt>TestPersistPartitionUnreliable()</tt>.
<pre>
$ go test -v
=== RUN   TestBasic
Test: One client ...
  ... Passed
--- PASS: TestBasic (15.22s)
=== RUN   TestConcurrent
Test: concurrent clients ...
  ... Passed
--- PASS: TestConcurrent (15.83s)
=== RUN   TestUnreliable
Test: unreliable ...
  ... Passed
--- PASS: TestUnreliable (16.68s)
=== RUN   TestUnreliableOneKey
Test: Concurrent Append to same key, unreliable ...
  ... Passed
--- PASS: TestUnreliableOneKey (1.40s)
=== RUN   TestOnePartition
Test: Progress in majority ...
  ... Passed
Test: No progress in minority ...
  ... Passed
Test: Completion after heal ...
  ... Passed
--- PASS: TestOnePartition (2.54s)
=== RUN   TestManyPartitionsOneClient
Test: many partitions ...
  ... Passed
--- PASS: TestManyPartitionsOneClient (24.08s)
=== RUN   TestManyPartitionsManyClients
Test: many partitions, many clients ...
  ... Passed
--- PASS: TestManyPartitionsManyClients (26.12s)
=== RUN   TestPersistOneClient
Test: persistence with one client ...
  ... Passed
--- PASS: TestPersistOneClient (18.68s)
=== RUN   TestPersistConcurrent
Test: persistence with concurrent clients ...
  ... Passed
--- PASS: TestPersistConcurrent (19.34s)
=== RUN   TestPersistConcurrentUnreliable
Test: persistence with concurrent clients, unreliable ...
  ... Passed
--- PASS: TestPersistConcurrentUnreliable (20.37s)
=== RUN   TestPersistPartition
Test: persistence with concurrent clients and repartitioning servers...
  ... Passed
--- PASS: TestPersistPartition (26.91s)
=== RUN   TestPersistPartitionUnreliable
Test: persistence with concurrent clients and repartitioning servers, unreliable...
  ... Passed
--- PASS: TestPersistPartitionUnreliable (26.89s)
PASS
ok  kvraft 214.069s</pre>
</p>

<h2>Resources and Advice</h2>

<p>
  This assignment doesn't require you to write much code, but
  you will most likely spend a substantial amount of time
  thinking and staring at debugging logs to figure out
  why your implementation doesn't work. Debugging will be
  more challenging than in the Raft lab because there are
  more components that work asynchronously of each other.
  Start early!
</p>

<p>
  You should implement the service without worrying about the 
  Raft log's growing without bound. You do not need to implement 
  snapshots (from Section 7 in the paper) to allow garbage collection
  of old log entries.
</p>

<p>
  As noted, a kvraft server should not complete a <tt>Get()</tt>
  RPC if it is not part of a majority (so that it does
  not serve stale data). A simple solution is to enter
  every <tt>Get()</tt> (as well as each <tt>Put()</tt>
  and <tt>Append()</tt>) in the Raft log. You don't have
  to implement the optimization for read-only operations
  that is described in Section 8.
</p>


<p>
  In Part I, you should probably modify your client
  Clerk to remember which server turned out to
  be the leader for the last RPC, and send the
  next RPC to that server first. This will avoid
  wasting time searching for the leader on every
  RPC.
</p>


## Point Distribution

<table>
<tr><th>Test</th><th>Points</th></tr>
<tr><td>TestBasic$</td><td>8</td></tr>
<tr><td>TestConcurrent$</td><td>8</td></tr>
<tr><td>TestUnreliable$</td><td>8</td></tr>
<tr><td>TestUnreliableOneKey$</td><td>8</td></tr>
<tr><td>TestOnePartition$</td><td>8</td></tr>
<tr><td>TestManyPartitionsOneClient$</td><td>8</td></tr>
<tr><td>TestManyPartitionsManyClients$</td><td>8</td></tr>
<tr><td>TestPersistOneClient$</td><td>8</td></tr>
<tr><td>TestPersistConcurrent$</td><td>9</td></tr>
<tr><td>TestPersistConcurrentUnreliable$</td><td>9</td></tr>
<tr><td>TestPersistPartition$</td><td>9</td></tr>
<tr><td>TestPersistPartitionUnreliable$</td><td>9</td></tr>
</table>

## Submitting Assignment

You hand in your assignment as before.

```bash
$ git commit -am "[you fill me in]"
$ git tag -a -m "i finished assignment 5" a5-handin
$ git push origin master a5-handin
```

<p>Recall, in order to overwrite a tag use the force flag as follows.</p>

```bash
$ git tag -fam "i finished assignment 5" a5-handin
$ git push -f --tags
```

<p>
  You will receive full credit for Part I if your software passes the tests mentioned for that section on the CS servers.
  You will receive full credit for Part II if your software passes the tests mentioned for that section on the CS servers.
</p>


<h2>Acknowledgements</h2>
<p>This assignment is adapted from MIT's 6.824 course. Thanks to Frans Kaashoek, Robert Morris, and Nickolai Zeldovich for their support.</p>
