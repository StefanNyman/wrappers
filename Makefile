
ROOTPKG=github.com/StefanNyman/kubectl

all:
	go build ${ROOTPKG}/cmd/kubectl
	go build ${ROOTPKG}/cmd/helm

install: all
	cp {kubectl,helm} ~/bin/

clean:
	rm -f {kubectl,helm}

