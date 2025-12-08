go build -ldflags "-s -w" -o k620 main.go
tar zcvf k620.tar.gz k620
mv k620.tar.gz ~/mozart229/www/
