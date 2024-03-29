#include "tools.h"
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <regex.h>

unsigned int countDelim(const char *str) { //  Seulement si DELIM ne fait qu'un caractère.
    unsigned int count = 0;
    while (*str != '\0') {
        if (*str == *DELIM) {
            ++count;
        }
        ++str;
    }
    return count + (count > 0);
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

// Check if the nth bit of the buffer map is 1.
int is_bit_set(BufferMap buffer_map, int n) {
    return (buffer_map.bit_sequence[n / BITS_PER_INT] >> (n % BITS_PER_INT)) & 1;
}

regex_t *announce_regex() {
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

regex_t *look_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^look \\[(([a-z]+(<|<=|!=|=|>|>=)\"[a-zA-Z0-9]*\"| )*)\\]$";//"^look \\[((?:[a-z]+(?:<|<=|!=|=|>|>=)\\\"[a-zA-Z0-9]*\\\"| )*)\\]$";//"^look \\[((?:[a-z]+(?:<|<=|!=|=|>|>=)\\\".*\\\"| ))*\\]$"; //"^look \[((?:[a-z]+(?:<|<=|!=|=|>|>=)\"[a-zA-Z0-9]*\"| )*)\]$"
    int ret = regcomp(regex, pattern, REG_EXTENDED); //^look \[((?:[a-z]+(?:<|<=|!=|=|>|>=)\"[a-zA-Z0-9]*\"| )*)\]$
    if (ret) {
        char error_message[100];
        regerror(ret, regex, error_message, sizeof(error_message));
        fprintf(stderr, "Failed to compile look regular expression. %s\n", error_message);
    }
    return regex;
}

regex_t *comparison_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^([a-z]+)((<|<=|!=|=|>|>=))\"([a-zA-Z0-9]*)\"$";
    if (regcomp(regex, pattern, REG_EXTENDED)) {
        fprintf(stderr, "Failed to compile regular expression\n");
    }
    return regex;
}

regex_t *getfile_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^getfile ([a-zA-Z0-9]{32})$";
    int ret = regcomp(regex, pattern, REG_EXTENDED);
    if (ret) {
        char error_message[100];
        regerror(ret, regex, error_message, sizeof(error_message));
        fprintf(stderr, "Failed to compile look regular expression. %s\n", error_message);
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
    for (int i = 0; i < data->nb_files; ++i)
        free_file(&data->files[i]);
    free(data->files);
    for (int i = 0; i < data->nb_leech_keys; ++i)
        free(data->leechKeys[i]);
}

void free_lookData(lookData *data) {
    for (int i = 0; i < data->nb_criterions; ++i) {
        if (data->criterions[i].value_type == STR)
            free(data->criterions[i].value.value_str);
    }
    free(data->criterions);
}

// Function to check announce message
announceData announceCheck(char *message) { // TODO : Valid announceStruct if error
    announceData announceStruct;
    announceStruct.files = NULL;

    regex_t *regex = announce_regex();
    regmatch_t matches[3];
    if (regexec(regex, message, 3, matches, 0)) {
        fprintf(stderr, "Failed to match regular expression\n");
        return announceStruct;
    }

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
            BufferMap tmp = {(files[index / 4].size - 1) / files[index / 4].pieceSize / BITS_PER_INT + 1};
            files[index / 4].buffer_map = tmp;
        }
        token = strtok(NULL, DELIM);
        ++index;
    }

    announceStruct.port = port;
    announceStruct.nb_files = nbFiles;
    announceStruct.files = files;
    announceStruct.nb_leech_keys = nb_leech_keys;
    announceStruct.leechKeys = malloc(nb_leech_keys * 33 * sizeof(char));

    free(filesData);

    return announceStruct;
}

// Function to check announce message
lookData lookCheck(char *message) {
    lookData lookStruct;
    lookStruct.criterions = NULL;

    regex_t *regex = look_regex();
    regmatch_t matches[3];
    if (regexec(regex, message, 3, matches, 0)) {
        fprintf(stderr, "Failed to match regular expression.\n");
        return lookStruct;
    }

    char *criterions_str = strndup(message + matches[1].rm_so, matches[1].rm_eo - matches[1].rm_so);
    int count = countDelim(criterions_str);
    printf("%s\n", criterions_str);
    count = count + (!count && matches[1].rm_eo - matches[1].rm_so);
    if (!count) {
        fprintf(stderr, "No criteria found.\n");
        return lookStruct;
    }
    regex_t *comp_regex = comparison_regex();
    regmatch_t comparison_match[4];


    criterion *criterions = malloc(count * sizeof(criterion));
    char *token = strtok(criterions_str, DELIM);

    char *endptr;
    int int_val;
    float float_val;
    int index = 0;
    int len_crit;

    char criteria[25];

    while (token != NULL) {
        if (regexec(comp_regex, token, 4, comparison_match, 0)) {
            fprintf(stderr, "Failed to match regular expression.\n");
            return lookStruct;
        }

        len_crit = comparison_match[1].rm_eo - comparison_match[1].rm_so;
        strncpy(criteria, token + comparison_match[1].rm_so, len_crit);
        criteria[len_crit] = '\0';

        if (streq(criteria, "filename")) { // Ajouter critères si besoin.
            criterions[index].criteria = FILENAME;
        } else if (streq(criteria, "filesize")) {
            criterions[index].criteria = FILESIZE;
        } else {
            fprintf(stderr, "Incorrect criteria.\n");
            return lookStruct;
        }

        len_crit = comparison_match[2].rm_eo - comparison_match[2].rm_so;
        strncpy(criteria, token + comparison_match[2].rm_so, len_crit);
        criteria[len_crit] = '\0';

        if (streq(criteria, "<")) {
            criterions[index].op = LT;
        } else if (streq(criteria, "<=")) {
            criterions[index].op = LE;
        } else if (streq(criteria, "=")) {
            criterions[index].op = EQ;
        } else if (streq(criteria, ">=")) {
            criterions[index].op = GE;
        } else if (streq(criteria, "<")) {
            criterions[index].op = GT;
        } else if (streq(criteria, "!=")) {
            criterions[index].op = DI;
        } else {
            fprintf(stderr, "Incorrect operator.\n");
            return lookStruct;
        }

        int_val = strtol(token + comparison_match[3].rm_so, &endptr, 10);
        if (*endptr == '\0') {
            criterions[index].value_type = INT;
            criterions[index].value.value_int = int_val;
        } else {
            float_val = strtof(token, &endptr);
            if (*endptr == '\0') {
                criterions[index].value_type = FLOAT;
                criterions[index].value.value_float = float_val;
            } else {
                criterions[index].value_type = STR;
                criterions[index].value.value_str = strndup(token, strlen(token));
            }
        }

        token = strtok(NULL, DELIM);
        ++index;
    }

    lookStruct.nb_criterions = count;
    lookStruct.criterions = criterions;

    free(criterions_str);

    return lookStruct;
}

getfileData getfileCheck(char *message) {
    getfileData getfileStruct;
    getfileStruct.is_valid = 0;

    regex_t *regex = getfile_regex();
    regmatch_t matches[2];
    if (regexec(regex, message, 2, matches, 0)) {
        fprintf(stderr, "Failed to match regular expression.\n");
        return getfileStruct;
    }
    // Check if key in files.
    for (int i = 0; i < 32; ++i)
        getfileStruct.key[i] = *(message + matches[1].rm_so + i);

    getfileStruct.is_valid = 1;
    return getfileStruct;
}

void print_criterion(criterion crit) {
    switch (crit.criteria) {
        case FILENAME:
            printf("filename ");
            break;
        case FILESIZE:
            printf("filesize ");
            break;
        default:
            printf("UNRECOGNISED_CRITERIA ");
    }
    switch (crit.op) {
        case LT:
            printf("lower than ");
            break;
        case LE:
            printf("lower than or equal to ");
            break;
        case EQ:
            printf("equal to ");
            break;
        case GE:
            printf("greater than or equal to");
            break;
        case GT:
            printf("greater than ");
            break;
        case DI:
            printf("different from ");
            break;
        default:
            printf("UNRECOGNISED_OPERATOR ");
    }
    switch (crit.value_type) {
        case INT:
            printf("%d\n", crit.value.value_int);
            break;
        case FLOAT:
            printf("%f\n", crit.value.value_float);
            break;
        case STR:
            printf("%s\n", crit.value.value_str);
            break;
        default:
            printf("UNRECOGNISED_VALUE\n");
    }
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

void printGetFileData(getfileData data) {
    if (data.is_valid) {
        printf("getfile : key : %s\n", data.key);
    } else {
        printf("getfileData is not valid.\n");
    }

}

void printLookData(lookData data) {
    printf("Nb criterions : %d\n", data.nb_criterions);
    for (int i = 0; i < data.nb_criterions; ++i) {
        print_criterion(data.criterions[i]);
    }
}


int main() {
    announceData data = announceCheck(
            "announce listen 2222 seed [teemo 14 5 jzichfnt8SBYA8NS8AZNY8SN9kzox83h]");
    announceData data2 = announceCheck(
            "announce listen 2522 seed [ferIV 120 13 9kOz8SB18SBYA8NS8AZNY8SN9kzox83h interruption 94 3 8jzkhfnt8SBYA8NS8AZNY8SN9kzox83h]");

    lookData data3 = lookCheck("look [filename=\"Enfin\"]");
    lookData data4 = lookCheck("look [filesize=\"9128\" filename=\"Alttab\"]");

    getfileData getfile = getfileCheck("getfile jzicsfnt8SBYA8NS8AZNY8SN9dkzo83h");


    printAnnounceData(data);
    printAnnounceData(data2);
    printLookData(data3);
    printLookData(data4);
    printGetFileData(getfile);

    free_regex(announce_regex());
    free_regex(look_regex());
    free_regex(comparison_regex());
    free_announceData(&data);
    free_announceData(&data2);
    free_lookData(&data3);
    free_lookData(&data4);

    return 0;
}
