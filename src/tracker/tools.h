#ifndef TOOLS_H
#define TOOLS_H

#include <regex.h>
#include "structs.h"

#define BITS_PER_INT 8*sizeof(int)
#define DELIM " "
#define PORT_MAX_LENGTH 5

void max(int *, int);
int streq(const char *, const char *);
int streqlim(const char *, const char *, int);

regex_t *announce_regex();
regex_t *look_regex();
regex_t *comparison_regex();
regex_t *getfile_regex();
regex_t *update_regex();

announceData announceCheck(char *);
lookData lookCheck(char *);
getfileData getfileCheck(char *);
updateData updateCheck(char *);

int peerCmp(Peer, Peer);
int announceStructCmp(announceData, announceData);
int criterionCmp(criterion, criterion);
int lookStructCmp(lookData, lookData);
int getfileStructCmp(getfileData, getfileData);
int updateStructCmp(updateData, updateData);

void printAnnounceData(announceData);
void print_criterion(criterion);
void printLookData(lookData);
void printGetFileData(getfileData);
void printUpdateData(updateData);

void free_peer(Peer *);
void free_announceData(announceData *);
void free_regex(regex_t *);
void free_all_regex();
void free_file(File *);
void free_announceData(announceData *);
void free_lookData(lookData *);
void free_updateData(updateData *);

#endif //TOOLS_H
