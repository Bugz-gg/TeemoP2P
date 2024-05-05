#ifndef STRUCTS_H
#define STRUCTS_H
#define MAX_FILE_NAME_SIZE 50
#define MAX_IP_ADDR_SIZE 17
#define ALLOC_FILES 100
#define ALLOC_PEERS 200


enum criterias {
    FILENAME, FILESIZE, KEY
};
enum operations {
    LT, LE, EQ, GE, GT, DI
};
enum types {
    INT, FLOAT, STR
};

typedef struct {
    int num_port;
    int peer_id;
    char addr_ip[MAX_IP_ADDR_SIZE];
} Peer;

typedef struct {
    unsigned long long size;
    unsigned long long pieceSize;
    char key[33];
    int nb_peers;
    int max_peer_ind;
    int alloc_peers;
    Peer **peers;
    char name[MAX_FILE_NAME_SIZE];
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
    Peer **peers;
    File **files;
    int nb_files;
    int max_file_ind;
    int nb_peers;
    int max_peer_ind;
    int alloc_files;
    int alloc_peers;
} Tracker;

#endif //STRUCTS_H
