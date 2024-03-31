#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <stdbool.h>
#include "thpool.h" 
#include "tools.h"
#ifndef NB_THREADS
#define NB_THREADS 2
#endif

void error(char *msg) {
    perror(msg);
    exit(1);
}

bool check_message(char* message){
    /* segmentation fault s
    announceData d= announceCheck(message);
    printAnnounceData(d);
    return  d.is_valid ;*/ 
    /* segmentation fault
    lookData d=lookCheck(message);
    printLookData(d);
    return d.is_valid;*/
    announceData aData = announceCheck(message);
    if (aData.is_valid)
        return true;
    lookData lData = lookCheck(message);
    if (lData.is_valid)
        return true;
    getfileData gfData = getfileCheck(message);
    if (gfData.is_valid)
        return true;
    return false;
}

// Fonction pour gérer les connexions clients dans des threads
void handle_client_connection(void* newsockfd_void_ptr) {
    int newsockfd = (int)(intptr_t)newsockfd_void_ptr;
    char buffer[256];
    int error_count=0;
    while (1) {
        memset(buffer, 0, 256);
        int n = read(newsockfd, buffer, 255);
        if (n < 0) {
            error("ERROR reading from socket");
            break; 
        }
        if (n == 0) {
            // Le client a fermé la connexion
            printf("Client disconnected\n");
            break; 
        }
        
        buffer[strcspn(buffer, "\r\n")] = 0;
        if (strcmp(buffer, "exit") == 0) {
            printf("Client requested to disconnect.\n");
            close(newsockfd);
            break;
        }
        // Vérifie si le message est bien formaté
        int is_formatted_correctly = check_message(buffer);
        if (!is_formatted_correctly) {
            // Message mal formaté
            error_count++;
            if (error_count >= 3) {
                // Trois erreurs de suite, fermer la connexion
                printf("Message mal formaté détecté 3 fois, fermeture de la connexion.\n");
                close(newsockfd);
                break;
            }
        } else {
            // Message bien formaté, réinitialiser le compteur d'erreurs
            error_count = 0;
        }
        printf("Here is the message: %s\n", buffer);
        n = write(newsockfd, "I got your message", 18);
        if (n < 0) {
            error("ERROR writing to socket");
            break; 
        }
}}

int main(int argc, char *argv[]) {
    int sockfd, newsockfd, portno;
    socklen_t clilen;
    struct sockaddr_in serv_addr, cli_addr;

    if (argc < 2) {
        fprintf(stderr,"ERROR, no port provided\n");
        exit(1);
    }

    // Initialiser le pool de threads
    threadpool thpool = thpool_init(NB_THREADS);

    sockfd = socket(AF_INET, SOCK_STREAM, 0);
    if (sockfd < 0) 
        error("ERROR opening socket");
    memset((char *) &serv_addr, 0, sizeof(serv_addr));
    portno = atoi(argv[1]);

    serv_addr.sin_family = AF_INET;
    serv_addr.sin_addr.s_addr = INADDR_ANY;
    serv_addr.sin_port = htons(portno);
    if (bind(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr)) < 0)
        error("ERROR on binding");
    
    listen(sockfd,5);
    clilen = sizeof(cli_addr);
    
    while (1) {
        newsockfd = accept(sockfd, (struct sockaddr *) &cli_addr, &clilen);
        if (newsockfd < 0)
            error("ERROR on accept");

        // Soumettre la gestion de chaque nouvelle connexion au pool de threads
        thpool_add_work(thpool, (void (*)(void*))handle_client_connection, (void*)(intptr_t)newsockfd);
    }

    thpool_destroy(thpool);

    close(sockfd);
    return 0;
}
