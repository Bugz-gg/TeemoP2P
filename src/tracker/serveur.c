#include <stdio.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

void error(char *msg)
{
perror(msg);
exit(1);
}

void handle_connexion(int socket){
    while(1){
    char buffer[256];
    memset(buffer, 0, sizeof(buffer));
    int n = read(socket,buffer,255);
    if (n < 0) error("ERROR reading from socket");
    printf("Here is the message: %s\n",buffer);
    n = write(socket,buffer,18);
    if (n < 0) error("ERROR writing to socket");}
    //TO DO
}

int main(int argc, char *argv[])
{
int sockfd, newsockfd, portno, clilen;

struct sockaddr_in serv_addr, cli_addr;


if (argc < 2) {
fprintf(stderr,"ERROR, no port provided\n");
exit(1);
}

sockfd = socket(AF_INET, SOCK_STREAM, 0);
if (sockfd < 0) error("ERROR opening socket");
memset((char *) &serv_addr, 0, sizeof(serv_addr));
portno = atoi(argv[1]);

serv_addr.sin_family = AF_INET;
serv_addr.sin_addr.s_addr = INADDR_ANY;
serv_addr.sin_port = htons(portno);
if (bind(sockfd, (struct sockaddr *) &serv_addr,
sizeof(serv_addr)) < 0)
error("ERROR on binding");
listen(sockfd,5);
clilen = sizeof(cli_addr);

while(1){ 

    newsockfd = accept(sockfd,
    (struct sockaddr *) &cli_addr,
    &clilen);
    if (newsockfd < 0){
        error("ERROR on accept");
    }
    pid_t pid=fork();
    if (pid < 0) {
            error("ERROR on fork");
    }
    else if( pid == 0 ){
        handle_connexion(newsockfd);
        //close(newsockfd);
        exit(0);
    }
    else{
        //close(newsockfd);
    }
}
return 0;
}