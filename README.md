# GoTorrent
BitTorrent Client Implementation

## Project Status : [![Build Status](https://travis-ci.com/IITH-SBJoshi/concurrency-8.svg?token=PzczDKzHVxyhM8id75xo&branch=master)](https://travis-ci.com/IITH-SBJoshi/concurrency-8)
- Able to parse torrent files.

## Setup
1. **Install Golang** :
	- Follow [this](https://golang.org/doc/install) link **OR**
	- Run ```sudo apt-get install golang```

2. **Dependencies** :
	- Install Bencode Parser : ```go get github.com/zeebo/bencode```
	- Install golint for formatting styles : ```go get -u golang.org/x/lint/golint```


## Guidelines for contribution :
1. Take open issues and ask for assignment in comment section.
2. **Working on seperate branch :**
	- Clone the repository : ```git clone <url>```
	- Create a issue specific branch in cloned repository : ```git checkout -b <branch name>```
	- You can now start working on your current branch.
3. **Testing the changes:**
	- Run the test cases if any : ```go test <test file>.go```
	- Check the linting : ```golint <file_name>```
4. **Commiting and pushing the changes :**
	- Update ```.gitignore``` if there is any need .
	- To add changes in your working directory : ```git add .```
	- Commit your changes : ```git commit -m "<message>"```
	- Follow a simple commit message guideline eg . ``` Fix <issue_id> : <small description> Author@<your name>```
	- Push your changes : ```git push origin <branch name>:<branch name>```
	- Make sure that ```Travis CI build``` is passed.
5. **Generating Pull requests :**
	- [Generate a pull request](https://help.github.com/articles/about-pull-requests/) from your ```branch``` branch to ```master``` branch.
	- Don't close the issue by your own.
6. Squash all your commits in a branch before submitting pull request.
7. **Commenting your Code**
	- Include your comments directly preceding an object for GoDoc to document it.
	- Refer to the [Guidelines](https://blog.golang.org/godoc-documenting-go-code) for commenting.


## Project Structure :
1. **Package parser** :
	- ```types.go``` : contain all the data structures used for parsing torrent files.
	- ```parser.go``` : parsing the torrent file
	- ```utis.go``` : utility cyrptographic functions
2. **Package tracker** :
	- ```utils.go``` : tracker utility functions.
3. **test_torrents** :  sample torrent files for testing.

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
