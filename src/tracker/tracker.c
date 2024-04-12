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
        streq(t->files[i]->name, "") ? printf("\033[0;34mLeech file key\033[39m:%s, \033[0;34mPeers' ids\033[39m: ", t->files[i]->key) :printf("\033[0;34mFilename\033[39m: %s, \033[0;34mSize\033[39m: %d(%d), \033[0;34mKey\033[39m:%s, \033[0;34mPeers' ids\033[39m: ",
               t->files[i]->name, t->files[i]->size, t->files[i]->pieceSize, t->files[i]->key);
        for (int j = 0; j < t->files[i]->nb_peers; ++j)
            printf("%d ", t->files[i]->peers[j]->peer_id);
        printf("\n");

    }
}

void print_peer(Peer *p) {
    printf("(%d) \033[0;33m%s:%d\033[39m.\033[39m\n", p->peer_id, p->addr_ip, p->num_port);
}

void print_tracker_peers(Tracker *t) {
    for (int i = 0; i < t->nb_peers; ++i) {
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

File *findfile(Tracker *t, char *k) {
    for (int i = 0; i < t->nb_files; ++i) {
        if (streq(t->files[i]->key, k)) {
            return t->files[i];
        }
    }
    return NULL;
}

Peer *getpeer(Peer **peers, int nb_peers, char *IP, int port) {
    for (int i = 0; i < nb_peers; ++i) {
        if (streq(peers[i]->addr_ip, IP) && peers[i]->num_port == port)
            return peers[i];
    }
    return NULL;
}

Peer *announce(Tracker *t, announceData *d, char *addr_ip, int socket_fd) {
    Peer *peer = getpeer(t->peers, t->nb_peers, addr_ip, d->port); // Vérifie si le peer a déjà communiqué.
    if (peer == NULL) { // Enregistre le peer.
        if (t->nb_peers + 1 > t->alloc_peers) { // Réalloue de la place dans t->peers
            t->alloc_peers *= 2;
            t->peers = realloc(t->peers, (t->alloc_peers) * sizeof(Peer));
        }
        t->peers[t->nb_peers] = malloc(sizeof(Peer)); // Alloue la place d'un Peer pour mettre son adresse dans t->peers
        peer = t->peers[t->nb_peers];
        peer->num_port = d->port;
        peer->peer_id = new_id(t, addr_ip, d->port);
        strcpy(peer->addr_ip, addr_ip);
        ++t->nb_peers;
    }
    // Le peer est maintenant enregistré.

    // TODO: Check coherence ? The following lines may not be necessary.
    /*peer->num_port = d->port;
    peer->peer_id = new_id(t, addr_ip, d->port);
    strcpy(peer->addr_ip, addr_ip);*/

    File *file;
    for (int i = 0; i < d->nb_files; ++i) {
        file = findfile(t, d->files[i].key); // Vérifie si le fichier est déjà enregistré.
        if (file == NULL) { // Enregistre le fichier.
            if (t->alloc_files < t->nb_files + 1) { // Réalloue de la place si besoin.
                t->alloc_files *= 2;
                t->files = realloc(t->files, t->alloc_files * sizeof(void *));
            }
            t->files[t->nb_files] = malloc(sizeof(File)); // Alloue un File pour mettre son adresse dans t->files
            file = t->files[t->nb_files];
            strcpy(file->name, d->files[i].name);
            file->size = d->files[i].size;
            file->pieceSize = d->files[i].pieceSize;
            strcpy(file->key, d->files[i].key);
            file->nb_peers = 0;
            file->alloc_peers = ALLOC_PEERS;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
        }

        // Le fichier est maintenant enregistré.
        // Check is file data is coherent ?
        if (streq(file->name, "")) { // Si le fichier a été ajouté en leech, on ne connaît pas ses informations.
            strcpy(file->name, d->files[i].name);
            file->size = d->files[i].size;
            file->pieceSize = d->files[i].pieceSize;
        }
        Peer *search_peer = getpeer(file->peers, file->nb_peers, peer->addr_ip, peer->num_port);
        if (search_peer == NULL) { // Ajout du peer pour le fichier si besoin.
            if (file->alloc_peers < file->nb_peers + 1) {
                file->alloc_peers *= 2;
                file->peers = realloc(file->peers, file->alloc_peers * sizeof(Peer));
            }
            file->peers[file->nb_peers] = peer;
            ++file->nb_peers;
        }
    }

    for (int i = 0; i < d->nb_leech_keys; ++i) {
        file = findfile(t, d->leechKeys[i]);
        if (file == NULL) {
            if (t->alloc_files < t->nb_files + 1) {
                t->alloc_files *= 2;
                t->files = realloc(t->files, t->alloc_files * sizeof(void *));
            }
            t->files[t->nb_files] = malloc(sizeof(File));
            file = t->files[t->nb_files];
            file->name[0] = '\0'; // Indique que le fichier a été découvert en leech.
            strcpy(file->key, d->leechKeys[i]);
            file->nb_peers = 0;
            file->alloc_peers = ALLOC_PEERS;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
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

    write(socket_fd, "OK\n", 3);
    return peer;
}

void look(Tracker *t, lookData *d, int socket_fd) {
    File **files = malloc(t->nb_files * sizeof(void *));
    memcpy(files, t->files, t->nb_files * sizeof(void *));
    select_files(t->nb_files, files, d->nb_criterions, d->criterions);
    write(socket_fd, "list [", 6);
    int after_first = 0;

    for (int i = 0; i < t->nb_files; ++i) {
        if (files[i] != NULL) {
            if (after_first)
                write(socket_fd, " ", 1);
            sprintf(tmp_buffer, "%s %d %d %s", files[i]->name, files[i]->size, files[i]->pieceSize, files[i]->key);
            write(socket_fd, tmp_buffer, strlen(tmp_buffer));
            after_first = 1;
        }
    }

    write(socket_fd, "]\n", 2);
    free(files);
}

void getfile(Tracker *t, getfileData *d, int socket_fd) {
    File *file = findfile(t, d->key);
    sprintf(tmp_buffer, "peers %s [", d->key);
    write(socket_fd, tmp_buffer, strlen(tmp_buffer));
    int after_first = 0;
    if (file!=NULL) {
        for (int i=0; i<file->nb_peers; ++i) {
            if (after_first)
                write(socket_fd, " ", 1);
            sprintf(tmp_buffer, "%s:%d", file->peers[i]->addr_ip, file->peers[i]->num_port);
            write(socket_fd, tmp_buffer, strlen(tmp_buffer));
            after_first = 1;
        }
    }
    write(socket_fd, "]\n", 2);
}

void updatedata(Tracker *t, updateData *d, int socket_fd) {

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

void select_by_name(File **f, criterion *c) {
    switch (c->op) {
        case LT:
            if (strcmp((*f)->name, c->value.value_str) >= 0)
                *f = NULL;
            break;
        case LE:
            if (strcmp((*f)->name, c->value.value_str) > 0)
                *f = NULL;
            break;
        case EQ:
            if (strcmp((*f)->name, c->value.value_str))
                *f = NULL;
            break;
        case GE:
            if (strcmp((*f)->name, c->value.value_str) < 0)
                *f = NULL;
            break;
        case GT:
            if (strcmp((*f)->name, c->value.value_str) <= 0)
                *f = NULL;
            break;
        case DI:
            if (!strcmp((*f)->name, c->value.value_str))
                *f = NULL;
            break;
        default:
            printf("UNRECOGNISED_OPERATOR ");
    }
}

void select_by_file_size(File **f, criterion *c) {
    switch (c->op) {
        case LT:
            if ((*f)->size >= c->value.value_int) {
                *f = NULL;
            }
            break;
        case LE:
            if ((*f)->size > c->value.value_int) {
                *f = NULL;
            }
            break;
        case EQ:
            if ((*f)->size != c->value.value_int) {
                *f = NULL;
            }
            break;
        case GE:
            if ((*f)->size < c->value.value_int) {
                *f = NULL;
            }
            break;
        case GT:
            if ((*f)->size <= c->value.value_int) {
                *f = NULL;
            }
            break;
        case DI:
            if ((*f)->size == c->value.value_int) {
                *f = NULL;
            }
            break;
        default:
            printf("UNRECOGNISED_OPERATOR ");
    }
}

void select_files(int nb_files, File **f, int nb_criterion, criterion *c) {
    for (int i = 0; i < nb_files; ++i) {
        for (int j = 0; j < nb_criterion; ++j) {
            if (f[i]->name[0]=='\0')
                f[i] = NULL;
            if (f[i] == NULL) // Déjà éliminé par un critérion
                break;
            switch (c[j].criteria) {
                case FILENAME:
                    select_by_name(&f[i], &c[j]);
                    break;
                case FILESIZE:
                    select_by_file_size(&f[i], &c[j]);
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
Peer * findfile(Tracker *t ,char * k ){
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
    printf("On exit 2 : %d\n", tracker.nb_files);
    (void) signo;
    for (int i = 0; i < tracker.nb_peers; ++i)
        free_peer(tracker.peers[i]);
    free(tracker.peers);
    printf("On exit : %d\n", tracker.nb_files);
    for (int i = 0; i < tracker.nb_files; ++i) {
        printf("i: %d\n", i);
        free_file(tracker.files[i]);
    }
    free(tracker.files);
    exit(0);
    return;
}
