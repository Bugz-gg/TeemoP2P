#include "tools.h"
#include "structs.h"

static Tracker tracker;

int new_id(Tracker *, char *);

void announce(Tracker *, announceData, char *);

void look(Tracker *, lookData);

void free_on_exit(int);

void init_tracker();