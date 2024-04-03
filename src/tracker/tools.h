#ifndef TOOLS_H
#define TOOLS_H

#include <regex.h>

#define BITS_PER_INT 8*sizeof(int)
#define DELIM " "
#define PORT_MAX_LENGTH 5

enum criterias {
    FILENAME, FILESIZE
};
enum operations {
    LT, LE, EQ, GE, GT, DI
};
enum types {
    INT, FLOAT, STR
};

// Define structures

typedef struct {
    char *addr_ip;
    int num_port;
    int peer_id;
} Peer;

typedef struct {
    char *name;
    int size;
    int pieceSize;
    char key[33];
    int nb_peers;
    Peer *peers;
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

typedef struct {
    char key[33];
    int is_valid;
} getfileData;

int streq(char *, char *);

regex_t *announce_regex();

regex_t *look_regex();

regex_t *comparison_regex();

regex_t *getfile_regex();

announceData announceCheck(char *);

lookData lookCheck(char *);

getfileData getfileCheck(char *);

int peerCmp(Peer, Peer);

int announceStructCmp(announceData, announceData);

int criterionCmp(criterion, criterion);

int lookStructCmp(lookData, lookData);

int getfileStructCmp(getfileData, getfileData);

void printAnnounceData(announceData);

void print_criterion(criterion);

void printLookData(lookData);

void printGetFileData(getfileData);

void free_peer(Peer *);

void free_announceData(announceData *);

void free_regex(regex_t *);

void free_all_regex();

void free_file(File *);

void free_announceData(announceData *);

void free_lookData(lookData *);


#endif //TOOLS_H
