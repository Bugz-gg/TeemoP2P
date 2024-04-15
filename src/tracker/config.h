#ifndef FREE_EIRB2SHARE_T4_CONFIG_H
#define FREE_EIRB2SHARE_T4_CONFIG_H

typedef struct {
    char IP[16];
    int IP_mode;
    int listen_port;
} config;

config read_config();

void create_log();
void write_log(char *);


#endif //FREE_EIRB2SHARE_T4_CONFIG_H
