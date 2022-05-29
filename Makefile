TARGETS=doggleship gfsalvo bbcsalvo cheater spyspy

all: $(TARGETS)

gfsalvo: fortran/spsalvo.for
	cd fortran && make && cp gfsalvo ..

cheater: swift/cheater.swift
	cd swift && make && cp cheater ..

doggleship: dogui.go doggleship.go
	go build -o doggleship dogui.go doggleship.go

bbcsalvo: basic/bbcsalvo.bas
	cd basic && make && cp bbcsalvo.bbc bbcsalvo ..

spyspy: c/spyspy.c
	cd c && make && cp spyspy ..

clean:
	cd fortran && make clean
	cd basic && make clean
	cd swift && make clean
	cd c && make clean
	rm -f bbcsalvo.bbc
	rm -f $(TARGETS)
	rm -f s????_[12].txt
