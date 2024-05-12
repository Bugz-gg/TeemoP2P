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

void max(int *a, int b) {
    *a = (*a<b) ? b: *a;
}

int streq(const char *str1, const char *str2) {
    return !strcmp(str1, str2);
}

int streqlim(const char *str1, const char *str2, int n) {
    return !strncmp(str1, str2, n);
}

regex_t *announce_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^announce listen ([0-9]+) seed \\[((.+ [0-9]+ [0-9]+ [a-zA-Z0-9]{32}| )*)\\]( leech \\[(([a-zA-Z0-9]{32}| )*)\\])?(\r\n|\n)$";
    if (regcomp(regex, pattern, REG_EXTENDED)) {
        fprintf(stderr, "Failed to compile `announce` regular expression\n");
    }
    return regex;
}

regex_t *look_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^look \\[(([a-z]+(<|<=|!=|=|>|>=)\"[a-zA-Z0-9]*\"| )*)\\](\r\n?|\n)$";
    int ret = regcomp(regex, pattern, REG_EXTENDED);
    if (ret) {
        char error_message[100];
        regerror(ret, regex, error_message, sizeof(error_message));
        fprintf(stderr, "Failed to compile `look` regular expression. %s\n", error_message);
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
        fprintf(stderr, "Failed to compile `comparison` regular expression\n");
    }
    return regex;
}

regex_t *getfile_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^getfile ([a-zA-Z0-9]{32})(\r\n?|\n)$";
    int ret = regcomp(regex, pattern, REG_EXTENDED);
    if (ret) {
        char error_message[100];
        regerror(ret, regex, error_message, sizeof(error_message));
        fprintf(stderr, "Failed to compile `getfile` regular expression. %s\n", error_message);
    }
    return regex;
}

regex_t *update_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    char *pattern = "^update seed \\[(([a-zA-Z0-9]{32}| )*)\\]( leech \\[(([a-zA-Z0-9]{32}| )*)\\])?(\r\n?|\n)$";
    int ret = regcomp(regex, pattern, REG_EXTENDED);
    if (ret) {
        char error_message[100];
        regerror(ret, regex, error_message, sizeof(error_message));
        fprintf(stderr, "Failed to compile `update` regular expression. %s\n", error_message);
    }
    return regex;
}

void free_peer(Peer *peer) {
    //free(peer->addr_ip);
    free(peer);
}

void free_regex(regex_t *regex) {
    regfree(regex);
    free(regex);
}

void free_all_regex() {
    regex_t *announce = announce_regex();
    free_regex(announce);
    regex_t *look = look_regex();
    free_regex(look);
    regex_t *getfile = getfile_regex();
    free_regex(getfile);
    regex_t *comparison = comparison_regex();
    free_regex(comparison);
    regex_t *update = update_regex();
    free_regex(update);
}

void free_file(File *file) {
    //for (int i=0; i<file->nb_peers; ++i)
    //    free(file->peers[i]);
    free(file->peers);
    free(file);
}

void free_announceData(announceData *data) {
    //for (int i = 0; i < data->nb_files; ++i)
    //    free_file(&data->files[i]);
    free(data->files);
    for (int i = 0; i < data->nb_leech_keys; ++i)
        free(data->leechKeys[i]);
    free(data->leechKeys);
}

void free_lookData(lookData *data) {
    for (int i = 0; i < data->nb_criterions; ++i) {
        if (data->criterions[i].value_type == STR)
            free(data->criterions[i].value.value_str);
    }
    free(data->criterions);
}

void free_updateData(updateData *data) {
    for (int i = 0; i < data->nb_keys; ++i)
        free(data->keys[i]);
    free(data->keys);
    for (int i = 0; i < data->nb_leech; ++i)
        free(data->leech[i]);
    free(data->leech);
}

// Function to check announce message
announceData announceCheck(char *message) {
    announceData announceStruct = {.port=0, .nb_files=0, .nb_leech_keys=0, .is_valid=0};

    regex_t *regex = announce_regex();
    regmatch_t matches[6];
    if (regexec(regex, message, 6, matches, 0)) {
        fprintf(stderr, "(announce) Message invalide :\033[0;33m%s\033[39m\n", message);
        return announceStruct;
    }

    //Port
    char port_str[PORT_MAX_LENGTH + 1];
    int len_port = matches[1].rm_eo - matches[1].rm_so;
    strncpy(port_str, message + matches[1].rm_so, len_port);
    port_str[len_port] = '\0';
    unsigned long long int port = strtoull(port_str, NULL, 10);

    char *filesData, *tofreefiles;
    filesData = tofreefiles = strndup(message + matches[2].rm_so, matches[2].rm_eo - matches[2].rm_so);
    int count = countDelim(filesData);

    int nbFiles = count / 4;
    File *files = malloc(nbFiles * sizeof(File));
    char *token = strsep(&filesData, DELIM);

    int index = 0;
    int remainder;
    while (token != NULL && *token != 0) {
        remainder = index % 4;
        if (!remainder) {
            strcpy(files[index / 4].name, token);
            files[index / 4].peers = NULL;
        } else if (remainder == 1) {
            files[index / 4].size = strtoull(token, NULL, 10);
        } else if (remainder == 2) {
            files[index / 4].pieceSize = strtoull(token, NULL, 10);
        } else {
            strcpy(files[index / 4].key, token);
        }
        token = strsep(&filesData, DELIM);
        ++index;
    }

    // leech
    int nb_leech_keys = 0;
    char *leechData, *tofreeleech;
    leechData = tofreeleech = strndup(message + matches[5].rm_so, matches[5].rm_eo - matches[5].rm_so);
    count = countDelim(leechData);
    count += (!count && (matches[5].rm_eo - matches[5].rm_so));
    char **leechKeys = malloc(count * sizeof(char *));
    token = strsep(&leechData, DELIM);
    int leech_index = 0;
    while (token != NULL && *token != 0) {
        leechKeys[leech_index] = malloc(33 * sizeof(char));
        strncpy(leechKeys[leech_index], token, 32);
        leechKeys[leech_index][32] = '\0';
        token = strsep(&leechData, DELIM);
        ++nb_leech_keys;
        ++leech_index;
    }

    announceStruct.port = port;
    announceStruct.nb_files = nbFiles;
    announceStruct.files = files;
    announceStruct.nb_leech_keys = nb_leech_keys;
    announceStruct.leechKeys = leechKeys;
    announceStruct.is_valid = 1;

    free(tofreeleech);
    free(tofreefiles);

    return announceStruct;
}

// Function to check announce message
lookData lookCheck(char *message) {
    lookData lookStruct = {.is_valid=0, .nb_criterions=0, .criterions=NULL};

    regex_t *regex = look_regex();
    regmatch_t matches[3];
    if (regexec(regex, message, 3, matches, 0)) {
        fprintf(stderr, "(look) Message invalide :\033[0;33m%s\033[39m\n", message);
        return lookStruct;
    }

    char *criterions_str, *tofreecrit;
    criterions_str = tofreecrit = strndup(message + matches[1].rm_so, matches[1].rm_eo - matches[1].rm_so);
    int count = countDelim(criterions_str);
    count += (!count && (matches[1].rm_eo - matches[1].rm_so));
    if (!count) {
        fprintf(stderr, "No criteria found in %s.\n", criterions_str);
        free(tofreecrit);
        lookStruct.nb_criterions = 0;
        lookStruct.criterions = NULL;
        lookStruct.is_valid = 1;
        return lookStruct;
    }
    regex_t *comp_regex = comparison_regex();
    regmatch_t comparison_match[5];

    criterion *criterions = malloc(count * sizeof(criterion));
    char *token = strsep(&criterions_str, DELIM);

    char *endptr;
    int int_val;
    float float_val;
    int index = 0;
    int len_crit;
    char value[100]; // Supposing the max length for the criterions' values is 100.

    char criteria[25];

    while (token != NULL && *token != 0) {
        if (regexec(comp_regex, token, 5, comparison_match, 0)) {
            fprintf(stderr, "(criterion) Message invalide :\033[0;33m%s\033[39m\n.\n", token);
            free(criterions);
            free(tofreecrit);
            return lookStruct;
        }

        len_crit = comparison_match[1].rm_eo - comparison_match[1].rm_so;
        strncpy(criteria, token + comparison_match[1].rm_so, len_crit);
        criteria[len_crit] = '\0';

        if (streq(criteria, "filename")) { // Ajouter critères si besoin.
            criterions[index].criteria = FILENAME;
        } else if (streq(criteria, "filesize")) {
            criterions[index].criteria = FILESIZE;
        } else if (streq(criteria, "key")) {
            criterions[index].criteria = KEY;
        } else {
            fprintf(stderr, "Incorrect criteria : %s.\n", criteria);
            free(criterions);
            free(tofreecrit);
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
        } else if (streq(criteria, ">")) {
            criterions[index].op = GT;
        } else if (streq(criteria, "!=")) {
            criterions[index].op = DI;
        } else {
            free(criterions);
            free(tofreecrit);
            fprintf(stderr, "Incorrect operator : %s.\n", criteria);
            return lookStruct;
        }
        strncpy(value, token + comparison_match[4].rm_so, comparison_match[4].rm_eo - comparison_match[4].rm_so);
        value[comparison_match[4].rm_eo - comparison_match[4].rm_so] = '\0';
        int_val = strtol(value, &endptr, 10);
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
                criterions[index].value.value_str = strndup(value,
                                                            comparison_match[4].rm_eo - comparison_match[4].rm_so);
            }
        }

        // Check is the type is coherent with the criteria
        if (criterions[index].criteria == FILENAME && criterions[index].value_type != STR) {
            free(criterions);
            free(tofreecrit);
            fprintf(stderr, "Incorrect value type for filesize : %s.\n", value);
            return lookStruct;
        }

        token = strsep(&criterions_str, DELIM);
        ++index;
    }
    free(tofreecrit);

    lookStruct.nb_criterions = count;
    lookStruct.criterions = criterions;
    lookStruct.is_valid = 1;

    return lookStruct;
}

getfileData getfileCheck(char *message) {
    getfileData getfileStruct = {.is_valid=0};

    regex_t *regex = getfile_regex();
    regmatch_t matches[2];
    if (regexec(regex, message, 2, matches, 0)) {
        fprintf(stderr, "(getfile) Message invalide :\033[0;33m%s\033[39m\n", message);
        return getfileStruct;
    }
    // Check if key in files.
    for (int i = 0; i < 32; ++i)
        getfileStruct.key[i] = *(message + matches[1].rm_so + i);
    getfileStruct.key[32] = '\0';
    getfileStruct.is_valid = 1;
    return getfileStruct;
}

updateData updateCheck(char *message) {
    updateData updateStruct = {.is_valid=0};

    regex_t *regex = update_regex();
    regmatch_t matches[5];
    if (regexec(regex, message, 5, matches, 0)) {
        fprintf(stderr, "(update) Message invalide :\033[0;33m%s\033[39m\n", message);
        return updateStruct;
    }

    int nb_keys = 0;
    char *keyData, *tofreekey;
    keyData = tofreekey = strndup(message + matches[1].rm_so, matches[1].rm_eo - matches[1].rm_so);
    int count = countDelim(keyData);
    count += (!count && matches[1].rm_eo - matches[1].rm_so);
    char **keys = malloc(count * sizeof(char *));
    char *token = strsep(&keyData, DELIM);
    int key_index = 0;
    while (token != NULL && *token != 0) {
        keys[key_index] = malloc(33 * sizeof(char));
        strncpy(keys[key_index], token, 32);
        keys[key_index][32] = '\0';
        token = strsep(&keyData, DELIM);
        ++nb_keys;
        ++key_index;
    }
    free(tofreekey);
    int nb_leech_keys = 0;
    char *leechData, *tofreeleech;
    leechData = tofreeleech = strndup(message + matches[4].rm_so, matches[4].rm_eo - matches[4].rm_so);
    count = countDelim(leechData);
    count += (count==0);
    char **leechKeys = malloc(count * sizeof(char *));
    token = strsep(&leechData, DELIM);
    int leech_index = 0;
    while (token != NULL && *token != 0) {
        leechKeys[leech_index] = malloc(33 * sizeof(char));
        strncpy(leechKeys[leech_index], token, 32);
        leechKeys[leech_index][32] = '\0';
        token = strsep(&leechData, DELIM);
        ++nb_leech_keys;
        ++leech_index;
    }
    free(tofreeleech);
    updateStruct.is_valid = 1;
    updateStruct.nb_keys = nb_keys;
    updateStruct.keys = keys;
    updateStruct.nb_leech = nb_leech_keys;
    updateStruct.leech = leechKeys;
    return updateStruct;
}

int peerCmp(Peer p1, Peer p2) {
    return streq(p1.addr_ip, p2.addr_ip) && p1.num_port == p2.num_port;
}

int fileCmp(File f1, File f2) { // The equality of the peers having the file is not checked.
    return streq(f1.name, f2.name) && f1.size == f2.size && f1.pieceSize == f2.pieceSize && streq(f1.key, f2.key);
}

int announceStructCmp(announceData a1, announceData a2) {
    if (a1.is_valid != a2.is_valid || a1.port != a2.port || a1.nb_files != a2.nb_files ||
        a1.nb_leech_keys != a2.nb_leech_keys)
        return 0;
    if (!a1.is_valid)
        return 1;

    for (int i = 0; i < a1.nb_files; ++i) {
        if (!fileCmp(a1.files[i], a2.files[i]))
            return 0;
    }
    for (int i = 0; i < a1.nb_leech_keys; ++i) {
        if (!streq(a1.leechKeys[i], a2.leechKeys[i]))
            return 0;
    }
    return 1;
}

int criterionCmp(criterion c1, criterion c2) {
    if (c1.value_type != c2.value_type || c1.criteria != c2.criteria || c1.op != c2.op)
        return 0;
    switch (c1.value_type) {
        case INT:
            return c1.value.value_int == c2.value.value_int;
        case FLOAT:
            return c1.value.value_float == c2.value.value_float;
        case STR:
            return streq(c1.value.value_str, c2.value.value_str);
        default:
            return 0;
    }
}

int lookStructCmp(lookData l1, lookData l2) {
    if ((l1.is_valid != l2.is_valid) || l1.nb_criterions != l2.nb_criterions)
        return 0;
    if (!l1.is_valid)
        return 1;
    for (int i = 0; i < l1.nb_criterions; ++i) {
        for (int j = 0; j < l1.nb_criterions; ++j) {
            if (criterionCmp(l1.criterions[i], l2.criterions[j]))
                break;
            if (j == l1.nb_criterions - 1)
                return 0;
        }

    }
    return 1;
}

int getfileStructCmp(getfileData gf1, getfileData gf2) {
    if (gf1.is_valid == gf2.is_valid && !gf1.is_valid)
        return 1;
    return streq(gf1.key, gf2.key);
}

int updateStructCmp(updateData u1, updateData u2) {
    if ((u1.is_valid != u2.is_valid) || u1.nb_keys != u2.nb_keys || u1.nb_leech != u2.nb_leech)
        return 0;
    if (!u1.is_valid)
        return 1;
    for (int i = 0; i < u1.nb_keys; ++i) {
        if (!streq(u1.keys[i], u2.keys[i]))
            return 0;
    }
    for (int i = 0; i < u1.nb_leech; ++i) {
        if (!streq(u1.leech[i], u2.leech[i]))
            return 0;
    }
    return 1;
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
            printf("greater than or equal to ");
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
            printf("%d (INT)\n", crit.value.value_int);
            break;
        case FLOAT:
            printf("%f (FLOAT)\n", crit.value.value_float);
            break;
        case STR:
            printf("%s (STR)\n", crit.value.value_str);
            break;
        default:
            printf("UNRECOGNISED_VALUE\n");
    }
}

void printAnnounceData(announceData data) {
    if (data.is_valid) {
        printf("Port: %d\n", data.port);
        for (int i = 0; i < data.nb_files; ++i) {
            printf("File %d: %s, Size: %lld, PieceSize: %lld, Key: %s\n", i + 1,
                   data.files[i].name, data.files[i].size, data.files[i].pieceSize,
                   data.files[i].key);
        }
        for (int i = 0; i < data.nb_leech_keys; ++i) {
            printf("Leech key %d: %s\n", i + 1, data.leechKeys[i]);
        }
    } else {
        printf("announceData is not valid.\n");
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
    if (data.is_valid) {
        printf("Nb criterions : %d\n", data.nb_criterions);
        for (int i = 0; i < data.nb_criterions; ++i) {
            print_criterion(data.criterions[i]);
        }
    } else {
        printf("lookData is not valid.\n");
    }
}

void printUpdateData(updateData data) {
    if (data.is_valid) {
        printf("Keys :");
        for (int i = 0; i < data.nb_keys; ++i) {
            printf(" %s", data.keys[i]);
        }
        printf("\nLeech :");
        for (int i = 0; i < data.nb_leech; ++i) {
            printf(" %s", data.leech[i]);
        }
        printf("\n");
    } else {
        printf("updateData is not valid.\n");
    }
}

