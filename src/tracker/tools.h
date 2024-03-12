#ifndef TOOLS_H
#define TOOLS_H
#define BITS_PER_INT 8*sizeof(int)
#define DELIM " "
#define PORT_MAX_LENGTH 5

typedef struct {
    int len;
    unsigned int *bit_sequence;
} BufferMap;

// Define structures
typedef struct {
    char *name;
    int size;
    int pieceSize;
    char *key;
    BufferMap buffer_map;
} File;

typedef struct {
    int port;
    unsigned int nb_files;
    File *files;
    unsigned int nb_leech_keys;
    char **leechKeys;
} announceData;


struct file {
    char *name;
    int size;
    int piece_size;
    char key[32];
    BufferMap buffer_map;
};

int streq(char *, char *);

#endif //TOOLS_H
