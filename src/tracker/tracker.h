#include "tools.h"
#include "structs.h"

static Tracker tracker;

int new_id(Tracker *, char *, int);

void announce(Tracker *, announceData *, char *, int);
void look(Tracker *, lookData *, int);

void select_files(int, File **, int, criterion *);

void free_on_exit(int);

void init_tracker(Tracker *);