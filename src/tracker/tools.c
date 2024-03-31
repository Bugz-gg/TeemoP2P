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

regex_t *announce_regex() {
    static regex_t *regex = NULL;
    if (regex != NULL)
        return regex;
    regex = malloc(sizeof(regex_t));
    //char *pattern = "^announce listen ([0-9]+) seed \\[([a-zA-Z0-9 ]*)\\]( leech \\[(( |[a-zA-Z0-9]{32})*)\\])?$";
    char *pattern = "^announce listen ([0-9]+) seed \\[(([a-zA-Z0-9]+ [0-9]+ [0-9]+ [a-zA-Z0-9]{32}| )*)\\]( leech \\[(([a-zA-Z0-9]{32}| )*)\\])?$";
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
announceData announceCheck(char *message) {
    announceData announceStruct = {.port=0, .nb_files=0, .nb_leech_keys=0, .is_valid=0};

    regex_t *regex = announce_regex();
    regmatch_t matches[6];
    if (regexec(regex, message, 6, matches, 0)) {
        fprintf(stderr, "Failed to match regular expression in %s\n", message);
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
            /*if (strlen(token) != 32) {
                fprintf(stderr, "Wrong md5sum hash size for %s.\n", files[index / 4].name);
                return announceStruct;
            }*/
            for (int i = 0; i < 32; ++i) {
                files[index / 4].key[i] = token[i];
            }
            files[index / 4].key[32] = '\0';
        }
        token = strtok(NULL, DELIM);
        ++index;
    }

    // leech
    int nb_leech_keys = 0;
    char *leechData = strndup(message + matches[5].rm_so, matches[5].rm_eo - matches[5].rm_so);
    count = countDelim(leechData);
    char **leechKeys = malloc(count * sizeof(char *));
    token = strtok(leechData, DELIM);
    int leech_index = 0;
    while (token != NULL) {
        leechKeys[leech_index] = malloc(33 * sizeof(char));
        strncpy(leechKeys[leech_index], token, 32);
        leechKeys[leech_index][32] = '\0';
        token = strtok(NULL, DELIM);
        ++nb_leech_keys;
        ++leech_index;
    }

    announceStruct.port = port;
    announceStruct.nb_files = nbFiles;
    announceStruct.files = files;
    announceStruct.nb_leech_keys = nb_leech_keys;
    announceStruct.is_valid = 1;
    announceStruct.leechKeys = leechKeys;

    free(filesData);

    return announceStruct;
}

// Function to check announce message
lookData lookCheck(char *message) {
    lookData lookStruct = {.is_valid=0, .nb_criterions=0, .criterions=NULL};

    regex_t *regex = look_regex();
    regmatch_t matches[3];
    if (regexec(regex, message, 3, matches, 0)) {
        fprintf(stderr, "Failed to match regular expression in %s\n", message);
        return lookStruct;
    }

    char *criterions_str = strndup(message + matches[1].rm_so, matches[1].rm_eo - matches[1].rm_so);
    int count = countDelim(criterions_str);
    //printf("%s\n", criterions_str);
    count = count + (!count && matches[1].rm_eo - matches[1].rm_so);
    if (!count) {
        fprintf(stderr, "No criteria found.\n");
        return lookStruct;
    }
    regex_t *comp_regex = comparison_regex();
    regmatch_t comparison_match[5];


    criterion *criterions = malloc(count * sizeof(criterion));
    char *token = strtok(criterions_str, DELIM);

    char *endptr;
    int int_val;
    float float_val;
    int index = 0;
    int len_crit;
    char value[100]; // Supposing the max length for the criterions' values is 100.

    char criteria[25];

    while (token != NULL) {
        if (regexec(comp_regex, token, 5, comparison_match, 0)) {
            fprintf(stderr, "Failed to match criterion expression.\n");
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

        token = strtok(NULL, DELIM);
        ++index;
    }

    lookStruct.nb_criterions = count;
    lookStruct.criterions = criterions;
    lookStruct.is_valid = 1;
    free(criterions_str);

    return lookStruct;
}

getfileData getfileCheck(char *message) {
    getfileData getfileStruct = {.is_valid=0};

    regex_t *regex = getfile_regex();
    regmatch_t matches[2];
    if (regexec(regex, message, 2, matches, 0)) {
        fprintf(stderr, "Failed to match regular expression in %s\n", message);
        return getfileStruct;
    }
    // Check if key in files.
    for (int i = 0; i < 32; ++i)
        getfileStruct.key[i] = *(message + matches[1].rm_so + i);
    getfileStruct.key[32] = '\0';
    getfileStruct.is_valid = 1;
    return getfileStruct;
}

int peerCmp(Peer p1, Peer p2) {
    return streq(p1.IP, p2.IP) && p1.port == p2.port;
}

int fileCmp(File f1, File f2) { // The equality of the peers having the file is not checked.
    return streq(f1.name, f2.name) && f1.size == f2.size && f1.pieceSize == f2.pieceSize && streq(f1.key, f2.key);
}

int announceStructCmp(announceData a1, announceData a2) {
    if (a1.is_valid != a2.is_valid || a1.port != a2.port || a1.nb_files != a2.nb_files || a1.nb_leech_keys != a2.nb_leech_keys)
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
            printf("File %d: %s, Size: %d, PieceSize: %d, Key: %s\n", i + 1,
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




