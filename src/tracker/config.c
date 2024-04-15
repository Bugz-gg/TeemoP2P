#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <string.h>
#include "config.h"

void create_log() {
    char buffer[15];
    time_t t = time(NULL);
    struct tm tm = *localtime(&t);
    sprintf(buffer, "log-%02d-%02d-%d", tm.tm_mday, tm.tm_mon + 1, tm.tm_year + 1900);
    FILE *file = fopen(buffer, "a");

    fclose(file);
}


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

void write_log(char *) {

}