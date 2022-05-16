TARGETS=doggleship gfsalvo bbcsalvo

all: $(TARGETS)

gfsalvo: fortran/spsalvo.for
	cd fortran && make && cp gfsalvo ..

doggleship: dogui.go doggleship.go
	go build -o doggleship dogui.go doggleship.go

bbcsalvo: basic/bbcsalvo.bas
	cd basic && make && cp bbcsalvo.bbc bbcsalvo ..

clean:
	cd fortran && make clean
	cd basic && make clean
	rm -f bbcsalvo.bbc
	rm -f $(TARGETS)
	rm -f s????_[12].txt
