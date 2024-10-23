#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

#define PORT 8080
#define BUFFER_SIZE 40

char* a = char[12]; // Variable to be overwritten

void handle_client(int client_socket) {
    char buffer[BUFFER_SIZE];
    int bytes_received;

    // Receive data from the client
    bytes_received = recv(client_socket, buffer, sizeof(buffer) - 1, 0);
    if (bytes_received < 0) {
        perror("recv failed");
        close(client_socket);
        return;
    }

    // Null terminate the received data
    buffer[bytes_received] = '\0';

    // Allocate memory for the overflow_buffer based on the length of the received data
    char *overflow_buffer = (char *)malloc(bytes_received + 1);
    if (overflow_buffer == NULL) {
        perror("malloc failed");
        close(client_socket);
        return;
    }

    // Copy the received data into the overflow_buffer
    strcpy(overflow_buffer, buffer); // Vulnerable to buffer overflow

    // Optionally, you can use overflow_buffer here for further processing

    // Free the allocated memory
    free(overflow_buffer);

    // Close the client socket
    close(client_socket);
}

void query_a(int client_socket) {
    char response[64];
    snprintf(response, sizeof(response), "Value of a: %sn", a);
    send(client_socket, response, strlen(response), 0);
}

int main() {
    int server_socket, client_socket;
    struct sockaddr_in server_addr, client_addr;
    socklen_t client_addr_len = sizeof(client_addr);

    // Create socket
    server_socket = socket(AF_INET, SOCK_STREAM, 0);
    if (server_socket < 0) {
        perror("socket failed");
        exit(EXIT_FAILURE);
    }

    // Set up the server address structure
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(PORT);

    // Bind the socket
    if (bind(server_socket, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
        perror("bind failed");
        close(server_socket);
        exit(EXIT_FAILURE);
    }

    // Listen for incoming connections
    if (listen(server_socket, 3) < 0) {
        perror("listen failed");
        close(server_socket);
        exit(EXIT_FAILURE);
    }

    printf("Server is listening on port %d\n", PORT);

    while (1) {
        // Accept a new connection
        client_socket = accept(server_socket, (struct sockaddr *)&client_addr, &client_addr_len);
        if (client_socket < 0) {
            perror("accept failed");
            continue;
        }

        // Handle the client in a simple way (you can extend this for API calls)
        handle_client(client_socket);
        query_a(client_socket);
    }

    // Close the server socket
    close(server_socket);
    return 0;
}`