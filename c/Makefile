# Linux
GCC=gcc
#FLAGS=-O2 -Wall
FLAGS=-O2 -Wall -DHASU64

# Historical Unix R7 on the PDP-11
#GCC=cc
#FLAGS=-O2 -DKANDR -DASCII

TARGETS=spyspy weyl32

all: $(TARGETS)

spyspy: spyspy.c weyl32.o
	$(GCC) $(FLAGS) -o spyspy spyspy.c weyl32.o

weyl32gold: weyl32.c weyl32.h
	$(GCC) -O2 -Wall -DHASU64 -DDEBUG -o weyl32gold weyl32.c

weyl32: weyl32.c weyl32.h
	$(GCC) $(FLAGS) -DDEBUG -o weyl32 weyl32.c

weyl32.o: weyl32.c weyl32.h
	$(GCC) $(FLAGS) -c weyl32.c

clean:
	rm -f $(TARGETS) *.o
