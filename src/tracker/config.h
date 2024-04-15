#ifndef FREE_EIRB2SHARE_T4_CONFIG_H
#define FREE_EIRB2SHARE_T4_CONFIG_H

extern FILE *log_file;

typedef struct {
    char IP[16];
    int IP_mode;
    int listen_port;
} config;

config read_config();

FILE *open_log();
void write_log(const char *format, ...);
#endif //FREE_EIRB2SHARE_T4_CONFIG_H
