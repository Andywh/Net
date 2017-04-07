//
//  main.c
//  tcpserver01
//
//  Created by MacBookPro on 17/4/7.
//  Copyright © 2017年 Hust. All rights reserved.
//

#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <netdb.h>
#include <sys/types.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <time.h>
#include <errno.h>

#define BUFSIZE 1024
#define MAXLINE 1024

int main(int argc, char **argv) {
    int     listenfd, connfd;
    int     childpid;
    socklen_t   clilen;
    ssize_t n;
    char    buf[MAXLINE];
    struct  sockaddr_in  cliaddr, servaddr;
    
    listenfd = socket(AF_INET, SOCK_STREAM, 0);
    
    bzero(&servaddr, sizeof(servaddr));
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    servaddr.sin_port = htons(9877);
    bind(listenfd, (struct sockaddr *) &servaddr, sizeof(servaddr));
    listen(listenfd, 5);
    
    for ( ; ; ) {
        clilen = sizeof(cliaddr);
        connfd = accept(listenfd, (struct sockaddr *) &cliaddr, &clilen);
        if ( (childpid = fork()) == 0) {
            close(listenfd);
            
        again:
            while ( (n = read(connfd, buf, MAXLINE)) > 0)
                write(connfd, buf, n);
            
            if (n < 0 && errno == EINTR)
                goto again;
            else if (n < 0)
                printf("str_echo: read error");
            
            exit(0);
        }
        close(connfd);
    }
}

