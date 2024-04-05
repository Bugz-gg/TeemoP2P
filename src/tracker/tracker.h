#include "tools.h"
#include "structs.h"

static Tracker tracker;

int new_id(Tracker *, char *);

char *announce(Tracker *, announceData, char *);

char *look(Tracker *, lookData);

void free_on_exit(int);

void init_tracker();