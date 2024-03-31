#include "tools.h"

typedef struct{
    Peer * peers;
    File * files;
    int nb_files;
    int nb_peers;
}Tracker;

int new_id(Tracker * t , char * addr_ip);

void announce( Tracker * t , char* message , char * addr_ip);