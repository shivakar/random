PRNGS=mt19937

default: test

test:
	go test -cover -v ./...

$(PRNGS):
	go test -cover -v ./prng/$@

$(PRNGS)-long:
	go test -cover -v ./prng/$(patsubst %-long,%,$@) -long -timeout 1h


$(PRNGS)-coverage:
	$(eval package := $(patsubst %-coverage,%,$@))
	go test -coverprofile=$(package).out -v ./prng/$(package)
	go tool cover -html=$(package).out -o $(package).html
	elinks $(package).html

clean:
	rm *.out *.html

bench:
	go test -bench=. ./...

