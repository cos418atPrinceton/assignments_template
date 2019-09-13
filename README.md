# Assignments for COS 418

### Environment Setup

Teaching support is only provided for work done using the cycles or courselab servers. The tests are known to work under Go 1.9. Thus later version 1 releases should also work. Learn more about semantic versioning [here](https://semver.org/). Please follow <a href="setup.md">these instructions</a> for requesting access to the servers if you have not already done so.

### Tools
<p>
 There are many commonly used tools in the Go ecosystem. The three most useful starting out are:
 <a href="https://golang.org/cmd/gofmt/">Go fmt</a> and <a href="https://golang.org/cmd/vet/">Go vet</a>, which are built-ins, and <a href="https://github.com/golang/lint">Golint</a>, which is similar to the <tt>splint</tt> tool you used in COS217.
</p>

### Editors
<p>
 For those of you in touch with your systems side (this <em>is</em> Distributed Systems, after all), there are quite a few resources for Go development in both <a href="https://github.com/dominikh/go-mode.el">emacs</a> (additional information available <a href="http://dominik.honnef.co/posts/2013/03/emacs-go-1/">here</a>) and <a href="https://github.com/fatih/vim-go">vim</a> (additional resources <a href="http://farazdagi.com/blog/2015/vim-as-golang-ide/">here</a>).
</p>

<p>
 As many Princeton COS students have become attached to Sublime, here are the two indispensible Sublime packages for Go development: <a href="https://github.com/DisposaBoy/GoSublime">GoSublime</a> and <a href="https://github.com/golang/sublime-build">Sublime-Build</a>. And -- learning from the ancient emacs-vi holy war -- it would be inviting trouble to offer Sublime information without likewise dispensing the must-have Atom plugin: <a href="https://atom.io/packages/go-plus">Go-Plus</a> (walkthrough and additional info <a href="https://rominirani.com/setup-go-development-environment-with-atom-editor-a87a12366fcf#.v49dtbadi">here</a>).
</p>

### Coding Style

<p>All of the code you turn in for this course should have good style.
Make sure that your code has proper indentation, descriptive comments,
and a comment header at the beginning of each file, which includes
your name, userid, and a description of the file.</p>

<p>It is recommended to use the standard tools <tt>gofmt</tt> and <tt>go
vet</tt>. You can also use the <a
href="https://github.com/qiniu/checkstyle">Go Checkstyle</a> tool for
advice on how to improve your code's style. It would also be advisable to
produce code that complies with <a
href="https://github.com/golang/lint">Golint</a> where possible. </p>

<h3>How do I git?</h3>
<p>Please read this <a href="https://git-scm.com/docs/gittutorial">Git Tutorial</a>.</p>

<p>The basic git workflow in the shell (assuming you already have a repo set up):</br>
<ul>
<li>git pull</li>
<li>do some work</li>
<li>git status (shows what has changed)</li>
<li>git add <i>all files you want to commit</i></li>
<li>git commit -m "brief message on your update"</li>
<li>git push</li>
</ul>
</p>

<p> All programming assignments, require Git for submission. <p> We are using Github for distributing and collecting your assignments. At the time of seeing this, you should have already joined the <a href="https://github.com/orgs/COS418F19">COS418F19</a> organization on Github and forked your private repository. Your Github page should have a link. Normally, you only need to clone the repository once, and you will have everything you need for all the assignments in this class.

```bash
$ git clone https://github.com/COS418F19/assignments-myusername.git 418
$ cd 418
$ ls
assignment1-1  assignment1-2  assignment1-3  assignment2  assignment3  assignment4  assignment5  README.md  setup.md
$ 
```

Now, you have everything you need for doing all assignments, i.e., instructions and starter code. Git allows you to keep track of the changes you make to the code. For example, if you want to checkpoint your progress, you can <emph>commit</emph> your changes by running:

```bash
$ git commit -am 'partial solution to assignment 1-1'
$ 
```

You should do this early and often!  You can _push_ your changes to Github after you commit with:

```bash
$ git push origin master
$ 
```

Please let us know that you've gotten this far in the assignment, by pushing a tag to Github.

```bash
$ git tag -a -m "i got git and cloned the assignments" gotgit
$ git push origin gotgit
$
```

As you complete parts of the assignments (and begin future assignments) we'll ask you push tags. You should also be committing and pushing your progress regularly.

### Stepping into Assignment 1-1

Now it's time to go to the [assignment 1-1](assignment1-1) folder to begin your adventure!
