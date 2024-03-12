#ifndef TOOLS_H
#define TOOLS_H
#include <regex.h>

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

void set_bit(BufferMap, int);
void clear_bit(BufferMap, int);
int is_bit_set(BufferMap, int);

regex_t *announce_regex();
announceData announceCheck(char *);
void printAnnounceData(announceData);
void free_announceData(announceData *);

void free_regex(regex_t *);
void free_file(File *);


#endif //TOOLS_H
