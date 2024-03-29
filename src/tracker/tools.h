#ifndef TOOLS_H
#define TOOLS_H
#include <regex.h>

#define BITS_PER_INT 8*sizeof(int)
#define DELIM " "
#define PORT_MAX_LENGTH 5

enum criterias {FILENAME, FILESIZE};
enum operations {LT, LE, EQ, GE, GT, DI};
enum types {INT, FLOAT, STR};

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
    int peer_id;
} File;

typedef struct {
    int port;
    unsigned int nb_files;
    File *files;
    unsigned int nb_leech_keys;
    char **leechKeys;
    int is_valid;
} announceData;

typedef struct {
    enum types value_type;
    enum criterias criteria;
    enum operations op;
    union {
        int value_int;
        float value_float;
        char *value_str;
    } value;
} criterion;

typedef struct {
    unsigned int nb_criterions;
    criterion *criterions;
    int is_valid;
} lookData;

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
regex_t *look_regex();
regex_t *comparison_regex();
announceData announceCheck(char *);
lookData lookCheck(char *);
void printAnnounceData(announceData);
void print_criterion(criterion);
void printLookData(lookData);
void free_announceData(announceData *);

void free_regex(regex_t *);
void free_file(File *);
void free_announceData(announceData *);
void free_lookData(lookData *);


#endif //TOOLS_H
