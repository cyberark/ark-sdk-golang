all:
	./scripts/build.sh

lint:
	./scripts/golint.sh

clean:
	rm -f ark
	rm -rf bin
