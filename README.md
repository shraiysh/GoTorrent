# GoTorrent
BitTorrent Client Implementation

## Project Status : [![Build Status](https://travis-ci.com/IITH-SBJoshi/concurrency-8.svg?token=PzczDKzHVxyhM8id75xo&branch=master)](https://travis-ci.com/IITH-SBJoshi/concurrency-8)
- Able to parse torrent files.

## Setup
1. **Install Golang**
	- Follow [this](https://golang.org/doc/install) link **OR**
	- Run ```sudo apt-get install golang```
	- Set the environment variables `GOPATH` and `GOBIN`
2. **Run the code**
	Get [dep](https://github.com/golang/dep) for installing the dependencies
	```
	$ cd $GOPATH/src/github.com # Come to the appropriate directory
	$ git clone https://github.com/IITH-SBJoshi/concurrency-8.git
	$ cd concurrency-8/
	$ dep ensure		# Get the dependencies
	$ go run main.go	# Run the code
	```

## Guidelines for contribution :
1. Take open issues and ask for assignment in comment section.
2. **Working on seperate branch**
	- Clone the repository : ```git clone https://github.com/IITH-SBJoshi/concurrency-8.git```
	- Create a issue specific branch in cloned repository : ```git checkout -b issue#<issue number>```
	- Run the code by following the steps above
	- You can now start working on your current branch
3. **Testing the changes**
	- Run the test cases if any: ```go test <test file>.go```
	- Check the linting (Install [golint](https://github.com/golang/lint), if not already installed): ```golint <file_name>```
4. **Commiting the changes**
	- Update ```.gitignore``` if there is any need .
	- To add changes in your working directory : ```git add .```
	- Commit your changes : ```git commit -m "<message>"```
	- Follow a simple commit message guideline eg . ``` Fix <issue_id> : <small description> Author@<your name>```
5. **Pushing the changes**
	- Get current master: `git fetch origin master`
	- Merge master with your branch: `git merge master`
	- Push your changes : ```git push origin <your branch name>:<your branch name>```
	- Make sure that ```Travis CI build``` is passed.
5. **Generating Pull requests :**
	- [Generate a pull request](https://help.github.com/articles/about-pull-requests/) from your ```branch``` branch to ```master``` branch.
	- Give the PR and apt title, and mention `Fixes #<issue_number>` in the comment to link it with the issue.
	- Don't close the issue by your own.
7. **Commenting your Code**
	- Include your comments directly preceding an object for GoDoc to document it.
	- Indent pre-formatted comments.
	- Refer to the [Guidelines](https://blog.golang.org/godoc-documenting-go-code) for more info on commenting.
8. **Documentation**
	- You can refer to [this](http://13.71.92.90/pkg/concurrency-8) to see the documentation generated for the master branch.
	- Documentation will be updated automatically upon a pull request into master.

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
