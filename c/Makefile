GCC=gcc
FLAGS=-O2 -Wall -DHASU64
#FLAGS=-O2 -Wall

TARGETS=weyl32.o spyspy weyl32

all: $(TARGETS)

spyspy: spyspy.c weyl32.o
	$(GCC) $(FLAGS) -o spyspy spyspy.c weyl32.o

weyl32: weyl32.c weyl32.h
	$(GCC) $(FLAGS) -DDEBUG -o weyl32 weyl32.c

weyl32.o: weyl32.c weyl32.h
	$(GCC) $(FLAGS) -c weyl32.c

clean:
	rm -f $(TARGETS)
