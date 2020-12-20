docker build . -f ./build.dockerfile -t go-server-with-python-handlers && docker run -it -p 8080:8080 go-server-with-python-handlers
