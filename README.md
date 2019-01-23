# GoTorrent
BitTorrent Client Implementation

### Dependencies :
- Install Bencode Parser : ```go get github.com/zeebo/bencode```
- Install golint for formatting styles : ```go get -u golang.org/x/lint/golint```

### Guidelines for contribution :
- Take open issues and ask for assignment in comment section
- After assignment , fork and work on a seperate branch
- Follow a simple commit message guideline eg . ``` Fix <issue_id> : <small description> Author@<your name>```
- Merge and squash the changes in your ```local master``` .
	- ```git checkout master```
	- ```git merge --squash <branch name>```
	- ```git commmit```
- Make sure that ```Travis CI build``` is passed.
- Generate a pull request from your ```local master``` branch to ```remote master``` branch.
- Don't close the issue by your own.

### Resources and References
- [Bittorent Specifications](http://jonas.nitro.dk/bittorrent/bittorrent-rfc.html)
- [Bittorent Specifications](http://www.bittorrent.org/beps/bep_0003.html)
- [Bittorrent in Javascript](https://allenkim67.github.io/programming/2016/05/04/how-to-make-your-own-bittorrent-client.html)
- [Bittorrent in C#](https://www.seanjoflynn.com/research/bittorrent.html)
- [Complete Bittorrent in Go](https://github.com/jackpal/Taipei-Torrent)
- [Network Programming in Go](https://ipfs.io/ipfs/QmfYeDhGH9bZzihBUDEQbCbTc5k5FZKURMUoUvfmc27BwL/index.html)
- [Concurrency in Go](https://github.com/golang/go/wiki/LearnConcurrency)
- [Go advanced testing tips & tricks](https://medium.com/@povilasve/go-advanced-tips-tricks-a872503ac859)
- [Travis CI tutorial](https://docs.travis-ci.com/user/tutorial/)
