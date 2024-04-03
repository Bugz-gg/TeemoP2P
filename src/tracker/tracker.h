#include "tools.h"
#include "structs.h"

static Tracker tracker;

int new_id(Tracker * t , char * addr_ip);

void announce( Tracker * t , char* message , char * addr_ip);

void free_on_exit(int);

void init_tracker();