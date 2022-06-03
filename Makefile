PRNGS := mt19937 splitmix64 xorshift128plus xorshift1024star xoroshiro128plus
PRNGS_LONG := $(addsuffix -long, $(PRNGS))
PRNGS_COVERAGE := $(addsuffix -coverage, $(PRNGS))

DISTRIBUTIONS := uniform normal lognormal cauchy
DISTRIBUTIONS_COVERAGE := $(addsuffix -coverage, $(DISTRIBUTIONS))

default: test

test:
	go test -race -cover -covermode=atomic -v ./...

$(PRNGS):
	go test -race -cover -covermode=atomic -v ./prng/$@

$(PRNGS_LONG):
	go test -race -cover -covermode=atomic -v ./prng/$(patsubst %-long,%,$@) -long -timeout 1h


$(PRNGS_COVERAGE):
	$(eval package := $(patsubst %-coverage,%,$@))
	go test -race -covermode=atomic -coverprofile=$(package).out -v ./prng/$(package)
	go tool cover -html=$(package).out -o $(package).html
	open $(package).html

$(DISTRIBUTIONS):
	go test -race -cover -covermode=atomic -v ./distribution/$@

$(DISTRIBUTIONS_COVERAGE):
	$(eval package := $(patsubst %-coverage,%,$@))
	go test -race -covermode=atomic -coverprofile=$(package).out -v ./distribution/$(package)
	go tool cover -html=$(package).out -o $(package).html
	open $(package).html

clean:
	rm *.out *.html

bench:
	go test -bench=. ./...

