# COS418 Assignment 1 (Part 2): Sequential Map/Reduce

<h2>Introduction</h2>
<p>
  In parts 2 and 3 of the first assignment, you will build a Map/Reduce library
  as a way to learn the Go programming language and as a way to learn
  about fault tolerance in distributed systems. For part 2, you
  will work with a sequential Map/Reduce implementation and write a
  sample program that uses it.

  The interface to the library is similar to the one described in
  the original <a href="http://research.google.com/archive/mapreduce-osdi04.pdf">MapReduce paper</a>.
</p>

<h2>Software</h2>
<p>
  You'll implement this assignment (and all the assignments) in
  <a href="http://www.golang.org/">Go</a>. The Go web site
  contains lots of tutorial information which you may want to
  look at.
</p>

<p>
  For the next two parts of this assignment, we will provide you with a significant amount of scaffolding code to get started.
  The relevant code is under this directory.
  We will ensure that all the code we supply works on the CS servers (cycles/courselab).
</p>


<p>
  In this assignment, we supply you with parts of a flexible
  MapReduce implementation. It has support for two
  modes of operation, <em>sequential</em> and
  <em>distributed</em>. Part 2 deals with the former. The map and reduce tasks
  are all executed in serial: the first map task is
  executed to completion, then the second, then the third, etc.
  When all the map tasks have finished, the first reduce task is
  run, then the second, etc. This mode, while not very fast, can
  be very useful for debugging, since it removes much of the
  noise seen in a parallel execution. The sequential mode also
  simplifies or eliminates various corner cases of a distributed system.
</p>


<h3>Getting familiar with the source</h3>
<p>
  The mapreduce package (located at <tt>src/mapreduce</tt>) provides a simple Map/Reduce library with
  a sequential implementation. Applications would normally call
  <tt>Distributed()</tt> &mdash; located in <tt>mapreduce/master.go</tt> &mdash; to start a job, but may
  instead call <tt>Sequential()</tt> &mdash; also in <tt>mapreduce/master.go</tt> &mdash; to get a
  sequential execution, which will be our approach in this assignment.
</p>

<p>
  The flow of the mapreduce implementation is as follows:
  <ol>
    <li>
      The application provides a number of input files, a map
      function, a reduce function, and the number of reduce
      tasks (<tt>nReduce</tt>).
    </li>
    <li>
      A master is created with this knowledge. It spins up an
      RPC server (see <tt>mapreduce/master_rpc.go</tt>), and waits for
      workers to register (using the RPC call
      <tt>Register()</tt> defined in <tt>mapreduce/master.go</tt>).
      As tasks become available, <tt>schedule()</tt> &mdash; located in <tt>mapreduce/schedule.go</tt> &mdash;
      decides how to assign those tasks to workers, and how to
      handle worker failures.
    </li>
    <li>
      The master considers each input file one map task, and
      makes a call to <tt>doMap()</tt>
      in <tt>mapreduce/common_map.go</tt> at least once for each task. It
      does so either directly (when using
      <tt>Sequential()</tt>) or by issuing the <tt>DoTask</tt>
      RPC &mdash; located in <tt>mapreduce/worker.go</tt> &mdash; on a worker. Each call to
      <tt>doMap()</tt> reads the appropriate file, calls the
      map function on that file's contents, and produces
      <tt>nReduce</tt> files for each map file. Thus, after all map
      tasks are done, the total number of files will be the product of the
      number of files given to map (<tt>nIn</tt>) and <tt>nReduce</tt>.
    


<pre>f0-0, ...,  f0-[nReduce-1],
...
f[nIn-1]-0, ..., f[nIn-1]-[nReduce-1].
</pre>
</li>
    
<li>
  The master next makes a call to <tt>doReduce()</tt>
  in <tt>mapreduce/common_reduce.go</tt> at least once for each
  reduce task.  As with <tt>doMap()</tt>, it does so either
  directly or through a worker. <tt>doReduce()</tt>
  collects corresponding files from each map result
  (e.g. <tt>f0-i, f1-i, ... f[nIn-1]-i</tt>), and runs the reduce function
  on each collection. This process produces <tt>nReduce</tt> result
  files.
</li>
<li>
  The master calls <tt>mr.merge()</tt>
  in <tt>mapreduce/master_splitmerge.go</tt>, which merges all the
  <tt>nReduce</tt> files produced by the previous step
  into a single output.
</li>
<li>
  The master sends a Shutdown RPC to each of its workers,
  and then shuts down its own RPC server.
</li>
  </ol>

  You should look through the files in the MapReduce implementation,
  as reading them might be useful to understand how the other methods
  fit into the overall architecture of the system hierarchy. However,
  for this assignment, you will write/modify <strong>only</strong> <tt>doMap</tt>
  in <tt>mapreduce/common_map.go</tt>,  <tt>doReduce</tt> in <tt>mapreduce/common_reduce.go</tt>,
  and <tt>mapF</tt> and <tt>reduceF</tt> in <tt>main/wc.go</tt>. You will not be able to submit
  other files or modules. In other words, any helper functions must reside within these listed files.
</p>


<h2>Part A: Map/Reduce input and output</h2>
<p>
  The Map/Reduce implementation you are given is missing some
  pieces. Before you can write your first Map/Reduce function
  pair, you will need to fix the sequential implementation. In
  particular, the code we give you is missing two crucial
  pieces: the function that divides up the output of a map task,
  and the function that gathers all the inputs for a reduce task.
  These tasks are carried out by the <tt>doMap()</tt> function in
  <tt>mapreduce/common_map.go</tt>, and the <tt>doReduce()</tt> function in
  <tt>mapreduce/common_reduce.go</tt> respectively. The comments in those
  files should point you in the right direction.
</p>

<p>
  To help you determine if you have correctly implemented
  <tt>doMap()</tt> and <tt>doReduce()</tt>, we have provided you
  with a Go test suite that checks the correctness of your
  implementation. These tests are implemented in the file
  <tt>test_test.go</tt>. To run the tests for the sequential
  implementation that you have now fixed, follow this (or a non-<tt>bash</tt> equivalent) sequence of shell commands,
  starting from the <tt>418/assignment1-2</tt> directory:

<pre>
$ cd 418/assignment1-2
$ ls
README.md src
# Go needs $GOPATH to be set to the directory containing "src"
$ export GOPATH="$PWD"
$ cd src
$ go test -run Sequential mapreduce/...
ok  mapreduce 4.515s
</pre>

<p>
  If the output did not show <em>ok</em> next to the tests, your
  implementation has a bug in it. To give more verbose output,
  set <tt>debugEnabled = true</tt> in <tt>mapreduce/common.go</tt>, and add
  <tt>-v</tt> to the test command above. You will get much more
  output along the lines of:

<pre>
$ go test -v -run Sequential
=== RUN   TestSequentialSingle
master: Starting Map/Reduce task test
Merge: read mrtmp.test-res-0
master: Map/Reduce task completed
--- PASS: TestSequentialSingle (2.30s)
=== RUN   TestSequentialMany
master: Starting Map/Reduce task test
Merge: read mrtmp.test-res-0
Merge: read mrtmp.test-res-1
Merge: read mrtmp.test-res-2
master: Map/Reduce task completed
--- PASS: TestSequentialMany (2.32s)
PASS
ok  mapreduce4.635s
</pre>


<h2>Part B: Single-worker word count</h3>
<p>
  Now that the map and reduce tasks are connected, we can start
  implementing some interesting Map/Reduce operations. For this
  assignment, we will be implementing word count &mdash; a simple and
  classic Map/Reduce example. Specifically, your task is to
  modify <tt>mapF</tt> and <tt>reduceF</tt> within <tt>main/wc.go</tt>
  so that the application reports the number of occurrences of each word.
  A word is any contiguous sequence of letters, as
  determined by
  <a href="http://golang.org/pkg/unicode/#IsLetter"><tt>unicode.IsLetter</tt></a>.
</p>

<p>
  There are some input files with pathnames of the form <tt>pg-*.txt</tt> in
  the <tt>main</tt> directory, downloaded from <a
                                       href="https://www.gutenberg.org/ebooks/search/%3Fsort_order%3Ddownloads">Project
    Gutenberg</a>.
  Remember to set <tt>GOPATH</tt> before continuing. This is the result when you initially try to compile the code we provide you
  and run it:
<pre>
$ cd "$GOPATH/src/main"
$ go run wc.go master sequential pg-*.txt
# command-line-arguments
./wc.go:14: missing return at end of function
./wc.go:21: missing return at end of function</pre>
</p>

<p>
  The compilation fails because we haven't written a complete map
  function (<tt>mapF()</tt>) nor a complete reduce function
  (<tt>reduceF()</tt>) in <tt>wc.go</tt> yet. Before you start
  coding read Section 2 of the
  <a href="http://research.google.com/archive/mapreduce-osdi04.pdf">MapReduce paper</a>.
  Your <tt>mapF()</tt> and <tt>reduceF()</tt> functions will
  differ a bit from those in the paper's Section 2.1. Your
  <tt>mapF()</tt> will be passed the name of a file, as well as
  that file's contents; it should split it into words, and return
  a Go slice of key/value pairs, of type
  <tt>mapreduce.KeyValue</tt>. Your <tt>reduceF()</tt> will be
  called once for each key, with a slice of all the values
  generated by <tt>mapF()</tt> for that key; it should return a
  single output value.
</p>

<p>
  You can test your solution using:
<pre>
$ cd "$GOPATH/src/main"
$ go run wc.go master sequential pg-*.txt
master: Starting Map/Reduce task wcseq
Merge: read mrtmp.wcseq-res-0
Merge: read mrtmp.wcseq-res-1
Merge: read mrtmp.wcseq-res-2
master: Map/Reduce task completed</pre>

The output will be in the file <tt>mrtmp.wcseq</tt>.
We will test your implementation's correctness with the following command,
which should produce the following top 10 words:
<pre>
$ sort -n -k2 mrtmp.wcseq | tail -10
he: 34077
was: 37044
that: 37495
I: 44502
in: 46092
a: 60558
to: 74357
of: 79727
and: 93990
the: 154024</pre>
(this sample result is also found in main/mr-testout.txt)
</p>

<p>
You can remove the output file and all intermediate files with:
  <pre>$ rm mrtmp.*</pre>
</p>

<p>
  To make testing easy for you, from the <tt>$GOPATH/src/main</tt> directory, run:
  <pre>$ sh ./test-wc.sh</pre>
  and it will report if your solution is correct or not.
</p>

<h2>Resources and Advice</h2>
<ul class="hints">
  <li>
    a good read on what strings are in Go is the
    <a href="http://blog.golang.org/strings">Go Blog on strings</a>.
  </li>
  <li>
    you can use
    <a href="http://golang.org/pkg/strings/#FieldsFunc"><tt>strings.FieldsFunc</tt></a>
    to split a string into components.
  </li>
  <li>
    the <a href="http://golang.org/pkg/strconv/">strconv package</a>
    is handy to convert strings to integers etc.
  </li>
</ul>

## Point Distribution

<table>
<tr><th>Test</th><th>Points</th></tr>
<tr><td>SequentialSingle</td><td>10</td></tr>
<tr><td>SequentialMany</td><td>10</td></tr>
<tr><td>test-wc.sh</td><td>10</td></tr>
</table>

## Submitting Assignment

You hand in your assignment exactly as you've been letting us know
your progress:

```bash
$ git commit -am "[you fill me in]"
$ git tag -a -m "i finished assignment 1-2" a12-handin
$ git push origin master a12-handin
```

You should verify that you are able to see your final commit and your
a12-handin tag on the Github page in your repository for this
assignment.

<p>Recall, in order to overwrite a tag use the force flag as follows.</p>

```bash
$ git tag -fam "i finished assignment 1-2" a12-handin
$ git push -f --tags
```


We will use the timestamp of your **last** tag for the
purpose of calculating late days, and we will only grade that version of the
code. (We'll also know if you backdate the tag, don't do that.)

<p>
  Our test script is not a substitute for doing your own testing on the full go test cases detailed above.
  Before submitting, please run the full tests given above for both parts one final time. <b>You</b> are responsible for making sure your code works.
</p>

<p>
  You will receive full credit for Part A if your software passes
  the Sequential tests in <tt>test_test.go</tt> and <tt>test-wc.sh</tt> (as run by the <tt>go test</tt> commands above) on the CS servers.
  You will receive full credit for Part B if your Map/Reduce word count
  output matches the correct output for the sequential
  execution above when run on the CS servers.
</p>

<h2>Acknowledgements</h2>
<p>This assignment is adapted from MIT's 6.824 course. Thanks to Frans Kaashoek, Robert Morris, and Nickolai Zeldovich for their support.</p>
