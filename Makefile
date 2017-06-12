PRNGS := mt19937 splitmix64 xorshift128plus xorshift1024star xoroshiro128plus
PRNGS_LONG := $(addsuffix -long, $(PRNGS))
PRNGS_COVERAGE := $(addsuffix -coverage, $(PRNGS))

DISTRIBUTIONS := uniform normal lognormal
DISTRIBUTIONS_COVERAGE := $(addsuffix -coverage, $(DISTRIBUTIONS))

default: test

test:
	go test -cover -v ./...

$(PRNGS):
	go test -cover -v ./prng/$@

$(PRNGS_LONG):
	go test -cover -v ./prng/$(patsubst %-long,%,$@) -long -timeout 1h


$(PRNGS_COVERAGE):
	$(eval package := $(patsubst %-coverage,%,$@))
	go test -coverprofile=$(package).out -v ./prng/$(package)
	go tool cover -html=$(package).out -o $(package).html
	elinks $(package).html

$(DISTRIBUTIONS):
	go test -cover -v ./distribution/$@

$(DISTRIBUTIONS_COVERAGE):
	$(eval package := $(patsubst %-coverage,%,$@))
	go test -coverprofile=$(package).out -v ./distribution/$(package)
	go tool cover -html=$(package).out -o $(package).html
	elinks $(package).html

clean:
	rm *.out *.html

bench:
	go test -bench=. ./...

