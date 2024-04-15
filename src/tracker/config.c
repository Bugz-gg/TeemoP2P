#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include <stdarg.h>
#include "config.h"

config read_config() {
    char buffer[50];
    config conf = {
            .IP = "127.0.0.1",
            .IP_mode = 0,
            .listen_port=9000
    };
    FILE *config_file = fopen("config.ini", "r");
    char *strToken;
    while (fgets(buffer, 50, config_file) != NULL) {
        strToken = strtok(buffer, " =");
        while (strToken != NULL) {
            if (!strcmp(strToken, "tracker-ip")) {
                strToken = strtok(NULL, " =");
                strcpy(conf.IP, strToken);
            }
            if (!strcmp(strToken, "tracker-ip-mode")) {
                strToken = strtok(NULL, " =");
                conf.IP_mode = atoi(strToken);
            }
            if (!strcmp(strToken, "tracker-port")) {
                strToken = strtok(NULL, " =");
                conf.listen_port = atoi(strToken);
            }
            strToken = strtok(NULL, " =");
        }
    }
    fclose(config_file);
    return conf;
}

FILE *open_log() {
    char buffer[15];
    time_t t = time(NULL);
    struct tm tm = *localtime(&t);
    sprintf(buffer, "%02d-%02d-%d.log", tm.tm_mday, tm.tm_mon + 1, tm.tm_year + 1900);
    return fopen(buffer, "a");
}

void write_log(const char *format, ...) {
    time_t t = time(NULL);
    struct tm tm = *localtime(&t);
    va_list args;
    va_start(args, format);
    fprintf(log_file, "%02d/%02d/%d %02d:%02d:%02d: ", tm.tm_mday, tm.tm_mon + 1, tm.tm_year + 1900, tm.tm_hour, tm.tm_min, tm.tm_sec);
    vfprintf(log_file, format, args);

    va_end(args);
}