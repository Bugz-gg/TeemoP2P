CC = gcc
CFLAGS = -Wall -Werror -g -fsanitize=address
LDFLAGS = -lpthread
OBJ = thpool.o tools.o tracker.o main.o
DEPS = thpool.h tools.h tracker.h
TARGET = server

test: src/tracker/tools.c tst/tracker/test_tools.c
	$(CC) $(CFLAGS) $^ -I src/tracker && ./a.out

.PHONY: clean

clean:
	rm -f $(OBJ) $(TARGET)
