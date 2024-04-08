#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include "tools.h"
#include "tracker.h"

#define MAX_FILES 100 // Check coherence with structs.h
#define MAX_PEERS 200

char tmp_buffer[200]; // Used to send messages back.

void init_tracker(Tracker *t) {
    t->nb_files = 0;
    t->nb_peers = 0;
    t->alloc_files = MAX_FILES;
    t->alloc_peers = MAX_PEERS;
    t->files = malloc(MAX_FILES * sizeof(File));
    t->peers = malloc(MAX_PEERS * sizeof(Peer));
}

void print_tracker_files(Tracker *t) {
    for (int i = 0; i < t->nb_files; i++) {
        printf("\033[0;34mFilename\033[39m: %s, \033[0;34mSize\033[39m: %d(%d), \033[0;34mKey\033[39m:%s, \033[0;34mPeers' ids\033[39m: ", t->files[i]->name, t->files[i]->size, t->files[i]->pieceSize, t->files[i]->key);
        for (int j = 0; j < t->files[i]->nb_peers; ++j)
            printf("%d ", t->files[i]->peers[j]->peer_id);
        printf("\n");

    }
}

void print_peer(Peer *p) {
    printf("(%d) \033[0;33m%s:%d\033[39m.\033[39m\n", p->peer_id, p->addr_ip, p->num_port);
}

void print_tracker_peers(Tracker *t) {
    for (int i=0; i<t->nb_peers; ++i) {
        print_peer(t->peers[i]);
    }
}

int new_id(Tracker *t, char *addr_ip, int port) {
    static int new_id = 0;
    for (int i = 0; i < t->nb_peers; i++) {
        if (streq(t->peers[i]->addr_ip, addr_ip) && t->peers[i]->num_port == port)
            return t->peers[i]->peer_id;
        else {
            if (new_id < t->peers[i]->peer_id)
                new_id = t->peers[i]->peer_id;
        }
    }
    return new_id + 1;
}

File *getfile(Tracker *t, char *k) {
    for (int i = 0; i < t->nb_files; ++i) {
        if (streq(t->files[i]->key, k)) {
            return t->files[i];
        }
    }
    return NULL;
}

Peer *getpeer(Peer **peers, int nb_peers, char *IP, int port) {
    printf("Searching for %s:%d\n", IP, port);
    for (int i = 0; i < nb_peers; ++i) {
        printf("Currently %s:%d\n", peers[i]->addr_ip, peers[i]->num_port);
        if (streq(peers[i]->addr_ip, IP) && peers[i]->num_port == port)
            return peers[i];
    }
    return NULL;
}

void announce(Tracker *t, announceData *d, char *addr_ip, int socket_fd) {
    Peer *peer = getpeer(t->peers, t->nb_peers, addr_ip, d->port);
    if (peer == NULL) {
        if (t->nb_peers+1 > t->alloc_peers) {
            t->alloc_peers *= 2;
            t->peers = realloc(t->peers, (t->alloc_peers) * sizeof(Peer));
        }
        t->peers[t->nb_peers] = malloc(sizeof(Peer));
        peer = t->peers[t->nb_peers];
        peer->num_port = d->port;
        peer->peer_id = new_id(t, addr_ip, d->port);
        strcpy(peer->addr_ip, addr_ip);
        ++t->nb_peers;
    }

    // Check coherence ? The following lines may not be necessary.
    peer->num_port = d->port;
    peer->peer_id = new_id(t, addr_ip, d->port);
    strcpy(peer->addr_ip, addr_ip);

    File *file;
    for (int i = 0; i < d->nb_files; ++i) {
        file = getfile(t, d->files[i].key);
        if (file == NULL) {

            if (t->alloc_files < t->nb_files+1) {
                t->alloc_files *= 2;
                t->files = realloc(t->files, t->alloc_files * sizeof(void *));
            }
            t->files[t->nb_files] = malloc(sizeof(File));
            file = t->files[t->nb_files];
            strcpy(file->name, d->files[i].name);
            file->size = d->files[i].size;
            strcpy(file->key, d->files[i].key);
            file->pieceSize = d->files[i].pieceSize;
            file->alloc_peers = ALLOC_PEERS;
            file->nb_peers = 0;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
        }
        // Check is file data is coherent ?
        if (streq(file->name, "")) { // If it was first added as leech.
            file = t->files[t->nb_files];
            strcpy(file->name, d->files[i].name);
            file->size = d->files[i].size;
            file->pieceSize = d->files[i].pieceSize;
        }
        Peer *search_peer = getpeer(file->peers, file->nb_peers, peer->addr_ip, peer->num_port);
        if (search_peer == NULL) {
            if (file->alloc_peers < file->nb_peers + 1) {
                file->alloc_peers *= 2;
                file->peers = realloc(file->peers, file->alloc_peers * sizeof(Peer));
            }
            file->peers[file->nb_peers] = peer;
            ++file->nb_peers;
        }
    }

    for (int i = 0; i < d->nb_leech_keys; ++i) {
        file = getfile(t, d->leechKeys[i]);
        if (file == NULL) {
            if (t->alloc_files < t->nb_files+1) {
                t->alloc_files *= 2;
                t->files = realloc(t->files, t->alloc_files * sizeof(void *));
            }

            file = t->files[t->nb_files];
            file->name[0] = '\0'; // Tell that we only got the key from leech.
            strcpy(file->key, d->files[i].key);
            file->alloc_peers = ALLOC_PEERS;
            file->nb_peers = 0;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
        }
        Peer *search_peer = getpeer(file->peers, file->nb_peers, peer->addr_ip, peer->num_port);
        if (search_peer == NULL) {
            if (file->alloc_peers < file->nb_peers + 1) {
                file->alloc_peers *= 2;
                file->peers = realloc(file->peers, file->alloc_peers * sizeof(Peer));
            }
            file->peers[file->nb_peers] = t->peers[t->nb_peers];
            ++file->nb_peers;
        }
    }

    write(socket_fd, "OK\n", 3);
}

void look(Tracker *t, lookData *data, int socket_fd) {
    File **files = malloc(t->nb_files * sizeof(void *));
    memcpy(files, t->files, t->nb_files * sizeof(void *));
    select_files(t->nb_files, files, data->nb_criterions, data->criterions);
    write(socket_fd, "list [", 6);

    for (int i = 0; i < t->nb_files - 1; ++i) {
        if (files[i] != NULL) {
            sprintf(tmp_buffer, "%s %d %d %s ", files[i]->name, files[i]->size, files[i]->pieceSize, files[i]->key);
            write(socket_fd, tmp_buffer, strlen(tmp_buffer));
        }
    }
    if (files[t->nb_files - 1] != NULL) {
        sprintf(tmp_buffer, "%s %d %d %s", files[t->nb_files - 1]->name, files[t->nb_files - 1]->size,
                files[t->nb_files - 1]->pieceSize, files[t->nb_files - 1]->key);
        write(socket_fd, tmp_buffer, strlen(tmp_buffer));
    }

    write(socket_fd, "]\n", 3);
}

void remove_file(File *fs, File f, int *nb) {
    int i, j;
    for (i = 0; i < *nb; i++) {
        if (streq(fs[i].name, f.name)) {
            for (j = i; j < *nb - 1; j++) {
                fs[j] = fs[j + 1];
            }
            (*nb)--;
            return;
        }
    }
}

void select_by_name(File *f, criterion *c) {
    switch (c->op) {
        case LT:
            if (strcmp(f->name, c->value.value_str) >= 0)
                f = NULL;
            break;
        case LE:
            if (strcmp(f->name, c->value.value_str) > 0)
                f = NULL;
            break;
        case EQ:
            if (strcmp(f->name, c->value.value_str))
                f = NULL;
            break;
        case GE:
            if (strcmp(f->name, c->value.value_str) < 0)
                f = NULL;
            break;
        case GT:
            if (strcmp(f->name, c->value.value_str) <= 0)
                f = NULL;
            break;
        case DI:
            if (!strcmp(f->name, c->value.value_str))
                f = NULL;
            break;
        default:
            printf("UNRECOGNISED_OPERATOR ");
    }
}

void select_by_file_size(File *f, criterion *c) {
    switch (c->op) {
        case LT:
            if (f->size >= c->value.value_int) {
                f = NULL;
            }
            break;
        case LE:
            if (f->size > c->value.value_int) {
                f = NULL;
            }
            break;
        case EQ:
            if (f->size != c->value.value_int) {
                f = NULL;
            }
            break;
        case GE:
            if (f->size < c->value.value_int) {
                f = NULL;
            }
            break;
        case GT:
            if (f->size <= c->value.value_int) {
                f = NULL;
            }
            break;
        case DI:
            if (f->size == c->value.value_int) {
                f = NULL;
            }
            break;
        default:
            printf("UNRECOGNISED_OPERATOR ");
    }
}

void select_files(int nb_files, File **f, int nb_criterion, criterion *c) {
    for (int i = 0; i < nb_files; ++i) {
        for (int j = 0; j < nb_criterion; ++j) {
            if (f[i] == NULL) // Déjà éliminé par un critérion
                break;
            switch (c[j].criteria) {
                case FILENAME:
                    select_by_name(f[i], &c[j]);
                    break;
                case FILESIZE:
                    select_by_file_size(f[i], &c[j]);
                    break;
                default:
                    printf("UNRECOGNISED_CRITERIA ");
            }

        }
    }

}

Peer *select_peer(Tracker *t, int id) {
    for (int i = 0; i < t->nb_peers; i++) {
        if (t->peers[i]->peer_id == id) {
            return t->peers[i];
        }
    }
    //Peer not_found = {.peer_id=-1, .addr_ip="", .num_port=-1};
    return NULL;
}

/*
Peer * getfile(Tracker *t ,char * k ){
    Peer * p=malloc(t->nb_peers * sizeof(Peer));
    int nb=0;
    for( int i=0;i<t->nb_files;i++){
        if(streq(t->files[i].key,k)){
            int id=t->files[i].peer_id;
            p[nb]=select_peer(t,id);
        }
    }
    return p;
}*/

void free_on_exit(int signo) {
    (void) signo;
    for (int i = 0; i < tracker.nb_peers; ++i)
        free_peer(tracker.peers[i]);
    free(tracker.peers);
    for (int i = 0; i < tracker.nb_files; ++i)
        free_file(tracker.files[i]);
    free(tracker.files);
    exit(0);
    return;
}
