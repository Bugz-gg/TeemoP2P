#ifndef STRUCTS_H
#define STRUCTS_H

enum criterias {
    FILENAME, FILESIZE
};
enum operations {
    LT, LE, EQ, GE, GT, DI
};
enum types {
    INT, FLOAT, STR
};

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

typedef struct {
    int nb_keys;
    char **keys;
    int nb_leech;
    char **leech;
    int is_valid;
} updateData;

typedef struct {
    Peer *peers;
    File *files;
    int nb_files;
    int nb_peers;
} Tracker;

#endif //STRUCTS_H
