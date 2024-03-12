#include "tools.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <regex.h>

unsigned int countDelim(const char *str) { //  Seulement si DELIM ne fait qu'un caractère.
    unsigned int count = 0;
    while (*str != '\0') {
        if (*str == *DELIM) {
            count++;
        }
        str++;
    }
    return count+(count>0);
}

int streq(char *str1, char *str2) {
    return !strcmp(str1, str2);
}

// Set the nth bit of the buffer map to 1.
void set_bit(BufferMap buffer_map, int n) {
    if (n >= buffer_map.len)
        printf("OUT OF RANGE !\n");
    buffer_map.bit_sequence[n / BITS_PER_INT] |= 1 << (n % BITS_PER_INT);
}

// Set the nth bit of the buffer map to 0.
void clear_bit(BufferMap buffer_map, int n) {
    buffer_map.bit_sequence[n / BITS_PER_INT] &= ~(1 << (n % BITS_PER_INT));
}

// Check if the nth bit of the buffer map 1.
int is_bit_set(BufferMap buffer_map, int n) {
    return (buffer_map.bit_sequence[n / BITS_PER_INT] >> (n % BITS_PER_INT)) & 1;
}

regex_t *announce_regex() { //Add free
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^announce listen ([0-9]+) seed \\[([a-zA-Z0-9 ]*)\\]$";
    if (regcomp(regex, pattern, REG_EXTENDED)) {
        fprintf(stderr, "Failed to compile regular expression\n");
    }
    return regex;
}

void free_regex(regex_t *regex) {
    regfree(regex);
    free(regex);
}

void free_file(File *file) {
    free(file->name);
    free(file->key);
}

void free_announceData(announceData *data) {
    for (int i=0; i<data->nb_files; ++i)
        free_file(&data->files[i]);
    free(data->files);
    for (int i=0; i<data->nb_leech_keys; ++i)
        free(data->leechKeys[i]);

}

// Function to check announce message
announceData announceCheck(char *message) {
    announceData announceStruct;
    announceStruct.files = NULL;

    regex_t *regex = announce_regex();
    //printf("%d.\n", regex);
    // Execute regular expression
    regmatch_t matches[3];
    if (regexec(regex, message, 3, matches, 0)) {
        fprintf(stderr, "Failed to match regular expression\n");
        return announceStruct;
    }
    printf("Here.\n");
    //Port
    char port_str[PORT_MAX_LENGTH + 1];
    int len_port = matches[1].rm_eo - matches[1].rm_so;
    strncpy(port_str, message + matches[1].rm_so, len_port);
    port_str[len_port] = '\0';
    int port = atoi(port_str);

    char *filesData = strndup(message + matches[2].rm_so, matches[2].rm_eo - matches[2].rm_so);
    int count = countDelim(filesData);
    int nb_leech_keys = 0; // Add leech key handling
    if (count % 4) {
        fprintf(stderr, "Wrong file data.\n");
        return announceStruct;
    }
    int nbFiles = count / 4;
    File *files = malloc(nbFiles * sizeof(File));
    char *token = strtok(filesData, DELIM);

    int index = 0;
    int remainder;
    while (token != NULL) {
        remainder = index % 4;
        if (!remainder) {
            files[index / 4].name = strndup(token, strlen(token));
        } else if (remainder == 1) {
            files[index / 4].size = atoi(token);
        } else if (remainder == 2) {
            files[index / 4].pieceSize = atoi(token);
        } else {
            if (strlen(token) != 32) {
                fprintf(stderr, "Wrong md5sum hash size for %s.\n", files[index / 4].name);
                return announceStruct;
            }
            files[index / 4].key = strndup(token, 32);
            BufferMap tmp = {(files[index/4].size-1)/files[index/4].pieceSize/BITS_PER_INT+1};
            files[index / 4].buffer_map = tmp;
        }
        //printf("%s\n", token);
        token = strtok(NULL, DELIM);
        ++index;
    }

    announceStruct.port = port;
    announceStruct.nb_files = nbFiles;
    announceStruct.files = files;
    announceStruct.nb_leech_keys = nb_leech_keys;
    announceStruct.leechKeys = malloc(nb_leech_keys * 33*sizeof(char));

    //regfree(regex);
    free(filesData);

    return announceStruct;
}

void printAnnounceData(announceData data) {
    printf("Port: %d\n", data.port);
    for (int i = 0; i < data.nb_files; ++i) {
        printf("File %d: %s, Size: %d, PieceSize: %d, Key: %s\n", i + 1,
               data.files[i].name, data.files[i].size, data.files[i].pieceSize,
               data.files[i].key);
    }
    for (int i = 0; i < data.nb_leech_keys; ++i) {
        printf("Leech key %d: %s\n", i + 1, data.leechKeys[i]);
    }
}

int main() {
    announceData data = announceCheck("announce listen 2222 seed [fe 12 1 duB18SB18SBYA8NS8AZNY8SN9kzox83h teemo 14 5 jzichfnt8SBYA8NS8AZNY8SN9kzox83h]");
    announceData data2 = announceCheck("announce listen 2522 seed [ferIV 120 13 9kOz8SB18SBYA8NS8AZNY8SN9kzox83h interruption 94 3 8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h]");

    printAnnounceData(data);
    printAnnounceData(data2);

    free_regex(announce_regex());
    free_announceData(&data);
    free_announceData(&data2);

    return 0;
}
