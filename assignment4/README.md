# COS418 Assignment 4: Raft Log Consensus

<h2>Introduction</h2>

<p>
  This is the second in the series of assignments in which you'll build a
  fault-tolerant key/value storage system. You've started off in Assignment 3
  assignment by implementing the leader election features of Raft. In this assignment,
  you will implement Raft's core features: log consensus agreement. From here, Assignment 4
  will be a key/value service that uses this Raft implementation as a foundational module.
</p>

<p>
  While being able to elect a leader is useful, we want to use
  Raft to keep a consistent, replicated log of operations. To do
  so, we need to have the servers accept client operations
  through <tt>Start()</tt>, and insert them into the log. In
  Raft, only the leader is allowed to append to the log, and
  should disseminate new entries to other servers by including
  them in its outgoing <tt>AppendEntries</tt> RPCs.
</p>

<p>
  If this sounds only vaguely familiar (or even if it's crystal clear), you are
  highly encouraged to go back to reread the
  <a href="https://raft.github.io/raft.pdf">extended Raft paper</a>,
  the Raft lecture notes, and the
  <a href="http://thesecretlivesofdata.com/raft/">illustrated Raft guide</a>.
  You should, of course, also review your work from Assignment 3, as this assignment
  directly builds off that.
</p>

<h2>Software</h2>

<p>
  You will continue to use the same <tt>cos418</tt> code bundle from the previous assignments.
  For this assignment, we will focus primarily on the code and tests for the Raft implementation in
  <tt>src/raft</tt> and the simple RPC-like system in <tt>src/labrpc</tt>. It is worth your while to
  read and digest the code in these packages again, including your implementation from Assignment 3.
</p>

<h2>Part I</h2>

<p>
  In this lab you'll implement most of the Raft design
  described in the extended paper, including saving
  persistent state and reading it after a node fails and
  then restarts. You will not implement cluster
  membership changes (Section 6) or log compaction /
  snapshotting (Section 7).
</p>

<p>
  A set of Raft instances talk to each other with
  RPC to maintain replicated logs. Your Raft interface will
  support an indefinite sequence of numbered commands, also
  called log entries. The entries are numbered with <em>index numbers</em>.
  The log entry with a given index will eventually
  be committed. At that point, your Raft should send the log
  entry to the larger service for it to execute.
</p>

<p>
  Your first major task is to implement the leader and follower code
  to append new log entries.
  This will involve implementing <tt>Start()</tt>, completing the
  <tt>AppendEntries</tt> RPC structs, sending them, and fleshing
  out the <tt>AppendEntry</tt> RPC handler. Your goal should
  first be to pass the <tt>TestBasicAgree()</tt> test (in
  <tt>test_test.go</tt>). Once you have that working, you should
  try to get all the tests before the "basic persistence" test to
  pass before moving on.
</p>

<p class="note">
  Only RPC may be used for interaction between different Raft
  instances. For example, different instances of your Raft
  implementation are not allowed to share Go variables.
  Your implementation should not use files at all.
</p>


<h2>Part II</h2>
<p>
  The next major task is to handle the fault tolerant aspects of the Raft protocol,
  making your implementation robust against various kinds of failures. These failures
  could include servers not receiving some RPCs and servers that crash and restart.
</p>

<p>
  A Raft-based server must be able to pick up where it left off,
  and continue if the computer it is running on reboots. This requires
  that Raft keep persistent state that survives a reboot (the
  paper's Figure 2 mentions which state should be persistent).
</p>

<p>
  A &ldquo;real&rdquo; implementation would do this by writing
  Raft's persistent state to disk each time it changes, and reading the latest saved
  state from
  disk when restarting after a reboot. Your implementation won't use
  the disk; instead, it will save and restore persistent state
  from a <tt>Persister</tt> object (see <tt>persister.go</tt>).
  Whoever calls <tt>Make()</tt> supplies a <tt>Persister</tt>
  that initially holds Raft's most recently persisted state (if
  any). Raft should initialize its state from that
  <tt>Persister</tt>, and should use it to save its persistent
  state each time the state changes. You can use the
  <tt>ReadRaftState()</tt> and <tt>SaveRaftState()</tt> methods
  for this respectively.
</p>

<p class="todo">
  Implement persistence by first adding code to serialize any
  state that needs persisting in <tt>persist()</tt>, and to
  unserialize that same state in <tt>readPersist()</tt>. You now
  need to determine at what points in the Raft protocol your
  servers are required to persist their state, and insert calls
  to <tt>persist()</tt> in those places. Once this code is
  complete, you should pass the remaining tests.  You may want to
  first try and pass the "basic persistence" test (<tt>go test
    -run 'TestPersist1$'</tt>), and then tackle the remaining ones.
</p>

<p class="note">
  You will need to encode the state as an array of bytes in order
  to pass it to the <tt>Persister</tt>; <tt>raft.go</tt> contains
  some example code for this in <tt>persist()</tt> and
  <tt>readPersist()</tt>.
</p>

<p>
  In order to pass some of the challenging tests towards the end, such as
  those marked "unreliable", you will need to implement the optimization to
  allow a follower to back up the leader's nextIndex by more than one entry
  at a time. See the description in the
  <a href="https://raft.github.io/raft.pdf">extended Raft paper</a> starting at
  the bottom of page 7 and top of page 8 (marked by a gray line).
</p>


<h2>Resources and Advice</h2>

<li>
  Remember that the field names any structures you will
  be sending over RPC (e.g. information about each log entry) must start with capital letters, as
  must the field names in any structure passed inside an RPC.
</li>

<li>
  Similarly to how the RPC system only sends structure
  field names that begin with upper-case letters, and
  silently ignores fields whose names start with
  lower-case letters, the GOB encoder you'll use to save
  persistent state only saves fields whose names start
  with upper case letters. This is a common source of
  mysterious bugs, since Go doesn't warn you.
</li>

<li>
  While the Raft leader is the only server that causes
  entries to be appended to the log, all the servers need
  to independently give newly committed entries to their local service
  replica (via their own <tt>applyCh</tt>). Because of this, you
  should try to keep these two activities as separate as
  possible.
</li>

<li>
  It is possible to figure out the minimum number of messages Raft should
  use when reaching agreement in non-failure cases. You should make your
  implementation use that minimum.
</li>


<li>
  In order to avoid running out of memory, Raft must periodically
  discard old log entries, but you <strong>do not</strong> have
  to worry about garbage collecting the log in this lab.</li>

## Point Distribution

<table>
<tr><th>Test</th><th>Points</th></tr>
<tr><td>BasicAgree</td><td>6</td></tr>
<tr><td>FailAgree</td><td>6</td></tr>
<tr><td>FailNoAgree</td><td>6</td></tr>
<tr><td>ConcurrentStarts</td><td>6</td></tr>
<tr><td>Rejoin</td><td>6</td></tr>
<tr><td>Backup</td><td>7</td></tr>
<tr><td>Count</td><td>7</td></tr>
<tr><td>Persist1</td><td>7</td></tr>
<tr><td>Persist2</td><td>7</td></tr>
<tr><td>Persist3</td><td>7</td></tr>
<tr><td>Figure8$</td><td>7</td></tr>
<tr><td>UnreliableAgree</td><td>7</td></tr>
<tr><td>Figure8Unreliable</td><td>7</td></tr>
<tr><td>ReliableChurn</td><td>7</td></tr>
<tr><td>UnreliableChurn</td><td>7</td></tr>
</table>

## Submitting Assignment

You hand in your assignment as before.

```bash
$ git commit -am "[you fill me in]"
$ git tag -a -m "i finished assignment 4" a4-handin
$ git push origin master a4-handin
```

<p>Recall, in order to overwrite a tag use the force flag as follows.</p>

```bash
$ git tag -fam "i finished assignment 4" a4-handin
$ git push -f --tags
```

You should verify that you are able to see your final commit and tags
on the Github page of your repository for this assignment.


<p>
  You will receive full credit for Part I if your software passes the tests mentioned for that section on the CS servers.
  You will receive full credit for Part II if your software passes the tests mentioned for that section on the CS servers.
</p>

<h2>Acknowledgements</h2>
<p>This assignment is adapted from MIT's 6.824 course. Thanks to Frans Kaashoek, Robert Morris, and Nickolai Zeldovich for their support.</p>
