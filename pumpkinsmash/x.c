#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>
#include <sys/stat.h>
#include <sys/sendfile.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <string.h>
#include <stdio.h>

int readfile(char* fname) {
    int fd = -1;
    struct stat fdstat;

    fd = open(fname, O_RDONLY);
    if (fd < 0) {
        perror("open");
        exit(EXIT_FAILURE);
    }

    if (fstat(fd, &fdstat) == -1) {
        perror("fstat");
        exit(EXIT_FAILURE);
    }

    if (sendfile(STDOUT_FILENO, fd, NULL, fdstat.st_size) < 0) {
        perror("write");
        exit(EXIT_FAILURE);
    }

    close(fd);
    return 0;
}

int play() {
    int a;
    int b;
    char buffer[10];
    a = 0x41414141;
    b = 0x42424242;

    if (write(STDOUT_FILENO, "Alright, let's crack into the 'Stack of Terror'… You're up first. What's your move?\n> ", strlen("Alright, let's crack into the 'Stack of Terror'… You're up first. What's your move?\n> ")) < 0) {
        perror("write");
    }

    sleep(1);

    if (read(STDIN_FILENO, &buffer, 12) < 0) {
        perror("read");
    }

    if (a == 31337) {
        system(buffer);
    } else if (b == 42) {
        readfile("flag.0");
    } else if (b == 23) {
        readfile("vuln1.txt");
    } else {
        write(STDOUT_FILENO, "Looks like your spells have drifted into the night. Better luck next haunting!\n", strlen("Looks like your spells have drifted into the night. Better luck next haunting!\n"));
    }

    return 0;
}

int start_server(int port) {
    int server_fd, new_socket;
    struct sockaddr_in address;
    int opt = 1;
    int addrlen = sizeof(address);

    // Creating socket file descriptor
    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("socket failed");
        exit(EXIT_FAILURE);
    }

    // Forcefully attaching socket to the port
    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR | SO_REUSEPORT, &opt, sizeof(opt))) {
        perror("setsockopt");
        exit(EXIT_FAILURE);
    }

    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(port);

    // Bind the socket to the specified port
    if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) < 0) {
        perror("bind failed");
        exit(EXIT_FAILURE);
    }

    // Start listening for incoming connections
    if (listen(server_fd, 3) < 0) {
        perror("listen");
        exit(EXIT_FAILURE);
    }

    printf("Server is listening on port %d...\n", port);

    // Accept a connection
    while ((new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t *)&addrlen)) >= 0) {
        printf("Connection established!\n");

        // Redirect STDIN, STDOUT to the new socket
        dup2(new_socket, STDIN_FILENO);
        dup2(new_socket, STDOUT_FILENO);

        // Run the game function
        play();

        // Close the socket
        close(new_socket);
    }

    if (new_socket < 0) {
        perror("accept");
        exit(EXIT_FAILURE);
    }

    return 0;
}

int main(int argc, char *argv[]) {
    int port = 2000;  // Default port
    if (argc > 1) {
        port = atoi(argv[1]);
    }
    start_server(port);
    exit(EXIT_SUCCESS);
}
