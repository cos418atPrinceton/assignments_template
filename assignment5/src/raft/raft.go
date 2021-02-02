package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"labrpc"
	"bytes"
	"encoding/gob"
        "fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	// Interval between each AppendEntries heartbeat sent from the leader to followers
	// Note: Because we synchronize accesses to the functions that send and receive
	// these heartbeats, this interval should not be too low (e.g. 10ms). Otherwise
	// we may end up starving the handlers for the request vote messages and not
	// electing any leaders.
	heartbeatInterval time.Duration = 50 * time.Millisecond

	// Base timeout before a follower or candidate starts a new election
	// This is a base timeout because the real timeout used is randomized based on this value
	baseElectionTimeout time.Duration = 150 * time.Millisecond

	// Timeout for sending RPCs
	rpcTimeout time.Duration = 100 * time.Millisecond
)

//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

type LogEntry struct {
	Term int
	Command interface{}
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu                 sync.Mutex
	peers              []*labrpc.ClientEnd
	persister          *Persister
	me                 int // index into peers[]

	// Persistent state
	currentTerm        int
	votedFor           int
	log                []*LogEntry

	// Volatile state
	isLeader           bool
	commitIndex        int
	lastApplied        int
	timer              *time.Timer
	stopTimer          chan bool
	applyCh            chan ApplyMsg

	// State that exists on the leader only
	nextIndex          []int // each entry is always > 0
	matchIndex         []int

	// Mutex for all of the above state
	lock               sync.Mutex
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (term int, isLeader bool) {
	term = rf.currentTerm
	isLeader = rf.isLeader
	return
}

// return a randomized election timeout for this server
func (rf *Raft) GetElectionTimeout() time.Duration {
	return time.Duration(
		int(float32(baseElectionTimeout.Nanoseconds()) * (1 + rand.Float32())))
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	w := new(bytes.Buffer)
	e := gob.NewEncoder(w)
	e.Encode(rf.currentTerm)
	e.Encode(rf.votedFor)
	e.Encode(rf.log)
	data := w.Bytes()
	rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	d.Decode(&rf.currentTerm)
	d.Decode(&rf.votedFor)
	d.Decode(&rf.log)
}

//
// example RequestVote RPC arguments structure.
//
type RequestVoteArgs struct {
	Term         int
	CandidateId  int
	LastLogIndex int
	LastLogTerm  int
}

//
// example RequestVote RPC reply structure.
//
type RequestVoteReply struct {
	Voter       int // just for debugging
	Term        int
	VoteGranted bool
}

//
// AppendEntries request sent from the leader to backups
//
type AppendEntriesArgs struct {
	Term int
	LeaderId int
	PrevLogIndex int
	PrevLogTerm int
	Entries []*LogEntry
	LeaderCommit int
}

//
// AppendEntries reply sent from backups to the leader
//
type AppendEntriesReply struct {
	Term           int
	Success        bool
	// Index of the latest entry on the follower's log that matches the leader's log
	// This is only > 0 if Success is false
	MatchIndex     int
	// Index of the first entry of the conflicting term, if any
	// This is only > 0 if Success is false
	RequestedIndex int
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args RequestVoteArgs, reply *RequestVoteReply) {
	rf.lock.Lock()
	defer rf.lock.Unlock()
	reply.Voter = rf.me
	reply.Term = rf.currentTerm
	reply.VoteGranted = true
	lastLogTerm := rf.log[len(rf.log) - 1].Term
	DPrintf("Server %v (term %v, votedFor = %v) received vote request from candidate %v " +
		"(term %v)\n", rf.me, rf.currentTerm, rf.votedFor, args.CandidateId, args.Term)
	// It's a new term! Reset some variables...
	if args.Term > rf.currentTerm {
		rf.DiscoveredNewTerm(args.Term)
	} else if args.Term < rf.currentTerm {
		// Deny vote if candidate term is lower than ours
		DPrintf("Server %v rejecting vote from %v because candidate term (%v) " +
			"< our term (%v)\n",
			rf.me, args.CandidateId, args.Term, rf.currentTerm)
		reply.VoteGranted = false
		return
	}
	// Deny vote if we voted for someone else before
	if rf.votedFor >= 0 && rf.votedFor != args.CandidateId {
		DPrintf("Server %v rejecting vote from %v because we voted for someone else %v\n",
			rf.me, args.CandidateId, rf.votedFor)
		reply.VoteGranted = false
		return
	}
	// Deny vote if last entry of candidate log has lower term than that of ours
	if args.LastLogTerm < lastLogTerm {
		DPrintf("Server %v rejecting vote from %v because candidate lastLogTerm " +
			"(%v) < ours (%v)\n",
			rf.me, args.CandidateId, args.LastLogTerm, lastLogTerm)
		reply.VoteGranted = false
		return
	}
	// Deny vote if last entry of candidate log has the same term than that of ours,
	// but candidate log is shorter than ours
	if args.LastLogTerm == lastLogTerm && args.LastLogIndex + 1 < len(rf.log) {
		DPrintf("Server %v rejecting vote from %v because candidate log " +
			"is shorter (%v) than ours (%v)\n",
			rf.me, args.CandidateId, args.LastLogIndex + 1, len(rf.log))
		reply.VoteGranted = false
		return
	}

	// We voted for this candidate, so keep track of it
	DPrintf("Server %v (term %v) voted for candidate %v (term %v)!\n",
		rf.me, rf.currentTerm, args.CandidateId, args.Term)
	rf.votedFor = args.CandidateId
	rf.timer.Reset(rf.GetElectionTimeout())
	rf.persist()
}

//
// AppendEntries RPC handler.
// Only called on the non-leaders.
//
func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {
	rf.lock.Lock()
	defer rf.lock.Unlock()
	reply.Term = rf.currentTerm
	reply.Success = false
	reply.RequestedIndex = -1
	reply.MatchIndex = -1

	if args.Term > rf.currentTerm {
		// It's a new term! Reset some variables...
		rf.DiscoveredNewTerm(args.Term)
	} else if args.Term == rf.currentTerm {
		// If a leader received a heartbeat from another leader of the same term...
		// something went wrong!
		if rf.isLeader {
			log.Fatalf("ERROR: Detected two servers (%v, %v) that claim to be " +
				"leaders in the same term (%v)\n",
				rf.me, args.LeaderId, rf.currentTerm)
		}
	} else {
		// Ignore if it's from an old term
		return
	}

	// Request is from valid leader; reset our election timer
	rf.timer.Reset(rf.GetElectionTimeout())

	// If the term of the previous entry doesn't match, return false
	// Start at index = 1 because the 0th index is a dummy entry
	if args.PrevLogIndex > 0 {
		if args.PrevLogIndex >= len(rf.log) {
			// There are some missing entries on the leader that we're not aware of
			// Return false so the leader will start sending us those entries
			reply.RequestedIndex = len(rf.log)
			return
		}
		prevLogTerm := rf.log[args.PrevLogIndex].Term
		if prevLogTerm != args.PrevLogTerm {
			// As an optimization, return the first index of the conflicting term so
			// the leader can go back a term at a time instead of an entry at a time
			i := args.PrevLogIndex
			for rf.log[i].Term == prevLogTerm {
				i--
			}
			reply.RequestedIndex = i + 1
			return
		}
	}

	// Make our log look like leader's log
	startIndex := args.PrevLogIndex + 1
	for i, entry := range args.Entries {
		logIndex := startIndex + i
		if logIndex == len(rf.log) {
			rf.log = append(rf.log, entry)
		} else if logIndex < len(rf.log) {
			rf.log[logIndex] = entry
		} else {
			// logIndex > len(rf.log); Should never happen
			log.Fatalf("ERROR: received starting index %v that is too " +
				"high for server %v, whose log only has %v entries\n",
				logIndex, rf.me, len(rf.log))
		}
	}

	// Delete extra entries, if any
	endIndex := startIndex + len(args.Entries)
	rf.log = rf.log[:endIndex]
	reply.MatchIndex = len(rf.log) - 1

	// Maybe commit some entries
	if args.LeaderCommit > rf.commitIndex {
		lastLogIndex := len(rf.log) - 1
		if args.LeaderCommit <= lastLogIndex {
			DPrintf("Server %v updating its commit index from %v to %v to match " +
				"leader %v's (prevLogIndex = %v, received %v entries)\n",
				rf.me, rf.commitIndex, args.LeaderCommit, args.LeaderId,
				args.PrevLogIndex, len(args.Entries))
			rf.commitIndex = args.LeaderCommit
		} else {
			// Will this ever happen?
			log.Fatalf("ERROR: Did this just happen on server %v? " +
				"Leader commit = %v, last log index = %v\n",
				rf.me, args.LeaderCommit, lastLogIndex)
		}
		rf.maybeApplyLog()
	}

	// If we added any new entries to our log, make sure we persist them
	if len(args.Entries) > 0 {
		rf.persist()
	}

	// Success!
	reply.Success = true
}

//
// Send a request vote message to the server specified asynchronously.
// Only called on non-leaders.
// MUST BE CALLED WITH rf.lock.
//
func (rf *Raft) sendRequestVote(server int, replyChan chan *RequestVoteReply) {
	lastLogIndex := len(rf.log) - 1
	lastLogTerm := rf.log[lastLogIndex].Term
	args := RequestVoteArgs{rf.currentTerm, rf.me, lastLogIndex, lastLogTerm}
	reply := RequestVoteReply{}
	DPrintf("Server %v (term %v) sending vote request to server %v...\n",
		rf.me, rf.currentTerm, server)
	go func() {
		if rf.sendRPC(server, "Raft.RequestVote", args, &reply) {
			DPrintf("Server %v got a response from server %v " +
				"(term = %v, vote granted = %v)\n",
				rf.me, server, reply.Term, reply.VoteGranted)
			replyChan <- &reply
		} else {
			DPrintf("Server %v never got a response from server %v :(\n", rf.me, server)
		}
	}()
}

//
// Synchronously send an RPC to the remote server, timing out if necessary.
// The caller is expected to pass a *pointer* of the reply struct.
//
func (rf *Raft) sendRPC(server int, name string, args interface{}, reply interface{}) bool {
	replyChan := make(chan bool)
	go func() {
		replyChan <- rf.peers[server].Call(name, args, reply)
	}()
	rpcTimer := time.NewTimer(rpcTimeout)
	select {
	case ok := <-replyChan:
		return ok
	case <-rpcTimer.C:
		return false
	}
}

//
// Send append entries to all followers.
// Only called on the leader.
// MUST BE CALLED WITH rf.lock!
//
func (rf *Raft) broadcastAppendEntries() {
	wg := sync.WaitGroup{}
	for i := 0; i < len(rf.peers); i++ {
		if i == rf.me {
			continue
		}
		wg.Add(1)
		go func(s *Raft, peerIndex int) {
			s.sendAppendEntries(peerIndex)
			wg.Done()
		}(rf, i)
	}
	wg.Wait()
}

//
// Send append entries to a specific follower.
// Return whether we received a response from the follower.
// Only called on the leader.
// MUST BE CALLED WITH rf.lock!
//
func (rf *Raft) sendAppendEntries(server int) {
	// Prepare args
	lastLogIndex := len(rf.log) - 1
	nextIndex := rf.nextIndex[server]
	entries := make([]*LogEntry, 0)
	prevLogIndex := nextIndex - 1
	prevLogTerm := -1
	if prevLogIndex > 0 {
		prevLogTerm = rf.log[prevLogIndex].Term
	}
	if lastLogIndex >= nextIndex {
		entries = rf.log[nextIndex:]
	}
	args := AppendEntriesArgs{
		rf.currentTerm, rf.me, prevLogIndex, prevLogTerm, entries, rf.commitIndex}
	reply := AppendEntriesReply{}
	go func() {
		ok := rf.sendRPC(server, "Raft.AppendEntries", args, &reply)
		if ok {
			rf.handleAppendEntriesReply(server, &reply)
		}
	}()
}

//
// Handler for the response to an AppendEntries request.
// Only called on the leader.
//
func (rf *Raft) handleAppendEntriesReply(server int, reply *AppendEntriesReply) {
	rf.lock.Lock()
	defer rf.lock.Unlock()
	if !rf.isLeader {
		DPrintf("Server %v attempted to handle AppendEntries reply even though " +
			"it is not the leader\n", rf.me)
		return
	}

	// If follower term is higher than ours, then it means we are no longer the leader
	if reply.Term > rf.currentTerm {
		rf.DiscoveredNewTerm(reply.Term)
		return
	}

	// If the logs were inconsistent, lower nextIndex and try again
	if !reply.Success {
		if reply.RequestedIndex == 0 {
			log.Fatalf("ERROR: Follower %v requested a log index of 0 from leader %v\n",
				server, rf.me)
		}
		if reply.RequestedIndex > 0 {
			// Drop our next index for this follower to the requested index
			if reply.RequestedIndex > rf.nextIndex[server] {
				// This could happen if replies come out of order? Not sure...
				log.Printf("WARNING: Follower %v requested a log index of %v, " +
					"which is higher than what leader %v first sent (%v)\n",
					server, reply.RequestedIndex, rf.me, rf.nextIndex[server])
				return
			}
			DPrintf("Follower %v requested leader %v to drop log index from %v to %v\n",
				server, rf.me, rf.nextIndex[server], reply.RequestedIndex)
			rf.nextIndex[server] = reply.RequestedIndex
			rf.sendAppendEntries(server)
		} else {
			log.Printf("WARNING: Follower %v failed the AppendEntries request " +
				"but did not provide a requested index to leader %v\n",
				server, rf.me)
		}
		return
	}

	// Follower log now matches ours, so we update follower state
	followerMatchIndex := reply.MatchIndex
	rf.nextIndex[server] = followerMatchIndex + 1
	rf.matchIndex[server] = followerMatchIndex

	// Maybe commit something
	for i := followerMatchIndex; i > rf.commitIndex; i-- {
		// Never directly commit entries from older terms
		if rf.log[i].Term < rf.currentTerm {
			continue
		}
		// Find the number of servers that contain this log entry
		matches := make([]int, 0)
		matches = append(matches, rf.me) // me
		for serv := 0; serv < len(rf.peers); serv++ {
			if serv == rf.me {
				continue
			}
			if rf.matchIndex[serv] >= i {
				matches = append(matches, serv)
			}
		}
		DPrintf("Leader %v (term %v) is trying to commit index %v: %v (matches = %v)\n",
			rf.me, rf.currentTerm, i, rf.log[i].Command, matches)
		// If the entry exists on a majority of servers, commit it!
		if len(matches) > len(rf.peers) / 2 {
			DPrintf("[COMMITTED] Leader %v (term %v) is committing index %v: %v\n",
				rf.me, rf.currentTerm, i, rf.log[i].Command)
			rf.commitIndex = i
			rf.maybeApplyLog()
			break
		}
	}
}

//
// If we have committed something, apply it to the state machine.
// This is idempotent.
// MUST BE CALLED WITH rf.lock!
//
func (rf *Raft) maybeApplyLog() {
	for rf.commitIndex > rf.lastApplied {
		indexToApply := rf.lastApplied + 1
		cmd := rf.log[indexToApply].Command
		rf.applyCh <- ApplyMsg{Index: indexToApply, Command: cmd}
		rf.lastApplied++
	}
}

//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (index int, term int, isLeader bool) {
	rf.lock.Lock()
	defer rf.lock.Unlock()
	isLeader = rf.isLeader
	if !isLeader {
		index = -1
		term = -1
		return
	}
	// We are the leader: append the entry to our log and it will get propagated to
	// the followers on the next AppendEntries heartbeat
	DPrintf("[NEW COMMAND] Leader %v received command %v, current log size = %v\n",
		rf.me, command, len(rf.log))
	rf.log = append(rf.log, &LogEntry{rf.currentTerm, command})
	rf.matchIndex[rf.me] = len(rf.log) - 1
	rf.persist()
	index = len(rf.log) - 1
	term = rf.currentTerm
	return
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	rf.stopTimer <- true
}

//
// Start a background goroutine that handles the expiration of the timer.
// The behavior of the timer depends on the role of the server.
// If this server is the leader, send append entries to all followers.
// If this server is not the leader, start a new round of election.
//
func (rf *Raft) StartTimerLoop() {
	go func (s *Raft) {
		keepGoing := true
		for keepGoing {
			select {
			case <-s.timer.C:
				rf.lock.Lock()
				if (s.isLeader) {
					s.broadcastAppendEntries()
					s.timer.Reset(heartbeatInterval)
				} else {
					s.StartElection()
					// If we became the leader, start the heartbeat timer
					// Otherwise just reset the election timer
					if (s.isLeader) {
						s.timer.Reset(heartbeatInterval)
					} else {
						s.timer.Reset(s.GetElectionTimeout())
					}
				}
				rf.lock.Unlock()
			case <-s.stopTimer:
				s.timer.Stop()
				keepGoing = false
			}
		}
	}(rf)
}

//
// Start a new round of election by requesting votes from all peers.
// If we get votes from a majority of the servers (including ourselves), become leader.
// MUST BE CALLED WITH rf.lock!
//
func (rf *Raft) StartElection() {
	DPrintf("[NEW ELECTION] Server %v started election for term %v!\n",
		rf.me, rf.currentTerm + 1)
	rf.currentTerm += 1
	rf.votedFor = rf.me
	rf.persist()
	// Send request vote message to everyone and collect replies
	majority := len(rf.peers) / 2 + 1
	replyChan := make(chan *RequestVoteReply, len(rf.peers))
	for i := 0; i < len(rf.peers); i++ {
		if i == rf.me {
			continue
		}
		rf.sendRequestVote(i, replyChan)
	}

	// Release lock while waiting for votes
	electionTerm := rf.currentTerm
	rf.lock.Unlock()

	// Wait until we have a majority of votes, timing out if necessary
	accepts := make([]int, 0)
	rejects := make([]int, 0)
	accepts = append(accepts, rf.me)
	voteTimer := time.NewTimer(rf.GetElectionTimeout())
	for len(accepts) < majority {
		select {
		case reply := <-replyChan:
			rf.lock.Lock()
			if reply.Term > rf.currentTerm {
				// If we discovered a higher term than ours, stop immediately!
				rf.DiscoveredNewTerm(reply.Term)
				return
			}
			if reply.VoteGranted {
				accepts = append(accepts, reply.Voter)
			} else {
				rejects = append(rejects, reply.Voter)
			}
			rf.lock.Unlock()
		case <-voteTimer.C:
			rf.lock.Lock()
			DPrintf("Server %v did not become leader in term %v :( " +
				"(accepts = %v, rejects = %v)\n",
				rf.me, electionTerm, accepts, rejects)
			return
		}
	}

	rf.lock.Lock()

	// Did my term change?
	if rf.currentTerm != electionTerm {
		DPrintf("Server %v term changed (%v -> %v) while running election, aborting " +
			"(accepts = %v, rejects = %v)\n",
			rf.me, electionTerm, rf.currentTerm, accepts, rejects)
		return
	}

	// We got a majority of votes, become leader
	DPrintf("Server %v became leader!! Term is %v! (accepts = %v, rejects = %v)\n",
		rf.me, rf.currentTerm, accepts, rejects)
	rf.isLeader = true
	for i := 0; i < len(rf.peers); i++ {
		rf.nextIndex[i] = len(rf.log) // last log index + 1
		if i == rf.me {
			rf.matchIndex[i] = len(rf.log) - 1
		} else {
			rf.matchIndex[i] = 0
		}
	}
	rf.broadcastAppendEntries()
}

//
// Called when a server discovers a term higher than its own.
// If this is called on the leader, the leader steps down.
// MUST BE CALLED WITH rf.lock!
//
func (rf *Raft) DiscoveredNewTerm(newTerm int) {
	if rf.isLeader {
		DPrintf("Leader %v discovered a higher term %v than its old term %v. " +
			"Stepping down!\n", rf.me, newTerm, rf.currentTerm)
	}
	rf.isLeader = false
	rf.currentTerm = newTerm
	rf.votedFor = -1
	rf.timer.Reset(rf.GetElectionTimeout())
	rf.persist()
}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.currentTerm = 0
	rf.votedFor = -1
	rf.log = make([]*LogEntry, 0)
	// dummy entry to make things 1-indexed
	rf.log = append(rf.log, &LogEntry{0, rf.currentTerm})
	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.timer = time.NewTimer(rf.GetElectionTimeout())
	rf.stopTimer = make(chan bool, 1)
	rf.applyCh = applyCh
	rf.StartTimerLoop()
	rf.nextIndex = make([]int, len(rf.peers))
	rf.matchIndex = make([]int, len(rf.peers))
	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
        fmt.Println("Creating RAFT instance from binary v2")
	return rf
}
