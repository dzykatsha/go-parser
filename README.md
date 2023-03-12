# Go-parser
Implementation of document parser for SPBU course Automated Systems of Data Collection and Processing 

### Requirements
preferably linux-based system or mac OS (haven't tested on Windows)

### Supported 
At this moment, parser supports the following formats:
- `html`
- `doc`
- `docx`
- `pdf`

If you need to retrieve text from `djvu`, check out command-line utility [DjVuLibre](https://djvu.sourceforge.net/features.html)
### How to use?
1. Clone repo and open it
2. Build docker-image: `docker build --tag golang-parser .`
3. Run container: `docker run --rm -it --name go-parser -v $PWD:/go/src/parser golang-parser`
4. `cd src/parser` and `go run client.go --help` to discover command-line arguments

### Manual test
To make sure all is fine, run from terminal `test.sh dir_with_test_input dir_for_test_output`
