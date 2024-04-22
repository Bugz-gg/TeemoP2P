CC = gcc
CFLAGS = -Wall -g
LDFLAGS = -lpthread
OBJ = thpool.o tools.o tracker.o config.o main.o
DEPS = thpool.h tools.h tracker.h config.h
TARGET = server

%.o: %.c $(DEPS)
	$(CC) -c -o $@ $< $(CFLAGS)

$(TARGET): $(OBJ)
	$(CC) -o $@ $^ $(LDFLAGS)

.PHONY: clean

clean:
	rm -f $(OBJ) $(TARGET)
