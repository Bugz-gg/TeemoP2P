CC=gcc
CFLAGS+= -std=c11 -Wall -Werror

all : main

main: main.c
	$(CC) $(CFLAGS) $^ -o $@
clean:
	rm -f *~ *.o main

