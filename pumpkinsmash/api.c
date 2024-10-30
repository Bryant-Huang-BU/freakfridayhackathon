#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <arpa/inet.h>

#define PORT 3010
#define BUFFER_SIZE 40
#define TEAM_NAME_LENGTH 12
#define MAX_TEAMS 10

char* team_names[MAX_TEAMS];
int num_teams = 0;

void load_teams(const char *filename) {
    FILE *file = fopen(filename, "r");
    if (!file) {
        perror("Failed to open teams file");
        exit(EXIT_FAILURE);
    }

    char line[TEAM_NAME_LENGTH + 2]; // Buffer for team names (up to 12 chars plus newline and null terminator)
    while (fgets(line, sizeof(line), file) && num_teams < MAX_TEAMS) {
        line[strcspn(line, "\n")] = '\0'; // Remove newline character if present
        team_names[num_teams] = strdup(line); // Allocate memory for team name
        //print team names
        printf("Team name: %s\n", team_names[num_teams]);
        if (team_names[num_teams] == NULL) {
            perror("Failed to allocate memory for team name");
            fclose(file);
            exit(EXIT_FAILURE);
        }
        num_teams++;
    }
    fclose(file);
}

int check_team(const char *input) {
    for (int i = 0; i < num_teams; i++) {
        //check in loop til null
        printf("Comparing %s and %s\n", team_names[i], input);
        int len = strlen(input) - 1;
        printf("%d\n", len);
        if (strncmp(input, team_names[i], len) == 0) {
            return i; // Match found
        }
    }
    return -1; // No match found
}

void handle_client(int client_socket) {
    volatile char target[12] = "NOTHING"; // The target array you want to overflow to
    volatile char buffer[BUFFER_SIZE];
    int bytes_received;
    send(client_socket, "Alright, let\'s crack into the \'Stack of Terror\'… You\'re up first, 52 in the chamber, team in the overflow. What\'s your move?\n", strlen("Alright, let\'s crack into the \'Stack of Terror\'… You\'re up first, 40 in the chamber, team in the overflow. What\'s your move?\n"), 0);
    // Vulnerable part: receiving data into a smaller buffer
    bytes_received = recv(client_socket, buffer, BUFFER_SIZE + 40, 0); // Intentionally oversized read
    if (bytes_received < 0) {
        perror("recv failed");
        close(client_socket);
        return;
    }
    buffer[bytes_received] = '\0'; // Null-terminate to avoid string-related issues
    printf("Buffer: %s\n", buffer);
    printf("Target: %s\n", target);
    // Check for team name match with possible overflowed input
    int team_index = check_team(target);
    if (team_index >= 0) {
        FILE *file = fopen("currflag.txt", "w");
        if (!file) {
            perror("Failed to open currflag.txt");
            close(client_socket);
            return;
        }
        fprintf(file, "Team: %s\n", team_names[team_index]);
        fclose(file);
        send(client_socket,"You dare disturb the stack of spirits?! Well, I’ll counter with 'Pumpkin Pulverize' and… SMASH that stack completely!\n", strlen("You dare disturb the stack of spirits?! Well, I’ll counter with 'Pumpkin Pulverize' and… SMASH that stack completely!\n"), 0);
    } else {
        char bufr[256]; // Adjust the size if necessary

    // Format the string into the buffer
        snprintf(bufr, sizeof(bufr), "Looks like your spells have drifted into the night. Better luck next haunting!\n Imagine being TEAM %s", target);
        //printf("has been wokrin");
        // No match found, so send out message and value of target
        send(client_socket, bufr , strlen(bufr), 0);
    }

    close(client_socket);
}

int main() {
    int server_socket, client_socket;
    struct sockaddr_in server_addr, client_addr;
    socklen_t client_addr_len = sizeof(client_addr);

    // Load team names from file
    load_teams("teams.txt");

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

        // Handle the client in a vulnerable way
        handle_client(client_socket);
    }

    // Clean up team names and close the server socket
    for (int i = 0; i < num_teams; i++) {
        free(team_names[i]);
    }
    close(server_socket);
    return 0;
}
