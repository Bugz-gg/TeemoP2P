#ifndef TOOLS_H
#define TOOLS_H
#define BITS_PER_INT 8*sizeof(int)

struct BufferMap {
    int len;
    unsigned int *bit_sequence;
};

struct file {
    char *name;
    int size;
    int piece_size;
    char key[32];
    struct BufferMap buffer_map;
};

int streq(char *, char*);

#endif //TOOLS_H
