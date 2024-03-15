#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <stdlib.h>
#include <string.h>
#include "tools.h"

#define MAX_LENGTH 100

void error(char *msg)
{
perror(msg);
exit(EXIT_FAILURE);
}

typedef struct {
    char *filename;
    int length;
    int piece_size;
    char *key;
    int *peer_id;
} File;

typedef struct {
    int peer_id;
    char *ip;
    int port;
} Peer;

/*
add_peer prend la liste des fichiers de la pair et les enregistre dans le fichier de configuration config.ini 
files contenient l'id de la pair p (files[0]) puis une succession nom taille taille du bloc clé pour chaque fichier de la pair puis NULL
*/
void add_peer(char * files[]){ 
   FILE *config_file = fopen("config.ini", "a");
    if (config_file == NULL) {
        error("Error opening file");
    }
    int i=1;
    while(files[i]!=NULL){
        fprintf(config_file, "\n[peer %s]\n" ,files[0]);
        fprintf(config_file, "File Name: %s\n",files[i] );
        i++;
        if (files[i] == NULL) {
            printf("Error : Length missing \n");
            break;
        }
        fprintf(config_file, "Length: %s\n",files[i]);
        i++;
        if (files[i] == NULL) {
            printf("Error: Block size missing.\n");
            break;
        }
        fprintf(config_file, "Block Size: %s\n",files[i]);
        i++;
        if (files[i] == NULL) {
            printf("Error: Key missing.\n");
            break;
        }
        fprintf(config_file, "Key: %s \n",files[i]);
        i++;
    }

    fclose(config_file);
    printf("OK\n");
}


/*
read_config_file lie le fichier de configuration et replie files par les informations des fichiers présents dans le réseau
*/
void read_config_file(File files[]) {
    FILE *config_file = fopen("config.ini", "r");
    if (config_file == NULL) {
        error("Error opening file\n");
    }

    char line[100];
    int file_count = -1; 

    while (fgets(line, sizeof(line), config_file) != NULL) {
        if (strstr(line, "[peer") != NULL) {
            file_count++;
            sscanf(line, "[peer %d]", &files[file_count].peer_id);
        } else if (strstr(line, "File Name:") != NULL) {
            char *filename_start = strstr(line, ":") + 2;
            size_t len = strcspn(filename_start, "\n");
            files[file_count].filename = malloc(len + 1);
            if (files[file_count].filename != NULL) {
                strncpy(files[file_count].filename, filename_start, len);
                files[file_count].filename[len] = '\0';
            } else {
                error("Error allocating memory for filename\n");
            }
        } else if (strstr(line, "Length:") != NULL) {
            sscanf(line, "Length: %d", &files[file_count].length);
        } else if (strstr(line, "Block Size:") != NULL) {
            sscanf(line, "Block Size: %d", &files[file_count].piece_size);
        } else if (strstr(line, "Key:") != NULL) {
            char *key_start = strstr(line, ":") + 2;
            size_t len = strlen(key_start);
            files[file_count].key = malloc(len + 1);
            if (files[file_count].key != NULL) {
                strncpy(files[file_count].key, key_start, len);
                files[file_count].key[len] = '\0';
            } else {
                error("Error allocating memory for key\n");
            }
        }
    }

    fclose(config_file);
}

void free_files_memory(File files[], int count) {
    for (int i = 0; i < count; i++) {
        free(files[i].filename);
        free(files[i].key);
    }
    
}



int is_his_name(char * name , File f){
    return (strcmp(name,f.filename)==0);
}

int is_his_key(char * key , File f){
    return (strcmp(key,f.key)==0);
}

int is_his_length(int l,File f){
    return f.length==l;
}

/*
look donne la liste des fichiers présents dans le réseau vérifiant un certain nombre de critères

cond est la liste des critères recherchès , le premier élément de cond doit étre le nom recherché ou "" ,
 le deuxiéme la taille des fichiers recherchés ou "" , le troisiéme la taille des blocs des fichiers ou "" 
 et le quatriéme la clé rechérchée ou "" et le ciquiéme l'id de la pair qui contient le fichier

*/
File * look(char* cond[5]){  
    File * f=malloc(MAX_LENGTH * sizeof(File)) ;
    File files[MAX_LENGTH];
    read_config_file(files);
    int j=0;
    for(int  i=0;i< MAX_LENGTH;i++){
           if(strcmp(cond[0], "") == 0 || is_his_name(cond[0], files[i])){
            if(strcmp(cond[1], "") == 0 || is_his_length(atoi(cond[1]), files[i])){
                if(strcmp(cond[2], "") == 0 || files[i].piece_size == atoi(cond[2])){
                    if(strcmp(cond[3], "") == 0 || is_his_key(cond[3], files[i])){
                        if(strcmp(cond[3], "") == 0 ||files[i].peer_id == atoi(cond[5])){
                            f[j] = files[i];
                            j++;
                        }
                        
                    }
                }
            }
        }
    }
    return f;
}

/* get_file cherche le fichier dont la clé est key*/

File * get_file(char * key){
    char* cond[5]={"","","",key,""};
    return look(cond);
}

int main(int argc, char *argv[]){
    
    /* test */

    char* files[] = {"1","file_a.dat","2097152",  "1024", "8905e92afeb80fc7722ec89eb0bf0966", "file_b.dat", "3145728", "1536" ,"330a57722ec8b0bf09669a2b35f88e9e",NULL};
    add_peer(files);
    char* files1[] = {"2","file_a.dat","2097152",  "1024", "8905e92afeb80fc7722ec89eb0bf0966", "file_b.dat", "3145728", "1536" ,"330a57722ec8b0bf09669a2b35f88e9e",NULL};
    add_peer(files1);

   File file[100]={0}; 

    read_config_file(file);

    for (int i = 0; i < 4; i++) {

        printf("Peer ID: %d\n", file[i].peer_id);
        printf("Filename: %s\n", file[i].filename);
        printf("Length: %d\n", file[i].length);
        printf("Block Size: %d\n", file[i].piece_size);
        printf("Key: %s\n", file[i].key);
        printf("\n");
    }
   
    free_files_memory(file,100);
    return 0;

}
