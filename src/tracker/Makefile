CC = gcc
CFLAGS = -Wall -Werror -g
LDFLAGS = -lpthread
OBJ = thpool.o tools.o tracker.o main.o
DEPS = thpool.h tools.h tracker.h
TARGET = server

%.o: %.c $(DEPS)
	$(CC) -c -o $@ $< $(CFLAGS)

$(TARGET): $(OBJ)
	$(CC) -o $@ $^ $(LDFLAGS)

.PHONY: clean

clean:
	rm -f $(OBJ) $(TARGET)
