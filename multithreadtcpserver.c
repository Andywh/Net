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
#include <pthread.h>

#define BUFSIZE 512
#define MAXLINE 1024

static void* doit(void *arg)
{
    ssize_t n;
    char    buf[MAXLINE];
    
    pthread_detach(pthread_self());
again:
    while ( (n = read((int)arg, buf, MAXLINE)) > 0)
        write((int)arg, buf, n);
    
    if (n < 0 && errno == EINTR)
        goto again;
    else if (n < 0)
        printf("str_echo: read error");
    
    close((int) arg);
    return (NULL);
}

int main(int argc, char **argv) {
    int     listenfd, connfd;
    //    int     childpid;
    pthread_t tid;
    socklen_t   addrlen, len;
    struct sockaddr *cliaddr;
    
    
    // step 1:设置地址结构体
    struct  sockaddr_in servaddr;
    bzero(&servaddr, sizeof(servaddr));
    servaddr.sin_family = AF_INET;
    servaddr.sin_addr.s_addr = htonl(INADDR_ANY);
    servaddr.sin_port = htons(9877);
    
    //  step 2: 创建套接字
    listenfd = socket(AF_INET, SOCK_STREAM, 0);
     
    //    if (argc == 2)
    //        listenfd = tcp_listen(NULL, argv[1], &addrlen);
    //    else if (argc == 3)
    //        listenfd = tcp_listen(argv[1], argv[2], &addrlen);
    //    else
    //        printf("usage: mytcpserver [<host>] <service or port>");
  
    // step 3: 绑定
    bind(listenfd, (struct sockaddr *) &servaddr, sizeof(servaddr));

    // stpe 4: 监听
    listen(listenfd, 128);
    
    cliaddr = malloc(addrlen);
    for ( ; ; ) {
        //        clilen = sizeof(cliaddr);
        len = addrlen;
        connfd = accept(listenfd, cliaddr, &len);
        pthread_create(&tid, NULL, &doit, (void*)connfd);
        //        close(connfd);
    }
}


