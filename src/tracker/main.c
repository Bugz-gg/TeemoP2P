#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <stdint.h>
#include <signal.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <errno.h>
#include <sys/select.h>
#include <arpa/inet.h>
#include "thpool.h"
#include "tools.h"
#include "tracker.h"
#include "config.h"

#ifndef NB_THREADS
#define NB_THREADS 2
#endif
#define MAX_PEERS 2000
#define TRUE 1
#define FALSE 0

static threadpool thpool;
static sig_t old_handler;
static int sockfd;
int client_socket[MAX_PEERS] = {0};
int bad_attempts[MAX_PEERS] = {0};
Peer *connected_peers[MAX_PEERS];
FILE *log_file;

void close_on_exit(int signo) {
    fclose(log_file);
    thpool_destroy(thpool);
    signal(SIGINT, old_handler);
    close(sockfd);
    free_all_regex();
    //free_on_exit(signo);
    (void) signo;
    for (int i = 0; i < tracker.nb_peers; ++i)
        free_peer(tracker.peers[i]);
    free(tracker.peers);
    for (int i = 0; i < tracker.nb_files; ++i) {
        free_file(tracker.files[i]);
    }
    free(tracker.files);
    exit(0);
}

void error(char *msg) {
    perror(msg);
    exit(1);
}

int handle_message(char *message, Tracker *tracker, char *addr_ip, int socket_fd, int index) {
    if (streq(message, "ping")) {
        write(socket_fd, "pong\n", 5);
        return 0;
    }
    if (streqlim(message, "announce", 8)) {
        announceData aData = announceCheck(message);
        if (aData.is_valid) {
            connected_peers[index] = announce(tracker, &aData, addr_ip, socket_fd);
            free_announceData(&aData);
            print_tracker_files(tracker);
            print_tracker_peers(tracker);
            return 0;
        }
        free_announceData(&aData);
    } else if (streqlim(message, "look", 4)) {
        lookData lData = lookCheck(message);
        if (lData.is_valid) {
            look(tracker, &lData, socket_fd);
            free_lookData(&lData);
            return 0;
        }
        free_lookData(&lData);
    } else if (streqlim(message, "getfile", 7)) {
        getfileData gfData = getfileCheck(message);
        if (gfData.is_valid) {
        getfile(tracker, &gfData, socket_fd);
            return 0;
        }
    } else if (streqlim(message, "update", 6)) {
        updateData uData = updateCheck(message);
        if (uData.is_valid) {
            // Handle data
            free_updateData(&uData);
            return 0;
        }
        free_updateData(&uData);
    }
    return 1;
}

// Fonction pour gérer les connexions clients dans des threads
void handle_client_connection(void *newsockfd_void_ptr) {
    int client_sockfd = (int) (intptr_t) newsockfd_void_ptr;
    char buffer[256] = {0};
    int index = 0;
    int n = 0;
    for (int i = 0; i < MAX_PEERS; ++i) {
        if (client_socket[i] == client_sockfd) {
            index = i;
            break;
        }
    }
    // memset(buffer, 0, 256);
    while (1) {
        n += read(client_sockfd, buffer + n, 255);
        if (n < 0) {
            error("ERROR reading from socket");
            break;
        }
        if (n == 0) {
            // Le client a fermé la connexion
            printf("Client disconnected\n");
            if (client_socket[index] == sockfd) {
                client_socket[index] = 0;
                bad_attempts[index] = 0;
                close(client_sockfd);
            }
            return;
        }
        if (buffer[n - 1] == '\n' || !strcmp(&buffer[n - 2], "\r\n"))
            break;
    }

    buffer[strcspn(buffer, "\r\n")] = 0;
    if (strcmp(buffer, "exit") == 0) {
        printf("Client requested to disconnect.\n");
        if (client_socket[index] == sockfd) {
            client_socket[index] = 0;
            bad_attempts[index] = 0;
            close(client_sockfd);
        }
        return;
    }
    struct sockaddr_in addr;
    socklen_t addr_size = sizeof(struct sockaddr_in);
    getpeername(client_sockfd, (struct sockaddr *) &addr, &addr_size);
    char clientip[MAX_IP_ADDR_SIZE];
    strcpy(clientip, inet_ntoa(addr.sin_addr));
    int port = ntohs(addr.sin_port);
    printf("Message reçu de \033[0;33m%s:%d\033[39m: %s\n", clientip, port, buffer);

    // Vérifie si le message est bien formaté
    int check = handle_message(buffer, &tracker, clientip, client_sockfd, index); // Replace NULL by addr_ip
    if (check == 1) { // Message mal formaté
        ++bad_attempts[index];
        if (bad_attempts[index] >= 3) {
            printf("\033[0;31mMessage mal formaté détecté 3 fois, fermeture de la connexion avec \033[0;33m%s:%d\033[39m.\033[39m\n",
                   clientip, port);
            client_socket[index] = 0;
            close(client_sockfd);
            return;
        }
    } else {
        // Message bien formaté, réinitialiser le compteur d'erreurs
        bad_attempts[index] = 0;
    }
    /*n = write(client_sockfd, "I got your message\n", 19);
    if (n < 0) {
        error("ERROR writing to socket");
        return;
    }*/
}

int main() {
    log_file = open_log();
    config conf = read_config();
    int opt = TRUE;
    old_handler = signal(SIGINT, close_on_exit);

    int sockfd, newsockfd, portno;
    int max_sd, sd, activity;
    socklen_t clilen;
    struct sockaddr_in serv_addr, cli_addr;

    init_tracker(&tracker);
    for (int i=0; i<MAX_PEERS; ++i)
        connected_peers[i] = NULL;

    // Initialiser le pool de threads
    thpool = thpool_init(NB_THREADS);

    sockfd = socket(AF_INET, SOCK_STREAM, 0);

    if (sockfd < 0)
        error("ERROR opening socket");
    memset((char *) &serv_addr, 0, sizeof(serv_addr));
    portno = conf.listen_port; //atoi(argv[1]);
    if (setsockopt(sockfd, SOL_SOCKET, SO_REUSEADDR, (char *) &opt,
                   sizeof(opt)) < 0) {
        error("setsockopt");
    }

    char tracker_ip[INET_ADDRSTRLEN];
    serv_addr.sin_family = AF_INET;
    if (!conf.IP_mode)
        serv_addr.sin_addr.s_addr = inet_addr(conf.IP);
    else if (conf.IP_mode==1)
        serv_addr.sin_addr.s_addr = htonl(INADDR_LOOPBACK);
    else // Et ouais
        serv_addr.sin_addr.s_addr = INADDR_ANY;

    serv_addr.sin_port = htons(portno);
    if (bind(sockfd, (struct sockaddr *) &serv_addr, sizeof(serv_addr)) < 0)
        error("ERROR on binding");

    listen(sockfd, 5);
    clilen = sizeof(cli_addr);
    fd_set readfds;
    inet_ntop(AF_INET, &(serv_addr.sin_addr), tracker_ip, INET_ADDRSTRLEN);
    printf("\033[92mReady on %s:%d\033[39m\n", tracker_ip, portno);
    write_log("Ready on %s:%d\n", tracker_ip, portno);

    while (1) {
        FD_ZERO(&readfds); // Clear the socket set
        FD_SET(sockfd, &readfds); // Add master socket to set
        max_sd = sockfd;

        for (int i = 0; i < MAX_PEERS; i++) { // Add child sockets to set
            sd = client_socket[i]; // socket descriptor

            if (sd > 0) // Only add valid socket descriptors to read list
                FD_SET(sd, &readfds);

            if (sd > max_sd) // Highest file descriptor number for the select function
                max_sd = sd;
        }

        // Wait indefinitely for an activity on one of the sockets
        activity = select(max_sd + 1, &readfds, NULL, NULL, NULL);

        if ((activity < 0) && (errno != EINTR)) {
            printf("select error");
        }

        // Incoming connection
        if (FD_ISSET(sockfd, &readfds)) {
            if ((newsockfd = accept(sockfd,
                                    (struct sockaddr *) &cli_addr, (socklen_t * ) & clilen)) < 0) {
                error("accept");
            }

            // Inform user of socket number - used in send and receive commands
            printf("New connection, socket fd: %d, \033[0;33m%s:%d\033[39m.\033[39m \n", newsockfd,
                   inet_ntoa(cli_addr.sin_addr), ntohs(cli_addr.sin_port));

            /*char *message = "ECHO Daemon v1.0 \r\n";
            // Send new connection greeting message
            if (send(newsockfd, message, strlen(message), 0) != strlen(message)) {
                perror("send");
            }*/

            // Add new socket to array of sockets
            for (int i = 0; i < MAX_PEERS; ++i) {
                if (client_socket[i] == 0) {
                    client_socket[i] = newsockfd;
                    printf("Adding to list of sockets as %d at %d\n", newsockfd, i);
                    break;
                }
            }
        }

        // Handle what happened to the active connections
        for (int i = 0; i < MAX_PEERS; i++) {
            sd = client_socket[i];

            if (sd > 0 && FD_ISSET(sd, &readfds)) {
                // Soumettre la gestion de chaque nouvelle connexion au pool de threads
                thpool_add_work(thpool, (void (*)(void *)) handle_client_connection, (void *) (intptr_t) sd);

            }
        }
        thpool_wait(thpool);
    }
    return 0;
}
