#ifndef TOOLS_H
#define TOOLS_H

struct BufferMap {
    int len;
    unsigned int *bit_sequence;
};

struct file {
    char *name;
    int size;
    char key[32];
    struct BufferMap buffer_map;
};

#endif //TOOLS_H
