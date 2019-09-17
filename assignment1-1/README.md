# COS418 Assignment 1 (Part 1): Intro to Go

<h2>Introduction</h2>
<p>
  In this assignment you will solve two short problems as a way to familiarize
  yourself with the Go programming language. We expect you to already have a
  basic knowledge of the language. If you're starting from nothing, we highly
  recommend going through the <a href="http://tour.golang.org/list">Golang tour</a>
  before you begin this assignment.
</p>
<h2>Software</h2>
<p>
  You will find the code in the same directory as this readme. The two problems that you need to solve are in <tt>q1.go</tt>
  and <tt>q2.go</tt>. You should only add code to places that say <tt>TODO: implement me</tt>. 
  Do not change any of the function signatures as our testing framework uses them.
</p>

<p>
  <b>Q1 - Top K words:</b> The task is to find the <tt>K</tt> most common words in a
  given document. To exclude common words such as "a" and "the", the user of your program
  should be able to specify the minimum character threshold for a word. Word matching is
  case insensitive and punctuation should be removed. You can find more details on what
  qualifies as a word in the comments in the code.
</p>

<p>
  <b>Q2 - Parallel sum:</b> The task is to implement a function that sums a list of
  numbers in a file in parallel. For this problem you are required to use goroutines (the
  <tt>go</tt> keyword) and channels to pass messages across the goroutines. While it is
  possible to just sum all the numbers sequentially, the point of this problem is to
  familiarize yourself with the synchronization mechanisms in Go.
</p>

<h3>Testing</h3>

<p>
  Our grading uses the tests in <tt>q1_test.go</tt> and <tt>q2_test.go</tt> provided to you.
  To test the correctness of your code, run the following:
</p>
<pre>
  $ cd assignment1-1
  $ go test
</pre>
<p>
  If all tests pass, you should see the following output:
</p>
<pre>
  $ go test
  PASS
  ok      /path/to/assignment1-1   0.009s
</pre>

### Point Distribution

<table>
<tr><th>Test</th><th>Points</th></tr>
<tr><td>Q1Simple</td><td>4</td></tr>
<tr><td>Q1DeclarationOfIndependence</td><td>4</td></tr>
<tr><td>Q2_1</td><td>3</td></tr>
<tr><td>Q2_2</td><td>3</td></tr>
<tr><td>Q2_3</td><td>3</td></tr>
<tr><td>Q2_4</td><td>3</td></tr>
</table>

### Submitting Assignment
<p> Now you need to submit your assignment. Commit your change and push it to the remote repository by doing the following: </p>

```bash
$ git commit -am "[you fill me in]"
$ git tag -a -m "i finished assignment 1-1" a11-handin
$ git push origin master
$ git push origin a11-handin
```
<p>Please fill in <a href="https://forms.gle/xQEMTYZALsCdz5657">this google form</a> with your GitHub username and your Princeton NetID. This allows us to map your GitHub to your Princeton information for grading. We'll run the grading scripts daily and send your grading information to your Princeton email.</p>

<p>In order to overwrite a tag use the force flag as follows.</p>

```bash
$ git tag -fam "i finished assignment 1-1" a11-handin
$ git push -f --tags
```
