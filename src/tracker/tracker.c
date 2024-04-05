#include <stdlib.h>
#include <stdio.h>
#include <unistd.h>
#include <string.h>
#include "tools.h"
#include "tracker.h"
#define MAX_FILES 100
#define MAX_PEERS 200

void init_tracker(Tracker *t){
    t->nb_files=0;
    t->nb_peers=0;
    t->files=malloc(MAX_FILES*sizeof(File));
    t->peers=malloc(MAX_PEERS*sizeof(Peer));
}

void print_tracker_files(Tracker *t){

    for (int i=0;i<t->nb_files;i++){

        printf("Filename: %s\n", t->files[i].name);
        printf("Size: %d\n", t->files[i].size);
        printf("Block Size: %d\n", t->files[i].pieceSize);
        printf("Key: %s\n", t->files[i].key);
        printf("Peers' ids :");
        for (int j=0; i<t->files[i].nb_peers; ++j)
            printf("%d ", t->files[i].peers[j].peer_id);
        printf("\n\n");

    }
}

int new_id(Tracker * t , char * addr_ip){
    static int new_id=0;
    for (int i=0;i<t->nb_peers;i++){
        if(streq(t->peers[i].addr_ip , addr_ip))
            return t->peers[i].peer_id;
        else{
            if(new_id<t->peers[i].peer_id)
                new_id=t->peers[i].peer_id;
        }
    }
    return new_id +1;
}

void announce( Tracker * t , announceData d, char * addr_ip){
    int nb_new_files=d.nb_files;
    if(t->nb_files +nb_new_files>MAX_FILES){
        t->files=realloc(t->files,(nb_new_files + t->nb_files)*sizeof(File));
    }
    for(int i=0;i<nb_new_files;i++){
        
        t->files[t->nb_files+i]=d.files[i];
    }
    t->nb_files+=nb_new_files;
    t->peers[t->nb_peers].num_port=d.port;
    t->peers[t->nb_peers].peer_id = new_id(t,addr_ip) ;
    t->peers[t->nb_peers].addr_ip=addr_ip;
    t->nb_peers+=1;
    printf("OK\n");
}

void remove_file(File * fs , File f ,int* nb){
    int i, j;
    for (i = 0; i < *nb; i++) {
        if (streq(fs[i].name , f.name)) {
            for (j = i; j < *nb - 1; j++) {
                fs[j] = fs[j + 1];
            }
            (*nb)--;
            return;
        }
    }
}

void select_by_name(File * f ,int nb, criterion c ){
    for(int i=0 ; i< nb ; i++){
        if (! streq(c.value.value_str,f[i].name)){
            remove_file(f,f[i],&nb);
        }
    }
}

void select_by_file_size( File * f ,int nb, criterion c){
    
    switch (c.op) {
        case LT:
            for(int i=0 ; i< nb ; i++){
            if ( f[i].size < c.value.value_int){
                remove_file(f,f[i],&nb);
            }
            }
            break;
        case LE:
            for(int i=0 ; i< nb ; i++){
            if (c.value.value_int <= f[i].size){
                remove_file(f,f[i],&nb);
            }
            }
            break;
        case EQ:
            for(int i=0 ; i< nb ; i++){
            if (f[i].size == c.value.value_int){
                remove_file(f,f[i],&nb);
            }
            }
            break;
        case GE:
            for(int i=0 ; i< nb ; i++){
            if (f[i].size >= c.value.value_int){
                remove_file(f,f[i],&nb);
            }
            }
            break;
        case GT:
            for(int i=0 ; i< nb ; i++){
            if (c.value.value_int < f[i].size){
                remove_file(f,f[i],&nb);
            }
            }
            break;
        case DI:
            for(int i=0 ; i< nb ; i++){
            if (c.value.value_int!=f[i].size){
                remove_file(f,f[i],&nb);
            }
            }
            break;
        default:
            printf("UNRECOGNISED_OPERATOR ");
    }
}

void select_files(File * f ,int nb, criterion c ){
        switch (c.criteria) {
        case FILENAME:
            select_by_name(f ,nb, c );
            break;
        case FILESIZE:
            select_by_file_size(f , nb, c );
        default:
            printf("UNRECOGNISED_CRITERIA ");
    }
}

void look(Tracker *t , lookData data){
    File * files=t->files;;
    criterion * l=data.criterions;
    unsigned int nb=data.nb_criterions;
    for(int i=0;i<nb;i++){
        select_files(files,t->nb_files,l[i]);    
    }
    return;
}

Peer select_peer(Tracker *t ,int id){
    for(int i=0;i<t->nb_peers;i++){
        if(t->peers[i].peer_id ==id){
            return t->peers[i];
        }
    }
    Peer not_found = {.peer_id=-1, .addr_ip="", .num_port=-1};
    return not_found;
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


Peer *getfile(Tracker *t ,char *k){
    Peer *p=NULL;
    for(int i=0; i<t->nb_files; ++i){
        if(streq(t->files[i].key,k)){
            return t->files[i].peers;
        }
    }
    return p;
}

void free_on_exit(int signo) {
    (void)signo;
    for (int i=0; i<tracker.nb_peers; ++i)
        free_peer(&tracker.peers[i]);
    for (int i=0; i<tracker.nb_files; ++i)
        free_file(&tracker.files[i]);
    exit(0);
    return;
}

void init_tracker() {
    tracker.nb_peers = 0;
    tracker.nb_files = 0;
}