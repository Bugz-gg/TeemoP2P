#include "tools.h"
#include "structs.h"

static Tracker tracker;

int new_id(Tracker *, char *, int);

Peer *announce(Tracker *, announceData *, char *, int, int);
void look(Tracker *, lookData *, int);
void getfile(Tracker *, getfileData *, int);
void update(Tracker *, updateData *, int, int);

void remove_peer_all_files(Tracker *, Peer *);
void select_files(int, File **, int, criterion *);

void free_on_exit(int);

void init_tracker(Tracker *);

void print_tracker_peers(Tracker *);
void print_tracker_files(Tracker *);