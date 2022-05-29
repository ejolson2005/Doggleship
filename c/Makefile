GCC=gcc -O2 -Wall

TARGETS=weyl32.o spyspy

all: $(TARGETS)

spyspy: spyspy.c weyl32.o
	$(GCC) -o spyspy spyspy.c weyl32.o

weyl32: weyl32.c
	$(GCC) -DDEBUG -o weyl32 weyl32.c

weyl32.o: weyl32.c
	$(GCC) -DHASU64 -c weyl32.c

clean:
	rm -f $(TARGETS)
