#include "tools.h"
#include <stdio.h>
#include <string.h>

int streq(char *str1, char *str2) {
    return !strcmp(str1, str2);
}

// Set the nth bit of the buffer map to 1.
void set_bit(struct BufferMap buffer_map, int n) {
    if (n>=buffer_map.len)
        printf("OUT OF RANGE !\n");
    buffer_map.bit_sequence[n / BITS_PER_INT] |= 1 << (n % BITS_PER_INT);
}

// Set the nth bit of the buffer map to 0.
void clear_bit(int *sequence, int n) {
    buffer_map.bit_sequence[n / BITS_PER_INT] &= ~(1 << (n % BITS_PER_INT));
}

// Check if the nth bit of the buffer map 1.
int is_bit_set(int *sequence, int n) {
    return (buffer_map.bit_sequence[n / BITS_PER_INT] >> (n % BITS_PER_INT)) & 1;
}