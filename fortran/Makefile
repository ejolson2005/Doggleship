all: gfsalvo

clean:
	rm -f gfsalvo gfsalvo.for rnd.o

gfsalvo.for: spsalvo.for
	./sp2gf

gfsalvo: gfsalvo.for rnd.o
	gfortran -O2 -o gfsalvo -std=legacy gfsalvo.for rnd.o

rnd.o: rnd.c
	gcc -c -O2 rnd.c
