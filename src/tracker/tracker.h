#include "tools.h"

typedef struct{
    char * addr_ip; 
    int num_port;
    int peer_id;
}Peer;

typedef struct{
    Peer * peers;
    File * files;
    int nb_files;
    int nb_peers;
}Tracker;

int new_id(Tracker * t , char * addr_ip);

void announce( Tracker * t , char* message , char * addr_ip);