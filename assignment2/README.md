# COS418 Assignment 2: Chandy-Lamport Distributed Snapshots

<h2>Introduction</h2>
<p>
  In this assignment you will implement the
  <a href="https://lamport.azurewebsites.net/pubs/chandy.pdf">Chandy-Lamport algorithm</a> for distributed snapshots.
  Your snapshot algorithm will be implemented on top of a token passing system, similar
  to the ones presented in precept and in
  the Chandy-Lamport paper.

  The algorithm makes the following assumptions:
  <ul>
    <li>There are no failures and all messages arrive intact and only once</li>
    <li>The communication channels are unidirectional and FIFO ordered</li>
    <li>There is a communication path between any two processes in the system</li>
    <li>Any process may initiate the snapshot algorithm</li>
    <li>The snapshot algorithm does not interfere with the normal execution of the processes</li>
    <li>Each process in the system records its local state and the state of its incoming channels</li>
  </ul>
</p>

<h2>Software</h2>
<p>
  You will find the code under this directory. The code is organized
  as follows:
  <ul>
    <li><tt>server.go</tt>: A process in the distributed algorithm</li>
    <li><tt>simulator.go</tt>: A discrete time simulator</li>
    <li><tt>logger.go</tt>: A logger that records events executed by system (useful for debugging)</li>
    <li><tt>common.go</tt>: Debug flag and common message types used in the server, logger, and simulator</li>
    <li><tt>snapshot_test.go</tt>: Tests you need to pass</li>
    <li><tt>syncmap.go</tt>: Implementation of thread-safe map</li>
    <li><tt>queue.go</tt>: Simple queue interface implemented using lists in Go</li>
    <li><tt>test_common.go</tt>: Helper functions for testing</li>
    <li><tt>test_data/</tt>: Test case inputs and expected results</li>
  </ul>
</p>
<p>
  Of these files, you only need to turn in <tt>server.go</tt> and <tt>simulator.go</tt>. However, some other
  files also contain information that will be useful for your implementation or debugging, such as the <tt>debug</tt>
  flag in <tt>common.go</tt> and the thread-safe map in <tt>syncmap.go</tt>. However, you do not have to use the provided SyncMap if you would prefer to implement its functionality yourself. Your task is to implement the functions
  that say <tt>TODO: IMPLEMENT ME</tt>, adding fields to the surrounding classes if necessary.
</p>

<h3>Testing</h3>

<p>
  Our grading uses the tests in <tt>snapshot_test.go</tt> provided to you. Test cases can be found in
  <tt>test_data/</tt>. To test the correctness of your code, simply run the following command:
</p>
<pre>
  $ cd chandy-lamport/
  $ go test
  Running test '2nodes.top', '2nodes-simple.events'
  Running test '2nodes.top', '2nodes-message.events'
  Running test '3nodes.top', '3nodes-simple.events'
  Running test '3nodes.top', '3nodes-bidirectional-messages.events'
  Running test '8nodes.top', '8nodes-sequential-snapshots.events'
  Running test '8nodes.top', '8nodes-concurrent-snapshots.events'
  Running test '10nodes.top', '10nodes.events'
  PASS
  ok      _/path/to/chandy-lamport 0.012s
</pre>
<p>
  To run individual tests, you can look up the test names in <tt>snapshot_test.go</tt> and run:
</p>
<pre>
  $ go test -run 2Node
  Running test '2nodes.top', '2nodes-simple.events'
  Running test '2nodes.top', '2nodes-message.events'
  PASS
  ok      chandy-lamport  0.006s
</pre>

## Point Distribution

<table>
<tr><th>Test</th><th>Points</th></tr>
<tr><td>2NodesSimple</td><td>13</td></tr>
<tr><td>2NodesSingleMessage</td><td>13</td></tr>
<tr><td>3NodesMultipleMessages</td><td>14</td></tr>
<tr><td>3NodesMultipleBidirectionalMessages</td><td>14</td></tr>
<tr><td>8NodesSequentialSnapshots</td><td>15</td></tr>
<tr><td>8NodesConcurrentSnapshots</td><td>15</td></tr>
<tr><td>10Nodes</td><td>16</td></tr>
</table>

## Submitting Assignment

You hand in your assignment as before.

```bash
$ git commit -am "[you fill me in]"
$ git tag -a -m "i finished assignment 2" a2-handin
$ git push origin master a2-handin
```

<p>Recall, in order to overwrite a tag use the force flag as follows.</p>

```bash
$ git tag -fam "i finished assignment 2" a2-handin
$ git push -f --tags
```

You should verify that you are able to see your final commit and tags
on the Github page of your repository for this assignment.


