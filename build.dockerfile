FROM python:3.7
ENV PATH=$PATH:/usr/local/go/bin
RUN wget -q https://golang.org/dl/go1.15.6.linux-amd64.tar.gz \
 && tar -C /usr/local -xzf go1.15.6.linux-amd64.tar.gz \
 && go version
WORKDIR /opt/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go test -v -bench=. .
#RUN go build -o /dev/null -gcflags="-m" .
RUN go build -o server -gcflags="-c=16" .
CMD ./server
