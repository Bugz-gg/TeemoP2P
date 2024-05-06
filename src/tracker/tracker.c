#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include <pthread.h>
#include "tools.h"
#include "tracker.h"
#include "structs.h"

#define MAX_FILES 100 // Check coherence with structs.h
#define MAX_PEERS 200

extern Peer *connected_peers[MAX_PEERS];
extern pthread_mutex_t mutex_for_connected_peers;

char tmp_buffer[256]; // Used to send messages back.

void init_tracker(Tracker *t) {
    t->nb_files = 0;
    t->nb_peers = 0;
    t->max_peer_ind = 0;
    t->alloc_files = MAX_FILES;
    t->alloc_peers = MAX_PEERS;
    t->files = malloc(MAX_FILES * sizeof(File));
    t->peers = malloc(MAX_PEERS * sizeof(Peer));
}

void print_tracker_files(Tracker *t) {
    for (int i = 0; i < t->max_file_ind; i++) {
        if (t->files[i] == NULL)
            continue;
        streq(t->files[i]->name, "") ? printf("\033[0;34mLeech file key\033[39m:%s, \033[0;34mPeers' ids\033[39m: ",
                                              t->files[i]->key) : printf(
                "\033[0;34mFilename\033[39m: %s, \033[0;34mSize\033[39m: %lld(%lld), \033[0;34mKey\033[39m:%s, \033[0;34mPeers' ids\033[39m: ",
                t->files[i]->name, t->files[i]->size, t->files[i]->pieceSize, t->files[i]->key);
        for (int j = 0; j < t->files[i]->max_peer_ind; ++j) {
            if (t->files[i]->peers[j] == NULL)
                continue;
            printf("%d ", t->files[i]->peers[j]->peer_id);
        }
        printf("\n");

    }
}

void print_peer(Peer *p) {
    printf("(%d) \033[0;33m%s:%d\033[39m.\033[39m\n", p->peer_id, p->addr_ip, p->num_port);
}

void print_tracker_peers(Tracker *t) {
    for (int i = 0; i < t->max_peer_ind; ++i) {
        if (t->peers[i] == NULL)
            continue;
        print_peer(t->peers[i]);
    }
}

int new_id(Tracker *t, char *addr_ip, int port) {
    static int new_id = 0;
    for (int i = 0; i < t->max_peer_ind; i++) {
        if (t->peers[i] == NULL) {
            continue;
        }
        if (streq(t->peers[i]->addr_ip, addr_ip) && t->peers[i]->num_port == port)
            return t->peers[i]->peer_id;
        else {
            if (new_id < t->peers[i]->peer_id)
                new_id = t->peers[i]->peer_id;
        }
    }
    return new_id + 1;
}

int peer_cmp(void *p1, void *p2) {
    return (streq(((Peer *) p1)->addr_ip, ((Peer *) p2)->addr_ip) &&
            ((Peer *) p1)->num_port == ((Peer *) p2)->num_port);
}

int file_cmp(void *f1, void *f2) {
    return streq(((File *) f1)->key, (((File *) f2)->key));
}

int key_cmp(void *k1, void *k2) {
    return streq(((char *) k1), (((char *) k2)));
}

// If found is NULL, index is the first NULL's index. Otherwise, it is the index of the found.
typedef struct {
    int index;
    void *found;
} findings;

findings find(void **tab, int length, void *search_struct, int (*cmp_func)(void *, void *)) {
    findings f = {.index = length, .found = NULL};
    for (int i = 0; i < length; ++i) {
        if (tab[i] == NULL) {
            if (f.index == length)
                f.index = i;
            continue;
        }
        if (cmp_func(tab[i], search_struct)) {
            f.index = i;
            f.found = tab[i];
            return f;
        }
    }
    return f;
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
        if (peers[i] == NULL)
            continue;
        if (streq(peers[i]->addr_ip, IP) && peers[i]->num_port == port)
            return peers[i];
    }
    return NULL;
}

void cond_realloc(void **ptr, int *nb_alloc, int nb, void *obj) {
    if (nb > *nb_alloc) { // Réalloue de la place dans t->peers
        *nb_alloc *= 2;
        *ptr = realloc(*ptr, *nb_alloc * sizeof(*obj));
        (void) obj;
    }
}

Peer *announce(Tracker *t, announceData *d, char *addr_ip, int socket_fd, int index) {
    Peer target_peer = {.num_port=d->port};
    strcpy(target_peer.addr_ip, addr_ip);
    findings found_peer = find((void **) t->peers, t->max_peer_ind, &target_peer, *peer_cmp);
    //Peer *peer = getpeer(t->peers, t->nb_peers, addr_ip, d->port); // Vérifie si le peer a déjà communiqué.
    Peer *peer = (Peer *) found_peer.found;
    if (peer == NULL) { // Enregistre le peer.
        if (found_peer.index == t->nb_peers)
            cond_realloc((void **) t->peers, &t->alloc_peers, t->nb_peers + 1, &target_peer);
        t->peers[found_peer.index] = malloc(
                sizeof(Peer)); // Alloue la place d'un Peer pour mettre son adresse dans t->peers
        peer = t->peers[found_peer.index];
        peer->num_port = d->port;
        peer->peer_id = index; //new_id(t, addr_ip, d->port);
        strcpy(peer->addr_ip, addr_ip);
        ++t->nb_peers;
        max(&t->max_peer_ind, t->nb_peers);
    }
    // Le peer est maintenant enregistré.

    // TODO: Check coherence ? The following lines may not be necessary.
    /*peer->num_port = d->port;
    peer->peer_id = new_id(t, addr_ip, d->port);
    strcpy(peer->addr_ip, addr_ip);*/

    File *file;
    File target_file;
    for (int i = 0; i < d->nb_files; ++i) {
        strcpy(target_file.key, d->files[i].key);
        findings found_file = find((void **) t->files, t->max_file_ind, &target_file, *file_cmp);
        //file = findfile(t, d->files[i].key); // Vérifie si le fichier est déjà enregistré.
        file = found_file.found;
        if (file == NULL) { // Enregistre le fichier.
            if (found_file.index == t->nb_files)
                cond_realloc((void **) t->files, &t->alloc_files, t->nb_files + 1, &target_file);

            t->files[found_file.index] = malloc(sizeof(File)); // Alloue un File pour mettre son adresse dans t->files
            file = t->files[found_file.index];
            strcpy(file->name, d->files[i].name);
            file->size = d->files[i].size;
            file->pieceSize = d->files[i].pieceSize;
            strcpy(file->key, d->files[i].key);
            file->nb_peers = 0;
            file->max_peer_ind = 0;
            file->alloc_peers = ALLOC_PEERS;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
            max(&t->max_file_ind, t->nb_files);
        }

        // Le fichier est maintenant enregistré.
        // Check if file data is coherent ?
        if (streq(file->name, "")) { // Si le fichier a été ajouté en leech, on ne connaît pas ses informations.
            strcpy(file->name, d->files[i].name);
            file->size = d->files[i].size;
            file->pieceSize = d->files[i].pieceSize;
        }
        findings found_peer = find((void **) file->peers, file->max_peer_ind, peer, *peer_cmp);
        //Peer *search_peer = getpeer(file->peers, file->nb_peers, peer->addr_ip, peer->num_port);
        Peer *found_file_peer = found_peer.found;
        if (found_file_peer == NULL) { // Ajout du peer pour le fichier si besoin.
            cond_realloc((void **) file->peers, &file->alloc_peers, file->nb_peers + 1, peer);
            file->peers[found_peer.index] = peer;
            ++file->nb_peers;
            max(&file->max_peer_ind, file->nb_peers);
        }
    }

    for (int i = 0; i < d->nb_leech_keys; ++i) {
        //file = findfile(t, d->leechKeys[i]);
        strcpy(target_file.key, d->leechKeys[i]);
        findings found_file = find((void **) t->files, t->max_file_ind, &target_file, *file_cmp);
        file = found_file.found;
        if (file == NULL) {
            if (found_file.index == t->nb_files)
                cond_realloc((void **) t->files, &t->alloc_files, t->nb_files + 1, &target_file);
            t->files[found_file.index] = malloc(sizeof(File));
            file = t->files[found_file.index];
            file->name[0] = '\0'; // Indique que le fichier a été découvert en leech.
            strcpy(file->key, d->leechKeys[i]);
            file->nb_peers = 0;
            file->max_peer_ind = 0;
            file->alloc_peers = ALLOC_PEERS;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
            max(&t->max_file_ind, t->nb_files);
        }
        //Peer *search_peer = getpeer(file->peers, file->nb_peers, peer->addr_ip, peer->num_port);
        findings found_peer = find((void **) file->peers, file->max_peer_ind, &target_peer, *peer_cmp);

        if (found_peer.found == NULL) {
            if (found_peer.index == t->nb_peers)
                cond_realloc((void **) file->peers, &t->alloc_peers, t->nb_peers + 1, &target_peer);
            file->peers[found_peer.index] = peer;
            ++file->nb_peers;
            max(&file->max_peer_ind, file->nb_peers);
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
            sprintf(tmp_buffer, "%s %lld %lld %s", files[i]->name, files[i]->size, files[i]->pieceSize, files[i]->key);
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
    if (file != NULL) {
        for (int i = 0; i < file->max_peer_ind; ++i) {
            if (after_first)
                write(socket_fd, " ", 1);
            if (file->peers[i] == NULL)
                continue;
            sprintf(tmp_buffer, "%s:%d", file->peers[i]->addr_ip, file->peers[i]->num_port);
            write(socket_fd, tmp_buffer, strlen(tmp_buffer));
            after_first = 1;
        }
    }
    write(socket_fd, "]\n", 2);
}

void update(Tracker *t, updateData *d, int socket_fd, int index) {
    pthread_mutex_lock(&mutex_for_connected_peers);
    Peer *peer = connected_peers[index];
    pthread_mutex_unlock(&mutex_for_connected_peers);
    if (peer == NULL) { // Pas enregistré donc pas fait d'announce pour annoncer son port.
        return;
    }

    File *file;
    File target_file;
    for (int i = 0; i < d->nb_keys; ++i) {
        strcpy(target_file.key, d->keys[i]);
        findings found_file = find((void **) t->files, t->max_file_ind, &target_file, *file_cmp);
        //file = findfile(t, d->files[i].key); // Vérifie si le fichier est déjà enregistré.
        file = found_file.found;
        if (file == NULL) { // Enregistre le fichier.
            if (found_file.index == t->nb_files)
                cond_realloc((void **) t->files, &t->alloc_files, t->nb_files + 1, &target_file);

            t->files[found_file.index] = malloc(sizeof(File)); // Alloue un File pour mettre son adresse dans t->files
            file = t->files[found_file.index];
            file->name[0] = '\0';
            file->size = 0;
            file->pieceSize = 0;
            strcpy(file->key, d->keys[i]);
            file->nb_peers = 0;
            file->max_peer_ind = 0;
            file->alloc_peers = ALLOC_PEERS;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
            max(&t->max_file_ind, t->nb_files);
        }

        // Le fichier est maintenant enregistré.

        findings found_peer = find((void **) file->peers, file->max_peer_ind, peer, *peer_cmp);
        //Peer *search_peer = getpeer(file->peers, file->nb_peers, peer->addr_ip, peer->num_port);
        Peer *found_file_peer = found_peer.found;
        if (found_file_peer == NULL) { // Ajout du peer pour le fichier si besoin.
            cond_realloc((void **) file->peers, &file->alloc_peers, file->nb_peers + 1, peer);
            file->peers[found_peer.index] = peer;
            ++file->nb_peers;
            max(&file->max_peer_ind, file->nb_peers);
        }
    }

    for (int i = 0; i < d->nb_leech; ++i) {
        //file = findfile(t, d->leechKeys[i]);
        strcpy(target_file.key, d->leech[i]);
        findings found_file = find((void **) t->files, t->max_file_ind, &target_file, *file_cmp);
        file = found_file.found;
        if (file == NULL) {
            if (found_file.index == t->nb_files)
                cond_realloc((void **) t->files, &t->alloc_files, t->nb_files + 1, &target_file);
            t->files[found_file.index] = malloc(sizeof(File));
            file = t->files[found_file.index];
            file->name[0] = '\0'; // Indique que le fichier a été découvert en leech.
            strcpy(file->key, d->leech[i]);
            file->nb_peers = 0;
            file->max_peer_ind = 0;
            file->alloc_peers = ALLOC_PEERS;
            file->peers = malloc(file->alloc_peers * sizeof(Peer));
            ++t->nb_files;
            max(&t->max_file_ind, t->nb_files);
        }
        findings found_peer = find((void **) file->peers, file->max_peer_ind, peer, *peer_cmp);
        if (found_peer.found == NULL) {
            if (found_peer.index == t->nb_peers)
                cond_realloc((void **) file->peers, &t->alloc_peers, t->nb_peers + 1, &peer);
            file->peers[found_peer.index] = peer;
            ++file->nb_peers;
            max(&file->max_peer_ind, file->nb_peers);
        }
    }

    // On cherche d'abord un fichier pas dans les seed ni dans les leech puis on regarde si le peer est parmi la liste des peers.
    // Il serait aussi possible de d'abord chercher si le peer est dans la liste des peers avant de vérifier que le fichier soit dans les seed ou leech.
    findings found_peer;
    findings found_key;
    for (int i = 0; i <
                    t->max_file_ind; ++i) { // Suppression des entrées du peer si le peer l'a pas indiqué dans update qu'il ait encore les fichiers.
        if (t->files[i] == NULL)
            continue;
        found_key = find((void **) d->keys, d->nb_keys, t->files[i]->key, *key_cmp);
        if (found_key.found == NULL) {// Pas dans les seed mais peut-être dans leech.
            found_key = find((void **) d->leech, d->nb_leech, t->files[i]->key, *key_cmp);
            if (found_key.found == NULL) { // Pas dans seed ni dans leech, peut-être à supprimer.
                found_peer = find((void **) t->files[i]->peers, t->files[i]->max_peer_ind, peer, *peer_cmp);
                if (found_peer.found == NULL) // Le peer n'a en fait jamais annoncé avoir ce fichier.
                    continue;

                // Ici on sait que le peer avait annoncé avoir le fichier mais ne l'a pas fait dans cette commande update.
                t->files[i]->peers[found_peer.index] = NULL;
                --t->files[i]->nb_peers;
                if (!t->files[i]->nb_peers) { // Plus aucun peer n'a le fichier.
                    free_file(t->files[i]);
                    t->files[i] = NULL;
                    --t->nb_files;
                }
            }
        }
    }
    write(socket_fd, "OK\n", 3);
}

void remove_peer_all_files(Tracker *t, Peer *peer) {
    findings found_peer;
    for (int i = 0; i < t->max_file_ind; ++i) {
        if (t->files[i] == NULL)
            continue;
        found_peer = find((void **) t->files[i]->peers, t->files[i]->max_peer_ind, peer, *peer_cmp);
        if (found_peer.found == NULL) // Le peer n'a en fait jamais annoncé avoir ce fichier.
            continue;
        t->files[i]->peers[found_peer.index] = NULL;
        --t->files[i]->nb_peers;
        if (!t->files[i]->nb_peers) { // Plus aucun peer n'a le fichier.
            free_file(t->files[i]);
            t->files[i] = NULL;
            --t->nb_files;
        }
    }
    found_peer = find((void **) t->peers, t->max_peer_ind, peer, *peer_cmp);
    if (found_peer.found == NULL) // Impossible d'arriver ici normalement.
        return;
    free_peer(peer);
    t->peers[found_peer.index] = NULL;
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

void select_by_key(File **f, criterion *c) {
    switch (c->op) {
        case LT:
            if (strcmp((*f)->key, c->value.value_str) >= 0)
                *f = NULL;
            break;
        case LE:
            if (strcmp((*f)->key, c->value.value_str) > 0)
                *f = NULL;
            break;
        case EQ:
            if (strcmp((*f)->key, c->value.value_str))
                *f = NULL;
            break;
        case GE:
            if (strcmp((*f)->key, c->value.value_str) < 0)
                *f = NULL;
            break;
        case GT:
            if (strcmp((*f)->key, c->value.value_str) <= 0)
                *f = NULL;
            break;
        case DI:
            if (!strcmp((*f)->key, c->value.value_str))
                *f = NULL;
            break;
        default:
            printf("UNRECOGNISED_OPERATOR ");
    }
}

void select_files(int nb_files, File **f, int nb_criterion, criterion *c) {
    for (int i = 0; i < nb_files; ++i) {
        for (int j = 0; j < nb_criterion; ++j) {
            if (f[i]->name[0] == '\0')
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
                case KEY:
                    select_by_key(&f[i], &c[j]);
                    break;
                default:
                    printf("UNRECOGNISED_CRITERIA ");
            }

        }
    }

}

Peer *select_peer(Tracker *t, int id) {
    for (int i = 0; i < t->max_peer_ind; i++) {
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
    for (int i = 0; i < tracker.max_peer_ind; ++i) {
        if (tracker.peers[i] == NULL)
            continue;
        free_peer(tracker.peers[i]);
    }
    free(tracker.peers);
    printf("On exit : %d\n", tracker.nb_files);
    for (int i = 0; i < tracker.nb_files; ++i) {
        printf("i: %d\n", i);
        free_file(tracker.files[i]);
    }
    free(tracker.files);
    exit(0);
}
